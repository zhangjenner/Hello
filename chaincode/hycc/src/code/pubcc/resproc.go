package main

import (
	"encoding/json"
	"github.com/jenner/chaincode/hycc/src/msic"
	"github.com/jenner/chaincode/hycc/src/res"
	"github.com/jenner/chaincode/hycc/src/res/public"
	"github.com/jenner/chaincode/hycc/src/utils"
)

//=============================================================================
//公司数据处理
func (cc *PubCC) CompProc(rsc res.ResIf, opt, arg string) (rst res.ResIf, err error) {
	stub := cc.GetStub()
	org, ok := rsc.(*pub.Comp)
	if !ok {
		return nil, utils.Errorf("The res's type is wrong")
	}
	switch opt {
	case "AddRes":
		//新建秘密并加密
		uEcc := utils.NewECC(cc.Nonce)
		rAes := utils.NewAES(cc.Nonce).GenKey()
		org.Encrypt(rAes)
		//为root用户赋权
		rKey := cc.getThenTKey(pub.SYSMG_COMP.Id, org.Cts)
		uEcc.LoadPemPubKey(rKey.TPubKey)
		org.Crypto[pub.SYSMG_COMP.Id] = uEcc.Encrypt(rAes.GetKey())
		//为所属组织赋权
		tKey := cc.getThenTKey(org.Id, org.Cts)
		uEcc.LoadPemPubKey(tKey.TPubKey)
		org.Crypto[org.Id] = uEcc.Encrypt(rAes.GetKey())
	case "DelRes":
	case "UpgRes":
		//获取数据密码
		oldComp := msic.GetData(stub, rsc).(*pub.Comp)
		rAes, err := cc.GetRAes(oldComp.Cts, oldComp.Crypto)
		if err != nil {
			return nil, utils.Error(err)
		}
		//更新数据处理
		oldComp.Decrypt(rAes)
		newRsc, err := UpgResProc(oldComp, arg)
		if err != nil {
			return nil, utils.Error(err)
		}
		newComp, _ := newRsc.(*pub.Comp)
		newComp.Encrypt(rAes)
		return newComp, nil
	case "GetRes":
		fallthrough
	case "QryRes":
		//获取数据密码
		rAes, err := cc.GetRAes(org.Cts, org.Crypto)
		if err != nil {
			return nil, utils.Error(err)
		}
		//解密数据
		org.Decrypt(rAes)
		return org, nil
	default:
		return nil, utils.Errorf("Unkown operation")
	}
	return org, nil
}

//权限数据处理
func (cc *PubCC) PermProc(rsc res.ResIf, opt, arg string) (rst res.ResIf, err error) {
	stub := cc.GetStub()
	perm, ok := rsc.(*pub.Perm)
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
	return perm, nil
}

//鉴权数据处理
func (cc *PubCC) AuthProc(rsc res.ResIf, opt, arg string) (rst res.ResIf, err error) {
	stub := cc.GetStub()
	auth, ok := rsc.(*pub.Auth)
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
	return auth, nil
}

//角色数据处理
func (cc *PubCC) RoleProc(rsc res.ResIf, opt, arg string) (rst res.ResIf, err error) {
	stub := cc.GetStub()
	role, ok := rsc.(*pub.Role)
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
	return role, nil
}

//用户数据处理
func (cc *PubCC) UserProc(rsc res.ResIf, opt, arg string) (rst res.ResIf, err error) {
	stub := cc.GetStub()
	user, ok := rsc.(*pub.User)
	if !ok {
		return nil, utils.Errorf("The res's type is wrong")
	}
	switch opt {
	case "AddRes":
		//新建秘密并加密
		uEcc := utils.NewECC(cc.Nonce)
		rAes := utils.NewAES(cc.Nonce).GenKey()
		user.Encrypt(rAes)
		//为root用户赋权
		rKey := cc.getThenTKey(pub.SYSMG_COMP.Id, user.Cts)
		uEcc.LoadPemPubKey(rKey.TPubKey)
		user.Crypto[pub.SYSMG_COMP.Id] = uEcc.Encrypt(rAes.GetKey())
		//为所属组织赋权
		tKey := cc.getThenTKey(user.Cid, user.Cts)
		uEcc.LoadPemPubKey(tKey.TPubKey)
		user.Crypto[user.Cid] = uEcc.Encrypt(rAes.GetKey())
	case "DelRes":
	case "UpgRes":
		//获取数据密码
		oldUser := msic.GetData(stub, rsc).(*pub.Comp)
		rAes, err := cc.GetRAes(oldUser.Cts, oldUser.Crypto)
		if err != nil {
			return nil, utils.Error(err)
		}
		//更新数据处理
		oldUser.Decrypt(rAes)
		newRsc, err := UpgResProc(oldUser, arg)
		if err != nil {
			return nil, utils.Error(err)
		}
		newUser, _ := newRsc.(*pub.User)
		newUser.Encrypt(rAes)
		return newUser, nil
	case "GetRes":
		fallthrough
	case "QryRes":
		//获取数据密码
		rAes, err := cc.GetRAes(user.Cts, user.Crypto)
		if err != nil {
			return nil, utils.Error(err)
		}
		//解密数据
		user.Decrypt(rAes)
		return user, nil
	default:
		return nil, utils.Errorf("Unkown operation")
	}
	return user, nil
}

//秘钥数据处理
func (cc *PubCC) TKeyProc(rsc res.ResIf, opt, arg string) (rst res.ResIf, err error) {
	stub := cc.GetStub()
	tKey, ok := rsc.(*pub.TKey)
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
	return tKey, nil
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

//获取数据密码
func (cc *PubCC) GetRAes(ts string, crypto map[string]string) (rAes *utils.AES, err error) {
	//获取提案者信息
	creator, err := cc.GetCreatorInfo()
	if err != nil {
		return nil, utils.Error(err)
	}
	//查询公司交易秘钥
	uEcc := cc.GetUEcc()
	tKey := cc.getThenTKey(creator.Cid, ts)
	tPriKey, ok := tKey.TPriKeys[uEcc.GetPemPubKey()]
	if !ok {
		return nil, utils.Errorf("Can't found the tPriKey of this user")
	}
	tEcc := utils.NewECC(cc.Nonce)
	tEcc.LoadBytePriKey(uEcc.Decrypt(tPriKey))
	//获取数据密码
	crAes, ok := crypto[creator.Cid]
	if !ok {
		return nil, utils.Errorf("Can't find the crypto to cid:%s", creator.Cid)
	}
	rAes = utils.NewAES(cc.Nonce)
	rAes.LoadKey(tEcc.Decrypt(crAes))
	return rAes, nil
}
