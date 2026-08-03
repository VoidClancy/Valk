package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"maps"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

var benchLine = regexp.MustCompile(`^Benchmark(\w+)/(\w+)-(\d+)\s+(\d+)\s+([0-9]+(?:\.[0-9]+)?)\s+(ns|µs|us|ms)/op\s+(\d+)\s+B/op\s+(\d+)\s+allocs/op`)

type benchSample struct {
	ns          int64
	bytesPerOp  int64
	allocsPerOp int64
	procs       int
}

type MetricResult struct {
	MsPerOp     float64 `json:"msPerOp"`
	OpsPerSec   int64   `json:"opsPerSec"`
	BytesPerOp  int64   `json:"bytesPerOp"`
	AllocsPerOp int64   `json:"allocsPerOp"`
}

type DatabaseResults struct {
	Meta       DatabaseInfo                       `json:"meta"`
	Machine    MachineInfo                        `json:"machine"`
	Benchmark  BenchmarkConfig                    `json:"benchmark"`
	Operations map[string]map[string]MetricResult `json:"operations"`
}

type DomainResult struct {
	Name        string                     `json:"name,omitempty"`
	Description string                     `json:"description,omitempty"`
	Databases   map[string]DatabaseResults `json:"databases"`
}

type MachineInfo struct {
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	GoVersion string `json:"goVersion"`
	CPU       string `json:"cpu"`
	Cores     int    `json:"cores"`
}

type DatabaseInfo struct {
	Name     string `json:"name"`
	Driver   string `json:"driver,omitempty"`
	Location string `json:"location"`
	Storage  string `json:"storage"`
	Endpoint string `json:"endpoint,omitempty"`
}

type BenchmarkConfig struct {
	BenchTime string `json:"benchTime"`
	Count     int    `json:"count"`
	Aggregate string `json:"aggregate"`
}

type BenchmarkReport struct {
	UpdatedAt string                  `json:"updatedAt"`
	Domains   map[string]DomainResult `json:"domains"`
}

func getReportPath() string {
	if out := os.Getenv("BENCH_OUT"); out != "" {
		return out
	}

	dir, err := os.Getwd()
	if err != nil {
		return "benchmarks.json"
	}

	current := dir
	for {
		if filepath.Base(current) == "benchmark" {
			return filepath.Join(current, "benchmarks.json")
		}
		benchDir := filepath.Join(current, "benchmark")
		if info, err := os.Stat(benchDir); err == nil && info.IsDir() {
			return filepath.Join(benchDir, "benchmarks.json")
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}

	return filepath.Join(dir, "benchmarks.json")
}

func saveBenchmarkJSON(parsed map[string]map[string]MetricResult, machine MachineInfo, cfg BenchmarkConfig) {
	if len(parsed) == 0 {
		return
	}

	targetPath := getReportPath()
	var report BenchmarkReport

	if data, err := os.ReadFile(targetPath); err == nil {
		_ = json.Unmarshal(data, &report)
	}

	if report.Domains == nil {
		report.Domains = make(map[string]DomainResult)
	}

	domainKey := os.Getenv("BENCH_DOMAIN")
	if domainKey == "" {
		domainKey = "orm"
	}

	domain, exists := report.Domains[domainKey]
	if !exists {
		domainName := "ORM Performance"
		if domainKey != "orm" {
			domainName = strings.Title(domainKey) + " Benchmarks"
		}
		domain = DomainResult{
			Name:        domainName,
			Description: fmt.Sprintf("Performance benchmarks for %s", domainKey),
			Databases:   make(map[string]DatabaseResults),
		}
	}
	if domain.Databases == nil {
		domain.Databases = make(map[string]DatabaseResults)
	}

	dbInfo := collectDatabaseInfo()
	dbKey := dbInfo.Name
	if dbKey == "" {
		dbKey = "sqlite"
	}

	dbRes, exists := domain.Databases[dbKey]
	if !exists || dbRes.Operations == nil {
		dbRes = DatabaseResults{
			Operations: make(map[string]map[string]MetricResult),
		}
	}
	dbRes.Meta = dbInfo
	dbRes.Machine = machine
	dbRes.Benchmark = cfg

	for opName, targets := range parsed {
		if dbRes.Operations[opName] == nil {
			dbRes.Operations[opName] = make(map[string]MetricResult)
		}
		maps.Copy(dbRes.Operations[opName], targets)
	}

	domain.Databases[dbKey] = dbRes
	report.Domains[domainKey] = domain
	report.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create directory for benchmarks.json: %v\n", err)
		return
	}

	outData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal benchmarks.json: %v\n", err)
		return
	}

	if err := os.WriteFile(targetPath, outData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write benchmarks.json: %v\n", err)
		return
	}

	fmt.Printf("\n── Benchmark Report ──────────────────────────────────\nSaved benchmark results to %s\n", targetPath)
}

func extractAggregateFlag(args []string) ([]string, string) {
	aggregate := "median"
	filtered := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case strings.HasPrefix(a, "-aggregate="):
			aggregate = strings.TrimPrefix(a, "-aggregate=")
		case a == "-aggregate" && i+1 < len(args):
			aggregate = args[i+1]
			i++
		default:
			filtered = append(filtered, a)
		}
	}
	switch aggregate {
	case "best", "last":
	default:
		aggregate = "median"
	}
	return filtered, aggregate
}

func aggregateSamples(samples map[string]map[string][]benchSample, mode string) map[string]map[string]MetricResult {
	parsed := make(map[string]map[string]MetricResult, len(samples))
	for op, orms := range samples {
		if parsed[op] == nil {
			parsed[op] = make(map[string]MetricResult)
		}
		for orm, list := range orms {
			groups := groupByProcs(list)
			for _, g := range groups {
				parsed[op][orm] = aggregateOne(g, mode)
			}
		}
	}
	return parsed
}

func groupByProcs(list []benchSample) [][]benchSample {
	var groups [][]benchSample
	index := make(map[int]int)
	for _, s := range list {
		if idx, ok := index[s.procs]; ok {
			groups[idx] = append(groups[idx], s)
		} else {
			index[s.procs] = len(groups)
			groups = append(groups, []benchSample{s})
		}
	}
	return groups
}

func aggregateOne(list []benchSample, mode string) MetricResult {
	sorted := append([]benchSample(nil), list...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].ns < sorted[j].ns })

	var chosen benchSample
	switch mode {
	case "best":
		chosen = sorted[0]
	case "last":
		chosen = list[len(list)-1]
	default: // median
		chosen = sorted[len(sorted)/2]
	}

	ms := float64(chosen.ns) / 1_000_000
	ops := int64(0)
	if chosen.ns > 0 {
		ops = 1_000_000_000 / chosen.ns
	}
	return MetricResult{
		MsPerOp:     ms,
		OpsPerSec:   ops,
		BytesPerOp:  chosen.bytesPerOp,
		AllocsPerOp: chosen.allocsPerOp,
	}
}

func toNanos(value float64, unit string) int64 {
	switch unit {
	case "ms":
		return int64(value * 1_000_000)
	case "µs", "us":
		return int64(value * 1_000)
	default:
		return int64(value)
	}
}

func collectMachineInfo() MachineInfo {
	goVersion := ""
	if v, err := exec.Command("go", "version").Output(); err == nil {
		goVersion = strings.TrimSpace(string(v))
	}

	return MachineInfo{
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		GoVersion: goVersion,
		CPU:       cpuModel(),
		Cores:     runtime.NumCPU(),
	}
}

func cpuModel() string {
	if runtime.GOOS != "linux" {
		return runtime.GOARCH + " CPU"
	}

	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return runtime.GOARCH + " CPU"
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if k, v, ok := strings.Cut(sc.Text(), ":"); ok && strings.TrimSpace(k) == "model name" {
			if model := strings.TrimSpace(v); model != "" {
				return model
			}
		}
	}
	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to read /proc/cpuinfo: %v\n", err)
	}

	return runtime.GOARCH + " CPU"
}

func collectBenchmarkConfig(args []string, aggregate string) BenchmarkConfig {
	flags := flagMap(args)
	return BenchmarkConfig{
		BenchTime: flagValue(flags, "benchtime", "1s"),
		Count:     atoiFlag(flags, "count", 1),
		Aggregate: aggregate,
	}
}

func flagMap(args []string) map[string]string {
	out := make(map[string]string)
	for _, a := range args {
		if !strings.HasPrefix(a, "-") {
			continue
		}
		name, value, has := strings.Cut(strings.TrimLeft(a, "-"), "=")
		if name == "" {
			continue
		}
		if !has {
			value = "true"
		}
		out[name] = value
	}
	return out
}

func flagValue(flags map[string]string, name, def string) string {
	if v, ok := flags[name]; ok {
		return v
	}
	return def
}

func atoiFlag(flags map[string]string, name string, def int) int {
	if v, ok := flags[name]; ok {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func collectDatabaseInfo() DatabaseInfo {
	dbName := strings.ToLower(os.Getenv("BENCH_DB"))
	if dbName == "" {
		dbName = "sqlite"
	}
	info := DatabaseInfo{Name: dbName, Location: "local"}

	switch dbName {
	case "postgres", "postgresql":
		info.Name = "postgres"
		info.Driver = "postgres"
		info.Storage = "server"
		host, port := "localhost", "5432"
		if dsn := os.Getenv("PG_DATABASE_URL"); dsn != "" {
			if u, err := url.Parse(dsn); err == nil && u.Hostname() != "" {
				host = u.Hostname()
				if p := u.Port(); p != "" {
					port = p
				}
			}
		}
		info.Endpoint = net.JoinHostPort(host, port)
		if !isLocalHost(host) {
			info.Location = "remote"
		}
	case "sqlite":
		info.Driver = "sqlite3"
		dsn := os.Getenv("BENCH_SQLITE_DSN")
		if dsn == "" {
			dsn = "file:benchmark?mode=memory&cache=shared&_fk=1"
		}
		info.Endpoint = dsn
		if strings.Contains(dsn, "mode=memory") {
			info.Storage = "memory"
		} else {
			info.Storage = "file"
		}
	default:
		info.Driver = dbName
		info.Storage = "unknown"
	}
	return info
}

func isLocalHost(host string) bool {
	switch strings.ToLower(host) {
	case "localhost", "127.0.0.1", "::1", "0.0.0.0":
		return true
	}
	return false
}
