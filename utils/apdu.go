package utils

import (
	"errors"
)

type APDUType byte

type APDU interface {
	GetType() APDUType
	Encode() ([]byte, error)
	Decode([]byte) error
	GenOutGoing(in APDU)
}

const (
	LINK_Request     APDUType = 0x01
	LINK_Response    APDUType = 0x81
	Connect_Request  APDUType = 0x02
	Connect_Response APDUType = 0x82
	Get_Request      APDUType = 0x05
	Get_Response     APDUType = 0x85
	Set_Request      APDUType = 0x06
	Set_Response     APDUType = 0x86
)

func GetAPDUType(t byte) (APDUType, error) {
	var at APDUType
	var err error
	switch t {
	case byte(LINK_Response):
		at = LINK_Response
	case byte(LINK_Request):
		at = LINK_Request
	case byte(Connect_Request):
		at = Connect_Request
	case byte(Connect_Response):
		at = Connect_Response
	case byte(Get_Request):
		at = Get_Request
	case byte(Get_Response):
		at = Get_Response
	case byte(Set_Request):
		at = Set_Request
	case byte(Set_Response):
		at = Set_Response

	default:
		err = errors.New("APDUType is not exist")
	}

	if err != nil {
		return at, err
	}
	return at, nil
}
