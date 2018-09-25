package main

import (
	//	"bytes"
	//	"encoding/gob"
	"encoding/json"
	"fmt"
	"utils"
	//	"github.com/jenner/chaincode/stock/src/res"
	//	"github.com/jenner/chaincode/stock/src/res/public"
	"reflect"
	//	"unsafe"
	//	"strings"
)

type TT struct {
	M map[string]string `json:"m,omitempty"`
	V int               `json:"v,omitempty"`
}

func NewCopy(rsc interface{}) interface{} {
	srcp := reflect.ValueOf(rsc)
	srcv := reflect.Indirect(srcp)
	dstp := reflect.New(srcv.Type())
	dstv := reflect.Indirect(dstp)
	srcjs, err := json.Marshal(rsc)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(srcjs, dstp.Interface())
	if err != nil {
		panic(err)
	}
	if reflect.TypeOf(rsc).Kind() != reflect.Ptr {
		return dstv.Interface()
	} else {
		return dstp.Interface()
	}
}

func main3() {
	aes := utils.NewAES(1).GenKey()
	key := aes.GetKey()
	fmt.Println(key)
	ecc := utils.NewECC(1).GenKey()
	cy := ecc.Encrypt(key)
	fmt.Println(cy)
}
