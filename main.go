// Copyright 2021 Chris Thunes
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Thread struct {
	Header  string
	Name    string
	CPU     time.Duration
	Elapsed time.Duration
	TID     string
	NID     string
	Stack   string
}

func getHeaderField(line string, name string) (string, error) {
	startMarker := name + "="
	startIdx := strings.Index(line, startMarker)
	if startIdx < 0 {
		return "", fmt.Errorf("field not found: %s", name)
	}

	suffix := line[startIdx+len(startMarker):]
	endIdx := strings.Index(suffix, " ")
	if endIdx < 0 {
		return suffix, nil
	}
	return suffix[:endIdx], nil
}

func parseThread(lines []string) (*Thread, error) {
	header := lines[0]
	stack := strings.Join(lines[1:], "\n")

	// Extract name and remove quotes
	startName := 1
	endName := strings.LastIndexByte(header, '"')
	if endName <= 0 {
		return nil, errors.New("invalid thread name (no closing quote)")
	}
	name := header[startName:endName]

	// Extract other header fields
	cpuString, err := getHeaderField(header, "cpu")
	if err != nil {
		return nil, err
	}

	cpu, err := time.ParseDuration(cpuString)
	if err != nil {
		return nil, err
	}

	elapsedString, err := getHeaderField(header, "elapsed")
	if err != nil {
		return nil, err
	}

	elapsed, err := time.ParseDuration(elapsedString)
	if err != nil {
		return nil, err
	}

	tid, err := getHeaderField(header, "tid")
	if err != nil {
		return nil, err
	}

	nid, err := getHeaderField(header, "nid")
	if err != nil {
		return nil, err
	}

	thread := &Thread{
		Header:  header,
		Name:    name,
		CPU:     cpu,
		Elapsed: elapsed,
		TID:     tid,
		NID:     nid,
		Stack:   stack,
	}

	return thread, nil
}

func ParseThreads(text string) (map[string]*Thread, error) {
	threads := make(map[string]*Thread)
	var thread []string
	for _, l := range strings.Split(text, "\n") {
		if len(l) > 0 && (l[0] == ' ' || l[0] == '\t') {
			if len(thread) > 0 {
				thread = append(thread, l)
			}
		} else {
			if len(thread) > 0 {
				parsed, err := parseThread(thread)
				if err != nil {
					return nil, err
				}
				threads[parsed.TID] = parsed
				thread = nil
			}
			if len(l) > 0 && l[0] == '"' && strings.Contains(l, "nid=") {
				thread = append(thread, l)
			}
		}
	}
	if len(thread) > 0 {
		parsed, err := parseThread(thread)
		if err != nil {
			return nil, err
		}
		threads[parsed.TID] = parsed
	}
	return threads, nil
}

func printHeader(cpuFrac float64, header string) {
	fi, err := os.Stdout.Stat()
	if err != nil || fi.Mode()&os.ModeCharDevice == 0 {
		// Error checking status or non-TTY
		fmt.Printf("%.6f\t%s\n", cpuFrac, header)
	} else {
		fmt.Printf("[%6.2f%%] %s\n", 100*cpuFrac, header)
	}
}

func printTopThreads(dump0, dump1 string, n int, summary bool) error {
	threads0, err := ParseThreads(dump0)
	if err != nil {
		panic(err)
	}

	threads1, err := ParseThreads(dump1)
	if err != nil {
		panic(err)
	}

	type withCPUFrac struct {
		frac   float64
		thread *Thread
	}

	var totalCPU time.Duration
	var maxElapsed time.Duration
	var top []*withCPUFrac
	for tid, t1 := range threads1 {
		var cpu time.Duration
		var elapsed time.Duration

		if t0, ok := threads0[tid]; ok {
			cpu = t1.CPU - t0.CPU
			elapsed = t1.Elapsed - t0.Elapsed
		} else {
			cpu = t1.CPU
			elapsed = t1.Elapsed
		}

		totalCPU += cpu
		if elapsed > maxElapsed {
			maxElapsed = elapsed
		}

		frac := float64(cpu) / float64(elapsed)
		top = append(top, &withCPUFrac{frac, t1})
	}

	// Sort by CPU time in descending order
	sort.Slice(top, func(i, j int) bool {
		if top[i].frac == top[j].frac {
			return top[i].thread.TID < top[j].thread.TID
		}
		return top[i].frac > top[j].frac
	})

	if n < 0 {
		n = len(top)
	}

	for _, t := range top[:n] {
		printHeader(t.frac, t.thread.Header)

		if !summary {
			if len(t.thread.Stack) > 0 {
				fmt.Println(t.thread.Stack)
			}
			fmt.Println()
		}
	}

	totalFrac := float64(totalCPU) / float64(maxElapsed)
	printHeader(totalFrac, fmt.Sprintf("Total (elapsed %s)", maxElapsed))

	return nil
}

func parseJavaPID(s string) (int, error) {
	// Try as pid and then search with JPS
	pid, err := strconv.Atoi(s)
	if err == nil {
		return pid, nil
	}

	cmd := exec.Command("jps", "-l")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	for _, line := range strings.Split(string(out), "\n") {
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 && parts[1] == s {
			return strconv.Atoi(parts[0])
		}
	}

	return 0, fmt.Errorf("no process found matching \"%s\"", s)
}

func jstack(pid int) (string, error) {
	cmd := exec.Command("jstack", strconv.Itoa(pid))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func main() {
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		usage := "usage: %s [flags] [stack0 stack1 | pid | name]\n"
		fmt.Fprintf(out, usage, os.Args[0])
		flag.PrintDefaults()
	}

	topN := -1
	duration := 5 * time.Second
	summary := false

	flag.IntVar(&topN, "n", topN, "limit output to the top `N` threads")
	flag.DurationVar(&duration, "d", duration, "sample duration")
	flag.BoolVar(&summary, "summary", summary, "omit stacks")
	flag.Parse()

	var dump0, dump1 string

	if flag.NArg() == 2 {
		bytes0, err := ioutil.ReadFile(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}

		bytes1, err := ioutil.ReadFile(flag.Arg(1))
		if err != nil {
			log.Fatal(err)
		}

		dump0 = string(bytes0)
		dump1 = string(bytes1)
	} else if flag.NArg() == 1 {
		pid, err := parseJavaPID(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}

		ch0 := make(chan string)
		go func() {
			dump, err := jstack(pid)
			if err != nil {
				log.Fatal(err)
			}
			ch0 <- dump
		}()

		ch1 := make(chan string)
		go func() {
			time.Sleep(duration)
			dump, err := jstack(pid)
			if err != nil {
				log.Fatal(err)
			}
			ch1 <- dump
		}()

		dump0 = <-ch0
		dump1 = <-ch1
	} else {
		fmt.Fprint(os.Stderr, "invalid arguments\n\n")
		flag.Usage()
		os.Exit(1)
	}

	if err := printTopThreads(dump0, dump1, topN, summary); err != nil {
		log.Fatal(err)
	}

	// TODO: Add option for "since process start"
	// TODO: Add option to read data from /proc
	// TODO: Make single dump mode the default
}
