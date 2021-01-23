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

package proc

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func setRuneField(name, value string, setter func(c rune)) error {
	runes := []rune(value)
	if len(runes) != 1 {
		return fmt.Errorf("invalid \"%s\" value \"%s\"", name, value)
	}
	setter(runes[0])
	return nil
}

func setIntField(name, value string, setter func(n int64)) error {
	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid \"%s\" value \"%s\" (%w)", name, value, err)
	}
	setter(n)
	return nil
}

func setUintField(name, value string, setter func(n uint64)) error {
	n, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid \"%s\" value \"%s\" (%w)", name, value, err)
	}
	setter(n)
	return nil
}

func Parse(s string) (*ProcStat, error) {
	stat := &ProcStat{}
	s = strings.TrimSuffix(s, "\n")

	pidEnd := strings.Index(s, " ")
	if pidEnd < 0 {
		return nil, errors.New("\"pid\" field not found")
	}

	err := setIntField("pid", s[:pidEnd], func(v int64) { stat.Pid = int(v) })
	if err != nil {
		return nil, err
	}

	commBegin := strings.Index(s, " (")
	commEnd := strings.LastIndex(s, ") ")
	if commBegin < 0 || commEnd < commBegin {
		return nil, errors.New("\"comm\" field not found")
	}
	stat.Comm = s[commBegin+2 : commEnd]

	rest := s[commEnd+2:]
	fields := strings.Split(rest, " ")

	if err := stat.parseRest(fields); err != nil {
		return nil, err
	}

	return stat, nil
}

func Duration(v uint64) time.Duration {
	tck := scClkTck()
	return time.Second * time.Duration(v) / time.Duration(tck)
}
