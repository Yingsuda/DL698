package model

import (
	"fmt"
	"strings"
	"testing"
)

func TestAD(t *testing.T) {
	name := "nihao"
	ss := strings.Split(name, "[")
	fmt.Println("len:", len(ss))
}
