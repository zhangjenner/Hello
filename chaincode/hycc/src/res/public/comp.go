package pub

import (
	"encoding/json"
	"github.com/jenner/chaincode/hycc/src/res"
	"github.com/jenner/chaincode/hycc/src/utils"
	"github.com/pkg/errors"
)

//=============================================================================
//加密数据
type CompECD struct {
	Username string `json:"username,omitempty"` //
	Sex      int32  `json:"sex,omitempty"`      //
	Birthday string `json:"birthday,omitempty"` //
	Mobile   string `json:"mobile,omitempty"`   //
	Phone    string `json:"phone,omitempty"`    //
	Email    string `json:"email,omitempty"`    //
	Position string `json:"position,omitempty"` //
}

//公司数据
type Comp struct {
	res.ResBase                         //基础数据
	Mspid             string            `json:"mspid,omitempty"`             //所属组织MspID
	Id                string            `json:"id,omitempty"`                //企业ID
	Pid               string            `json:"pid,omitempty"`               //
	Sourceid          string            `json:"sourceid,omitempty"`          //发展成为会员的企业id,源头
	Name              string            `json:"name,omitempty"`              //
	Type              int32             `json:"type,omitempty"`              //医药客户1,专线物流2,第三方物流服务-类似华药物流3
	IsBranch          int32             `json:"isBranch,omitempty"`          //0
	Disable           int32             `json:"disable,omitempty"`           //0
	DeleteFlag        int32             `json:"deleteFlag,omitempty"`        //0
	Status            int32             `json:"status,omitempty"`            //
	LockStatus        int32             `json:"lockStatus,omitempty"`        //
	Proid             int32             `json:"proid,omitempty"`             //
	Cityid            int32             `json:"cityid,omitempty"`            //
	Areaid            int32             `json:"areaid,omitempty"`            //
	Address           string            `json:"address,omitempty"`           //
	Creator           string            `json:"creator,omitempty"`           //
	CreateDateTime    string            `json:"createDateTime,omitempty"`    //
	Updator           string            `json:"updator,omitempty"`           //
	UpdateDateTime    string            `json:"updateDateTime,omitempty"`    //
	BusLicensePhoto   string            `json:"busLicensePhoto,omitempty"`   //
	VerifyStatus      int32             `json:"verifyStatus,omitempty"`      //审核状态,默认0未审核,1审核通过,-1未通过审核
	Verifyid          string            `json:"verifyid,omitempty"`          //审核人id
	VerifyDateTime    string            `json:"verifyDateTime,omitempty"`    //
	VerifyDescription string            `json:"verifyDescription,omitempty"` //
	AdminUserid       string            `json:"adminUserid,omitempty"`       //企业负责人id
	Crypto            map[string]string `json:"crypto,omitempty"`            //秘钥
	EcdStr            string            `json:"ecdStr,omitempty"`            //
	CompECD
}

//-----------------------------------------------------------------------------
//系统管理公司
var SYSMG_COMP = &Comp{ResBase: res.ResBase{Idx: "comp", Ver: "1.0", Cts: "1514736000"},
	Id: "-1", AdminUserid: "-1", Mspid: "SysMgMSP", Type: -1, Name: "SysMgComp"}

//-----------------------------------------------------------------------------
//构造函数
func NewComp() *Comp {
	comp := Comp{ResBase: res.ResBase{Idx: "comp", Ver: "1.0"}}
	comp.Crypto = make(map[string]string)
	return &comp
}

//获取索引
func (self *Comp) GetIdx() string {
	return self.Idx
}

//获取字段
func (self *Comp) GetField(ftype res.FieldType) (fields []string, err error) {
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
func (self *Comp) Encrypt(rAes *utils.AES) {
	nilECD := CompECD{}
	if self.CompECD != nilECD {
		pdata, err := json.Marshal(&self.CompECD)
		if err != nil {
			panic(utils.Error(err))
		}
		self.EcdStr = rAes.Encrypt(pdata)
		self.CompECD = nilECD
	}
}

//资源解密
func (self *Comp) Decrypt(rAes *utils.AES) {
	if self.EcdStr != "" {
		pdata := rAes.Decrypt(self.EcdStr)
		err := json.Unmarshal(pdata, &self.CompECD)
		if err != nil {
			panic(utils.Error(err))
		}
	}
	self.EcdStr = ""
}
