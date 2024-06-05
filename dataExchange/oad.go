package dataExchange

import (
	"dev.magustek.com/bigdata/dass/iotdriver/OP2_DL_698/utils"
	"gitee.com/iotdrive/tools/logs"
	"time"
)

const (
	OAD_40000200 = "40000200" //获取时钟
	OAD_40010200 = "40010200"
)

type OadHandle func(oad string) (utils.DL698Data, byte)

var OAD_Map = make(map[string]OadHandle)

func init() {
	//默认支持的oad func
	OAD_Map[OAD_40000200] = func(oad string) (utils.DL698Data, byte) {
		return &utils.DTDateTimeS{
			Value: time.Now(),
		}, 0
	}
}

func RegisterOAD(oad string, f func(oad string) (utils.DL698Data, byte)) error {
	OAD_Map[oad] = f
	return nil
}

func GetDL698DataTypeByOAD(oad string) (utils.DL698Data, byte) {
	//通过OAD获取到对应的消息，生成对应的数据结构
	if f, ok := OAD_Map[oad]; ok {
		return f(oad)
	} else {
		logs.Error("oad ", oad, "is not exist")
	}
	return nil, 0xFF
}
