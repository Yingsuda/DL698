package protocol

import (
	"encoding/binary"
	"fmt"
	"sync"
)

type PacketCodec interface {
	Decode([]byte) (*PacketData, error)
}

type DL698PacketCodec struct{}

var dl698PacketCodec *DL698PacketCodec
var codecOnce sync.Once

func GetJT808PacketCodec() *DL698PacketCodec {
	codecOnce.Do(func() {
		dl698PacketCodec = &DL698PacketCodec{}
	})
	return dl698PacketCodec
}

// 反转义 -> 校验 -> 反序列化
func (pc *DL698PacketCodec) Decode(data []byte) (*PacketData, error) {
	if len(data) < 10 {
		return nil, fmt.Errorf("minimum length err")
	}
	if data[0] != 0x68 || data[len(data)-1] != 0x16 {
		return nil, fmt.Errorf("head or tail err")
	}
	//长度域
	length := append([]byte{}, data[1:3]...)
	length[1] = 0b00111111 & length[1]
	la := LengthArea{
		length: binary.LittleEndian.Uint16(length),
	}
	fmt.Println("la.length:", la.length)
	//长度校验
	if len(data) != int(la.length+2) {
		return nil, fmt.Errorf("total length err")
	}

	//控制域
	ca := CTRLArea{
		dir:      data[3] >> 7,
		prm:      (data[3] << 1) >> 7,
		framing:  (data[3] << 2) >> 7,
		blur:     (data[3] << 4) >> 7,
		funcCode: data[3] & 0b00000111,
	}

	//地址域
	addrType := data[4] >> 6
	logicAddr := (data[4] << 2) >> 6
	addressLen := data[4]&0b00001111 + 1

	address := "" // fmt.Sprintf("%02x", data[4])
	for i := int(addressLen) - 1; i >= 0; i-- {
		address += fmt.Sprintf("%02x", data[5+i])
	}

	switch addrType {
	case 0: //单地址
		fmt.Println("单地址")
	case 1:
		fmt.Println("组地址")
	case 2:
		fmt.Println("通配地址")
	case 3:
		fmt.Println("广播地址")
	default:
		return nil, fmt.Errorf("address type err")
	}

	aa := AddressArea{
		sa: AddressSA{
			addressType:  addrType,
			logicAddress: logicAddr,
			addressLen:   addressLen,
			//extendLogic:  data[8],
			address: address,
		},
		ca: data[5+addressLen],
	}

	fmt.Println("logic:", aa.sa.logicAddress)
	fmt.Println("addressLen:", aa.sa.addressLen)
	fmt.Println("address:", aa.sa.address)
	fmt.Println("CA:", aa.ca)
	//HCS

	//APDU
	apduData := data[5+addressLen+3 : len(data)-3]
	fmt.Printf("APDU Data:% 02X\n", apduData)
	//校验尾CS
	if ca.blur == 1 {
		//扰码还原
		apduData = pc.unescape(apduData)
	}

	//encode APDU
	apduType, err := GetAPDUType(apduData[0])
	if err != nil {
		return nil, err
	}

	var apdu APDU
	switch apduType {
	case LINK_Request:
		apdu = &Link_Request{}
	case LINK_Response:
		apdu = &Link_Response{}

	}

	p := new(PacketData)
	p.length = la
	p.control = ca
	p.address = aa
	p.apdu = apdu
	return p, nil
}

// 还原原始数据包
func (pc *DL698PacketCodec) unescape(src []byte) []byte {
	dst := make([]byte, 0)
	return dst
}
