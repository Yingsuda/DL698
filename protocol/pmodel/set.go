package pmodel

import (
	"dl698/utils"
	"fmt"
)

type SetRequest struct {
}

func (s *SetRequest) GetType() utils.APDUType {
	return utils.Set_Request
}

func (s *SetRequest) Encode() ([]byte, error) {
	return nil, nil
}

func (s *SetRequest) Decode(bytes []byte) error {
	fmt.Println("Set Request")
	return nil
}

func (s *SetRequest) GenOutGoing(in utils.APDU) {

}

type SetResponse struct {
}

func (s *SetResponse) GetType() utils.APDUType {
	return utils.Set_Response
}

func (s *SetResponse) Encode() ([]byte, error) {
	data := make([]byte, 0)
	fmt.Println("Set Response")

	return data, nil
}

func (s *SetResponse) Decode(bytes []byte) error {
	return nil
}

func (s *SetResponse) GenOutGoing(in utils.APDU) {
	if in.GetType() == utils.Set_Request {
		fmt.Println("")
	} else {
		panic("TODO")
	}
}
