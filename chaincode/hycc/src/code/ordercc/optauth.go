package main

import (
	"encoding/json"
	"fmt"
	"github.com/jenner/chaincode/hycc/src/utils"
	"net/http"
)

//跨链操作鉴权
func (cc *OrderCC) SpanOptAuth(opt string) (rst bool, err error) {
	stub := cc.GetStub()
	args := [][]byte{[]byte("OptAuth"), []byte(fmt.Sprintf("%d", cc.Nonce)), []byte(opt)}
	rsp := stub.InvokeChaincode(utils.PUBCC_NAME, args, utils.PUBCHAN_NAME)
	if rsp.Status == http.StatusOK {
		return true, nil
	}
	return false, utils.Errorf(rsp.Message)
}

//跨链获取用户权限
func (cc *OrderCC) SpanGetPerms() (perms []string, err error) {
	stub := cc.GetStub()
	args := [][]byte{[]byte("GetPerms"), []byte(fmt.Sprintf("%d", cc.Nonce))}
	rsp := stub.InvokeChaincode(utils.PUBCC_NAME, args, utils.PUBCHAN_NAME)
	if rsp.Status == http.StatusOK {
		perms := []string{}
		if err := json.Unmarshal(rsp.Payload, &perms); err != nil {
			return nil, utils.Error(err)
		}
		return perms, nil
	}
	return nil, utils.Errorf(rsp.Message)
}
