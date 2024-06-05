package protocol

import (
	"fmt"
	"net"
	"testing"
)

func TestClientLogin(t *testing.T) {

	conn, err := net.Dial("tcp", "127.0.0.1:7001")
	//conn, err := net.Dial("tcp", "192.168.60.35:6001")
	if err != nil {
		fmt.Println("Dial err:", err)
		return
	}

	pipe := NewPipeLine(conn, "12345678")

	pipe.Login()
	//want:68 1c 00 81 03 78 56 34 12 00 58 89 01 00 00 01 2c 07 e0 09 0b 07 0a 15 28 02 f5 0b 65 16
	//real:68 1C 00 81 07 78 56 34 12 10 75 89 01 00 00 01 2C 07 E8 06 03 01 0F 15 32 01 99 DC 3E 16

	pipe.Start()

	//Connect
	//pipe.Connect()
	//want:
	//02 00 00 10 FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF
	//04 00 04 00 01 04 00 00 00 00 64 00 00
	//FF FF 16

	//real:
	//68 36 00 81 03 78 56 34 12 10 D0 8A
	//02 00 00 16 FF FF FF FF C0 00 00 00 00 01 FF FE 00 00 00 00 00 00 00 00 00 00 00 00
	//04 00 04 00 01 1B 58 00 00 1C 20
	//00 00
	//B7 47 16

	//heart
}

func TestLinkLogin(t *testing.T) {
	pc := GetJT808PacketCodec()

	//Login
	data := []byte{0x68, 0x1E, 0x00, 0x81, 0x05, 0x07, 0x09, 0x19, 0x05, 0x16, 0x20, 0x00, 0xAA, 0xAA, 0x01, 0x00, 0x00, 0x00, 0xB4, 0x07, 0xE0, 0x05, 0x13, 0x04, 0x08, 0x05, 0x00, 0x00, 0xA4, 0xAA, 0xAA, 0x16}
	pkd, err := pc.Decode(data)
	if err != nil {
		fmt.Println("Decode err:", err)
		return
	}

	mp := GetJT808MsgProcessor()
	pcd, err := mp.Process(pkd)
	if err != nil {
		fmt.Println("Process err:", err)
		return
	}

	res, err := pc.Encode(pkd, pcd)
	if err != nil {
		fmt.Println("packetCoder Encode err:", err)
		return
	}
	fmt.Printf("Data: % 02X\n", res)

	//send data
}

func TestLinkHeartBeat(t *testing.T) {
	pc := GetJT808PacketCodec()

	//HeartBeat

	//68 1e 00 81 05 01 00 00 00 00 00 00 d2 b6 01 10 01 01 2c 07 e0 09 0c 01 00 00 1f 02 ab 9d 39 16
	data := []byte{0x68, 0x1E, 0x00, 0x81, 0x05, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xD2, 0xB6, 0x01, 0x10, 0x01, 0x01,
		0x2C, 0x07, 0xE0, 0x09, 0x0C, 0x01, 0x00, 0x00, 0x1F, 0x02, 0xAB, 0x9D, 0x39, 0x16}
	pkd, err := pc.Decode(data)
	if err != nil {
		fmt.Println("Decode err:", err)
		return
	}

	mp := GetJT808MsgProcessor()
	pcd, err := mp.Process(pkd)
	if err != nil {
		fmt.Println("Process err:", err)
		return
	}

	res, err := pc.Encode(pkd, pcd)
	if err != nil {
		fmt.Println("packetCoder Encode err:", err)
		return
	}
	fmt.Printf("Data: % 02X\n", res)

	//send data
}

func TestConnect(t *testing.T) {
	pc := GetJT808PacketCodec()

	//Connect

	/*
		68 2C 00 81 05 07 09 19 05 16 20 00 FF FF
		02 00 00 10 FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF FF
		FF 04 00 04 00 01 04 00 00 00 00 64 00 00
		FF FF 16

	*/
	data := []byte{0x68, 0x2C, 0x00, 0x81, 0x05, 0x07, 0x09, 0x19, 0x05, 0x16, 0x20, 0x00, 0xFF, 0xFF,
		0x02, 0x00, 0x00, 0x10, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0x04, 0x00, 0x04, 0x00, 0x01, 0x04, 0x00, 0x00, 0x00, 0x00, 0x64, 0x00, 0x00, 0xFF, 0xFF, 0x16}

	fmt.Println("len:", len(data))
	pkd, err := pc.Decode(data)
	if err != nil {
		fmt.Println("Decode err:", err)
		return
	}

	mp := GetJT808MsgProcessor()
	pcd, err := mp.Process(pkd)
	if err != nil {
		fmt.Println("Process err:", err)
		return
	}

	res, err := pc.Encode(pkd, pcd)
	if err != nil {
		fmt.Println("packetCoder Encode err:", err)
		return
	}
	fmt.Printf("Data: % 02X\n", res)

	//send data
}

func TestGet(t *testing.T) {
	pc := GetJT808PacketCodec()

	//Connect

	/*
		68 2C 00 81 05 07 09 19 05 16 20 00 FF FF
		05 01 01 40 01 02 00 00
		FF FF 16
	*/
	data := []byte{
		0x68, 0x17, 0x00, 0x81, 0x05, 0x07, 0x09, 0x19, 0x05, 0x16, 0x20, 0x00, 0xFF, 0xFF,
		0x05, 0x01, 0x01, 0x40, 0x01, 0x02, 0x00, 0x00,
		0xFF, 0xFF, 0x16,
	}

	fmt.Println("len:", len(data))
	pkd, err := pc.Decode(data)
	if err != nil {
		fmt.Println("Decode err:", err)
		return
	}

	mp := GetJT808MsgProcessor()
	pcd, err := mp.Process(pkd)
	if err != nil {
		fmt.Println("Process err:", err)
		return
	}

	res, err := pc.Encode(pkd, pcd)
	if err != nil {
		fmt.Println("packetCoder Encode err:", err)
		return
	}
	fmt.Printf("Data: % 02X\n", res)

	//send data

}

func TestGetTime(t *testing.T) {
	pc := GetJT808PacketCodec()

	//Connect

	/*
		68 17 00 43 05 01 00 00 00 00 00 10 26 f6 05 01 05 40 00 02 00 00 d1 0b 16
	*/
	data := []byte{
		0x68, 0x17, 0x00, 0x43, 0x05, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x26, 0xf6, 0x05, 0x01, 0x05, 0x40, 0x00, 0x02, 0x00, 0x00, 0xd1, 0x0b, 0x16,
	}

	fmt.Println("len:", len(data))
	pkd, err := pc.Decode(data)
	if err != nil {
		fmt.Println("Decode err:", err)
		return
	}

	mp := GetJT808MsgProcessor()
	pcd, err := mp.Process(pkd)
	if err != nil {
		fmt.Println("Process err:", err)
		return
	}

	res, err := pc.Encode(pkd, pcd)
	if err != nil {
		fmt.Println("packetCoder Encode err:", err)
		return
	}
	fmt.Printf("Data: % 02X\n", res)

	//send data
	//real:68 21 00 C3 05 01 00 00 00 00 00 10 FF FF 85 01 05 40 00 02 00 01 1C 07 E8 06 03 0B 1C 00 00 00 FF FF 16
	//want:68 21 00 c3 05 01 00 00 00 00 00 10 49 54 85 01 05 40 00 02 00 01 1c 07 e0 09 0b 12 1e 15 00 00 2a fb 16
}

func TestSet(t *testing.T) {
	pc := GetJT808PacketCodec()

	//Connect

	/*
		68 2C 00 81 05 07 09 19 05 16 20 00 FF FF
		06 01 02 40 00 02 00 1C 07 E0 01 14 10 1B 0B 00
		FF FF 16
	*/
	data := []byte{
		0x68, 0x1F, 0x00, 0x81, 0x05, 0x07, 0x09, 0x19, 0x05, 0x16, 0x20, 0x00, 0xFF, 0xFF,
		0x06, 0x01, 0x02, 0x40, 0x00, 0x02, 0x00, 0x1C, 0x07, 0xE0, 0x01, 0x14, 0x10, 0x1B, 0x0B, 0x00,
		0xFF, 0xFF, 0x16,
	}

	fmt.Println("len:", len(data))
	pkd, err := pc.Decode(data)
	if err != nil {
		fmt.Println("Decode err:", err)
		return
	}

	mp := GetJT808MsgProcessor()
	pcd, err := mp.Process(pkd)
	if err != nil {
		fmt.Println("Process err:", err)
		return
	}

	res, err := pc.Encode(pkd, pcd)
	if err != nil {
		fmt.Println("packetCoder Encode err:", err)
		return
	}
	fmt.Printf("Data: % 02X\n", res)

	//send data
}
