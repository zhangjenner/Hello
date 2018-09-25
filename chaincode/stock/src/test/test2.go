package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unsafe"
	//	"time"
	//	"github.com/jenner/chaincode/stock/src/base"
)

type MyDI interface {
	Add(i int) MyDI
}

type MyDT struct {
	Val1 int `json:"val1"` //数据索引
	Val2 int `json:"val2"` //数据版本
}

func (d *MyDT) Add(i int) MyDI {
	d.Val1 += i
	d.Val2 += i
	return d
}

func NewSlice(res MyDI) interface{} {
	srcp := reflect.ValueOf(res)
	srcv := reflect.Indirect(srcp)
	slice := reflect.MakeSlice(reflect.SliceOf(srcv.Type()), 0, 0)
	return slice.Interface()
}

func myGet(res MyDI, rsts interface{}) interface{} {
	vals := []string{`{"val1":1,"val2":2}`, `{"val1":3,"val2":4}`}
	rstsv := reflect.ValueOf(rsts)
	for _, val := range vals {
		err := json.Unmarshal([]byte(val), res)
		if err != nil {
			panic(err)
		} else {
			dv := reflect.Indirect(reflect.ValueOf(res))
			fmt.Println("dv:", dv.Type())
			fmt.Println("dsv:", rstsv.Type())
			rstsv = reflect.Append(rstsv, dv)
		}
	}
	return rstsv.Interface()
}

func test(d MyDI, rsts interface{}, a int) interface{} {
	vals := reflect.ValueOf(rsts)
	slice := reflect.MakeSlice(vals.Type(), 0, 0)
	for i := 0; i < vals.Len(); i++ {
		val := vals.Index(i)
		pval := reflect.NewAt(val.Type(), unsafe.Pointer(val.UnsafeAddr()))
		add := pval.MethodByName("Add")
		if add.IsValid() {
			add.Call([]reflect.Value{reflect.ValueOf(a)})
			slice = reflect.Append(slice, val)
		} else {
			fmt.Println("fuck")
		}
	}
	return slice.Interface()
}

func Run(d MyDI) {
	ds := NewSlice(d)
	nds := myGet(d, ds)
	nd := test(d, nds, 2)
	fmt.Println(nd)
}

func main2() {
	d := &MyDT{}
	nd := d.Add(1)
	fmt.Println(reflect.TypeOf(nd))
	Run(d)
}
