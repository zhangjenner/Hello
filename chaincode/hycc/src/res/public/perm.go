package pub

import (
	"github.com/jenner/chaincode/hycc/src/res"
	"github.com/pkg/errors"
)

//=============================================================================
//权限数据
type Perm struct {
	res.ResBase        //基础数据
	Id          string `json:"id,omitempty"`          //
	Pid         string `json:"pid,omitempty"`         //
	Name        string `json:"name,omitempty"`        //
	Description string `json:"description,omitempty"` //
	Type        string `json:"type,omitempty"`        //0：菜单，1按钮
	Permission  string `json:"permission,omitempty"`  //权限代码
	Url         string `json:"url,omitempty"`         //菜单对应url
	Icon        string `json:"icon,omitempty"`        //图标
	Sorted      string `json:"sorted,omitempty"`      //排序
	Disable     int32  `json:"disable,omitempty"`     //1禁用,0启用
}

//-----------------------------------------------------------------------------
//系统管理权限
var SYSMG_PERM = &Perm{ResBase: res.ResBase{Idx: "perm", Ver: "1.0", Cts: "1514736000"},
	Id: "-1", Pid: "0", Permission: "SysMgPerm"}

//-----------------------------------------------------------------------------
//构造函数
func NewPerm() *Perm {
	return &Perm{ResBase: res.ResBase{Idx: "perm", Ver: "1.0"}}
}

//获取字段
func (self *Perm) GetField(ftype res.FieldType) (fields []string, err error) {
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
