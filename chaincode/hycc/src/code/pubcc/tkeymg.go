package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/jenner/chaincode/hycc/src/msic"
	"github.com/jenner/chaincode/hycc/src/res"
	"github.com/jenner/chaincode/hycc/src/res/public"
	"github.com/jenner/chaincode/hycc/src/utils"
	"net/http"
)

//=============================================================================
//新建交易秘钥
func (cc *PubCC) tKeyNew(rsc res.ResIf, arg string) peer.Response {
	stub := cc.GetStub()
	logger := cc.GetLogger()
	logger.Info("==============================tKeyNew Start=======================================")
	//root用户校验
	rEcc := cc.GetUEcc()
	rTKey := cc.getLastTkey("-1")
	if cc.rootVerify(&rTKey, rEcc) == false {
		return utils.ErrorRspf("Only the root user can invoke TKeyNew")
	}
	//解析参数并检查
	args := cc.checkArg(arg, "cid", "sTime", "eTime")
	tKeys := cc.getAllTkey(args["cid"])
	if len(tKeys) != 0 {
		return utils.ErrorRspf("Already has the TKey of cid: %s", args["cid"])
	}
	//新建交易秘钥
	tKey := pub.NewTKey()
	tKey.Cid = args["cid"]
	tKey.STime = args["sTime"]
	tKey.ETime = args["eTime"]
	tEcc := utils.NewECC(cc.Nonce).GenKey()
	tKey.TPubKey = tEcc.GetPemPubKey()
	tKey.TPriKeys[rEcc.GetPemPubKey()] = rEcc.Encrypt(tEcc.PriKey.D.Bytes())
	tKey.RootSign = cc.rootSign(tKey, rEcc)
	msic.PutData(stub, tKey)
	logger.Info("==============================tKeyNew End=========================================")
	return shim.Success(nil)
}

//添加交易秘钥
func (cc *PubCC) tKeyAdd(rsc res.ResIf, arg string) peer.Response {
	stub := cc.GetStub()
	logger := cc.GetLogger()
	logger.Info("==============================tKeyAdd Start=======================================")
	//root用户校验
	rEcc := cc.GetUEcc()
	rTKey := cc.getLastTkey("-1")
	if cc.rootVerify(&rTKey, rEcc) == false {
		return utils.ErrorRspf("Only the root user can invoke TKeyAdd")
	}
	//解析参数并检查
	args := cc.checkArg(arg, "cid", "sTime", "eTime")
	lTKey := cc.getLastTkey(args["cid"])
	if args["sTime"] != lTKey.ETime {
		return utils.ErrorRspf("The STime must equal to recent ETime")
	}
	//新建交易秘钥
	tKey := pub.NewTKey()
	tKey.Cid = args["cid"]
	tKey.STime = args["sTime"]
	tKey.ETime = args["eTime"]
	uEcc := utils.NewECC(cc.Nonce)
	tEcc := utils.NewECC(cc.Nonce).GenKey()
	tKey.TPubKey = tEcc.GetPemPubKey()
	for uPubKey, _ := range lTKey.TPriKeys {
		uEcc.LoadPemPubKey(uPubKey)
		tKey.TPriKeys[uPubKey] = uEcc.Encrypt(tEcc.PriKey.D.Bytes())
	}
	tKey.RootSign = cc.rootSign(tKey, rEcc)
	msic.PutData(stub, tKey)
	logger.Info("==============================tKeyAdd End=========================================")
	return shim.Success(nil)
}

//获取组织交易秘钥
func (cc *PubCC) tKeyGet(rsc res.ResIf, arg string) peer.Response {
	logger := cc.GetLogger()
	logger.Info("==============================tKeyQry Start=======================================")
	//root用户校验
	rEcc := cc.GetUEcc()
	rTKey := cc.getLastTkey("-1")
	if cc.rootVerify(&rTKey, rEcc) == false {
		return utils.ErrorRspf("Only the root user can invoke TKeyUpg")
	}
	//获取交易秘钥
	var err error
	args := cc.checkArg(arg, "cid")
	lTKey := cc.getLastTkey(args["cid"])
	rsp := peer.Response{Status: http.StatusOK, Message: "1"}
	rsp.Payload, err = json.Marshal(lTKey)
	if err != nil {
		panic(utils.ErrorRsp(err))
	}
	logger.Info("==============================tKeyQry End=========================================")
	return rsp
}

//更新交易秘钥
func (cc *PubCC) tKeyUpg(rsc res.ResIf, arg string) peer.Response {
	stub := cc.GetStub()
	logger := cc.GetLogger()
	logger.Info("==============================tKeyUpg Start=======================================")
	//root用户校验
	rEcc := cc.GetUEcc()
	rTKey := cc.getLastTkey("-1")
	if cc.rootVerify(&rTKey, rEcc) == false {
		return utils.ErrorRspf("Only the root user can invoke TKeyUpg")
	}
	//数据检查
	args := cc.checkArg(arg, "cid", "sTime", "eTime")
	lTKey := cc.getLastTkey(args["cid"])
	if lTKey.STime != args["sTime"] {
		return utils.ErrorRspf("Can just upgrade the recent tKey")
	}
	//数据更新
	lTKey.ETime = args["eTime"]
	lTKey.RootSign = cc.rootSign(&lTKey, rEcc)
	msic.PutData(stub, &lTKey)
	logger.Info("==============================tKeyUpg End=========================================")
	return shim.Success(nil)
}

//查询组织交易秘钥
func (cc *PubCC) tKeyQry(rsc res.ResIf, arg string) peer.Response {
	logger := cc.GetLogger()
	logger.Info("==============================tKeyQry Start=======================================")
	//查询组织交易秘钥
	args := cc.checkArg(arg, "cid", "ts", "type")
	tKey := cc.getThenTKey(args["cid"], args["ts"])
	//类型判断
	rsp := peer.Response{Status: http.StatusOK}
	switch args["type"] {
	case "pub":
		rsp.Payload = []byte(tKey.TPubKey)
	case "pri":
		uEcc := cc.GetUEcc()
		tPriKey, ok := tKey.TPriKeys[uEcc.GetPemPubKey()]
		if !ok {
			return utils.ErrorRspf("Can't found the tPriKey of this user")
		}
		rsp.Payload = uEcc.Decrypt(tPriKey)
	default:
		return utils.ErrorRspf("The type of tKey is error: %s", args["type"])
	}
	logger.Info("==============================tKeyQry End=========================================")
	return rsp
}

//添加秘钥用户
func (cc *PubCC) tKeyAddUser(rsc res.ResIf, arg string) peer.Response {
	stub := cc.GetStub()
	logger := cc.GetLogger()
	logger.Info("==============================tKeyAddUser Start===================================")
	//root用户校验
	rEcc := cc.GetUEcc()
	rTKey := cc.getLastTkey("-1")
	if cc.rootVerify(&rTKey, rEcc) == false {
		return utils.ErrorRspf("Only the root user can invoke TKeyAddUser")
	}
	//解析参数
	args := cc.checkArg(arg, "cid", "uid")
	user := pub.NewUser()
	user.Id = args["uid"]
	user = msic.GetData(stub, user).(*pub.User)
	if user.Cid != args["cid"] {
		return utils.ErrorRspf("The user:%s is not bylong to the company:%d", args["uid"], args["cid"])
	}
	//查找交易秘钥
	tKeys := cc.getAllTkey(args["cid"])
	if len(tKeys) == 0 {
		return utils.ErrorRspf("There is no TKey of cid: %s", args["cid"])
	}
	//添加秘钥用户
	rPubKey := rEcc.GetPemPubKey()
	uPubKey := utils.GetPubKeyFromCert(user.Cert)
	uEcc := utils.NewECC(cc.Nonce)
	uEcc.LoadPemPubKey(uPubKey)
	for _, tKey := range tKeys {
		enTPriKey, ok := tKey.TPriKeys[rPubKey]
		if !ok {
			return utils.ErrorRspf("The TPriKeys has no root user's TKey")
		} else {
			tPriKey := rEcc.Decrypt(enTPriKey)
			tKey.TPriKeys[uPubKey] = uEcc.Encrypt(tPriKey)
			msic.PutData(stub, &tKey)
		}
	}
	logger.Info("==============================tKeyAddUser End=====================================")
	return shim.Success(nil)
}

//删除秘钥用户
func (cc *PubCC) tKeyDelUser(rsc res.ResIf, arg string) peer.Response {
	stub := cc.GetStub()
	logger := cc.GetLogger()
	logger.Info("==============================tKeyDelUser Start===================================")
	//root用户校验
	rEcc := cc.GetUEcc()
	rTKey := cc.getLastTkey("-1")
	if cc.rootVerify(&rTKey, rEcc) == false {
		return utils.ErrorRspf("Only the root user can invoke TKeyAddUser")
	}
	//解析参数
	args := cc.checkArg(arg, "cid", "uid")
	user := pub.NewUser()
	user.Id = args["uid"]
	user = msic.GetData(stub, user).(*pub.User)
	if user.Cid != args["cid"] {
		return utils.ErrorRspf("The user:%s is not bylong to the company:%d", args["uid"], args["cid"])
	}
	//查找交易秘钥
	tKeys := cc.getAllTkey(args["cid"])
	if len(tKeys) == 0 {
		return utils.ErrorRspf("There is no TKey of cid: %s", args["cid"])
	}
	//删除秘钥用户
	uPubKey := utils.GetPubKeyFromCert(user.Cert)
	for _, tKey := range tKeys {
		if _, ok := tKey.TPriKeys[uPubKey]; ok {
			delete(tKey.TPriKeys, uPubKey)
			msic.PutData(stub, &tKey)
		}
	}
	logger.Info("==============================tKeyDelUser End=====================================")
	return shim.Success(nil)
}

//-----------------------------------------------------------------------------
//参数检查
func (cc *PubCC) checkArg(arg string, feilds ...string) map[string]string {
	args := make(map[string]string)
	err := json.Unmarshal([]byte(arg), &args)
	if err != nil {
		panic(utils.ErrorRsp(err))
	}
	for _, feild := range feilds {
		if _, ok := args[feild]; !ok {
			panic(utils.ErrorRspf("Can't found '%s' feild in arg", feild))
		}
	}
	return args
}

//获取当时的交易秘钥
func (cc *PubCC) getThenTKey(cid, ts string) pub.TKey {
	stub := cc.GetStub()
	tKey := pub.NewTKey()
	qstr := fmt.Sprintf("{\"selector\":{\"$and\":[{\"idx\":\"%s\"},{\"cid\":\"%s\"},"+
		"{\"sTime\":{\"$lte\":\"%s\"}},{\"eTime\":{\"$gt\":\"%s\"}}]}}", tKey.Idx, cid, ts, ts)
	tKeys, ok := msic.QryData(stub, tKey, qstr).([]pub.TKey)
	if !ok {
		panic(utils.Errorf("The result of QryData is wrong"))
	} else if len(tKeys) == 0 {
		panic(utils.Errorf("Can't found the then tKey of %s-%s", cid, ts))
	} else if len(tKeys) > 1 {
		panic(utils.Errorf("Find too many then tKey of %s-%s", cid, ts))
	}
	return tKeys[0]
}

//获取最后的交易秘钥
func (cc *PubCC) getLastTkey(cid string) pub.TKey {
	tKeys := cc.getAllTkey(cid)
	if len(tKeys) == 0 {
		panic(utils.ErrorRspf("Can't find any TKey of cid: %s", cid))
	}
	num, eTime := 0, "0"
	for i := 0; i < len(tKeys); i++ {
		if tKeys[i].ETime > eTime {
			num = i
			eTime = tKeys[i].ETime
		}
	}
	return tKeys[num]
}

//获取所有交易秘钥
func (cc *PubCC) getAllTkey(cid string) []pub.TKey {
	stub := cc.GetStub()
	tKey := pub.NewTKey()
	qstr := fmt.Sprintf("{\"selector\":{\"idx\":\"%s\",\"cid\":\"%s\"}}", tKey.Idx, cid)
	tKeys, ok := msic.QryData(stub, tKey, qstr).([]pub.TKey)
	if !ok {
		panic(utils.ErrorRspf("The result of QryData is wrong"))
	}
	return tKeys
}

//root用户签名
func (cc *PubCC) rootSign(tKey *pub.TKey, rEcc *utils.ECC) (sign string) {
	tKey.RootSign = ""
	text, err := json.Marshal(tKey)
	if err != nil {
		panic(utils.Error(err))
	}
	return rEcc.Sign(text)
}

//root用户校验
func (cc *PubCC) rootVerify(tKey *pub.TKey, rEcc *utils.ECC) (valid bool) {
	signature := tKey.RootSign
	tKey.RootSign = ""
	text, err := json.Marshal(tKey)
	if err != nil {
		panic(utils.Error(err))
	}
	return rEcc.Verify(text, signature)
}
