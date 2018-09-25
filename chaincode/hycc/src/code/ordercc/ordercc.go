package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/jenner/chaincode/hycc/src/code"
	"github.com/jenner/chaincode/hycc/src/proc"
	"github.com/jenner/chaincode/hycc/src/res"
	"github.com/jenner/chaincode/hycc/src/res/order"
	//	"github.com/jenner/chaincode/hycc/src/res/public"
	"github.com/jenner/chaincode/hycc/src/utils"
	"net/http"
	"strconv"
)

//=============================================================================
type OrderCC struct {
	code.CcBase
}

//=============================================================================
//初始化
func (cc *OrderCC) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

//=============================================================================
//调用接口
func (cc *OrderCC) Invoke(stub shim.ChaincodeStubInterface) (rsp peer.Response) {
	cc.SetStub(stub)
	cc.Logger.Info("==============================Invoke Start======================================")
	//错误处理
	defer func() {
		if e := recover(); e != nil {
			rsp = peer.Response{Status: http.StatusInternalServerError}
			rsp.Message = fmt.Sprintf("%v", e)
		}
	}()
	//获取参数
	cc.Logger.Info("--------------------GetArgs--------------------------------")
	var arglen int
	var rsc res.ResIf
	mod, args := stub.GetFunctionAndParameters()
	switch mod {
	case "DrugMg": //药品管理
		arglen = 3
		rsc = order.NewDrug()
		cc.ResPFun = cc.DrugProc
	case "ExceptMg": //异常管理
		arglen = 3
		rsc = order.NewExcept()
		cc.ResPFun = cc.ExceptProc
	case "OrderMg": //订单管理
		arglen = 3
		rsc = order.NewOrder()
	default:
		return utils.ErrorRspf("Unknown method")
	}
	//参数处理
	if len(args) != arglen {
		return utils.ErrorRspf("Should have %d param: %+v", arglen, args)
	}
	var err error
	cc.Nonce, err = strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return utils.ErrorRsp(err)
	}
	//操作鉴权
	//	cc.Logger.Info("--------------------OptAuth-------------------------------")
	//	act := pub.Action{Mod: mod, Cmd: args[1], Args: args[2]}
	//	rst, err := cc.SpanOptAuth(act.GetOpt())
	//	if err != nil {
	//		return utils.ErrorRsp(err)
	//	} else if rst == false {
	//		return utils.ErrorRspf("Opt permission denied")
	//	}
	//执行命令
	cc.Logger.Info("--------------------RunCmd---------------------------------")
	switch mod {
	case "DrugMg": //药品管理
	case "ExceptMg": //异常管理
		rsp = cc.CommResMg(rsc, args[1], args[2])
	case "OrderMg": //订单管理
		rsp = cc.OrderResMg(rsc, args[1], args[2])
	default:
		return utils.ErrorRspf("Unknown method")
	}
	cc.Logger.Info("==============================Invoke End========================================")
	return rsp
}

//----------------------------------------------------------------------------
//一般资源管理
func (cc *OrderCC) CommResMg(rsc res.ResIf, cmd, arg string) peer.Response {
	rsp := shim.Success(nil)
	switch cmd {
	case "add": //添加资源
		rsp = proc.AddRes(cc, rsc, arg)
	case "del": //删除资源
		rsp = proc.DelRes(cc, rsc, arg)
	case "upg": //更新资源
		rsp = proc.UpgRes(cc, rsc, arg)
	case "get": //获取资源
		rsp = proc.GetRes(cc, rsc, arg)
	case "qry": //查询资源
		rsp = proc.QryRes(cc, rsc, arg)
	default:
		return utils.ErrorRspf("Unknown command")
	}
	return rsp
}

//订单资源管理
func (cc *OrderCC) OrderResMg(rsc res.ResIf, cmd string, arg string) peer.Response {
	rsp := shim.Success(nil)
	switch cmd {
	case "new": //新建订单
		rsp = cc.orderNew(rsc, arg)
	case "upg": //更新数据
	
	default:
		return utils.ErrorRspf("Unknown command")
	}
	return rsp
}

//=============================================================================
//主函数
func main() {
	ordercc := &OrderCC{}
	ordercc.Setting("ordercc", shim.LogDebug)
	err := shim.Start(ordercc)
	if err != nil {
		fmt.Printf("Error starting OrderMg chaincode: %s", err)
	}
}
