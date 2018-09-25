package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"encoding/base64"
	butils "github.com/hyperledger/fabric/bccsp/utils"
	"math/big"
	"strings"
)

//=============================================================================
//ECC加密
type ECC struct {
	Nonce          int64                 //随机数
	HasPub, HasPri bool                  //是否有公钥/私钥
	PubKey         *ecdsa.PublicKey      //公钥
	PriKey         *ecdsa.PrivateKey     //私钥
	Curve          elliptic.Curve        //椭圆曲线
	cParams        *elliptic.CurveParams //椭圆曲线参数
	cLen           int                   //坐标参数长度
}

//-----------------------------------------------------------------------------
//新建加密
func NewECC(nonce int64) *ECC {
	impl := &ECC{Nonce: nonce, HasPub: false, HasPri: false}
	impl.Curve = elliptic.P256()
	impl.cParams = impl.Curve.Params()
	impl.cLen = impl.cParams.BitSize / 8
	return impl
}

//生成密钥对
func (self *ECC) GenKey() *ECC {
	var err error
	randbyte := GetRand(256/8+8, self.Nonce)
	rander := strings.NewReader(string(randbyte))
	self.PriKey, err = ecdsa.GenerateKey(self.Curve, rander)
	self.PubKey = &self.PriKey.PublicKey
	self.HasPub, self.HasPri = true, true
	if err != nil {
		panic(Error(err))
	}
	return self
}

//=============================================================================
//获取PEM格式的公钥
func (self *ECC) GetPemPubKey() (pemPubKey string) {
	if !self.HasPub {
		panic(Errorf("No PubKey"))
	}
	PubKey, err := butils.PublicKeyToPEM(self.PubKey, nil)
	if err != nil {
		panic(Error(err))
	}
	return string(PubKey)
}

//导入PEM格式公钥
func (self *ECC) LoadPemPubKey(pemPubKey string) {
	PubKey, err := butils.PEMtoPublicKey([]byte(pemPubKey), nil)
	if err != nil {
		panic(Error(err))
	}
	ecPubKey, ok := PubKey.(*ecdsa.PublicKey)
	if !ok {
		panic(Errorf("Invalid key type. It must be *ecdsa.PublicKey"))
	}
	if self.PriKey == nil {
		self.PriKey = &ecdsa.PrivateKey{PublicKey: *ecPubKey}
		self.PubKey = &self.PriKey.PublicKey
	} else {
		self.PriKey.PublicKey = *ecPubKey
	}
	self.HasPub = true
}

//=============================================================================
//导入PEM格式私钥
func (self *ECC) LoadPemPriKey(pemPriKey string) {
	PriKey, err := butils.PEMtoPrivateKey([]byte(pemPriKey), nil)
	if err != nil {
		panic(Error(err))
	}
	ecPriKey, ok := PriKey.(*ecdsa.PrivateKey)
	if !ok {
		panic(Errorf("Invalid key type. It must be *ecdsa.PublicKey"))
	}
	PubKeyBK := self.PriKey.PublicKey
	self.PriKey = ecPriKey
	self.PriKey.PublicKey = PubKeyBK
	self.PubKey = &self.PriKey.PublicKey
	self.HasPri = true
}

//导入数字私钥
func (self *ECC) LoadBytePriKey(bytes []byte) {
	d := new(big.Int).SetBytes(bytes)
	if self.PriKey == nil {
		self.PriKey = &ecdsa.PrivateKey{D: d}
		self.PriKey.Curve = self.Curve
		self.PubKey = &self.PriKey.PublicKey
	} else {
		self.PriKey.D = d
	}
	self.HasPri = true
}

//=============================================================================
//加密信息
func (self *ECC) Encrypt(pdata []byte) (ctext string) {
	var cdata []byte
	blockNum := len(pdata) / self.cLen
	r := GetRand(self.cParams.BitSize/8+8, self.Nonce)
	if len(pdata)%self.cLen == 0 {
		cdata = make([]byte, 4*blockNum*self.cLen)
		for i := 0; i < blockNum; i++ {
			self.enBlock(cdata[4*i*self.cLen:], pdata[i*self.cLen:(i+1)*self.cLen], r)
		}
	} else {
		cdata = make([]byte, 4*(blockNum+1)*self.cLen)
		for i := 0; i < blockNum; i++ {
			self.enBlock(cdata[4*i*self.cLen:], pdata[i*self.cLen:(i+1)*self.cLen], r)
		}
		self.enBlock(cdata[4*blockNum*self.cLen:], pdata[blockNum*self.cLen:], r)
	}
	return base64.StdEncoding.EncodeToString(cdata)
}

//解密信息
func (self *ECC) Decrypt(ctext string) (pdata []byte) {
	cdata, err := base64.StdEncoding.DecodeString(ctext)
	if err != nil {
		panic(Error(err))
	}
	blockLen := self.cLen * 4
	blockNum := len(cdata) / blockLen
	pdata = make([]byte, blockNum*self.cLen)
	deLen := self.cLen
	for i := 0; i < blockNum; i++ {
		deLen = self.deBlock(pdata[i*self.cLen:], cdata[i*blockLen:(i+1)*blockLen])
	}
	if deLen != self.cLen {
		return pdata[:len(pdata)-(self.cLen-deLen)]
	}
	return pdata
}

//数字签名
func (self *ECC) Sign(text []byte) (signature string) {
	digest := md5.Sum([]byte(text))
	randbyte := GetRand(256/8+8, self.Nonce)
	rander := strings.NewReader(string(randbyte))
	r, s, err := ecdsa.Sign(rander, self.PriKey, digest[:])
	if err != nil {
		panic(Error(err))
	}
	s, _, err = butils.ToLowS(self.PubKey, s)
	if err != nil {
		panic(Error(err))
	}
	signdata, err := butils.MarshalECDSASignature(r, s)
	if err != nil {
		panic(Error(err))
	}
	return base64.StdEncoding.EncodeToString(signdata)
}

//签名验证
func (self *ECC) Verify(text []byte, signature string) (valid bool) {
	digest := md5.Sum([]byte(text))
	signdata, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		panic(Error(err))
	}
	r, s, err := butils.UnmarshalECDSASignature(signdata)
	if err != nil {
		panic(Error(err))
	}
	lowS, err := butils.IsLowS(self.PubKey, s)
	if err != nil {
		panic(Error(err))
	}
	if !lowS {
		panic(Errorf("Invalid S. Must be smaller than half the order"))
	}
	return ecdsa.Verify(self.PubKey, digest[:], r, s)
}

//-----------------------------------------------------------------------------
//加密信息块
func (self *ECC) enBlock(cBlock, pBlock, r []byte) {
	//信息映射
	Mx := big.NewInt(0)
	Mx.SetBytes(pBlock)
	My1 := big.NewInt(0)
	My1.Mul(Mx, Mx)
	My1.Mul(My1, Mx)
	My2 := big.NewInt(0)
	My2.Mul(Mx, self.cParams.N)
	My := big.NewInt(0)
	My.Add(My1, My2)
	My.Add(My, self.cParams.B)
	My.Mod(My, self.cParams.P)
	My.Sqrt(My)
	//加密信息
	Kx, Ky := self.PubKey.X, self.PubKey.Y
	rKx, rKy := self.Curve.ScalarMult(Kx, Ky, r)
	C1x, C1y := self.Curve.Add(Mx, My, rKx, rKy)
	C2x, C2y := self.Curve.ScalarBaseMult(r)
	copy(cBlock[0:self.cLen], self.getParam(C1x))
	copy(cBlock[1*self.cLen:2*self.cLen], self.getParam(C1y))
	copy(cBlock[2*self.cLen:3*self.cLen], self.getParam(C2x))
	copy(cBlock[3*self.cLen:], self.getParam(C2y))
}

//解密信息块
func (self *ECC) deBlock(pBlock, cBlock []byte) int {
	//获取加密数据
	C1x := big.NewInt(0)
	C1y := big.NewInt(0)
	C2x := big.NewInt(0)
	C2y := big.NewInt(0)
	size := len(cBlock) / 4
	C1x.SetBytes(cBlock[0*size : 1*size])
	C1y.SetBytes(cBlock[1*size : 2*size])
	C2x.SetBytes(cBlock[2*size : 3*size])
	C2y.SetBytes(cBlock[3*size : 4*size])
	//解密信息
	k := self.PriKey.D
	kC2x, kC2y := self.Curve.ScalarMult(C2x, C2y, k.Bytes())
	kC2y.Neg(kC2y)
	dMx, _ := self.Curve.Add(C1x, C1y, kC2x, kC2y)
	return copy(pBlock, dMx.Bytes())
}

//添加加密后的参数
func (self *ECC) getParam(param *big.Int) []byte {
	bytes := param.Bytes()
	if len(bytes) < self.cLen {
		nbytes := make([]byte, self.cLen-len(bytes))
		return append(nbytes, bytes...)
	}
	return bytes
}
