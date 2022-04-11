package cpu

import (
	"testing"

	"github.com/prometheus/procfs"
)

type UsageCase struct {
	name  string
	start procfs.CPUStat
	end   procfs.CPUStat
}

func TestCalculateUsage(t *testing.T) {

}
