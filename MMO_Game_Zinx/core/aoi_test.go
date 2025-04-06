package core

import (
"fmt"
"testing"
)

func TestNewAOIManager(t *testing.T) {
	aoi := NewAOIManager(0, 250, 0, 250, 5, 5)
	if aoi == nil {
		fmt.Println("aoi is nil")
	}
	a := aoi.GetSurroundGridsByGid(15)
	for _, grid := range a {
		fmt.Println(grid.GID)
	}
}
