package order

import (
	"github.com/jenner/chaincode/hycc/src/res"
	"github.com/jenner/chaincode/hycc/src/utils"
	"github.com/pkg/errors"
)

//=============================================================================
//药品数据
type Except struct {
	res.ResBase                 //基础数据
	Id                   string `json:"id,omitempty"`                   //异常单号
	Ordercode            string `json:"ordercode,omitempty"`            //
	Level                int32  `json:"level,omitempty"`                //预警级别,1黄,2橙,3红
	Type                 int32  `json:"type,omitempty"`                 //异常类型,1灭失,2货损,3货不对单,4超时未送达,5变更提货要求
	Status               int32  `json:"status,omitempty"`               //处理状态,0未处理,1处理中,2已处理
	SourceChannel        int32  `json:"sourceChannel,omitempty"`        //来源渠道,0系统超时未送达,1客户投诉,2专线,3华药,
	SourceModule         int32  `json:"sourceModule,omitempty"`         //来源系统模块,1投诉
	SourceModuleId       string `json:"sourceModuleId,omitempty"`       //来源系统模块的id,比如投诉模块的id
	Address              string `json:"address,omitempty"`              //提交地址
	IfSendError          int32  `json:"ifSendError,omitempty"`          //是否配送过程中产生的异常,1是,0不是(仓储)
	DutyType             int32  `json:"dutyType,omitempty"`             //责任主体类型,0个人,1公司
	DutyCompanyId        string `json:"dutyCompanyId,omitempty"`        //责任公司id
	DutyName             string `json:"dutyName,omitempty"`             //责任人姓名
	NewCompensationOrder int32  `json:"newCompensationOrder,omitempty"` //是否勾选赔偿单,1是,0否
	NewReturnOrder       int32  `json:"newReturnOrder,omitempty"`       //是否勾选退换货单,1是,0否
	NewChangeOrder       int32  `json:"newChangeOrder,omitempty"`       //是否勾选原始订单,1是,0否
	CompensationCode     string `json:"compensationCode,omitempty"`     //赔偿单code
	ReturnCode           string `json:"returnCode,omitempty"`           //退换货单code
	ChangeCode           string `json:"changeCode,omitempty"`           //订单变更code
	Description          string `json:"description,omitempty"`          //处理描述
	Creator              string `json:"creator,omitempty"`              //
	CreateDateTime       string `json:"createDateTime,omitempty"`       //
	Updator              string `json:"updator,omitempty"`              //
	UpdateDateTime       string `json:"updateDateTime,omitempty"`       //
	SendCompanyId        string `json:"sendCompanyId,omitempty"`        //原始订单发货公司id
}

//-----------------------------------------------------------------------------
//构造函数
func NewExcept() *Except {
	return &Except{ResBase: res.ResBase{Idx: "except"}}
}

//获取索引
func (self *Except) GetIdx() string {
	return self.Idx
}

//获取字段
func (self *Except) GetField(ftype res.FieldType) (fields []string, err error) {
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
func (self *Except) Encrypt(rAes *utils.AES) {
}

//资源解密
func (self *Except) Decrypt(rAes *utils.AES) {
}
