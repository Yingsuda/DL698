package pmodel

import (
	"dev.magustek.com/bigdata/dass/iotdriver/OP2_DL_698/DLContorl"
	"dev.magustek.com/bigdata/dass/iotdriver/OP2_DL_698/utils"
	"encoding/binary"
	"fmt"
	"gitee.com/iotdrive/tools/logs"
	"strconv"
)

type SetRequest struct {
	piid      byte
	oadLength byte
	oads      []string
	setType   byte
	data      utils.DL698Data
}

func (s *SetRequest) GetType() utils.APDUType {
	return utils.Set_Request
}

func (s *SetRequest) Encode() ([]byte, error) {
	return nil, nil
}

func (s *SetRequest) Decode(data []byte) error {
	var err error

	s.oads = make([]string, 0)
	s.setType = data[1]
	switch s.setType {
	case 1:
		//fmt.Println("SetRequestNormal")
		s.piid = data[2]
		s.oadLength = 1
		oad := ""
		for i := 0; i < 4; i++ {
			oad += fmt.Sprintf("%02x", data[3+i])
		}
		s.oads = append(s.oads, oad)
		//encode datatype
		fmt.Printf("DataType:% 02X\n", data[7:len(data)-1])
		dt, dErr := utils.DecodeDL698Data(data[7 : len(data)-1])
		if dErr == nil {
			s.data = dt
		} else {
			logs.Error("DL698 DataType Decode err:", dErr)
		}
	case 2:
		//fmt.Println("SetRequestNormalList")
	case 3:
		//fmt.Println("SetThenGetRequestNormalList")
	default:
		err = fmt.Errorf("unSupport set request type %v", s.setType)

	}
	return err
}

func (s *SetRequest) GenOutGoing(in utils.APDU) {

}

type SetResponse struct {
	piid_cad  byte
	oadLength byte
	oads      []string
	setType   byte
	values    []byte
}

func (s *SetResponse) GetType() utils.APDUType {
	return utils.Set_Response
}

func (s *SetResponse) Encode() ([]byte, error) {
	data := []byte{byte(utils.Set_Response), s.setType, s.piid_cad}
	if s.oadLength == 1 {

	} else {
		data = append(data, s.oadLength)
	}
	for i := 0; i < int(s.oadLength); i++ {
		oi, err := strconv.ParseUint(s.oads[i], 16, 64)
		if err != nil {
			return nil, err
		}
		ob := make([]byte, 4)
		binary.BigEndian.PutUint32(ob, uint32(oi))
		data = append(data, ob...)
		data = append(data, s.values[i])
	}
	data = append(data, 0x00, 0x00)
	return data, nil
}

func (s *SetResponse) Decode(bytes []byte) error {
	return nil
}

func (s *SetResponse) GenOutGoing(in utils.APDU) {
	if in.GetType() == utils.Set_Request {
		ind := in.(*SetRequest)
		s.piid_cad = ind.piid
		s.oadLength = ind.oadLength
		s.oads = make([]string, 0)
		s.values = make([]byte, 0)
		s.setType = ind.setType
		s.oads = append(s.oads, ind.oads...)
		for _, oad := range s.oads {
			//执行控制
			//根据OAD执行控制
			var errCode byte
			if ind.data != nil {
				//fmt.Println("Control OAD:", oad)
				err := DLContorl.DoControl(oad, ind.data.GetValue())
				if err != nil {
					logs.Error(fmt.Sprintf("OAD %s Control err %s", oad, err.Error()))
					errCode = 0xff
				}
			} else {
				logs.Error("Not found DL698DataType for control")
				errCode = 0xff
			}
			s.values = append(s.values, errCode)
		}
	} else {
		panic("TODO")
	}
}
