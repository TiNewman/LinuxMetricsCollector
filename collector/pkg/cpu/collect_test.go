package cpu

import (
	"errors"
	"os"
	"testing"

	"github.com/prometheus/procfs"
)

type UsageCase struct {
	name     string
	start    procfs.CPUStat
	end      procfs.CPUStat
	expected float32
}

type CollectCase struct {
	name string
}

func TestCollect(t *testing.T) {
	collector := NewCPUCollectorWithoutRepo()

	_, err := collector.Collect()
	if err != nil {
		t.Errorf("Collect method returned an error: %v\n", err.Error())
	}

}

func TestCollectError(t *testing.T) {
	ts := []CollectCase{
		{
			name: "nofiles",
		},
		{
			name: "withCPUInfo",
		},
	}

	for _, tc := range ts {
		t.Run(tc.name, func(t *testing.T) {
			collector := newTestCollector(tc.name)
			_, err := collector.Collect()
			_, ok := err.(*os.PathError)
			if !ok {
				t.Errorf("Test: %v; Unexpected error: %v; Expected os.PathError\n", tc.name, err.Error())
			}
		})
	}
}

func TestCollectInvalidMount(t *testing.T) {
	tc := CollectCase{name: "mount_does_not_exist"}
	collector := newTestCollector(tc.name)
	_, err := collector.Collect()
	e := errors.Unwrap(err)
	_, ok := e.(*os.PathError)
	if !ok {
		t.Errorf("Test: %v; Unexpected error: %v; Expected os.PathError\n", tc.name, err.Error())
	}

}

func TestCalculateUsage(t *testing.T) {

	ts := []UsageCase{
		{
			name:     "case1",
			start:    procfs.CPUStat{User: 108.01, Nice: 0.27, System: 46.68, Idle: 5671.43, Iowait: 2.68, IRQ: 34.18, SoftIRQ: 4.02, Steal: 0, Guest: 0, GuestNice: 0},
			end:      procfs.CPUStat{User: 108.09, Nice: 0.27, System: 46.71, Idle: 5679.27, Iowait: 2.68, IRQ: 34.23, SoftIRQ: 4.02, Steal: 0, Guest: 0, GuestNice: 0},
			expected: 1.3836478,
		},
		{
			name:     "case2",
			start:    procfs.CPUStat{User: 108.44, Nice: 0.27, System: 46.9, Idle: 5710.67, Iowait: 2.68, IRQ: 34.31, SoftIRQ: 4.04, Steal: 0, Guest: 0, GuestNice: 0},
			end:      procfs.CPUStat{User: 108.51, Nice: 0.27, System: 46.92, Idle: 5718.56, Iowait: 2.68, IRQ: 34.32, SoftIRQ: 4.05, Steal: 0, Guest: 0, GuestNice: 0},
			expected: 1.1278195,
		},
		{
			name:     "case3",
			start:    procfs.CPUStat{User: 108.97, Nice: 0.27, System: 47.08, Idle: 5749.96, Iowait: 2.68, IRQ: 34.36, SoftIRQ: 4.06, Steal: 0, Guest: 0, GuestNice: 0},
			end:      procfs.CPUStat{User: 109, Nice: 0.27, System: 47.09, Idle: 5757.9, Iowait: 2.68, IRQ: 34.36, SoftIRQ: 4.07, Steal: 0, Guest: 0, GuestNice: 0},
			expected: 0.5012531,
		},
		{
			name:     "case4",
			start:    procfs.CPUStat{User: 109.39, Nice: 0.27, System: 47.27, Idle: 5789.35, Iowait: 2.69, IRQ: 34.4, SoftIRQ: 4.08, Steal: 0, Guest: 0, GuestNice: 0},
			end:      procfs.CPUStat{User: 109.41, Nice: 0.27, System: 47.28, Idle: 5797.32, Iowait: 2.69, IRQ: 34.41, SoftIRQ: 4.09, Steal: 0, Guest: 0, GuestNice: 0},
			expected: 0.375,
		},
		{
			name:     "case5",
			start:    procfs.CPUStat{User: 110.17, Nice: 0.27, System: 47.68, Idle: 5827.77, Iowait: 2.69, IRQ: 34.8, SoftIRQ: 4.12, Steal: 0, Guest: 0, GuestNice: 0},
			end:      procfs.CPUStat{User: 110.22, Nice: 0.27, System: 47.76, Idle: 5835.43, Iowait: 2.69, IRQ: 34.99, SoftIRQ: 4.12, Steal: 0, Guest: 0, GuestNice: 0},
			expected: 1.6688062,
		},
	}

	for _, tc := range ts {
		// run each case in a sub-test
		t.Run(tc.name, func(t *testing.T) {
			result := calculateUsage(tc.start, tc.end)
			if result != tc.expected {
				t.Errorf("test %v; Expected: %v; Got: %v\n", tc.name, tc.expected, result)
			}
		})
	}

}
