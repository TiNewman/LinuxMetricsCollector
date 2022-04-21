package disk

import (
	"fmt"
	"testing"
)

func TestCollect(t *testing.T) {
	c := NewDiskCollector()
	m, err := c.Collect()
	if err != nil {
		t.Errorf("Error: %v\n", err.Error())
	}
	fmt.Printf("info: %+v\n", m)
}
