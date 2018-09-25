package pub

import (
	"github.com/jenner/chaincode/hycc/src/res"
	"github.com/pkg/errors"
)

//=============================================================================
//动作数据
type Action struct {
	Mod  string `json:"mod,omitempty"`  //调用模块
	Cmd  string `json:"cmd,omitempty"`  //方法命令
	Args string `json:"args,omitempty"` //命令参数
}

//获取操作码
func (self *Action) GetOpt() string {
	return self.Mod + "-" + self.Cmd
}

//=============================================================================
//鉴权数据
type Auth struct {
	res.ResBase              //基础数据
	Opt        string `json:"opt,omitempty"`     //操作码(Mod-Cmd)
	Pexp       string `json:"pexp,omitempty"`    //权限表达式
	Disable    int32  `json:"disable,omitempty"` //是否可用
}

//-----------------------------------------------------------------------------
//构造函数
func NewAuth() *Auth {
	return &Auth{ResBase: res.ResBase{Idx: "auth", Ver: "1.0"}}
}

//获取索引
func (self *Auth) GetIdx() string {
	return self.Idx
}

//获取字段
func (self *Auth) GetField(ftype res.FieldType) (fields []string, err error) {
	switch ftype {
	case res.FIELD_KEY: //子建字段
		fields = []string{"Opt"}
	case res.FIELD_UNQ: //独有字段
		fields = []string{"Opt"}
	case res.FIELD_REQ: //必选字段
		fields = []string{"Opt", "Pexp"}
	case res.FIELD_CST: //常量字段
		fields = []string{"Idx"}
	default:
		return []string{}, errors.Errorf("Unknown field type:%v", ftype)
	}
	return fields, nil
}
