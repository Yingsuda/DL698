package pmodel

import (
	"dl698/utils"
	"fmt"
	"testing"
	"time"
)

func TestStr2DataTime(t *testing.T) {
	st := "2016-09-12 00:00:31.683"
	fmt.Printf("% 02X\n", utils.Str2DataTime(st))
	//07 E0 09 0C 01 00 00 1F 02 AB
}

func TestTime(t *testing.T) {
	var tt time.Time

	fmt.Println("tt", tt.Unix())
}
