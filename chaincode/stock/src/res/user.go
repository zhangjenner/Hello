package res

import (
	"github.com/pkg/errors"
)

//=============================================================================
//用户数据
type User struct {
	ResBase        //基础数据
	Cert    string `json:"cert,omitempty"`  //用户证书
	User    string `json:"user,omitempty"`  //用户名
	Phone   string `json:"phone,omitempty"` //手机号
	Pwd     string `json:"pwd,omitempty"`   //用户密码
	Role    string `json:"role,omitempty"`  //用户角色
}

//-----------------------------------------------------------------------------
//root用户
var SYSMG_USER = &User{ResBase: ResBase{Idx: "User", Ver: "1.0", Cts: "1514736000"},
	Cert: "", User: "admin", Phone: "18202860046", Pwd: "21232F297A57A5A743894A0E4A801FC3", Role: "admin"}

//-----------------------------------------------------------------------------
//构造函数
func NewUser() *User {
	user := User{ResBase: ResBase{Idx: "User", Ver: "1.0"}}
	return &user
}

//获取字段
func (self *User) GetField(ftype FieldType) (fields []string, err error) {
	switch ftype {
	case FIELD_NIL: //空字段
		fields = []string{}
	case FIELD_KEY: //子建字段
		fields = []string{"Phone"}
	case FIELD_UNQ: //独有字段
		fields = []string{"Phone", "Cert"}
	case FIELD_REQ: //必选字段
		fields = []string{"Phone", "User", "Phone", "Pwd"}
	case FIELD_CST: //常量字段
		fields = []string{"Idx", "Phone"}
	default:
		return []string{}, errors.Errorf("Unknown field type:%v", ftype)
	}
	return fields, nil
}
