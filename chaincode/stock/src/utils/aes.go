package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"io"
	"strings"
)

//=============================================================================
//A加密
type AES struct {
	nonce int64  //随机数
	key   []byte //秘钥
}

//-----------------------------------------------------------------------------
//新建加密
func NewAES(nonce int64) *AES {
	return &AES{nonce: nonce}
}

//生成密钥对
func (self *AES) GenKey() *AES {
	self.key = GetRand(16, self.nonce)
	return self
}

//=============================================================================
//获取秘钥
func (self *AES) GetKey() []byte {
	return self.key
}

//导入秘钥
func (self *AES) LoadKey(key []byte) {
	self.key = key
}

//=============================================================================
//加密信息
func (self *AES) Encrypt(pdata []byte) (ctext string) {
	//新建加密器
	block, err := aes.NewCipher(self.key)
	if err != nil {
		panic(Error(err))
	}
	//加密信息
	pdata = self.padd(pdata)
	cdata := make([]byte, aes.BlockSize+len(pdata))
	iv := cdata[:aes.BlockSize]
	prng := strings.NewReader(string(GetRand(16, self.nonce)))
	if _, err := io.ReadFull(prng, iv); err != nil {
		panic(Error(err))
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cdata[aes.BlockSize:], pdata)
	return base64.StdEncoding.EncodeToString(cdata)
}

//解密信息
func (self *AES) Decrypt(ctext string) (pdata []byte) {
	//数据准备
	cdata, err := base64.StdEncoding.DecodeString(ctext)
	if err != nil {
		panic(Error(err))
	}
	block, err := aes.NewCipher(self.key)
	if err != nil {
		panic(Error(err))
	}
	if len(cdata) < aes.BlockSize {
		panic(Errorf("Invalid ciphertext. It must be a multiple of the block size"))
	}
	//解密数据
	iv := cdata[:aes.BlockSize]
	cdata = cdata[aes.BlockSize:]
	if len(cdata)%aes.BlockSize != 0 {
		panic(Errorf("Invalid ciphertext. It must be a multiple of the block size"))
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cdata, cdata)
	return self.unPadd(cdata)
}

//-----------------------------------------------------------------------------
//填补信息块
func (self *AES) padd(pdata []byte) []byte {
	padding := aes.BlockSize - len(pdata)%aes.BlockSize
	paddata := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(pdata, paddata...)
}

//取出信息块
func (self *AES) unPadd(cdata []byte) []byte {
	length := len(cdata)
	unpadding := int(cdata[length-1])
	if unpadding > aes.BlockSize || unpadding == 0 {
		panic(Errorf("Invalid pkcs7 padding (unpadding > aes.BlockSize || unpadding == 0)"))
	}
	pad := cdata[len(cdata)-unpadding:]
	for i := 0; i < unpadding; i++ {
		if pad[i] != byte(unpadding) {
			panic(Errorf("Invalid pkcs7 padding (pad[i] != unpadding)"))
		}
	}
	return cdata[:(length - unpadding)]
}
