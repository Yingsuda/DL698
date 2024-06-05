package pmodel

import (
	"dev.magustek.com/bigdata/dass/iotdriver/OP2_DL_698/utils"
	"encoding/binary"
)

type ConnectMechanismInfo struct {
	cmType   byte //0:公共连接  1:一般密码  2:对称加密  3:数字加密
	passWord string
	key      string
	value    string
}

// 建立应用连接
type ConnectRequest struct {
	Piid                byte
	ProtocolVersion     uint16
	ProtocolConformance [8]byte
	FunctionConformance [16]byte
	MaxSend             uint16
	MaxReceive          uint16
	MaxSize             byte
	MaxAPDULength       uint16
	WantTimeout         uint32
	//认证类型
	cmInfo   ConnectMechanismInfo
	timeTag  byte //0 没有时间Tag
	dateTime string
}

func (c *ConnectRequest) GetType() utils.APDUType {
	return utils.Connect_Request
}

func (c *ConnectRequest) Encode() ([]byte, error) {
	data := []byte{byte(utils.Connect_Request), 0x00, 0x00, 0x00}
	binary.BigEndian.PutUint16(data[2:], c.ProtocolVersion)
	data = append(data, []byte{0xff, 0xff, 0xff, 0xff, 0xc0, 0x00, 0x00, 0x00}...)
	data = append(data, []byte{0x00, 0x01, 0xff, 0xfe, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}...)
	u16 := make([]byte, 2)
	binary.BigEndian.PutUint16(u16, c.MaxSend)
	data = append(data, u16...)
	binary.BigEndian.PutUint16(u16, c.MaxReceive)
	data = append(data, u16...)
	data = append(data, c.MaxSize)
	binary.BigEndian.PutUint16(u16, c.MaxAPDULength)
	data = append(data, u16...)
	u32 := make([]byte, 4)
	binary.BigEndian.PutUint32(u32, c.WantTimeout)
	data = append(data, u32...)
	data = append(data, []byte{0x00, 0x00}...)
	return data, nil
}

func (c *ConnectRequest) Decode(data []byte) error {
	//fmt.Println("ConnectRequest")
	c.Piid = data[1]
	//fmt.Println("PIID:", c.Piid)

	return nil
}

func (c *ConnectRequest) GenOutGoing(in utils.APDU) {
	//TODO implement me
	panic("implement me")
}

type ConnectResponseInfo struct {
	connectResult byte   //0:允许  1：密码错误  2：对称解密错误  3：非对称解密错误  4：签名错误  5：版本不匹配  255：其他
	securityData  []byte //todo
}

type ConnectResponse struct {
	ppid_acd            byte
	factoryVersion      [32]byte
	protocolVersion     uint16
	protocolConformance [8]byte
	functionConformance [16]byte
	maxSend             uint16
	maxReceive          uint16
	maxSize             byte
	maxAPDULength       uint16
	wantTimeout         uint32
	crInfo              ConnectResponseInfo
	timeTag             byte //0 没有时间Tag
	dateTime            string
}

func (c *ConnectResponse) GetType() utils.APDUType {
	return utils.Connect_Response
}

func (c *ConnectResponse) Encode() ([]byte, error) {
	data := make([]byte, 0)
	return data, nil
}

func (c *ConnectResponse) Decode(bytes []byte) error {
	//TODO implement me
	panic("implement me")
}

func (c *ConnectResponse) GenOutGoing(in utils.APDU) {
	if in.GetType() == utils.Connect_Request {
		//fmt.Println("in piid:", in.(*ConnectRequest).Piid)
	} else {
		panic("不该来这里的")
	}
}
