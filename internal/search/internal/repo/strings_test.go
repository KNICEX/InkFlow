package repo

import (
	"fmt"
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {
	fmt.Println(strings.Split("", ","))
	fmt.Println(strings.Split("a,b,", ","))
}
