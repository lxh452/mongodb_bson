package model

import (
	"fmt"
	builder "refeject_bson/jsy_builder"
	"testing"
)

func TestName(t *testing.T) {
	b := builder.New()
	b.IsAdd("test", "eq", "test")
	b.IsAdd(1, "eq", "test1")
	fmt.Println(b)
}
