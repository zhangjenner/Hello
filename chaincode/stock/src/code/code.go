package code

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"res"
	"utils"
)

//=============================================================================
type ResPFunT func(rsc res.ResIf, opt, arg string) (rst res.ResIf, err error) //资源处理函数

//=============================================================================
//链码通用接口
type CcIf interface {
	Setting(name string, loglevel shim.LoggingLevel)              //初始化
	GetStub() shim.ChaincodeStubInterface                         //获取stub
	GetLogger() *shim.ChaincodeLogger                             //获取日志记录器
	GetUEcc() (uEcc *utils.ECC)                                   //获取用户秘钥
	ResProc(rsc res.ResIf, opt, arg string) (rst res.ResIf, err error) //资源处理
}

//=============================================================================
//链码通用数据
type CcBase struct {
	CcIf
	Name    string                      //链码名字
	Mod		string						//执行模块
	Cmd		string						//执行命令
	Logger  *shim.ChaincodeLogger       //日志记录器
	Stub    shim.ChaincodeStubInterface //通讯存根
	ResPFun ResPFunT                    //资源处理函数
	UEcc    *utils.ECC                  //用户秘钥
}

//-----------------------------------------------------------------------------
//初始化
func (cc *CcBase) Setting(name string, loglevel shim.LoggingLevel) {
	cc.Name = name
	if cc.Logger == nil {
		cc.Logger = shim.NewLogger(name)
	}
	cc.Logger.SetLevel(loglevel)
}

//设置stub
func (cc *CcBase) SetStub(stub shim.ChaincodeStubInterface) {
	cc.Stub = stub
}

//-----------------------------------------------------------------------------
//获取stub
func (cc *CcBase) GetStub() shim.ChaincodeStubInterface {
	return cc.Stub
}

//获取日志记录器
func (cc *CcBase) GetLogger() *shim.ChaincodeLogger {
	return cc.Logger
}

//获取用户秘钥对
func (cc *CcBase) GetUEcc() (uEcc *utils.ECC) {
	//参数检查
	stub := cc.GetStub()
	if cc.UEcc != nil {
		return cc.UEcc
	}
	//加载公钥
	//	stub := cc.GetStub()
	cc.UEcc = utils.NewECC(123)
	_, uCert := utils.GetCreator(stub)
	uPubKey := utils.GetPubKeyFromCert(uCert)
	cc.UEcc.LoadPemPubKey(uPubKey)
	//加载私钥
	RootPriKey := `
-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgMM3BUEPPY14K1mtr
3JHd5nfT54LhPsaqzXOLroNzRpKhRANCAATryENZCkCAp+OzU1+xkWkEtbSVTIgM
O+fC7WEyuGAXfDeXnAYLbMEV4xHVzkFI6ex2EZLacbPI9SpBiRcZklSl
-----END PRIVATE KEY-----`
	cc.UEcc.LoadPemPriKey(RootPriKey)
	//		tMap, err := stub.GetTransient()
	//		if err != nil {
	//			return nil, utils.Error(err)
	//		}
	//		uPriKey, ok := tMap["PriKey"]
	//		if !ok {
	//			return nil, utils.Errorf("Not found 'PriKey' feild in Transient")
	//		}
	//		uEcc.LoadPemPriKey(string(uPriKey))	
	return cc.UEcc
}

//资源处理
func (cc *CcBase) ResProc(rsc res.ResIf, opt, arg string) (rst res.ResIf, err error) {
	return cc.ResPFun(rsc, opt, arg)
}
