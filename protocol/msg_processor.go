package protocol

import (
	"dev.magustek.com/bigdata/dass/iotdriver/OP2_DL_698/protocol/pmodel"
	"dev.magustek.com/bigdata/dass/iotdriver/OP2_DL_698/utils"
	"fmt"
	"sync"
)

// MsgProcessor 处理消息的Handler接口
type MsgProcessor interface {
	Process(pkt *utils.PacketData) (utils.APDU, error)
}

// 消息处理方法调用表, <msgId, action>
type processOptions map[utils.APDUType]*action

type action struct {
	genAPDU func() *utils.ProcessData      // 定义生成消息的类型。由于go不支持type作为参数，所以这里直接初始化结构体
	process func(*utils.ProcessData) error // 处理消息的逻辑。可以设置消息字段、根据消息做相应处理逻辑
}

// 表驱动，初始化消息处理方法组
func initProcessOption() processOptions {
	options := make(map[utils.APDUType]*action)

	options[utils.LINK_Request] = &action{
		genAPDU: func() *utils.ProcessData {
			return &utils.ProcessData{
				InComing:  &pmodel.LinkRequest{},
				OutComing: &pmodel.LinkResponse{},
			}
		},
		process: nil,
	}

	options[utils.Connect_Request] = &action{
		genAPDU: func() *utils.ProcessData {
			return &utils.ProcessData{
				InComing:  &pmodel.ConnectRequest{},
				OutComing: &pmodel.ConnectResponse{},
			}
		},
		process: nil,
	}

	options[utils.Get_Request] = &action{
		genAPDU: func() *utils.ProcessData {
			return &utils.ProcessData{
				InComing:  &pmodel.GetRequest{},
				OutComing: &pmodel.GetResponse{},
			}
		},
		process: nil,
	}

	options[utils.Set_Request] = &action{
		genAPDU: func() *utils.ProcessData {
			return &utils.ProcessData{
				InComing:  &pmodel.SetRequest{},
				OutComing: &pmodel.SetResponse{},
			}
		},
		process: nil,
	}

	return options
}

type DL698MsgProcessor struct {
	opts processOptions
}

var dl698MsgProcessor *DL698MsgProcessor

var ProcessorOnce sync.Once

func GetJT808MsgProcessor() *DL698MsgProcessor {
	ProcessorOnce.Do(func() {
		dl698MsgProcessor = &DL698MsgProcessor{
			initProcessOption(),
		}
	})
	return dl698MsgProcessor
}

func (mp *DL698MsgProcessor) Process(pkt *utils.PacketData) (*utils.ProcessData, error) {
	apduType, err := utils.GetAPDUType(pkt.Body[0])
	if err != nil {
		return nil, err
	}

	if _, ok := mp.opts[apduType]; !ok {
		return nil, fmt.Errorf("not support apdu type %x", apduType)
	}

	gct := mp.opts[apduType]

	pData := gct.genAPDU()

	in := pData.InComing
	err = in.Decode(pkt.Body)
	if err != nil {
		return nil, err
	}

	if pData.OutComing != nil {
		//根据输入信息生成输出信息
		pData.OutComing.GenOutGoing(in)
	}
	return pData, nil
}
