package sup

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"reflect"
	"res"
	"strings"
	"utils"
)

//=============================================================================
//存入数据
func PutData(stub shim.ChaincodeStubInterface, rsc res.ResIf) {
	key, err := res.GenKey(stub, rsc)
	if err != nil {
		panic(utils.Error(err))
	}
	val, err := json.Marshal(rsc)
	if err != nil {
		panic(utils.Error(err))
	}
	err = stub.PutState(key, val)
	if err != nil {
		panic(utils.Error(err))
	}
}

//获取数据
func GetData(stub shim.ChaincodeStubInterface, ref res.ResIf) res.ResIf {
	key, err := res.GenKey(stub, ref)
	if err != nil {
		panic(utils.Error(err))
	}
	val, err := stub.GetState(key)
	if err != nil {
		panic(utils.Error(err))
	} else if val == nil {
		panic(utils.Errorf("%s not exist", strings.Replace(key, "\x00", "|", -1)))
	}
	rst := res.NewRes(ref, false)
	err = json.Unmarshal(val, rst)
	if err != nil {
		panic(utils.Error(err))
	}
	return rst
}

//删除数据
func DelData(stub shim.ChaincodeStubInterface, ref res.ResIf) {
	key, err := res.GenKey(stub, ref)
	if err != nil {
		panic(utils.Error(err))
	}
	err = stub.DelState(key)
	if err != nil {
		panic(utils.Error(err))
	}
}

//查询数据
func QryData(stub shim.ChaincodeStubInterface, ref res.ResIf, qstr string) interface{} {
	rstIt, err := stub.GetQueryResult(qstr)
	if err != nil {
		panic(utils.Error(err))
	}
	defer rstIt.Close()
	var rst res.ResIf
	rsts := res.NewSlice(ref)
	rstsv := reflect.ValueOf(rsts)
	for rstIt.HasNext() {
		rst = res.NewRes(ref, false)
		rsp, err := rstIt.Next()
		if err != nil {
			panic(utils.Error(err))
		}
		err = json.Unmarshal(rsp.GetValue(), rst)
		if err != nil {
			panic(utils.Error(err))
		} else {
			v := reflect.Indirect(reflect.ValueOf(rst))
			rstsv = reflect.Append(rstsv, v)
		}
	}
	return rstsv.Interface()
}

//选择数据
func SelectData(stub shim.ChaincodeStubInterface, ref res.ResIf) interface{} {
	cond, err := json.Marshal(ref)
	if err != nil {
		panic(utils.Error(err))
	}
	qstr := fmt.Sprintf("{\"selector\":%s}", string(cond))
	return QryData(stub, ref, qstr)
}

//=============================================================================
//数据查重
func HasDuplicate(stub shim.ChaincodeStubInterface, rsc res.ResIf) (rst bool, msg string) {
	fields, err := rsc.GetField(res.FIELD_UNQ)
	if err != nil {
		panic(utils.Error(err))
	}
	for _, field := range fields {
		jkv, err := res.GetJKVByName(rsc, field, "json")
		if err != nil {
			panic(utils.Error(err))
		}
		fmt.Println("jkv:%s", jkv);
		qstr := fmt.Sprintf("{\"selector\":{\"idx\":\"%s\",%s}}", rsc.GetIdx(), jkv)
		rstIt, err := stub.GetQueryResult(qstr)
		if err != nil {
			panic(utils.Error(err))
		}
		defer rstIt.Close()
		if rstIt.HasNext() {
			rsp, err := rstIt.Next()
			if err != nil {
				panic(utils.Error(err))
			}
			ktype, _, err := stub.SplitCompositeKey(rsp.GetKey())
			if err != nil {
				panic(utils.Error(err))
			} else if ktype == rsc.GetIdx() {
				return true, fmt.Sprintf("%s has duplicated", jkv)
			}
		}
	}
	return false, ""
}
