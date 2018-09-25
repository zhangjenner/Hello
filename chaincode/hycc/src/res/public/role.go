package pub

import (
	"github.com/jenner/chaincode/hycc/src/res"
	"github.com/pkg/errors"
)

//=============================================================================
//角色数据
type Role struct {
	res.ResBase             //基础数据
	PermIds        []string `json:"permIds,omitempty"`        //权限ID
	Id             string   `json:"id,omitempty"`             //
	CompanyId      string   `json:"companyId,omitempty"`      //企业id
	Name           string   `json:"name,omitempty"`           //名称
	Description    string   `json:"description,omitempty"`    //描述
	Creator        string   `json:"creator,omitempty"`        //创建者
	CreateDateTime string   `json:"createDateTime,omitempty"` //创建时间
	Updator        string   `json:"updator,omitempty"`        //更新者
	UpdateDateTime string   `json:"updateDateTime,omitempty"` //更新时间
	RoleType       int32    `json:"roleType,omitempty"`       //0：系统模板角色，1，用户创建角色
	RoleChannel    int32    `json:"roleChannel,omitempty"`    //角色适用渠道,1华药\2客户\3专线
	Disable        int32    `json:"disable,omitempty"`        //1启用，0禁用
	Status         int32    `json:"status,omitempty"`         //
}

//-----------------------------------------------------------------------------
//系统管理员角色
var SYSMG_ROLE = &Role{ResBase: res.ResBase{Idx: "role", Ver: "1.0", Cts: "1514736000"},
	PermIds: []string{"-1"}, Id: "-1", CompanyId: "-1", Name: "SysMgRole"}

//-----------------------------------------------------------------------------
//构造函数
func NewRole() *Role {
	return &Role{ResBase: res.ResBase{Idx: "role", Ver: "1.0"}}
}

//获取字段
func (self *Role) GetField(ftype res.FieldType) (fields []string, err error) {
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
