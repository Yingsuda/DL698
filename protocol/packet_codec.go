package protocol

import (
	"dl698/utils"
	"encoding/binary"
	"fmt"
	"github.com/sigurn/crc16"
	"sync"
)

type PacketCodec interface {
	Decode([]byte) (*utils.PacketData, error)
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

// 还原原始数据包
func (pc *DL698PacketCodec) unescape(src []byte) []byte {
	dst := make([]byte, 0)
	return dst
}

// 反转义 -> 校验 -> 反序列化
func (pc *DL698PacketCodec) Decode(data []byte) (*utils.PacketData, error) {
	if len(data) < 10 {
		return nil, fmt.Errorf("minimum length err")
	}
	if data[0] != 0x68 || data[len(data)-1] != 0x16 {
		return nil, fmt.Errorf("head or tail err")
	}
	//长度域
	length := append([]byte{}, data[1:3]...)
	length[1] = 0b00111111 & length[1]
	la := utils.LengthArea{
		Length: binary.LittleEndian.Uint16(length),
	}
	//fmt.Println("la.length:", la.Length)
	//长度校验
	if len(data) != int(la.Length+2) {
		return nil, fmt.Errorf("total length err")
	}

	//控制域
	ca := utils.CTRLArea{
		Dir:      data[3] >> 7,
		Prm:      (data[3] << 1) >> 7,
		Framing:  (data[3] << 2) >> 7,
		Blur:     (data[3] << 4) >> 7,
		FuncCode: data[3] & 0b00000111,
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

	aa := utils.AddressArea{
		Sa: utils.AddressSA{
			AddressType:  addrType,
			LogicAddress: logicAddr,
			AddressLen:   addressLen,
			//extendLogic:  data[8],
			Address: address,
		},
		Ca: data[5+addressLen],
	}

	//fmt.Println("logic:", aa.Sa.LogicAddress)
	//fmt.Println("addressLen:", aa.Sa.AddressLen)
	//fmt.Println("address:", aa.Sa.Address)
	//fmt.Println("CA:", aa.Ca)
	//HCS

	//APDU
	apduData := data[5+addressLen+3 : len(data)-3]
	//fmt.Printf("APDU Data:% 02X\n", apduData)
	//校验尾CS
	if ca.Blur == 1 {
		//扰码还原
		apduData = pc.unescape(apduData)
	}

	p := new(utils.PacketData)
	p.Length = la
	p.Control = ca
	p.Address = aa
	p.Body = apduData
	return p, nil
}

func (pc *DL698PacketCodec) Encode(pkd *utils.PacketData, pcd *utils.ProcessData) ([]byte, error) {
	data := []byte{0x68,
		0x00, 0x00,
		0x00,
	}
	//funcCode
	//pkd.Control.Dir = 0
	//fmt.Println("control:", pkd.Control)
	data[3] = pkd.EncodeControl()

	pkd.Address.Ca = 0x10

	data = append(data, pkd.EncodeAddress()...)
	hcs := len(data)
	data = append(data, 0xFF, 0xFF)

	od, err := pcd.OutComing.Encode()
	if err != nil {
		return nil, err
	}
	//
	data = append(data, od...)
	data = append(data, []byte{0xFF, 0xFF}...)
	data = append(data, 0x16)
	binary.LittleEndian.PutUint16(data[1:], uint16(len(data)-2))                                                              //长度
	binary.LittleEndian.PutUint16(data[hcs:], crc16.Checksum(data[1:hcs], crc16.MakeTable(crc16.CRC16_X_25)))                 //HCS
	binary.LittleEndian.PutUint16(data[len(data)-3:], crc16.Checksum(data[1:len(data)-3], crc16.MakeTable(crc16.CRC16_X_25))) //FCS
	return data, nil
}
