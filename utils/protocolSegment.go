package utils

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"time"
)

// 2 BYte
type LengthArea struct {
	unit   byte   //长度单位 0：Byte  1： KByte ==> bit14
	Length uint16 //bit13-bit0
}

// 1 Byte
type CTRLArea struct {
	Dir      byte //方向       bit7
	Prm      byte //启动标志位  bit6
	Framing  byte //分帧标志    bit5
	Blur     byte //混淆       bit3
	FuncCode byte //功能码     bit2-bit0 1：链路管理   3：应用连接管理 数据交换服务
}

// n Byte
type AddressArea struct {
	Sa AddressSA //服务器地址
	Ca byte
}

type AddressSA struct {
	AddressType  byte //bit7-bit6 0：单地址  1：通配地址  2：组地址  3：广播地址
	LogicAddress byte //bit5-bit4 有无扩展逻辑地址
	AddressLen   byte //bit3-bit0
	extendLogic  byte //
	Address      string
}

type PacketData struct {
	Length  LengthArea  //长度域
	Control CTRLArea    //控制域
	Address AddressArea //地址域
	Body    []byte
}

type ProcessData struct {
	InComing  APDU
	OutComing APDU
}

func (pkd *PacketData) EncodeControl() byte {
	var b byte
	b = pkd.Control.FuncCode & 0b00000111
	b = b | pkd.Control.Blur<<3
	b = b | pkd.Control.Framing<<5
	b = b | pkd.Control.Prm<<6
	b = b | pkd.Control.Dir<<7
	return b
}

func (pkd *PacketData) EncodeLength() uint16 {
	var l uint16

	return l
}

func (pkd *PacketData) EncodeAddress() []byte {
	val := make([]byte, 1)
	b1 := (pkd.Address.Sa.AddressLen - 1) & 0b00001111
	b1 = b1 | pkd.Address.Sa.LogicAddress<<4
	b1 = b1 | pkd.Address.Sa.AddressType<<6
	val[0] = b1

	var b byte
	//address 可能包含字母
	for i := 0; i < len(pkd.Address.Sa.Address); i++ {
		if i%2 == 0 {
			v, err := strconv.Atoi(pkd.Address.Sa.Address[i : i+1])
			if err != nil {
				break
			}
			b = byte(v) << 4
			continue
		}
		v, err := strconv.Atoi(pkd.Address.Sa.Address[i : i+1])
		if err != nil {
			b = b | 0xA
		} else {
			b = b | byte(v)
		}
		val = append(val, b)
	}
	j := len(val) - 1
	//fmt.Printf("Val % 02X\n", val)
	for i := 1; i < (len(val)+1)/2; i++ {
		//if i-1 == j {
		//	break
		//}
		x := val[i]
		y := val[j]
		val[i] = y
		val[j] = x
		j = j - 1
		//fmt.Printf("Val % 02X\n", val)
	}
	val = append(val, pkd.Address.Ca)
	return val
}

func DataTime2Str(dt []byte) string {
	fmt.Printf("Date:% 02X\n", dt)
	if len(dt) != 10 {
		return ""
	}
	var st string
	st += strconv.Itoa(int(binary.BigEndian.Uint16(dt[0:2]))) + "-" //07 E0 年
	st += fmt.Sprintf("%02d", dt[2]) + "-"                          //09 月
	st += fmt.Sprintf("%02d", dt[3]) + " "                          //0C

	//01 weak  0：周末
	st += fmt.Sprintf("%02d", dt[5]) + ":"                     //00
	st += fmt.Sprintf("%02d", dt[6]) + ":"                     //00
	st += fmt.Sprintf("%02d", dt[7]) + "."                     //1F
	st += strconv.Itoa(int(binary.BigEndian.Uint16(dt[8:10]))) //02 AB //毫秒
	return st
}

func Str2DataTime(dt string) []byte {
	t := make([]byte, 10)
	if dt == "" {
		t[0] = 0xFF
		t[1] = 0xFF
		return t
	}
	if len(dt) == 23 {
		rt, err := time.Parse("2006-01-02 15:04:05.000", dt)
		if err != nil {
			t[0] = 0xFF
			t[1] = 0xFF
			return t
		}
		binary.BigEndian.PutUint16(t[0:], uint16(rt.Year()))
		t[2] = byte(rt.Month())
		t[3] = byte(rt.Day())
		t[4] = byte(rt.Weekday())
		t[5] = byte(rt.Hour())
		t[6] = byte(rt.Minute())
		t[7] = byte(rt.Second())
		binary.BigEndian.PutUint16(t[8:], uint16(rt.Nanosecond()/1e6))
	} else {
		t[0] = 0xFF
		t[1] = 0xFF
		return t
	}

	return t
}
