package sup

import (
	"code"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"net/http"
	"reflect"
	"res"
	"strconv"
	"strings"
	"unsafe"
	"utils"
)

//----------------------------------------------------------=================
//添加资源
func AddRes(cc code.CcIf, rsc res.ResIf, arg string) peer.Response {
	stub := cc.GetStub()
	logger := cc.GetLogger()
	logger.Info("AddRes Start----------------------------------------------------------")
	//参数验证
	if arg != "" {
		ParseRsc(&arg, res.FIELD_REQ, rsc)
	}
	//数据处理
	rsc, err := cc.ResProc(rsc, "AddRes", arg)
	if err != nil {
		return utils.ErrorRsp(err)
	}
	//数据查重
	rst, msg := HasDuplicate(stub, rsc)
	if rst == true {
		return utils.ErrorRspf(msg)
	}
	//存入数据
	PutData(stub, rsc)
	logger.Info("AddRes End------------------------------------------------------------")
	return shim.Success(nil)
}

//-----------------------------------------------------------------------------
//删除资源
func DelRes(cc code.CcIf, rsc res.ResIf, arg string) peer.Response {
	stub := cc.GetStub()
	logger := cc.GetLogger()
	logger.Info("DelRes Start----------------------------------------------------------")
	//参数验证
	if arg != "" {
		ParseRsc(&arg, res.FIELD_KEY, rsc)
	}
	//数据处理
	rsc, err := cc.ResProc(rsc, "DelRes", arg)
	if err != nil {
		return utils.ErrorRsp(err)
	}
	//删除数据
	DelData(stub, rsc)
	logger.Info("DelRes End------------------------------------------------------------")
	return shim.Success(nil)
}

//-----------------------------------------------------------------------------
//更新资源
func UpgRes(cc code.CcIf, rsc res.ResIf, arg string) peer.Response {
	stub := cc.GetStub()
	logger := cc.GetLogger()
	logger.Info("UpgRes Start----------------------------------------------------------")
	//参数验证
	if arg != "" {
		ParseRsc(&arg, res.FIELD_KEY, rsc)
	}
	//数据处理
	rsc, err := cc.ResProc(rsc, "UpgRes", arg)
	if err != nil {
		return utils.ErrorRsp(err)
	}
	//合并资源数据
	nrsc, err2 := CompRsc(rsc, &arg)
	if err2 != nil {
		return utils.ErrorRsp(err2)
	}
	//更新数据
	PutData(stub, nrsc)
	logger.Info("UpgRes End------------------------------------------------------------")
	return shim.Success(nil)
}

//-----------------------------------------------------------------------------
//获取资源
func GetRes(cc code.CcIf, rsc res.ResIf, arg string) peer.Response {
	stub := cc.GetStub()
	logger := cc.GetLogger()
	logger.Info("GetRes Start----------------------------------------------------------")
	//参数验证
	var err error
	if arg != "" {
		ParseRsc(&arg, res.FIELD_KEY, rsc)
	}
	//数据处理
	rsc = GetData(stub, rsc)
	rsc, err = cc.ResProc(rsc, "GetRes", arg)
	if err != nil {
		return utils.ErrorRsp(err)
	}
	//回应数据
	rsp := &peer.Response{Status: http.StatusOK, Message: ""}
	rsp.Payload, err = json.Marshal(rsc)
	if err != nil {
		return utils.ErrorRsp(err)
	}
	logger.Info("GetRes End------------------------------------------------------------")
	return *rsp
}

//-----------------------------------------------------------------------------
//查询资源
func QryRes(cc code.CcIf, rsc res.ResIf, qstr string) peer.Response {
	stub := cc.GetStub()
	logger := cc.GetLogger()
	logger.Info("QryRes Start----------------------------------------------------------")
	//查询数据
	if strings.Index(qstr, "\"selector\":{") == -1 {
		return utils.ErrorRspf("Can't find selector or selector syntax error")
	} else {
		qstr = strings.Replace(qstr, "\"selector\":{", fmt.Sprintf("\"selector\":{\"idx\":\"%s\",", rsc.GetIdx()), 1)
	}
	//数据处理
	rsts := QryData(stub, rsc, qstr)
	rstsval := reflect.ValueOf(rsts)
	nrsts := reflect.MakeSlice(rstsval.Type(), 0, 0)
	for i := 0; i < rstsval.Len(); i++ {
		val := rstsval.Index(i)
		pval := reflect.NewAt(val.Type(), unsafe.Pointer(val.UnsafeAddr()))
		prsc := pval.Interface().(res.ResIf)
		rsc, err := cc.ResProc(prsc, "QryRes", qstr)
		if err != nil {
			return utils.ErrorRsp(err)
		}
		nrsts = reflect.Append(nrsts, reflect.Indirect(reflect.ValueOf(rsc)))
	}
	//回应数据
	var err error
	rsp := &peer.Response{Status: http.StatusOK, Message: strconv.Itoa(rstsval.Len())}
	if rsp.Payload, err = json.Marshal(rsts); err != nil {
		return utils.ErrorRsp(err)
	}
	logger.Info("QryRes End------------------------------------------------------------")
	return *rsp
}
