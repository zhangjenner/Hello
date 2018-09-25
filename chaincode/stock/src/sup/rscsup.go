package sup

import (
	"encoding/json"
	"res"
	"utils"
)

//=============================================================================
//解析资源
func ParseRsc(arg *string, field res.FieldType, rsc res.ResIf) {
	//参数验证
	rst, msg, err := res.HasField(rsc, field, *arg)
	if err != nil {
		panic(utils.Error(err))
	} else if rst == false {
		panic(utils.Errorf(msg))
	}
	//解析参数
	err = json.Unmarshal([]byte(*arg), rsc)
	if err != nil {
		panic(utils.Error(err))
	}
}

//合并资源数据
func CompRsc(oldRsc res.ResIf, arg *string) (newRsc res.ResIf, err error) {
	if *arg == "" {
		return oldRsc, nil
	}
	//合并数据
	newRsc = res.NewRes(oldRsc, true)
	err = json.Unmarshal([]byte(*arg), newRsc)
	if err != nil {
		return nil, utils.Error(err)
	}
	//数据验证
	rst, msg, err := res.CanMutable(newRsc, oldRsc)
	if err != nil {
		return nil, utils.Error(err)
	} else if rst == false {
		return nil, utils.Errorf(msg)
	}
	return newRsc, nil
}
