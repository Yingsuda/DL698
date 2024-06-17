package pmodel

import (
	"dev.magustek.com/bigdata/dass/iotdriver/OP2_DL_698/utils"
	"encoding/binary"
	"fmt"
	"time"
)

// 预连接

type LinkRequest struct {
	Piid_acd    byte
	RequestType byte   //0:登录  1：心跳  2：退出登录
	Heartbeat   uint16 //心跳周期
	requestDate string
	receiveDate string
}

func (l *LinkRequest) GenOutGoing(in utils.APDU) {
}

func (l *LinkRequest) GetType() utils.APDUType {
	return utils.LINK_Request
}

func (l *LinkRequest) Encode() ([]byte, error) {
	//	fmt.Println("Link Request Data")
	data := []byte{byte(utils.LINK_Request), l.Piid_acd, l.RequestType, 0x00, 0x00}
	binary.BigEndian.PutUint16(data[3:], l.Heartbeat)
	data = append(data, utils.Str2DataTime(time.Now().Format("2006-01-02 15:04:05.000"))...)
	return data, nil
}

func (l *LinkRequest) Decode(data []byte) error {
	l.Piid_acd = data[1]
	l.RequestType = data[2]
	l.Heartbeat = binary.BigEndian.Uint16(data[3:])
	//fmt.Println("heartBeat:", l.heartbeat)
	st := utils.DataTime2Str(data[5:15])
	l.requestDate = st //[]byte to data
	//fmt.Println("requestDate:", l.requestDate)
	l.receiveDate = time.Now().Format("2006-01-02 15:04:05.000") //时间
	return nil
}

type LinkResponse struct {
	piid         byte
	result       byte   //结果 bit7:始终可信标识 1：可信  0：不可信   bit2-bit0：结果  0：成功  1：地址重复  2：非法设备  3：容量不足
	requestDate  string //请求时间
	receiveDate  string //收到的时间
	responseDate string //回复时间
}

func (l *LinkResponse) GenOutGoing(in utils.APDU) {
	if in.GetType() == utils.LINK_Request {
		l.piid = in.(*LinkRequest).Piid_acd //根据这个ID 获取对应设备的Str，并进行相关逻辑处理
		//逻辑处理
		switch in.(*LinkRequest).RequestType {
		case 0:
			//fmt.Println("Login")
		case 1:
			//fmt.Println("HeartBeat")
		case 2:
			//fmt.Println("Exit Login")
		default:
			fmt.Println("UnSupport TypeCode")
		}
		l.result = 0x80
		l.requestDate = in.(*LinkRequest).requestDate
		l.receiveDate = in.(*LinkRequest).receiveDate
		l.responseDate = time.Now().Format("2006-01-02 15:04:05.000")
	} else {
		panic(" 不该来这儿的！")
	}
}

func (l *LinkResponse) GetType() utils.APDUType {
	return utils.LINK_Response
}

func (l *LinkResponse) Encode() ([]byte, error) {
	data := []byte{byte(utils.LINK_Response), l.piid, l.result}
	//time string to []byte
	//l.requestDate
	data = append(data, utils.Str2DataTime(l.requestDate)...)
	//l.receiveDate
	data = append(data, utils.Str2DataTime(l.receiveDate)...)
	//l.responseDate
	data = append(data, utils.Str2DataTime(l.responseDate)...)
	return data, nil
}

func (l *LinkResponse) Decode([]byte) error {
	return nil
}
