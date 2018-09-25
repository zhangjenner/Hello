package main

import (
	"encoding/json"
	"github.com/jenner/chaincode/hycc/src/msic"
	"github.com/jenner/chaincode/hycc/src/res"
	"github.com/jenner/chaincode/hycc/src/res/order"
	"github.com/jenner/chaincode/hycc/src/utils"
)

//药品数据处理
func (cc *OrderCC) DrugProc(rsc res.ResIf, opt, arg string) (rst res.ResIf, err error) {
	stub := cc.GetStub()
	org, ok := rsc.(*order.Drug)
	if !ok {
		return nil, utils.Errorf("The res's type is wrong")
	}
	switch opt {
	case "AddRes":
	case "DelRes":
	case "UpgRes":
		oldRsc := msic.GetData(stub, rsc)
		return UpgResProc(oldRsc, arg)
	case "GetRes":
	case "QryRes":
	default:
		return nil, utils.Errorf("Unkown operation")
	}
	return org, nil
}

//异常数据处理
func (cc *OrderCC) ExceptProc(rsc res.ResIf, opt, arg string) (rst res.ResIf, err error) {
	stub := cc.GetStub()
	except, ok := rsc.(*order.Except)
	if !ok {
		return nil, utils.Errorf("The res's type is wrong")
	}
	switch opt {
	case "AddRes":
	case "DelRes":
	case "UpgRes":
		oldRsc := msic.GetData(stub, rsc)
		return UpgResProc(oldRsc, arg)
	case "GetRes":
	case "QryRes":
	default:
		return nil, utils.Errorf("Unkown operation")
	}
	return except, nil
}

//=============================================================================
//更新数据处理
func UpgResProc(oldRsc res.ResIf, arg string) (newRsc res.ResIf, err error) {
	//合并数据
	newRsc = res.NewRes(oldRsc, true)
	err = json.Unmarshal([]byte(arg), newRsc)
	if err != nil {
		return nil, utils.Error(err)
	}
	//数据验证
	rst, msg, err := res.CanMutable(newRsc, oldRsc)
	if err != nil {
		return nil, utils.Error(err)
	} else if rst == false {
		return nil, utils.Errorf(msg)
	}
	return newRsc, nil
}
