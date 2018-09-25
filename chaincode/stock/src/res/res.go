package res

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"reflect"
	"strings"
	"utils"
)

//=============================================================================
//字段类型
type FieldType int

const (
	FIELD_NIL FieldType = 0 //空字段
	FIELD_KEY FieldType = 1 //键码字段
	FIELD_UNQ FieldType = 2 //独有字段
	FIELD_REQ FieldType = 3 //必选字段
	FIELD_CST FieldType = 4 //常量字段
)

//=============================================================================
//资源通用接口
type ResIf interface {
	GetIdx() string                                        //获取数据索引
	GetVer() string                                        //获取数据版本
	GetCts() string                                        //获取创建时间
	GetField(ftype FieldType) (fields []string, err error) //获取指定字段
}

//=============================================================================
//资源通用数据
type ResBase struct {
	Idx string `json:"idx,omitempty"` //数据索引
	Ver string `json:"ver,omitempty"` //数据版本
	Cts string `json:"cts,omitempty"` //创建时间
}

//获取数据索引
func (self ResBase) GetIdx() string {
	return self.Idx
}

//获取数据版本
func (self ResBase) GetVer() string {
	return self.Ver
}

//获取创建时间
func (self ResBase) GetCts() string {
	return self.Cts
}

//=============================================================================
//获取标签值
func GetTagByName(rsc ResIf, fname, sign string) (tag string, err error) {
	restype := reflect.TypeOf(rsc).Elem()
	field, ok := restype.FieldByName(fname)
	if !ok {
		return "", utils.Errorf("Can't filed %s in %s", fname, restype.String())
	}
	tag = field.Tag.Get(sign)
	if tag == "" {
		return "", utils.Errorf("The filed:%s in %s has no tag named %s", fname, restype.String(), sign)
	}
	tag = strings.Split(tag, ",")[0]
	return tag, nil
}

//获取字段值
func GetValByName(rsc ResIf, fname string) (val interface{}, err error) {
	resval := reflect.ValueOf(rsc).Elem()
	field := resval.FieldByName(fname)
	if field.Kind() == reflect.Invalid {
		return "", utils.Errorf("Can't filed %s in %s", fname, reflect.TypeOf(rsc).String())
	}
	return field.Interface(), nil
}

//获取键值对
func GetKVByName(rsc ResIf, field string, sign string) (rst string, err error) {
	key, err := GetTagByName(rsc, field, sign)
	if err != nil {
		return "", utils.Error(err)
	}
	val, err := GetValByName(rsc, field)
	if err != nil {
		return "", utils.Error(err)
	}
	return fmt.Sprintf("%s:%s", key, fmt.Sprint(val)), nil
}

//获取键值对(Json格式)
func GetJKVByName(rsc ResIf, field string, sign string) (rst string, err error) {
	key, err := GetTagByName(rsc, field, sign)
	if err != nil {
		return "", utils.Error(err)
	}
	val, err := GetValByName(rsc, field)
	if err != nil {
		return "", utils.Error(err)
	}
	kv := map[string]interface{}{key: val}
	jkv, err := json.Marshal(kv)
	if err != nil {
		return "", utils.Error(err)
	}
	return string(jkv[1 : len(jkv)-1]), nil
}

//=============================================================================
//生成键码
func GenKey(stub shim.ChaincodeStubInterface, res ResIf) (key string, err error) {
	attrs := []string{}
	fields, err := res.GetField(FIELD_KEY)
	if err != nil {
		return "", utils.Error(err)
	}
	for _, field := range fields {
		subkv, err := GetKVByName(res, field, "json")
		if err != nil {
			return "", utils.Error(err)
		}
		attrs = append(attrs, subkv)
	}
	key, err = stub.CreateCompositeKey(res.GetIdx(), attrs)
	if err != nil {
		return "", utils.Error(err)
	}
	return key, nil
}

//字段缺失判断
func HasField(res ResIf, ftype FieldType, jsonstr string) (rst bool, msg string, err error) {
	fields, err := res.GetField(ftype)
	if err != nil {
		return false, "", utils.Error(err)
	}
	for _, field := range fields {
		tagval, err := GetTagByName(res, field, "json")
		if err != nil {
			return false, "", utils.Error(err)
		}
		if !strings.Contains(jsonstr, "\""+tagval+"\":") {
			return false, fmt.Sprintf("Missing '%s' field", tagval), nil
		}
	}
	return true, "", nil
}

//字段可更改判断
func CanMutable(src, dst ResIf) (rst bool, msg string, err error) {
	fields, err := src.GetField(FIELD_CST)
	if err != nil {
		return false, "", err
	}
	for _, field := range fields {
		sval, err := GetValByName(src, field)
		if err != nil {
			return false, "", err
		}
		dval, err := GetValByName(dst, field)
		if err != nil {
			return false, "", err
		}
		if sval != dval {
			return false, fmt.Sprintf("The field:%s is immutable ", field), nil
		}
	}
	return true, "", nil
}

//=============================================================================
//新建资源
func NewRes(ref ResIf, cp bool) ResIf {
	srcp := reflect.ValueOf(ref)
	srcv := reflect.Indirect(srcp)
	dstp := reflect.New(srcv.Type())
	dstv := reflect.Indirect(dstp)
	if cp == true {
		srcjs, err := json.Marshal(ref)
		if err != nil {
			panic(utils.Error(err))
		}
		err = json.Unmarshal(srcjs, dstp.Interface())
		if err != nil {
			panic(utils.Error(err))
		}
	}
	if reflect.TypeOf(ref).Kind() == reflect.Ptr {
		return dstp.Interface().(ResIf)
	} else {
		return dstv.Interface().(ResIf)
	}
}

//新建资源组
func NewSlice(rsc ResIf) interface{} {
	srcp := reflect.ValueOf(rsc)
	srcv := reflect.Indirect(srcp)
	slice := reflect.MakeSlice(reflect.SliceOf(srcv.Type()), 0, 0)
	return slice.Interface()
}
