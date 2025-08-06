package model

import (
	"fmt"
	"github.com/qiniu/qmgo/operator"
	builder "refeject_bson/builder"
	"testing"
)

func TestName(t *testing.T) {
	b := builder.New()
	var num int

	var str string
	str = "我是你爹"
	var nums = []int{1, 2, 3}
	var cp = []int{1, 3}
	var strs = []string{"lxh", "fwf", "sdf"}
	b.IsAdd(nums, "in", "nums")
	b.IsAdd(num, operator.In, "num")
	b.IsAdd(strs, operator.In, "name").IsAdd(str, "like", "who")
	b.IsAdd(cp, "between", "cp")
	fmt.Println(b)
}

func TestBuilder(t *testing.T) {

}
