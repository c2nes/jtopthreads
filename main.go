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
	"math/bits"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/c2nes/jtopthreads/internal/proc"
)

// Combined output of jstack with data collected from /proc.
type StackDump struct {
	Text      string
	ProcStats map[int]string
	Uptime    time.Duration
}

type Thread struct {
	Header  string
	Name    string
	CPU     time.Duration
	Elapsed time.Duration
	TID     string
	NID     int
	Stack   string
}

func getHeaderField(line string, name string) string {
	startMarker := name + "="
	startIdx := strings.Index(line, startMarker)
	if startIdx < 0 {
		return ""
	}

	suffix := line[startIdx+len(startMarker):]
	endIdx := strings.Index(suffix, " ")
	if endIdx < 0 {
		endIdx = len(suffix)
	}
	return suffix[:endIdx]
}

func (dump *StackDump) parseThread(lines []string) (*Thread, error) {
	header := lines[0]
	stack := strings.Join(lines[1:], "\n")

	// Extract name and remove quotes
	startName := 1
	endName := strings.LastIndexByte(header, '"')
	if endName <= 0 {
		return nil, errors.New("invalid thread name (no closing quote)")
	}
	name := header[startName:endName]

	tid := getHeaderField(header, "tid")
	if tid == "" {
		return nil, errors.New("tid= field missing from header")
	}

	nidString := getHeaderField(header, "nid")
	if nidString == "" {
		return nil, errors.New("nid= field missing from header")
	}

	nid64, err := strconv.ParseInt(nidString, 0, bits.UintSize)
	if err != nil {
		return nil, fmt.Errorf("unable to parse nid: %w", err)
	}
	nid := int(nid64)

	// Parse prop data for thread if we have it
	var stat *proc.ProcStat
	if procString, ok := dump.ProcStats[nid]; ok {
		stat, err = proc.Parse(procString)
		if err != nil {
			return nil, fmt.Errorf("error parsing /proc/[pid]/task/[tid]/stat: %w", err)
		}
	}

	// In recent Java versions stack dumps include CPU and elapsed time. We use
	// this data if it is available, but will fall back to using data from /proc
	// (again, if available).

	var cpu, elapsed time.Duration

	cpuString := getHeaderField(header, "cpu")
	if cpuString != "" {
		cpu, err = time.ParseDuration(cpuString)
		if err != nil {
			return nil, err
		}
	} else if stat != nil {
		cpu = proc.Duration(stat.Utime + stat.Stime)
	} else {
		cpu = 0 * time.Second
	}

	elapsedString := getHeaderField(header, "elapsed")
	if elapsedString != "" {
		elapsed, err = time.ParseDuration(elapsedString)
		if err != nil {
			return nil, err
		}
	} else if stat != nil {
		elapsed = dump.Uptime - proc.Duration(stat.Starttime)
	} else {
		elapsed = 0 * time.Second
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

func (dump *StackDump) ParseThreads() (map[string]*Thread, error) {
	threads := make(map[string]*Thread)
	var thread []string
	for _, l := range strings.Split(dump.Text, "\n") {
		if len(l) > 0 && (l[0] == ' ' || l[0] == '\t') {
			if len(thread) > 0 {
				thread = append(thread, l)
			}
		} else {
			if len(thread) > 0 {
				parsed, err := dump.parseThread(thread)
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
		parsed, err := dump.parseThread(thread)
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

func printTopThreads(dump0, dump1 *StackDump, n int, summary bool) error {
	threads0, err := dump0.ParseThreads()
	if err != nil {
		panic(err)
	}

	threads1, err := dump1.ParseThreads()
	if err != nil {
		panic(err)
	}

	type withCPUFrac struct {
		frac    float64
		thread  *Thread
		cpu     time.Duration
		elapsed time.Duration
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
		top = append(top, &withCPUFrac{frac, t1, cpu, elapsed})
	}

	// Sort by CPU time in descending order
	sort.Slice(top, func(i, j int) bool {
		if top[i].frac == top[j].frac {
			return top[i].thread.TID < top[j].thread.TID
		}
		return top[i].frac > top[j].frac
	})

	if n <= 0 {
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

func collectProcStats(pid int) (map[int]string, error) {
	bytes, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	tasks, err := ioutil.ReadDir(fmt.Sprintf("/proc/%d/task", pid))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	res := make(map[int]string)
	res[pid] = string(bytes)

	// From proc(5), on /proc/[pid]/task
	//
	// "This is a directory that contains one subdirectory for each thread in
	//  the process.  The name of each subdirectory is the numerical thread ID
	//  ([tid]) of the thread (see gettid(2))."
	for _, task := range tasks {
		tid, err := strconv.Atoi(task.Name())
		// Ignore any invalid directory entries
		if err != nil {
			continue
		}

		// We already read this one
		if tid == pid {
			continue
		}

		bytes, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/task/%d/stat", pid, tid))
		if err != nil {
			if os.IsNotExist(err) {
				continue
			} else {
				return nil, err
			}
		}
		res[tid] = string(bytes)
	}

	return res, nil
}

func readProcUptime() (time.Duration, error) {
	v, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}

	var up, idle float64
	_, err = fmt.Sscanf(string(v), "%f %f", &up, &idle)
	if err != nil {
		return 0, err
	}

	return time.Duration(up * float64(time.Second)), nil
}

func jstack(pid int) (*StackDump, error) {
	// Read /proc/uptime so we can calculate elapsed process time
	uptime, err := readProcUptime()
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	// Collect process stats
	procResCh := make(chan map[int]string, 1)
	procErrCh := make(chan error, 1)
	go func() {
		res, err := collectProcStats(pid)
		procResCh <- res
		procErrCh <- err
	}()

	cmd := exec.Command("jstack", strconv.Itoa(pid))
	out, err := cmd.CombinedOutput()
	// Wait for go routine to complete before returning any errors
	procErr := <-procErrCh
	if err != nil {
		return nil, err
	}
	if procErr != nil {
		return nil, procErr
	}
	return &StackDump{string(out), <-procResCh, uptime}, nil
}

func main() {
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		usage := "usage: %s [options] <stack-file> [stack-file]\n"
		usage += "   or: %s [options] [-sample <duration>] <pid | main-class>\n\n"
		fmt.Fprintf(out, usage, os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	usageError := func(format string, a ...interface{}) {
		fmt.Fprintf(os.Stderr, "error: "+format, a...)
		fmt.Fprint(os.Stderr, "\n\n")
		flag.Usage()
		os.Exit(1)
	}

	topN := 0
	duration := time.Duration(0)
	summary := false

	flag.IntVar(&topN, "n", topN, "limit output to the top `N` threads")
	flag.DurationVar(&duration, "sample", duration, "sample process for `duration`")
	flag.BoolVar(&summary, "summary", summary, "omit stacks")
	flag.Parse()

	var dump0, dump1 *StackDump

	if flag.NArg() == 2 {
		if duration > 0 {
			usageError("-sample not supported with file arguments")
		}

		bytes0, err := ioutil.ReadFile(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}

		bytes1, err := ioutil.ReadFile(flag.Arg(1))
		if err != nil {
			log.Fatal(err)
		}

		dump0 = &StackDump{Text: string(bytes0)}
		dump1 = &StackDump{Text: string(bytes1)}
	} else if flag.NArg() == 1 {
		arg := flag.Arg(0)

		// A single argument can be a file, pid or main-class. Process as a file
		// if a matching file exists, otherwise assume the arg is a pid/main-class.
		if _, err := os.Stat(arg); err == nil {
			if duration > 0 {
				usageError("-sample not supported with file argument")
			}

			bytes1, err := ioutil.ReadFile(arg)
			if err != nil {
				log.Fatal(err)
			}

			dump0 = &StackDump{}
			dump1 = &StackDump{Text: string(bytes1)}
		} else {
			pid, err := parseJavaPID(arg)
			if err != nil {
				log.Fatal(err)
			}

			ch0 := make(chan *StackDump)
			go func() {
				dump, err := jstack(pid)
				if err != nil {
					log.Fatal(err)
				}
				ch0 <- dump
			}()

			if duration > 0 {
				ch1 := make(chan *StackDump)
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
				dump0 = &StackDump{}
				dump1 = <-ch0
			}
		}
	} else if flag.NArg() < 1 {
		usageError("argument missing")
	} else {
		usageError("too many arguments")
	}

	if err := printTopThreads(dump0, dump1, topN, summary); err != nil {
		log.Fatal(err)
	}
}
