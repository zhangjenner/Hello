package main

import (
	"crypto/md5"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"res"
	"sup"
	"time"
	"utils"
)

//=============================================================================
//用户管理
func (cc *StockCC) UserMg(cmd, arg string) peer.Response {
	stub := cc.GetStub()
	rsp := shim.Success(nil)
	logger := cc.GetLogger()
	logger.Info("UserMg Start----------------------------------------------------------")
	user := res.NewUser()
	switch cmd {
	case "register": //用户注册
		_, cert := utils.GetCreator(stub)
		sup.ParseRsc(&arg, res.FIELD_REQ, user)
		user.Cts = fmt.Sprint(time.Now().Unix())
		user.Cert = cert
		user.Pwd = fmt.Sprintf("%X", md5.Sum([]byte(user.Pwd)))
		user.Role = "user"
		rsp = sup.AddRes(cc, user, "")
	case "login": //用户登录
		sup.ParseRsc(&arg, res.FIELD_NIL, user)
		creator, err := cc.GetCreator()
		if err != nil {
			return utils.ErrorRsp(err)
		} else if creator.Phone != user.Phone {
			return utils.ErrorRspf("You are not '%s'", user.Phone)
		}
		rsp = sup.GetRes(cc, user, arg)
	case "query": //用户查询
		creator, err := cc.GetCreator()
		if err != nil {
			return utils.ErrorRsp(err)
		} else if creator.Role != "admin" {
			return utils.ErrorRspf("Permission denied")
		}
		rsp = sup.QryRes(cc, user, `{"selector":{"user":{"$gt":null}}}`)
	default:
		return utils.ErrorRspf("Unknown command")
	}
	logger.Info("UserMg End------------------------------------------------------------")
	return rsp
}

//用户数据处理
func (cc *StockCC) UserMgProc(rsc res.ResIf, opt, arg string) (rst res.ResIf, err error) {
	userbk := res.NewUser()
	user, ok := rsc.(*res.User)
	if !ok {
		return nil, utils.Errorf("The res's type is wrong")
	}
	switch opt {
	case "AddRes":
	case "DelRes":
	case "UpgRes":
	case "GetRes":
		if cc.Cmd == "login" {
			sup.ParseRsc(&arg, res.FIELD_KEY, userbk)
			if fmt.Sprintf("%X", md5.Sum([]byte(userbk.Pwd))) != user.Pwd {
				return nil, utils.Errorf("wrong password")
			}
			newUser := res.User{}
			newUser.User = user.User
			newUser.Phone = user.Phone
			newUser.Role = user.Role
			*user = newUser
		}
	case "QryRes":
		if cc.Cmd == "query" {
			newUser := res.User{}
			newUser.User = user.User
			newUser.Phone = user.Phone
			newUser.Role = user.Role
			*user = newUser
		}
	default:
		return nil, utils.Errorf("Unkown operation")
	}
	return user, nil
}
