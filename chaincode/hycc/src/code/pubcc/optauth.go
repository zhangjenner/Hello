package main

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/jenner/chaincode/hycc/src/msic"
	"github.com/jenner/chaincode/hycc/src/res/public"
	"github.com/jenner/chaincode/hycc/src/utils"
	"net/http"
	"sort"
	"strconv"
)

//操作鉴权
func (cc *PubCC) OptAuth(opt string) (rst bool, err error) {
	//获取用户权限
	stub := cc.GetStub()
	perms, err := cc.GetPerms()
	if err != nil {
		return false, err
	}
	//查看是否为系统管理员
	sort.Strings(perms)
	idx := sort.SearchStrings(perms, pub.SYSMG_PERM.Permission)
	if idx < len(perms) && perms[idx] == pub.SYSMG_PERM.Permission {
		return true, nil
	}
	//获取操作所需权
	auth := pub.NewAuth()
	auth.Opt = opt
	auth = msic.GetData(stub, auth).(*pub.Auth)
	//权限鉴别
	exp := msic.NewCondExp(auth.Pexp)
	rst, err = exp.Calc(perms)
	if err != nil {
		return false, utils.Error(err)
	}
	return rst, nil
}

//获取用户权限
func (cc *PubCC) GetPerms() (perms []string, err error) {
	stub := cc.GetStub()
	//获取用户信息
	user, err := cc.GetCreatorInfo()
	if err != nil {
		return nil, err
	}
	//查找角色
	roles := []pub.Role{}
	role := pub.NewRole()
	for _, role.Id = range user.RoleIds {
		role = msic.GetData(stub, role).(*pub.Role)
		roles = append(roles, *role)
	}
	//查找权限
	for _, role := range roles {
		perm := pub.NewPerm()
		for _, perm.Id = range role.PermIds {
			perm = msic.GetData(stub, perm).(*pub.Perm)
			sort.Strings(perms)
			if sort.SearchStrings(perms, perm.Permission) == len(perms) {
				perms = append(perms, perm.Permission)
			}
		}
	}
	//返回结果
	rst := peer.Response{Status: http.StatusOK, Message: strconv.Itoa(len(perms))}
	rst.Payload, err = json.Marshal(perms)
	if err != nil {
		return nil, utils.Error(err)
	}
	return perms, nil
}

//获取提案者信息
func (cc *PubCC) GetCreatorInfo() (user *pub.User, err error) {
	stub := cc.GetStub()
	_, cert := utils.GetCreator(stub)
	user = pub.NewUser()
	user.Cert = cert
	users, ok := msic.SelectRes(stub, user).([]pub.User)
	if !ok {
		return nil, utils.Errorf("The result of SelectRes is wrong")
	} else if len(users) != 1 {
		return nil, utils.Errorf("The number of user whose cert is '%s' is wrong ", user.Cert)
	}
	return &users[0], nil
}

//操作鉴权(跨链调用)
func (cc *PubCC) SpanOptAuth(opt string) peer.Response {
	rst, err := cc.OptAuth(opt)
	if err != nil {
		return utils.ErrorRsp(err)
	} else if rst == false {
		return utils.ErrorRspf("Opt permission denied")
	}
	return shim.Success(nil)
}

//获取用户权限(跨链调用)
func (cc *PubCC) SpanGetPerms() peer.Response {
	perms, err := cc.GetPerms()
	if err != nil {
		return utils.ErrorRsp(err)
	}
	jperms, err := json.Marshal(perms)
	if err != nil {
		return utils.ErrorRsp(err)
	}
	return peer.Response{Status: http.StatusOK, Message: strconv.Itoa(len(perms)), Payload: jperms}
}
