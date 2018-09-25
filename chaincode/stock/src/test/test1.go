package main

import (
	//	"crypto"
	//	"crypto/ecdsa"
	//	"crypto/elliptic"
	//	"crypto/md5"
	"math/rand"
	"utils"
	//	"encoding/base64"
	//	"encoding/json"
	//	"crypto/x509"
	//	"encoding/pem"
	"fmt"
	//	"github.com/hyperledger/fabric/bccsp"
	//	"github.com/hyperledger/fabric/bccsp/factory"
	//	"github.com/hyperledger/fabric/bccsp/mocks"
	//	"github.com/hyperledger/fabric/bccsp/sw"
	//	"io"
	//	"strings"
	"time"
)

var st time.Time

func ptime(msg string) {
	et := time.Now()
	fmt.Printf("%s:%v\n", msg, et.Sub(st))
	st = et
}

//=============================================================================
func main() {
	PemPub := `
-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEEVMI5amup8Nnk5gEIaL3Y5b9+d9z
e1ve9of7bgPJ7EOyu4ZHxxuPETS4xMi7/3O6cWoSwOfQ2kewRMace7XoTA==
-----END PUBLIC KEY-----`
	PemPri := `
-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQg163J2WPCpZ2U3aax
hkUf6XuypQ9TPJMz9GCsk/cHrsqhRANCAAQRUwjlqa6nw2eTmAQhovdjlv3533N7
W972h/tuA8nsQ7K7hkfHG48RNLjEyLv/c7pxahLA59DaR7BExpx7tehM
-----END PRIVATE KEY-----`

	fmt.Println("============statr===========")
	st = time.Now()
	ecc := utils.NewECC(12).GenKey()
	ptime("NewECC")
	ecc.LoadPemPubKey(PemPub)
	ptime("LoadPubKey")
	ecc.LoadPemPriKey(PemPri)
	ptime("LoadPriKey")
	//	ecc2 := utils.NewECC(10).GenKey()
	//	ptime("2NewECC")
	//	ecc2.LoadPemPubKey(PemPub)
	//	ptime("2LoadPubKey")
	//	ecc2.LoadPemPriKey(PemPri)
	//	ptime("2LoadPriKey")

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	msg := utils.GetRand(1024, r.Int63())
	ctext := ecc.Encrypt(msg)

	sst := time.Now()
	for i := 0; i < 1000; i++ {

		//		ptime("Encrypt")

		ecc.Decrypt(ctext)
		//		ptime("Decrypt")

		//		sign, err := ecc.Sign(msg)
		//		if err != nil {
		//			panic(err)
		//		}
		//		ptime("Sign")
		//		rst, err := ecc.Verify(msg, sign)
		//		if err != nil {
		//			panic(err)
		//		} else if rst != true {
		//			panic("Verify Wrong")
		//		}
		//		ptime("Verify")
	}
	eet := time.Now()
	fmt.Printf("Total time:%v\n", eet.Sub(sst))
	fmt.Println("=============end============")
}
