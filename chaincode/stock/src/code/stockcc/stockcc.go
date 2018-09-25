package main

import (
	"code"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"math/rand"
	"net/http"
	"res"
	"sup"
	"time"
	"utils"
)

//=============================================================================
type StockCC struct {
	code.CcBase
}

//=============================================================================
//初始化
func (cc *StockCC) Init(stub shim.ChaincodeStubInterface) (rsp peer.Response) {
	cc.SetStub(stub)
	cc.Logger.Info("==============================Init Start========================================")
	//错误处理
	defer func() {
		if e := recover(); e != nil {
			rsp = peer.Response{Status: http.StatusInternalServerError}
			rsp.Message = fmt.Sprintf("%v", e)
		}
	}()
	//添加系统管理用户
	_, cert := utils.GetCreator(stub)
	res.SYSMG_USER.Cert = cert
	sup.PutData(stub, res.SYSMG_USER)
	cc.Logger.Info("==============================Init End==========================================")
	return shim.Success(nil)
}

//=============================================================================
//调用接口
func (cc *StockCC) Invoke(stub shim.ChaincodeStubInterface) (rsp peer.Response) {
	cc.SetStub(stub)
	cc.Logger.Info("==============================Invoke Start======================================")
	//错误处理
	defer func() {
		if e := recover(); e != nil {
			rsp = peer.Response{Status: http.StatusInternalServerError}
			rsp.Message = fmt.Sprintf("%v", e)
		}
	}()
	//运行命令
	mod, args := stub.GetFunctionAndParameters()
	cc.Logger.Info(fmt.Sprintf("[%s] %s %s\r\n", mod, args[0], args[1]))
	rand.Seed(time.Now().Unix())
	cc.Mod = mod
	switch mod {
	case "UserMg": //用户管理
		cc.Cmd = args[0]
		cc.ResPFun = cc.UserMgProc
		rsp = cc.UserMg(args[0], args[1])
	case "StockMg": //股权管理
		cc.Cmd = args[0]
		cc.ResPFun = cc.StockMgProc
		rsp = cc.StockMg(args[0], args[1])
	default:
		return utils.ErrorRspf("Unknown method")
	}
	cc.Logger.Info(fmt.Sprintf("[RSP] Status:%d Message:%s Payload:%s",
		rsp.Status, rsp.Message, string(rsp.Payload)))
	cc.Logger.Info("==============================Invoke End========================================")
	return rsp
}

//=============================================================================
//主函数
func main() {
	stockcc := &StockCC{}
	stockcc.Setting("stockcc", shim.LogDebug)
	err := shim.Start(stockcc)
	if err != nil {
		fmt.Printf("Error starting StockCC chaincode: %s", err)
	}
}
