package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"math/rand"
	"net/http"
	"res"
	"strconv"
	"sup"
	"time"
	"utils"
)

//=============================================================================
//交易记录
type StockRecord struct {
	Num     uint32  `json:"num"`     //超市编号
	Name    string  `json:"name"`    //超市名称
	Address string  `json:"address"` //超市地址
	Stock   uint32  `json:"stock"`   //交易股权
	Price   float32 `json:"price"`   //交易股权价格
	Date    int64   `json:"date"`    //交易日期
	Type    uint32  `json:"type"`    //交易类型
}

//=============================================================================
//股权管理
func (cc *StockCC) StockMg(cmd, arg string) peer.Response {
	stub := cc.GetStub()
	rsp := shim.Success(nil)
	logger := cc.GetLogger()
	logger.Info("StockMg Start---------------------------------------------------------")
	//获取用户信息
	user, err := cc.GetCreator()
	if err != nil {
		return utils.ErrorRsp(err)
	}
	switch cmd {
	case "newSuperm": //新建超市
		if user.Role != "admin" {
			return utils.ErrorRspf("Permission denied")
		}
		superm := res.NewSuperm()
		sup.ParseRsc(&arg, res.FIELD_NIL, superm)
		superm.Id = rand.Uint32()
		superm.Owner = user.Phone
		superm.Hold = superm.Stock
		superm.State = res.SNST_DONE
		rsp = sup.AddRes(cc, superm, arg)
	case "stockQuery": //股权占有查询
		superm := res.NewSuperm()
		sup.ParseRsc(&arg, res.FIELD_NIL, superm)
		if user.Phone != superm.Owner {
			return utils.ErrorRspf("You are not '%s'", superm.Owner)
		}
		qstr := fmt.Sprintf("{\"selector\":{\"owner\":\"%s\"}}", superm.Owner)
		rsp = sup.QryRes(cc, superm, qstr)
	case "sellStock": //股权卖出操作
		sell := res.NewStockSell()
		sup.ParseRsc(&arg, res.FIELD_NIL, sell)
		suprm := res.NewSuperm()
		suprm.Id = sell.Id
		suprm = sup.GetData(stub, suprm).(*res.Superm)
		//参数验证
		if suprm.State != res.SNST_DONE {
			return utils.ErrorRspf("This surpmarket can't be sold")
		} else if user.Phone != sell.SellPhone {
			return utils.ErrorRspf("You are not '%s'", sell.SellPhone)
		} else if suprm.Owner != sell.SellPhone {
			return utils.ErrorRspf("You are not have supermarket:%d", suprm.Id)
		} else if sell.SellStock <= 0 {
			return utils.ErrorRspf("Stock num must be greater than zero")
		} else if sell.SellStock > suprm.Hold {
			return utils.ErrorRspf("There aren't enough stocks")
		}
		//修改超市数据
		suprm.State = res.SNST_SELLING
		sup.PutData(stub, suprm)
		//新建股权卖出交易
		sell.Id = suprm.Id
		sell.Num = suprm.Num
		sell.Name = suprm.Name
		sell.Address = suprm.Address
		sell.Capital = suprm.Capital
		sell.Stock = suprm.Stock
		sell.SellId = rand.Uint32()
		sell.SellState = res.SELLST_SELING
		sell.SellUser = user.User
		sell.SellPhone = user.Phone
		sell.SellStock = sell.SellStock
		sell.SellPrice = sell.SellPrice
		sell.SoldStock = 0
		sell.LeftStock = sell.SellStock
		sell.SellDate = time.Now().Unix()
		sell.Deadline = sell.Deadline
		sup.PutData(stub, sell)
	case "sellQuery": //股权卖出查询
		//参数解析
		var parm map[string]string
		err = json.Unmarshal([]byte(arg), &parm)
		if err != nil {
			panic(utils.Error(err))
		}
		//参数验证
		if (parm["type"] == "sell") && (user.Phone != parm["phone"]) {
			return utils.ErrorRspf("You are not '%s'", parm["phone"])
		}
		//查询数据
		var qstr string
		sell := res.NewStockSell()
		if parm["type"] == "sell" {
			qstr = fmt.Sprintf("{\"selector\":{\"sellId\":{\"$gt\":null},\"sellPhone\":\"%s\"}}", parm["phone"])
		} else {
			qstr = `{"selector":{"sellId":{"$gt":null}}}`
		}
		rsp = sup.QryRes(cc, sell, qstr)
	case "buyStock": //股权买入操作
		buy := res.NewStockBuy()
		sup.ParseRsc(&arg, res.FIELD_NIL, buy)
		//查找股权卖出交易
		sell := res.NewStockSell()
		sell.SellId = buy.SellId
		sell = sup.GetData(stub, sell).(*res.StockSell)
		//参数验证
		if user.Phone == sell.SellPhone {
			return utils.ErrorRspf("You can't buy your own stocks")
		} else if buy.BuyStock <= 0 {
			return utils.ErrorRspf("Stock num must be greater than zero")
		} else if buy.BuyStock > (sell.SellStock - sell.SoldStock) {
			return utils.ErrorRspf("There aren't enough stocks")
		} else if sell.SellState != res.SELLST_SELING {
			return utils.ErrorRspf("The sale of stocks has been halted")
		}
		//卖出股权交易更新
		sell.SoldStock += buy.BuyStock
		sup.PutData(stub, sell)
		//买入股权
		buy.Id = rand.Uint32()
		buy.Num = sell.Num
		buy.Name = sell.Name
		buy.Address = sell.Address
		buy.Capital = sell.Capital
		buy.Stock = sell.Stock
		buy.SellId = sell.SellId
		buy.SellUser = sell.SellUser
		buy.SellPhone = sell.SellPhone
		buy.BuyId = rand.Uint32()
		buy.BuyState = res.BUYST_BUYING
		buy.BuyUser = user.User
		buy.BuyPhone = user.Phone
		buy.Stock = buy.Stock
		buy.BuyPrice = sell.SellPrice
		buy.BuyDate = time.Now().Unix()
		sup.PutData(stub, buy)
		//新建超市数据
		suprm := res.NewSuperm()
		suprm.Id = buy.Id
		suprm.Num = buy.Num
		suprm.Name = buy.Name
		suprm.Address = buy.Address
		suprm.Capital = buy.Capital
		suprm.Stock = buy.Stock
		suprm.Owner = buy.BuyPhone
		suprm.Hold = buy.BuyStock
		suprm.State = res.SNST_BUYING
		sup.PutData(stub, suprm)
	case "sellDetial": //股权卖出详情
		//参数解析
		var parm map[string]uint32
		err = json.Unmarshal([]byte(arg), &parm)
		if err != nil {
			panic(utils.Error(err))
		}
		//查询数据
		buy := res.NewStockBuy()
		qstr := fmt.Sprintf("{\"selector\":{\"sellId\":%d}}", parm["sellId"])
		rsp = sup.QryRes(cc, buy, qstr)
	case "sellConfirm": //股权卖出确认
		//参数解析
		var parm map[string]uint32
		err = json.Unmarshal([]byte(arg), &parm)
		if err != nil {
			panic(utils.Error(err))
		}
		//更新买入股权交易
		buy := res.NewStockBuy()
		buy.BuyId = parm["buyId"]
		buy = sup.GetData(stub, buy).(*res.StockBuy)
		buy.BuyState = res.BUYST_DONE
		sup.PutData(stub, buy)
		//更新买入超市数据
		var hold uint32 = 0
		suprm := res.NewSuperm()
		suprm.Num = buy.Num
		suprm.Owner = buy.BuyPhone
		suprm.State = res.SNST_DONE
		suprms := sup.SelectData(stub, suprm).([]res.Superm)
		for i := range suprms {
			hold += suprms[i].Hold
			sup.DelData(stub, &suprms[i])
		}
		suprm.Id = buy.Id
		suprm = sup.GetData(stub, suprm).(*res.Superm)
		suprm.State = res.SNST_DONE
		suprm.Hold += hold
		sup.PutData(stub, suprm)
		//更新买出股权交易
		sell := res.NewStockSell()
		sell.SellId = buy.SellId
		sell = sup.GetData(stub, sell).(*res.StockSell)
		sell.LeftStock -= buy.BuyStock
		if sell.LeftStock <= 0 {
			sell.SellState = res.SELLST_DONE
		}
		sup.PutData(stub, sell)
		//更新卖出超市
		if sell.SellState == res.SELLST_DONE {
			suprm.Id = sell.Id
			suprm = sup.GetData(stub, suprm).(*res.Superm)
			suprm.Hold -= sell.SoldStock
			suprm.State = res.SNST_DONE
			if suprm.Hold <= 0 {
				sup.DelData(stub, suprm)
			} else {
				sup.PutData(stub, suprm)
			}
		}
	case "dealRecord": //股权买卖记录
		//参数解析
		var parm map[string]string
		err = json.Unmarshal([]byte(arg), &parm)
		if err != nil {
			panic(utils.Error(err))
		} else if user.Phone != parm["owner"] {
			return utils.ErrorRspf("You are not '%s'", parm["owner"])
		}
		//查询交易记录
		var record StockRecord
		var records []StockRecord
		buy := res.NewStockBuy()
		buys := sup.QryData(stub, buy, `{"selector":{"buyId":{"$gt":null}}}`).([]res.StockBuy)
		for _, b := range buys {
			record.Num = b.Num
			record.Name = b.Name
			record.Address = b.Address
			record.Stock = b.BuyStock
			record.Price = b.BuyPrice
			record.Date = b.BuyDate
			if b.SellPhone == parm["owner"] {
				record.Type = 1
			} else if b.BuyPhone == parm["owner"] {
				record.Type = 2
			} else {
				continue
			}
			records = append(records, record)
		}
		//回应数据
		rsp = peer.Response{Status: http.StatusOK, Message: strconv.Itoa(len(records))}
		if rsp.Payload, err = json.Marshal(records); err != nil {
			return utils.ErrorRsp(err)
		}
	default:
		return utils.ErrorRspf("Unknown command")
	}
	logger.Info("StockMg End-----------------------------------------------------------")
	return rsp
}

//用户数据处理
func (cc *StockCC) StockMgProc(rsc res.ResIf, opt, arg string) (rst res.ResIf, err error) {
	switch opt {
	case "AddRes":
	case "DelRes":
	case "UpgRes":
	case "GetRes":
	case "QryRes":
	default:
		return nil, utils.Errorf("Unkown operation")
	}
	return rsc, nil
}
