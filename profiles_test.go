// Copyright 2025 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package main

import (
	"os"
	"testing"

	"github.com/google/pprof/profile"
)

func TestMergeProfilesByName(t *testing.T) {
	dir := t.TempDir()
	bs := benchSuite{artDir: dir}
	if err := bs.ensureProfileDirs(); err != nil {
		t.Fatalf("ensure profile dirs: %v", err)
	}
	runMem := bs.profileRunPath(memProfileName)
	runExtra := bs.profileRunPath("mem_BenchmarkFoo.pb.gz")
	mergedMem := bs.profileMergedPath(memProfileName)
	mergedExtra := bs.profileMergedPath("mem_BenchmarkFoo.pb.gz")

	if err := writeTestProfile(runMem, 1); err != nil {
		t.Fatalf("write run mem profile: %v", err)
	}
	if err := writeTestProfile(runExtra, 2); err != nil {
		t.Fatalf("write run extra profile: %v", err)
	}
	if err := writeTestProfile(mergedMem, 3); err != nil {
		t.Fatalf("write merged mem profile: %v", err)
	}

	if err := bs.mergeProfiles(false, true, false); err != nil {
		t.Fatalf("merge profiles: %v", err)
	}
	if _, err := os.Stat(mergedMem); err != nil {
		t.Fatalf("expected merged mem profile: %v", err)
	}
	if _, err := os.Stat(mergedExtra); err != nil {
		t.Fatalf("expected merged extra profile: %v", err)
	}
}

func writeTestProfile(path string, value int64) error {
	p := &profile.Profile{
		TimeNanos:     1,
		DurationNanos: 1,
		SampleType: []*profile.ValueType{
			{Type: "alloc_objects", Unit: "count"},
		},
		Sample: []*profile.Sample{
			{Value: []int64{value}},
		},
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := p.Write(f); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}
