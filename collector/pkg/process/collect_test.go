package process

import (
	"testing"
)

func TestCollect(t *testing.T) {
	numProcs := 1
	collector := NewProcessCollectorWithoutRepo()

	list := collector.Collect()
	if len(list) != numProcs {
		t.Errorf("Expected %v processes; Got %v\n", numProcs, len(list))
	}
}

func TestCalcCPUUtilization(t *testing.T) {
	ts := []struct {
		name          string
		cpuTime       float64
		executionTime float64
		expected      float64
	}{
		{name: "full utilization", cpuTime: 450.6, executionTime: 450.6, expected: float64(100)},
		{name: "execution zero", cpuTime: 450.6, executionTime: 0, expected: float64(0)},
		{name: "cpu zero", cpuTime: 0, executionTime: 2043.56, expected: float64(0)},
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
