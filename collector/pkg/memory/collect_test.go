package memory

import (
	"fmt"
	"testing"
)

func TestCollect(t *testing.T) {
	mc := NewMemoryCollector()

	mem, err := mc.Collect()
	if err != nil {
		t.Errorf("Error: %v\n", err.Error())
	}
	fmt.Printf("%+v\n", mem)
}
