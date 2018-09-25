package main

import (
	"res"
	"sup"
	"utils"
)

//=============================================================================
//获取用户信息
func (cc *StockCC) GetCreator() (user *res.User, err error) {
	stub := cc.GetStub()
	_, cert := utils.GetCreator(stub)
	user = res.NewUser()
	user.Cert = cert
	users, ok := sup.SelectData(stub, user).([]res.User)
	if !ok {
		return nil, utils.Errorf("The result of SelectRes is wrong")
	} else if len(users) != 1 {
		return nil, utils.Errorf("The number of user whose cert is '%s' is wrong ", user.Cert)
	}
	return &users[0], nil
}
