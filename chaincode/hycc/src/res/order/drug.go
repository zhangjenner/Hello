package order

import (
	"github.com/jenner/chaincode/hycc/src/res"
	"github.com/jenner/chaincode/hycc/src/utils"
	"github.com/pkg/errors"
)

//=============================================================================
//药品数据
type Drug struct {
	res.ResBase              //基础数据
	Id                  string  `json:"id,omitempty"`                  //
	Cid                 string  `json:"cid,omitempty"`                 //公司id
	Code                string  `json:"code,omitempty"`                //药品编码
	InfoName            string  `json:"infoName,omitempty"`            //药品通用名称
	InfoNamed           string  `json:"infoNamed,omitempty"`           //药品商品名称
	InfoSpecification   string  `json:"infoSpecification,omitempty"`   //规格
	InfoForm            string  `json:"infoForm,omitempty"`            //剂型:如胶囊片剂
	InfoPrice           float32 `json:"infoPrice,omitempty"`           //单价
	InfoUnit            int32   `json:"infoUnit,omitempty"`            //单位,1盒,2包
	IfPrescription      int32   `json:"ifPrescription,omitempty"`      //是否处方药,1是,0不是
	IfRawMaterial       int32   `json:"ifRawMaterial,omitempty"`       //是否原料药,1是,0不是
	TrpInnerPacking     string  `json:"trpInnerPacking,omitempty"`     //包装规格，例如:18粒/板*2板/盒
	TrpOuterNumber      int32   `json:"trpOuterNumber,omitempty"`      //运输/存储外包装数量
	TrpOuterPackingUnit int32   `json:"trpOuterPackingUnit,omitempty"` //运输/存储外包装单位,1盒.2包
	TrpUnit             int32   `json:"trpUnit,omitempty"`             //运输/存储包装单位,1箱,2包
	TrpUnitPrice        float32 `json:"trpUnitPrice,omitempty"`        //运输/存储包装单价
	TrpDamagePrice      float32 `json:"trpDamagePrice,omitempty"`      //货损单价
	TrpWeight           float32 `json:"trpWeight,omitempty"`           //包装重量kg
	TrpBulk             float32 `json:"trpBulk,omitempty"`             //包装体积M3
	TrpLength           float32 `json:"trpLength,omitempty"`           //长,单位mm
	TrpWidth            float32 `json:"trpWidth,omitempty"`            //宽,单位mm
	TrpHeight           float32 `json:"trpHeight,omitempty"`           //高,单位mm
	TrpStorageCond      string  `json:"trpStorageCond,omitempty"`      //药品储存条件,常温,冷藏,阴冷
	FacCompanyName      string  `json:"facCompanyName,omitempty"`      //生产企业名称
	FacBatchNumber      string  `json:"facBatchNumber,omitempty"`      //产品批号
	FacApprovalNumber   string  `json:"facApprovalNumber,omitempty"`   //批准文号
	FacProducedDate     string  `json:"facProducedDate,omitempty"`     //生产日期
	FacValidityDate     string  `json:"facValidityDate,omitempty"`     //有效期至
	Creator             string  `json:"creator,omitempty"`             //
	CreateDateTime      string  `json:"createDateTime,omitempty"`      //
	Updator             string  `json:"updator,omitempty"`             //
	UpdateDateTime      string  `json:"updateDateTime,omitempty"`      //
	Disable             int32   `json:"disable,omitempty"`             //0
	DeleteFlag          int32   `json:"deleteFlag,omitempty"`          //0
}

//-----------------------------------------------------------------------------
//构造函数
func NewDrug() *Drug {
	return &Drug{ResBase: res.ResBase{Idx: "drug"}}
}

//获取索引
func (self *Drug) GetIdx() string {
	return self.Idx
}

//获取字段
func (self *Drug) GetField(ftype res.FieldType) (fields []string, err error) {
	switch ftype {
	case res.FIELD_KEY: //子建字段
		fields = []string{"Id"}
	case res.FIELD_UNQ: //独有字段
		fields = []string{"Id"}
	case res.FIELD_REQ: //必选字段
		fields = []string{"Id"}
	case res.FIELD_CST: //常量字段
		fields = []string{"Idx", "Id"}
	default:
		return []string{}, errors.Errorf("Unknown field type:%v", ftype)
	}
	return fields, nil
}

//资源加密
func (self *Drug) Encrypt(rAes *utils.AES) {
}

//资源解密
func (self *Drug) Decrypt(rAes *utils.AES) {
}

