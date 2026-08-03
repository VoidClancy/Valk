package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	defaultArgs := []string{"test", "-bench=.", "-benchmem", "-count=1", "-timeout=30m"}
	extraArgs := os.Args[1:]
	args := append(defaultArgs, extraArgs...)

	args, aggregate := extractAggregateFlag(args)

	started := time.Now()
	machine := collectMachineInfo()
	cfg := collectBenchmarkConfig(args, aggregate)

	cmd := exec.Command("go", args...)
	cmd.Dir = "."
	cmd.Stderr = os.Stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	benchSamples := make(map[string]map[string][]benchSample)
	scanner := bufio.NewScanner(stdout)
	var lastOp string
	for scanner.Scan() {
		line := scanner.Text()
		if m := benchLine.FindStringSubmatch(line); m != nil {
			op := m[1]
			if op != lastOp {
				fmt.Printf("\n── %s ────────────────────────────────────\n", op)
				lastOp = op
			}
			orm := m[2]
			procs, _ := strconv.Atoi(m[3])
			value, _ := strconv.ParseFloat(m[5], 64)
			ns := toNanos(value, m[6])
			ms := float64(ns) / 1_000_000
			ops := int64(0)
			if ns > 0 {
				ops = 1_000_000_000 / ns
			}
			bPerOp, _ := strconv.ParseInt(m[7], 10, 64)
			allocs, _ := strconv.ParseInt(m[8], 10, 64)
			fmt.Printf("%-15s  %8.3f ms/op  %8d ops/s  %7d B/op  %4d allocs/op\n", orm, ms, ops, bPerOp, allocs)

			if benchSamples[op] == nil {
				benchSamples[op] = make(map[string][]benchSample)
			}
			benchSamples[op][orm] = append(benchSamples[op][orm], benchSample{
				ns:          ns,
				bytesPerOp:  bPerOp,
				allocsPerOp: allocs,
				procs:       procs,
			})
		} else if strings.HasPrefix(line, "ok ") || strings.HasPrefix(line, "FAIL") || strings.HasPrefix(line, "? ") {
			fmt.Println(line)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	cmd.Wait()

	parsedOperations := aggregateSamples(benchSamples, aggregate)

	fmt.Printf("\nTotal run time: %s\n", time.Since(started).Round(time.Millisecond))
	saveBenchmarkJSON(parsedOperations, machine, cfg)
}
