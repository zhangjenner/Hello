package order

import (
	"encoding/json"
	"github.com/jenner/chaincode/hycc/src/res"
	"github.com/jenner/chaincode/hycc/src/utils"
	"github.com/pkg/errors"
)

//=============================================================================
//签名数据
type OrderSign struct {
	UPubKey  string `json:"uPubKey,omitempty"`  //用户公钥
	UserSign string `json:"userSign,omitempty"` //用户签名
}

//附加数据
type OrderAtch struct {
	Crypto map[string]string    `json:"crypto,omitempty"` //秘钥
	Signs  map[string]OrderSign `json:"signs,omitempty"`  //签名
}

//药品目录
type OrderItem struct {
	Id             string  `json:"id,omitempty"`             //
	Ordercode      string  `json:"ordercode,omitempty"`      //订单code
	WarehouseId    string  `json:"warehouseId,omitempty"`    //仓库ID
	Locationid     string  `json:"locationid,omitempty"`     //库位id
	Drugid         string  `json:"drugid,omitempty"`         //药品id
	Num            int32   `json:"num,omitempty"`            //下单数量(箱)
	CodeSection    string  `json:"codeSection,omitempty"`    //区域码区间，例如：1-100,120-150,180-250
	CompensateCost float32 `json:"compensateCost,omitempty"` //赔偿单,赔偿金额
}

//加密数据
type OrderECD struct {
	TotalBox             int32       `json:"totalBox,omitempty"`             //订单总箱数
	TotalWeight          float32     `json:"totalWeight,omitempty"`          //订单总重量
	TotalBulk            float32     `json:"totalBulk,omitempty"`            //订单总体积
	DrugTotalPrice       float32     `json:"drugTotalPrice,omitempty"`       //药品总价
	CostTotal            float32     `json:"costTotal,omitempty"`            //订单总额
	CostStorage          float32     `json:"costStorage,omitempty"`          //仓储费
	CostFetch            float32     `json:"costFetch,omitempty"`            //取货费
	CostMainLine         float32     `json:"costMainLine,omitempty"`         //干线运输费
	CostBranchLine       float32     `json:"costBranchLine,omitempty"`       //中转费
	CostSend             float32     `json:"costSend,omitempty"`             //配送费
	CostDamage           float32     `json:"costDamage,omitempty"`           //货损赔偿金
	CostManage           float32     `json:"costManage,omitempty"`           //管理费
	CostProtect          float32     `json:"costProtect,omitempty"`          //保价费
	CostLoadUnload       float32     `json:"costLoadUnload,omitempty"`       //装卸费
	TrpBeginCityid       string      `json:"trpBeginCityid,omitempty"`       //起点:城市
	TrpEndCityid         string      `json:"trpEndCityid,omitempty"`         //终点:城市
	SendProid            string      `json:"sendProid,omitempty"`            //发送省id
	SendCityid           string      `json:"sendCityid,omitempty"`           //发送市id
	SendAreaid           string      `json:"sendAreaid,omitempty"`           //发送区域id
	SendFullAreaName     string      `json:"sendFullAreaName,omitempty"`     //
	SendName             string      `json:"sendName,omitempty"`             //发货人姓名
	SendCompanyName      string      `json:"sendCompanyName,omitempty"`      //发货人公司
	SendMobile           string      `json:"sendMobile,omitempty"`           //发货人手机
	SendPhone            string      `json:"sendPhone,omitempty"`            //发货人座机
	SendAddress          string      `json:"sendAddress,omitempty"`          //发货人详细地址
	SendLng              float32     `json:"sendLng,omitempty"`              //发送人经度
	SendLat              float32     `json:"sendLat,omitempty"`              //发送人纬度
	ReceiveName          string      `json:"receiveName,omitempty"`          //收货人姓名
	ReceiveProid         string      `json:"receiveProid,omitempty"`         //接收省id
	ReceiveCityid        string      `json:"receiveCityid,omitempty"`        //接收市id
	ReceiveAreaid        string      `json:"receiveAreaid,omitempty"`        //接收区域id
	ReceiveFullAreaName  string      `json:"receiveFullAreaName,omitempty"`  //
	ReceiveCompanyName   string      `json:"receiveCompanyName,omitempty"`   //收货人公司
	ReceiveWarehouseName string      `json:"receiveWarehouseName,omitempty"` //收货人仓库名称
	ReceiveMobile        string      `json:"receiveMobile,omitempty"`        //收货人手机
	ReceivePhone         string      `json:"receivePhone,omitempty"`         //收货人座机
	ReceiveAddress       string      `json:"receiveAddress,omitempty"`       //收货人详细地址
	ReceiveLng           float32     `json:"receiveLng,omitempty"`           //收获人经度
	ReceiveLat           float32     `json:"receiveLat,omitempty"`           //收货人纬度
	OrderItems           []OrderItem `json:"orderItems,omitempty"`           //关联药品
}

//订单数据
type Order struct {
	res.ResBase              //基础数据
	Code              string `json:"code,omitempty"`              //订单号
	Cid               string `json:"cid,omitempty"`               //企业id
	SourceCompanyId   string `json:"sourceCompanyId,omitempty"`   //归属于哪个华药分公司
	Remark            string `json:"remark,omitempty"`            //订单备注
	OrderType         int32  `json:"orderType,omitempty"`         //订单类型，1销售订单，2仓储订单,3临时订单,4退换货单
	OrderStatus       int32  `json:"orderStatus,omitempty"`       //订单状态:0待审核1未派单2派单中3出库中4运输中5已签收并回单中6已回单
	TrpDeliveryWay    int32  `json:"trpDeliveryWay,omitempty"`    //交货方式1送货上门,2营业点自提,3代客卸货(平地卸货),4代客卸货(有电梯),5代客卸货(无电梯)
	TrpVas            int32  `json:"trpVas,omitempty"`            //是否需要增值服务(货物保价),1是,0否
	TrpFreightPayment int32  `json:"trpFreightPayment,omitempty"` //运费付款方式,1发货方付2收货方付(到付)
	TrpNeedStevedore  int32  `json:"trpNeedStevedore,omitempty"`  //是否需要装卸工.1是,0不是
	TrpNeedTrucks     int32  `json:"trpNeedTrucks,omitempty"`     //是否需要调派车辆,1是,0不是
	TrpReturnClaim    int32  `json:"trpReturnClaim,omitempty"`    //回单要求,1原件(签字+身份证)个人,2原件(签字+盖签收章)3,无需返单(电子回单)
	TrpReturnDateTime string `json:"trpReturnDateTime,omitempty"` //客户要求最迟返单时间
	Driverid          string `json:"driverid,omitempty"`          //司机id
	Creator           string `json:"creator,omitempty"`           //创建人id
	CreateDateTime    string `json:"createDateTime,omitempty"`    //创建时间
	Updator           string `json:"updator,omitempty"`           //修改人
	UpdateDateTime    string `json:"updateDateTime,omitempty"`    //修改时间
	Verifyid          string `json:"verifyid,omitempty"`          //审核人
	VerifyStatus      int32  `json:"verifyStatus,omitempty"`      //审核状态
	VerifyDateTime    string `json:"verifyDateTime,omitempty"`    //审核时间
	VerifyDescription string `json:"verifyDescription,omitempty"` //审核说明
	PhotoSalesInvoice string `json:"photoSalesInvoice,omitempty"` //销售发票附件
	PhotoQuality      string `json:"photoQuality,omitempty"`      //质检报告附件
	PhotoOutWh        string `json:"photoOutWh,omitempty"`        //出库清单附件
	TrpExArrivalDate  string `json:"trpExArrivalDate,omitempty"`  //预计到达时间
	IfMatching        int32  `json:"ifMatching,omitempty"`        //是否匹配专线,1匹配,0不匹配
	EcdStr            string `json:"ecd,omitempty"`               //加密字符串
	OrderAtch
	OrderECD
}

//-----------------------------------------------------------------------------
//构造函数
func NewOrder() *Order {
	return &Order{ResBase: res.ResBase{Idx: "order", Ver: "1.0"}}
}

//获取字段
func (self *Order) GetField(ftype res.FieldType) (fields []string, err error) {
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
func (self *Order) Encrypt(rAes *utils.AES) {
	nilECD := OrderECD{}
	pdata, err := json.Marshal(&self.OrderECD)
	if err != nil {
		panic(utils.Error(err))
	}
	self.EcdStr = rAes.Encrypt(pdata)
	self.OrderECD = nilECD
}

//资源解密
func (self *Order) Decrypt(rAes *utils.AES) {
	if self.EcdStr != "" {
		pdata := rAes.Decrypt(self.EcdStr)
		err := json.Unmarshal(pdata, &self.OrderECD)
		if err != nil {
			panic(utils.Error(err))
		}
	}
	self.EcdStr = ""
}
