package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/jenner/chaincode/hycc/src/code"
	"github.com/jenner/chaincode/hycc/src/msic"
	"github.com/jenner/chaincode/hycc/src/proc"
	"github.com/jenner/chaincode/hycc/src/res"
	"github.com/jenner/chaincode/hycc/src/res/public"
	"github.com/jenner/chaincode/hycc/src/utils"
	"net/http"
	"strconv"
	"time"
)

//=============================================================================
type PubCC struct {
	code.CcBase
}

//=============================================================================
//初始化
func (cc *PubCC) Init(stub shim.ChaincodeStubInterface) (rsp peer.Response) {
	cc.SetStub(stub)
	cc.Logger.Info("==============================Init Start========================================")
	//错误处理
	defer func() {
		if e := recover(); e != nil {
			rsp = peer.Response{Status: http.StatusInternalServerError}
			rsp.Message = fmt.Sprintf("%v", e)
		}
	}()
	//验证用户信息
	mspid, cert := utils.GetCreator(stub)
	if mspid != pub.SYSMG_COMP.Mspid {
		//return utils.ErrorRspf("Init permission denied:%s", cert)
	}
	//查找是否有交易秘钥
	it, err := stub.GetStateByPartialCompositeKey("tkey", []string{})
	if err == nil && it.HasNext() {
		return shim.Success(nil)
	}
	//添加系统管理交易秘钥
	cc.Nonce = time.Now().Unix()
	rEcc := cc.GetUEcc()
	tEcc := utils.NewECC(cc.Nonce).GenKey()
	pub.SYSMG_TKEY.TPubKey = tEcc.GetPemPubKey()
	pub.SYSMG_TKEY.TPriKeys[rEcc.GetPemPubKey()] = rEcc.Encrypt(tEcc.PriKey.D.Bytes())
	pub.SYSMG_TKEY.RootSign = cc.rootSign(pub.SYSMG_TKEY, rEcc)
	msic.PutData(stub, pub.SYSMG_TKEY)
	//添加系统管理组织
	msic.PutData(stub, pub.SYSMG_COMP)
	//添加系统管理权限
	msic.PutData(stub, pub.SYSMG_PERM)
	//添加系统管理员角色
	msic.PutData(stub, pub.SYSMG_ROLE)
	//添加系统管理用户
	pub.SYSMG_USER.Cert = cert
	msic.PutData(stub, pub.SYSMG_USER)
	cc.Logger.Info("==============================Init End==========================================")
	return shim.Success(nil)
}

//=============================================================================
//调用接口
func (cc *PubCC) Invoke(stub shim.ChaincodeStubInterface) (rsp peer.Response) {
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
	case "CompMg": //组织管理
		arglen = 3
		rsc = pub.NewComp()
		cc.ResPFun = cc.CompProc
	case "PermMg": //权限管理
		arglen = 3
		rsc = pub.NewPerm()
		cc.ResPFun = cc.PermProc
	case "AuthMg": //鉴权管理
		arglen = 3
		rsc = pub.NewAuth()
		cc.ResPFun = cc.AuthProc
	case "RoleMg": //角色管理
		arglen = 3
		rsc = pub.NewRole()
		cc.ResPFun = cc.RoleProc
	case "UserMg": //用户管理
		arglen = 3
		rsc = pub.NewUser()
		cc.ResPFun = cc.UserProc
	case "TKeyMg": //交易秘钥管理
		arglen = 3
		rsc = pub.NewTKey()
		cc.ResPFun = cc.TKeyProc
	case "GetPerms": //获取权限
		arglen = 1
	case "OptAuth": //操作鉴权
		arglen = 2
	default:
		return utils.ErrorRspf("Unknown method")
	}
	//参数处理
	if len(args) != arglen {
		return utils.ErrorRspf("Should have %d param: %v", arglen, args)
	}
	var err error
	cc.Nonce, err = strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return utils.ErrorRsp(err)
	}
	//操作鉴权
	cc.Logger.Info("--------------------OptAuth--------------------------------")
	act := pub.Action{Mod: mod, Cmd: args[1], Args: args[2]}
	rst, err := cc.OptAuth(act.GetOpt())
	if err != nil {
		return utils.ErrorRsp(err)
	} else if rst == false {
		return utils.ErrorRspf("Opt permission denied")
	}
	//运行命令
	cc.Logger.Info("--------------------RunCmd---------------------------------")
	switch mod {
	case "CompMg": //公司管理
		fallthrough
	case "PermMg": //权限管理
		fallthrough
	case "AuthMg": //鉴权管理
		fallthrough
	case "RoleMg": //角色管理
		fallthrough
	case "UserMg": //用户管理
		rsp = cc.CommResMg(rsc, args[1], args[2])
	case "TKeyMg": //交易秘钥管理
		rsp = cc.TKeyResMg(rsc, args[1], args[2])
	case "GetPerms": //获取权限
		rsp = cc.SpanGetPerms()
	case "OptAuth": //操作鉴权
		rsp = cc.SpanOptAuth(args[1])
	default:
		return utils.ErrorRspf("Unknown method")
	}
	cc.Logger.Info("==============================Invoke End========================================")
	return rsp
}

//----------------------------------------------------------------------------
//一般资源管理
func (cc *PubCC) CommResMg(rsc res.ResIf, cmd, arg string) peer.Response {
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

//秘钥资源管理
func (cc *PubCC) TKeyResMg(rsc res.ResIf, cmd string, arg string) peer.Response {
	rsp := shim.Success(nil)
	switch cmd {
	case "new": //新建交易秘钥
		rsp = cc.tKeyNew(rsc, arg)
	case "add": //添加交易秘钥
		rsp = cc.tKeyAdd(rsc, arg)
	case "get": //获取交易秘钥
		rsp = cc.tKeyGet(rsc, arg)
	case "upg": //更新交易秘钥
		rsp = cc.tKeyUpg(rsc, arg)
	case "qry": //查询交易秘钥
		rsp = cc.tKeyQry(rsc, arg)
	case "adduser": //添加秘钥用户
		rsp = cc.tKeyAddUser(rsc, arg)
	case "deluser": //删除秘钥用户
		rsp = cc.tKeyDelUser(rsc, arg)
	case "upguser": //更新秘钥用户
		rsp = cc.tKeyAddUser(rsc, arg)
	default:
		return utils.ErrorRspf("Unknown command")
	}
	return rsp
}

//=============================================================================
//主函数
func main() {
	pubcc := &PubCC{}
	pubcc.Setting("pubcc", shim.LogDebug)
	err := shim.Start(pubcc)
	if err != nil {
		fmt.Printf("Error starting PubMg chaincode: %s", err)
	}
}
