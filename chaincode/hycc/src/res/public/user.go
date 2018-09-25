package pub

import (
	"encoding/json"
	"github.com/jenner/chaincode/hycc/src/res"
	"github.com/jenner/chaincode/hycc/src/utils"
	"github.com/pkg/errors"
)

//=============================================================================
//加密数据
type UserECD struct {
	Username string `json:"username,omitempty"` //
	Sex      int32  `json:"sex,omitempty"`      //
	Birthday string `json:"birthday,omitempty"` //
	Mobile   string `json:"mobile,omitempty"`   //
	Phone    string `json:"phone,omitempty"`    //
	Email    string `json:"email,omitempty"`    //
	Position string `json:"position,omitempty"` //
}

//用户数据
type User struct {
	res.ResBase                      //基础数据
	Cert           string            `json:"cert,omitempty"`           //用户证书
	RoleIds        []string          `json:"roleIds,omitempty"`        //用户角色ID
	Id             string            `json:"id,omitempty"`             //
	Cid            string            `json:"cid,omitempty"`            //企业id
	Disable        int32             `json:"disable,omitempty"`        //是否可以使用，1禁用，0可用
	DeleteFlag     int32             `json:"deleteFlag,omitempty"`     //
	Creator        string            `json:"creator,omitempty"`        //
	CreateDateTime string            `json:"createDateTime,omitempty"` //
	Updator        string            `json:"updator,omitempty"`        //
	UpdateDateTime string            `json:"updateDateTime,omitempty"` //
	Crypto         map[string]string `json:"crypto,omitempty"`         //秘钥
	EcdStr         string            `json:"ecdStr,omitempty"`            //
	UserECD
}

//-----------------------------------------------------------------------------
//root用户
var SYSMG_USER = &User{ResBase: res.ResBase{Idx: "user", Ver: "1.0", Cts: "1514736000"},
	RoleIds: []string{"-1"}, Id: "-1", Cid: "-1"}

//-----------------------------------------------------------------------------
//构造函数
func NewUser() *User {
	user := User{ResBase: res.ResBase{Idx: "user", Ver: "1.0"}}
	user.Crypto = make(map[string]string)
	return &user
}

//获取字段
func (self *User) GetField(ftype res.FieldType) (fields []string, err error) {
	switch ftype {
	case res.FIELD_KEY: //子建字段
		fields = []string{"Id"}
	case res.FIELD_UNQ: //独有字段
		fields = []string{"Id"}
	case res.FIELD_REQ: //必选字段
		fields = []string{"Id"}
	case res.FIELD_CST: //常量字段
		fields = []string{"Idx", "Id"}
	default:
		return []string{}, errors.Errorf("Unknown field type:%v", ftype)
	}
	return fields, nil
}

//资源加密
func (self *User) Encrypt(rAes *utils.AES) {
	nilECD := UserECD{}
	if self.UserECD != nilECD {
		pdata, err := json.Marshal(&self.UserECD)
		if err != nil {
			panic(utils.Error(err))
		}
		self.EcdStr = rAes.Encrypt(pdata)
		self.UserECD = nilECD
	}
}

//资源解密
func (self *User) Decrypt(rAes *utils.AES) {
	if self.EcdStr != "" {
		pdata := rAes.Decrypt(self.EcdStr)
		err := json.Unmarshal(pdata, &self.UserECD)
		if err != nil {
			panic(utils.Error(err))
		}
	}
	self.EcdStr = ""
}
