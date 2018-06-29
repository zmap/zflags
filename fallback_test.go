package flags

import (
	"os"
	"reflect"
	"testing"
	"time"
)

// Test that when `long` is absent, it attempts to fall back to `json`, and
// when `env` is absent, it attempts to fall back to `long` (and to `json`)
func TestFallback(t *testing.T) {
	type Options struct {
		Int   int            `long:"int" json:"json-int" default:"1"`
		Time  time.Duration  `json:"time" default:"1m"`
		Map   map[string]int `json:"map,omitempty" default:"a:1" env-delim:";"`
		Slice []int          `long:"slice" default:"1" default:"2" env:"OVERRIDE_SLICE" env-delim:","`
	}

	var tests = []struct {
		msg      string
		args     []string
		expected Options
		env      map[string]string
	}{
		{
			msg:  "JSON override",
			args: []string{},
			expected: Options{
				Int: 23,
				Time: time.Minute * 3,
				Map: map[string]int{"key1": 1},
				Slice: []int{3,4,5},
			},
			env: map[string]string{
				// since both `json` and `long` are present, `long` ("int") wins
				"json-int": "4",
				"int": "23",
				"time": "3m",
				"map": "key1:1",
				// since both `env` and `long` are present, `env` ("OVERRIDE_SLICE") wins
				"slice": "3,2,1",
				"OVERRIDE_SLICE": "3,4,5",
			},
		},
		{
			msg:  "no arguments, no env, expecting default values",
			args: []string{},
			expected: Options{
				Int:   1,
				Time:  time.Minute,
				Map:   map[string]int{"a": 1},
				Slice: []int{1, 2},
			},
		},
		{
			msg:  "no arguments, env defaults, expecting env default values",
			args: []string{},
			expected: Options{
				Int:   2,
				Time:  2 * time.Minute,
				Map:   map[string]int{"a": 2, "b": 3},
				Slice: []int{4, 5, 6},
			},
			env: map[string]string{
				"Int": "2",
				"Time": "2m",
				"Map": "a:2;b:3",
				"OVERRIDE_SLICE": "4,5,6",
			},
		},
		{
			msg:  "non-zero value arguments, expecting overwritten arguments",
			args: []string{"--int=3", "--time=3ms", "--map=c:3", "--slice=3", "--map=d:4", "--slice=1"},
			expected: Options{
				Int:   3,
				Time:  3 * time.Millisecond,
				Map:   map[string]int{"c": 3,"d":4},
				Slice: []int{3,1},
			},
			env: map[string]string{
				"Int": "2",
				"Time": "2m",
				"Map": "a:2;b:3",
				"OVERRIDE_SLICE": "4,5,6",
			},
		},
		{
			msg:  "zero value arguments, expecting overwritten arguments",
			args: []string{"--int=0", "--time=0ms", "--map=:0", "--slice=0"},
			expected: Options{
				Int:   0,
				Time:  0,
				Map:   map[string]int{"": 0},
				Slice: []int{0},
			},
			env: map[string]string{
				"Int": "2",
				"Time": "2m",
				"Map": "a:2;b:3",
				"OVERRIDE_SLICE": "4,5,6",
			},
		},
		{
			msg:  "`long` used for env name even though `env` was present",
			args: []string{"--int=0", "--time=0ms", "--map=:0", "--slice=0"},
			expected: Options{
				Int:   0,
				Time:  0,
				Map:   map[string]int{"": 0},
				Slice: []int{0},
			},
			env: map[string]string{
				"Int": "2",
				"Time": "2m",
				"Map": "a:2;b:3",
				"Slice": "4,5,6",
			},
		},
	}

	oldEnv := EnvSnapshot()
	defer oldEnv.Restore()

	for _, test := range tests {
		var opts Options
		oldEnv.Restore()
		for envKey, envValue := range test.env {
			os.Setenv(envKey, envValue)
		}
		_, _, _, err := NewParser(&opts, Default | EnvironmentFallback).ParseCommandLine(test.args)
		if err != nil {
			t.Fatalf("%s:\nUnexpected error: %v", test.msg, err)
		}

		if opts.Slice == nil {
			opts.Slice = []int{}
		}

		if !reflect.DeepEqual(opts, test.expected) {
			t.Errorf("%s:\nUnexpected options with arguments %+v\nexpected\n%+v\nbut got\n%+v\n", test.msg, test.args, test.expected, opts)
		}
	}
}

// Test that when `long` is absent, it attempts to fall back to `json`, and
// when `env` is absent, it attempts to fall back to `long` (and to `json`)
func TestNoFallback(t *testing.T) {
	type Options struct {
		Int   int            `long:"int" json:"json-int" default:"1"`
		Time  time.Duration  `json:"time" default:"1m"`
		Map   map[string]int `json:"map,omitempty" default:"a:1" env-delim:";"`
		Slice []int          `long:"slice" default:"1" default:"2" env:"OVERRIDE_SLICE" env-delim:","`
	}

	var tests = []struct {
		msg      string
		args     []string
		expected Options
		env      map[string]string
	}{
		{
			msg:  "JSON override",
			args: []string{},
			expected: Options{
				Int: 1,
				Time: time.Minute * 1,
				Map: map[string]int{"a": 1},
				Slice: []int{3,4,5},
			},
			env: map[string]string{
				"json-int": "4",
				"int": "23",
				"time": "3m",
				"map": "key1:1",
				"slice": "3,2,1",
				"OVERRIDE_SLICE": "3,4,5",
			},
		},
		{
			msg:  "no arguments, no env, expecting default values",
			args: []string{},
			expected: Options{
				Int:   1,
				Time:  time.Minute,
				Map:   map[string]int{"a": 1},
				Slice: []int{1, 2},
			},
		},
		{
			msg:  "no arguments, env defaults, no fallback, expecting default values (except slice)",
			args: []string{},
			expected: Options{
				Int:   1,
				Time:  time.Minute,
				Map:   map[string]int{"a": 1},
				Slice: []int{4, 5, 6},
			},
			env: map[string]string{
				"Int": "2",
				"Time": "2m",
				"Map": "a:2;b:3",
				"OVERRIDE_SLICE": "4,5,6",
			},
		},
		{
			msg:  "non-zero value arguments, expecting overwritten arguments",
			args: []string{"--int=3", "--time=3ms", "--map=c:3", "--slice=3", "--map=d:4", "--slice=1"},
			expected: Options{
				Int:   3,
				Time:  3 * time.Millisecond,
				Map:   map[string]int{"c": 3,"d":4},
				Slice: []int{3,1},
			},
			env: map[string]string{
				"Int": "2",
				"Time": "2m",
				"Map": "a:2;b:3",
				"OVERRIDE_SLICE": "4,5,6",
			},
		},
		{
			msg:  "zero value arguments, expecting overwritten arguments",
			args: []string{"--int=0", "--time=0ms", "--map=:0", "--slice=0"},
			expected: Options{
				Int:   0,
				Time:  0,
				Map:   map[string]int{"": 0},
				Slice: []int{0},
			},
			env: map[string]string{
				"Int": "2",
				"Time": "2m",
				"Map": "a:2;b:3",
				"OVERRIDE_SLICE": "4,5,6",
			},
		},
	}

	oldEnv := EnvSnapshot()
	defer oldEnv.Restore()

	for _, test := range tests {
		var opts Options
		oldEnv.Restore()
		for envKey, envValue := range test.env {
			os.Setenv(envKey, envValue)
		}
		_, _, _, err := NewParser(&opts, Default & (^EnvironmentFallback)).ParseCommandLine(test.args)
		if err != nil {
			t.Fatalf("%s:\nUnexpected error: %v", test.msg, err)
		}

		if opts.Slice == nil {
			opts.Slice = []int{}
		}

		if !reflect.DeepEqual(opts, test.expected) {
			t.Errorf("%s:\nUnexpected options with arguments %+v\nexpected\n%+v\nbut got\n%+v\n", test.msg, test.args, test.expected, opts)
		}
	}
}