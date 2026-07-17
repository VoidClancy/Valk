package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var benchLine = regexp.MustCompile(`^Benchmark(\w+)/(\w+)-\d+\s+\d+\s+(\d+)\s+ns/op\s+(\d+)\s+B/op\s+(\d+)\s+allocs/op`)

func main() {
	args := []string{"test", "-bench=.", "-benchmem", "-count=1", "-timeout=30m"}
	args = append(args, os.Args[1:]...)

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
			ns, _ := strconv.ParseInt(m[3], 10, 64)
			ms := float64(ns) / 1_000_000
			ops := 1_000_000_000 / ns
			bPerOp, _ := strconv.ParseInt(m[4], 10, 64)
			allocs, _ := strconv.ParseInt(m[5], 10, 64)
			fmt.Printf("%-15s  %8.3f ms/op  %8d ops/s  %7d B/op  %4d allocs/op\n", orm, ms, ops, bPerOp, allocs)
		} else if strings.HasPrefix(line, "ok ") || strings.HasPrefix(line, "FAIL") || strings.HasPrefix(line, "? ") {
			fmt.Println(line)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	cmd.Wait()
}
