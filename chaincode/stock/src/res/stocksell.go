package res

import (
	"github.com/pkg/errors"
)

//=============================================================================
//卖出状态
type SellST uint32

const (
	SELLST_SELING SellST = 0 //正在卖出
	SELLST_DONE   SellST = 1 //已经卖出
)

//卖出交易
type StockSell struct {
	ResBase           //基础数据
	Id        uint32  `json:"id"`        //超市编号
	Num       uint32  `json:"num"`       //超市编号
	Name      string  `json:"name"`      //超市名称
	Address   string  `json:"address"`   //超市地址
	Capital   float32 `json:"capital"`   //注册资金
	Stock     uint32  `json:"stock"`     //股权数量
	SellId    uint32  `json:"sellId"`    //卖出交易ID
	SellState SellST  `json:"sellState"` //卖出状态
	SellUser  string  `json:"sellUser"`  //卖出者
	SellPhone string  `json:"sellPhone"` //卖出者电话号码
	SellStock uint32  `json:"sellStock"` //卖出股权数量
	SellPrice float32 `json:"sellPrice"` //卖出股权价格
	SoldStock uint32  `json:"soldStock"` //已卖出股权数量
	LeftStock uint32  `json:"leftStock"` //剩余股权数量
	SellDate  int64   `json:"sellDate"`  //卖出日期
	Deadline  int64   `json:"deadline"`  //到期日期
}

//-----------------------------------------------------------------------------
//构造函数
func NewStockSell() *StockSell {
	trade := StockSell{ResBase: ResBase{Idx: "StockSell", Ver: "1.0"}}
	return &trade
}

//获取字段
func (self *StockSell) GetField(ftype FieldType) (fields []string, err error) {
	switch ftype {
	case FIELD_NIL: //空字段
		fields = []string{}
	case FIELD_KEY: //子建字段
		fields = []string{"SellId"}
	case FIELD_UNQ: //独有字段
		fields = []string{"SellId"}
	case FIELD_REQ: //必选字段
		fields = []string{"SellId"}
	case FIELD_CST: //常量字段
		fields = []string{"Idx", "SellId"}
	default:
		return []string{}, errors.Errorf("Unknown field type:%v", ftype)
	}
	return fields, nil
}
