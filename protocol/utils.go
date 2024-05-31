package protocol

import "errors"

// 2 BYte
type LengthArea struct {
	unit   byte   //长度单位 0：Byte  1： KByte ==> bit14
	length uint16 //bit13-bit0
}

// 1 Byte
type CTRLArea struct {
	dir      byte //方向       bit7
	prm      byte //启动标志位  bit6
	framing  byte //分帧标志    bit5
	blur     byte //混淆       bit3
	funcCode byte //功能码     bit2-bit0 1：链路管理   3：应用连接管理 数据交换服务
}

// n Byte
type AddressArea struct {
	sa AddressSA //服务器地址
	ca byte
}

type AddressSA struct {
	addressType  byte //bit7-bit6 0：单地址  1：通配地址  2：组地址  3：广播地址
	logicAddress byte //bit5-bit4 有无扩展逻辑地址
	addressLen   byte //bit3-bit0
	extendLogic  byte //
	address      string
}

type PacketData struct {
	length  LengthArea  //长度域
	control CTRLArea    //控制域
	address AddressArea //地址域
	apdu    APDU        //应用层数据单元
}

type APDUType byte

const (
	LINK_Request  APDUType = 0x01
	LINK_Response APDUType = 0x81
)

func GetAPDUType(t byte) (APDUType, error) {
	var at APDUType
	var err error
	switch t {
	case byte(LINK_Response):
		at = LINK_Response
	case byte(LINK_Request):
		at = LINK_Request
	default:
		err = errors.New("APDUType is not exist")
	}

	if err != nil {
		return at, err
	}
	return at, nil
}
