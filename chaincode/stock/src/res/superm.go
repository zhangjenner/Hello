package res

import (
	"github.com/pkg/errors"
)

//=============================================================================
//超市状态
type SupermST uint32

const (
	SNST_DONE    SupermST = 0 //已完成
	SNST_BUYING  SupermST = 1 //正在买入
	SNST_SELLING SupermST = 2 //正在卖出
)

//超市数据
type Superm struct {
	ResBase          //基础数据
	Id      uint32   `json:"id"`      //超市编号
	Num     uint32   `json:"num"`     //超市编号
	Name    string   `json:"name"`    //超市名称
	Address string   `json:"address"` //超市地址
	Capital float32  `json:"capital"` //注册资金
	Stock   uint32   `json:"stock"`   //股权数量
	Owner   string   `json:"owner"`   //拥有者账号
	Hold    uint32   `json:"hold"`    //占有股权数量
	State   SupermST `json:"state"`   //买入状态
}

//-----------------------------------------------------------------------------
//构造函数
func NewSuperm() *Superm {
	supermarket := Superm{ResBase: ResBase{Idx: "Superm", Ver: "1.0"}}
	return &supermarket
}

//获取字段
func (self *Superm) GetField(ftype FieldType) (fields []string, err error) {
	switch ftype {
	case FIELD_NIL: //空字段
		fields = []string{}
	case FIELD_KEY: //子建字段
		fields = []string{"Id"}
	case FIELD_UNQ: //独有字段
		fields = []string{"Id", "Num"}
	case FIELD_REQ: //必选字段
		fields = []string{"Num", "Name"}
	case FIELD_CST: //常量字段
		fields = []string{"Idx", "Id", "Num"}
	default:
		return []string{}, errors.Errorf("Unknown field type:%v", ftype)
	}
	return fields, nil
}
