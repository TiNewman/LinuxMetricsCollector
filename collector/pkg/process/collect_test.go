package process

import (
	"errors"
	"os"
	"testing"
)

type CollectCase struct {
	name string
}

func TestCollect(t *testing.T) {
	numProcs := 4
	tc := CollectCase{name: "valid"}
	collector := newTestCollector(tc.name)

	list, err := collector.Collect()
	if err != nil {
		t.Errorf("Collect method returned an error: %v\n", err.Error())
	}
	if len(list) != numProcs {
		t.Errorf("Expected %v processes; Got %v\n", numProcs, len(list))
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

func TestCollectError(t *testing.T) {
	ts := []CollectCase{
		{name: "nostatus"},
		{name: "nostat"},
		{name: "noio"},
	}

	for _, tc := range ts {
		t.Run(tc.name, func(t *testing.T) {
			collector := newTestCollector(tc.name)
			_, e := collector.Collect()
			if e != nil {
				// fmt.Println(e.Error())
				_, ok := e.(*os.PathError)
				if !ok {
					t.Errorf("Test: %v; Unexpected error: %v; Expected os.PathError\n", tc.name, e.Error())
				}
			} else {
				t.Errorf("Test: %v; Expected os.PathError\n", tc.name)
			}
		})
	}
}

type CPUUtilizationCase struct {
	name          string
	cpuTime       float64
	executionTime float64
	expected      float64
}

func TestCalcCPUUtilization(t *testing.T) {
	ts := []CPUUtilizationCase{
		{name: "full utilization", cpuTime: 450.6, executionTime: 450.6, expected: float64(100)},
		{name: "execution zero", cpuTime: 450.6, executionTime: 0, expected: float64(0)},
		{name: "cpu zero", cpuTime: 0, executionTime: 2043.56, expected: float64(0)},
		{name: "all zero", cpuTime: 0, executionTime: 0, expected: float64(0)},
	}
	for _, tc := range ts {
		// run each case in a sub-test
		t.Run(tc.name, func(t *testing.T) {
			result := calcCPUUtilization(float64(tc.cpuTime), float64(tc.executionTime))
			if result != tc.expected {
				t.Errorf("test %v; Expected: %v, Got: %v\n", tc.name, tc.expected, result)
			}
		})
	}
}
