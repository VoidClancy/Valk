package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type MetricResult struct {
	NsPerOp     int64   `json:"nsPerOp"`
	MsPerOp     float64 `json:"msPerOp"`
	OpsPerSec   int64   `json:"opsPerSec"`
	BytesPerOp  int64   `json:"bytesPerOp"`
	AllocsPerOp int64   `json:"allocsPerOp"`
}

type DatabaseResults struct {
	Operations map[string]map[string]MetricResult `json:"operations"`
}

type DomainResult struct {
	Name        string                     `json:"name,omitempty"`
	Description string                     `json:"description,omitempty"`
	Databases   map[string]DatabaseResults `json:"databases"`
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

func saveBenchmarkJSON(parsed map[string]map[string]MetricResult) {
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

	dbKey := os.Getenv("BENCH_DB")
	if dbKey == "" {
		dbKey = "sqlite"
	}

	dbRes, exists := domain.Databases[dbKey]
	if !exists || dbRes.Operations == nil {
		dbRes = DatabaseResults{
			Operations: make(map[string]map[string]MetricResult),
		}
	}

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
