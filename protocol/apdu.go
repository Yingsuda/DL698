package protocol

type APDU interface {
	GetType() APDUType
	Encode([]byte) (interface{}, error)
	Decode()
}
