package res

import (
	"github.com/pkg/errors"
)

//=============================================================================
//买入状态
type BuyST uint32

const (
	BUYST_BUYING BuyST = 0 //正在买入
	BUYST_DONE   BuyST = 1 //已经买入
)

//买入交易
type StockBuy struct {
	ResBase           //基础数据
	Id        uint32  `json:"id"`                  //超市编号
	Num       uint32  `json:"num"`                 //超市编号
	Name      string  `json:"name"`                //超市名称
	Address   string  `json:"address"`             //超市地址
	Capital   float32 `json:"capital"`             //注册资金
	Stock     uint32  `json:"stock"`               //股权数量
	SellId    uint32  `json:"sellId,omitempty"`    //关联卖出ID
	SellUser  string  `json:"buyUser,omitempty"`   //买入者
	SellPhone string  `json:"sellPhone,omitempty"` //卖出者电话号码
	BuyId     uint32  `json:"buyId"`               //买入交易ID
	BuyState  BuyST   `json:"buyState"`            //买入状态
	BuyUser   string  `json:"buyUser"`             //买入者
	BuyPhone  string  `json:"buyPhone"`            //买入者电话号码
	BuyStock  uint32  `json:"buyStock"`            //买入股权数量
	BuyPrice  float32 `json:"buyPrice"`            //买入股权价格
	BuyDate   int64   `json:"buyDate"`             //买入日期
}

//-----------------------------------------------------------------------------
//构造函数
func NewStockBuy() *StockBuy {
	trade := StockBuy{ResBase: ResBase{Idx: "StockBuy", Ver: "1.0"}}
	return &trade
}

//获取字段
func (self *StockBuy) GetField(ftype FieldType) (fields []string, err error) {
	switch ftype {
	case FIELD_NIL: //空字段
		fields = []string{}
	case FIELD_KEY: //子建字段
		fields = []string{"BuyId"}
	case FIELD_UNQ: //独有字段
		fields = []string{"BuyId"}
	case FIELD_REQ: //必选字段
		fields = []string{"BuyId"}
	case FIELD_CST: //常量字段
		fields = []string{"Idx", "BuyId"}
	default:
		return []string{}, errors.Errorf("Unknown field type:%v", ftype)
	}
	return fields, nil
}
