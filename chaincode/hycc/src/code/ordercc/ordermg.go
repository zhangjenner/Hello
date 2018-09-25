package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/jenner/chaincode/hycc/src/msic"
	"github.com/jenner/chaincode/hycc/src/res"
	"github.com/jenner/chaincode/hycc/src/res/order"
	//	"github.com/jenner/chaincode/hycc/src/res/public"
	"github.com/jenner/chaincode/hycc/src/utils"
	"net/http"
)

//=============================================================================
//新建交易秘钥
func (cc *OrderCC) orderNew(rsc res.ResIf, arg string) peer.Response {
		stub := cc.GetStub()
	logger := cc.GetLogger()
	logger.Info("==============================orderNew Start======================================")
	//参数验证
	rst, msg, err := res.HasField(rsc, res.FIELD_REQ, arg)
	if err != nil {
		return utils.ErrorRsp(err)
	} else if rst == false {
		return utils.ErrorRspf(msg)
	}
	//解析参数
	err = json.Unmarshal([]byte(arg), rsc)
	if err != nil {
		return utils.ErrorRsp(err)
	}
	ord, ok := rsc.(*order.Order)
	if !ok {
		return utils.ErrorRspf("The rsc type is not match order.Order")
	}
	//添加组织秘钥
	dAes := utils.NewAES(cc.Nonce).GenKey()
	tEcc := utils.NewECC(cc.Nonce)
	rKeyPub := string(cc.spanQryTKey("-1", ord.Cts, "pub"))
	tEcc.LoadPemPubKey(rKeyPub)
	ord.Crypto[rKeyPub] = tEcc.Encrypt(dAes.GetKey())
	tKeyPub := string(cc.spanQryTKey(ord.Cid, ord.Cts, "pub"))
	tEcc.LoadPemPubKey(tKeyPub)
	ord.Crypto[tKeyPub] = tEcc.Encrypt(dAes.GetKey())
	//数据签名
	uEcc := cc.GetUEcc()
	if !uEcc.HasPri {
		return utils.ErrorRspf("Not find private key, can't signature")
	}
	//	ord.Signs[uEcc.GetPemPubKey()] = uEcc.Sign([]byte("test"))
	//数据加密存储
	ord.Encrypt(dAes)
	msic.PutData(stub, ord)
	logger.Info("==============================orderNew End========================================")
	return shim.Success(nil)
}

//-----------------------------------------------------------------------------
//查询组织交易秘钥(跨链调用)
func (cc *OrderCC) spanQryTKey(cid, ts, typ string) []byte {
	stub := cc.GetStub()
	args := [][]byte{[]byte("TKeyMg"), []byte(fmt.Sprintf("%d", cc.Nonce)),
		[]byte(fmt.Sprintf("{\"cid\":\"%s\",\"ts\":\"%s\",\"type\":\"%s\"}", cid, ts, typ))}
	rsp := stub.InvokeChaincode(utils.PUBCC_NAME, args, utils.PUBCHAN_NAME)
	if rsp.Status != http.StatusOK {
		panic(utils.Errorf(rsp.Message))
	}
	return rsp.Payload
}
