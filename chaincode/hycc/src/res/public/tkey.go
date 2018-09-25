package pub

import (
	"github.com/jenner/chaincode/hycc/src/res"
	"github.com/pkg/errors"
)

//=============================================================================
//交易秘钥
type TKey struct {
	res.ResBase                   //基础数据
	Cid         string            `json:"cid,omitempty"`      //所属公司ID
	STime       string            `json:"sTime,omitempty"`    //秘钥开始时间
	ETime       string            `json:"eTime,omitempty"`    //秘钥结束时间
	TPubKey     string            `json:"tPubKey,omitempty"`  //交易公钥
	TPriKeys    map[string]string `json:"tPriKeys,omitempty"` //交易私钥
	RootSign    string            `json:"rootSign,omitempty"` //root签名
}

//-----------------------------------------------------------------------------
//系统管理组织
var SYSMG_TKEY = &TKey{ResBase: res.ResBase{Idx: "tkey", Ver: "1.0", Cts: "1514736000"},
	Cid: "-1", STime: "1514736000", ETime: "1830268799", TPriKeys: make(map[string]string)}

//-----------------------------------------------------------------------------
//构造函数
func NewTKey() *TKey {
	impl := &TKey{ResBase: res.ResBase{Idx: "tkey", Ver: "1.0"}}
	impl.TPriKeys = make(map[string]string)
	return impl
}

//获取字段
func (self *TKey) GetField(ftype res.FieldType) (fields []string, err error) {
	switch ftype {
	case res.FIELD_KEY: //子建字段
		fields = []string{"Cid", "STime"}
	case res.FIELD_UNQ: //独有字段
		fields = []string{}
	case res.FIELD_REQ: //必选字段
		fields = []string{"Cid", "STime", "ETime"}
	case res.FIELD_CST: //常量字段
		fields = []string{"Idx", "Ver", "Cid", "STime", "TPubKey"}
	default:
		return []string{}, errors.Errorf("Unknown field type:%v", ftype)
	}
	return fields, nil
}
