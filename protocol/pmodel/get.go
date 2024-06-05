package pmodel

import (
	"dev.magustek.com/bigdata/dass/iotdriver/OP2_DL_698/dataExchange"
	"dev.magustek.com/bigdata/dass/iotdriver/OP2_DL_698/utils"
	"encoding/binary"
	"fmt"
	"strconv"
)

type GetRequest struct {
	getType   byte //
	piid      byte
	oadLength byte
	oads      []string
	timeTag   byte
	dateTime  string
}

func (g *GetRequest) GetType() utils.APDUType {
	return utils.Get_Request
}

func (g *GetRequest) Encode() ([]byte, error) {
	data := make([]byte, 0)
	return data, nil
}

func (g *GetRequest) Decode(data []byte) error {
	var err error
	g.oads = make([]string, 0)
	g.getType = data[1]
	switch data[1] {
	case 1:
		//fmt.Println("GetRequestNormal")
		if len(data) >= 8 {
			g.piid = data[2]
			g.oadLength = 1
			oad := ""
			for i := 0; i < 4; i++ {
				oad += fmt.Sprintf("%02x", data[3+i])
			}
			g.oads = append(g.oads, oad)
			g.timeTag = data[7]
		} else {
			err = fmt.Errorf("GetRequestNormal Data length error")
		}
	case 2:
		//fmt.Println("GetRequestNormalList")
		g.piid = data[2]
		g.oadLength = data[3]
		if len(data[4:]) == int(4*g.oadLength)+1 {
			for i := 0; i < int(g.oadLength); i++ {
				oad := ""
				for l := 0; l < 4; l++ {
					oad += fmt.Sprintf("%02x", data[4+i*4+l])
				}
				g.oads = append(g.oads, oad)
			}
			g.timeTag = data[4+4*g.oadLength]
		} else {
			err = fmt.Errorf("oad length err ")
		}
	case 3:
		//fmt.Println("GetRequestRecord")
		err = fmt.Errorf("unSupport Get request type %v", g.getType)
	case 4:
		//fmt.Println("GetRequestRecordList")
		err = fmt.Errorf("unSupport Get request type %v", g.getType)
	case 5:
		//fmt.Println("GetRequestNext")
		err = fmt.Errorf("unSupport Get request type %v", g.getType)
	case 6:
		//fmt.Println("GetRequestMD5")
		err = fmt.Errorf("unSupport Get request type %v", g.getType)
	default:
		err = fmt.Errorf("unSupport Get request type %v", g.getType)
	}
	return err
}

func (g *GetRequest) GenOutGoing(in utils.APDU) {

}

type GetResponse struct {
	getType   byte
	piid_acd  byte
	oadLength byte
	oads      []string
	errInfos  []byte
	values    [][]byte
}

func (g *GetResponse) GetType() utils.APDUType {
	return utils.Get_Response
}

func (g *GetResponse) Encode() ([]byte, error) {
	data := []byte{byte(utils.Get_Response), g.getType, g.piid_acd}
	if g.oadLength == 1 {

	} else {
		data = append(data, g.oadLength)
	}

	for i := 0; i < int(g.oadLength); i++ {
		oi, err := strconv.ParseUint(g.oads[i], 16, 64)
		if err != nil {
			return nil, err
		}
		ob := make([]byte, 4)
		binary.BigEndian.PutUint32(ob, uint32(oi))
		data = append(data, ob...)
		data = append(data, g.errInfos[i])
		data = append(data, g.values[i]...)
	}
	data = append(data, 0x00, 0x00)

	return data, nil
}

func (g *GetResponse) Decode(bytes []byte) error {
	return nil
}

func (g *GetResponse) GenOutGoing(in utils.APDU) {
	if in.GetType() == utils.Get_Request {
		g.oads = make([]string, 0)
		g.errInfos = make([]byte, 0)
		g.values = make([][]byte, 0)
		ind := in.(*GetRequest)
		g.piid_acd = ind.piid
		g.getType = ind.getType
		g.oadLength = ind.oadLength
		g.oads = append(g.oads, ind.oads...)
		for _, oad := range ind.oads {
			//根据oad获取数据， 获取到一个数据类型，但后Encode成[]byte

			dt, errCode := dataExchange.GetDL698DataTypeByOAD(oad)
			if errCode != 0 {
				g.errInfos = append(g.errInfos, 0x00)
				g.values = append(g.values, []byte{0x04}) //DAR  接口未定义
			} else {
				val, err := dt.Encode()
				if err != nil {
					g.errInfos = append(g.errInfos, 0x00)
					g.values = append(g.values, []byte{0x07}) //DAR  类型不匹配
				} else {
					g.errInfos = append(g.errInfos, 0x01) //数据
					g.values = append(g.values, val)
				}
			}
		}
	} else {
		panic("TODO")
	}
}
