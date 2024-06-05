package utils

import (
	"encoding/binary"
	"time"
)

type DL698Data interface {
	Encode() ([]byte, error)
}

type DL698DataType byte

const (
	DT_NULL       DL698DataType = 0
	DT_ARRAY      DL698DataType = 0x1
	DT_Int32      DL698DataType = 0x05 //double-long
	DT_Uint32     DL698DataType = 0x06 //double-long-unsigned
	DT_OCTET_STR  DL698DataType = 0x09 //octet-string
	DT_Int8       DL698DataType = 0x0f //integer
	DT_Int16      DL698DataType = 0x10 //long
	DT_Uint8      DL698DataType = 0x11 //unsigned
	DT_Uint16     DL698DataType = 0x12 //long-unsigned
	DT_Int64      DL698DataType = 0x14 //long64
	DT_Uint64     DL698DataType = 0x15 //long64-unsigned
	DT_Float32    DL698DataType = 0x17 //float32
	DT_Float64    DL698DataType = 0x18 //float64
	DT_DateTime_S DL698DataType = 0x1C //date_time_s
)

type DTOctetString struct {
	dataType DL698DataType
	Value    string
}

func (dt *DTOctetString) Encode() ([]byte, error) {
	res := []byte{0x01, byte(DT_OCTET_STR), byte(len(dt.Value))}
	res = append(res, []byte(dt.Value)...)
	return res, nil
}

type DTNull struct {
	dataType DL698DataType
}

func (dt *DTNull) Encode() ([]byte, error) {
	return []byte{byte(DT_NULL)}, nil
}

type DTDateTimeS struct {
	dataType DL698DataType
	Value    time.Time
}

func (dt *DTDateTimeS) Encode() ([]byte, error) {
	val := []byte{0x01, byte(DT_DateTime_S)}
	if dt.Value.Unix() < 0 {
		dt.Value = time.Now()
	}
	bt := make([]byte, 7)
	binary.BigEndian.PutUint16(bt[0:], uint16(dt.Value.Year()))
	bt[2] = byte(dt.Value.Month())
	bt[3] = byte(dt.Value.Day())
	bt[4] = byte(dt.Value.Hour())
	bt[5] = byte(dt.Value.Minute())
	bt[6] = byte(dt.Value.Second())
	val = append(val, bt...)
	return val, nil
}

type DTArray struct {
	dataType DL698DataType
	Value    []DL698Data
}

func (dt *DTArray) Encode() ([]byte, error) {
	res := []byte{0x01, byte(DT_ARRAY), byte(len(dt.Value))}
	for _, v := range dt.Value {
		vi, _ := v.Encode()
		res = append(res, vi...)
	}
	return res, nil
}

type DTUint8 struct {
	dataType DL698DataType
	Value    float64
}

func (dt *DTUint8) Encode() ([]byte, error) {
	res := []byte{byte(DT_Uint8)}
	return res, nil
}

type DTInt8 struct {
	dataType DL698DataType
	Value    float64
}

func (dt *DTInt8) Encode() ([]byte, error) {
	res := []byte{byte(DT_Int8)}
	return res, nil
}

type DTInt16 struct {
	dataType DL698DataType
	Value    float64
}

func (dt *DTInt16) Encode() ([]byte, error) {
	res := []byte{byte(DT_Int16)}
	return res, nil
}

type DTUint16 struct {
	dataType DL698DataType
	Value    float64
}

func (dt *DTUint16) Encode() ([]byte, error) {
	val := []byte{byte(DT_Uint16), 0xFF, 0xFF}
	binary.BigEndian.PutUint16(val[1:], uint16(dt.Value))
	return val, nil
}

type DTInt32 struct {
	dataType DL698DataType
	Value    float64
}

func (dt *DTInt32) Encode() ([]byte, error) {
	val := []byte{0x01, byte(DT_Int32), 0x00, 0x00, 0x00, 0x00}
	binary.BigEndian.PutUint32(val[1:], uint32(dt.Value))
	return val, nil
}

type DTUint32 struct {
	dataType DL698DataType
	Value    float64
}

func (dt *DTUint32) Encode() ([]byte, error) {
	res := []byte{byte(DT_Uint32)}

	return res, nil
}

type DTFloat32 struct {
	dataType DL698DataType
	Value    float64
}

func (dt *DTFloat32) Encode() ([]byte, error) {
	val := []byte{byte(DT_Float32)}

	return val, nil
}

type DTFloat64 struct {
	dataType DL698DataType
	Value    float64
}

func (dt *DTFloat64) Encode() ([]byte, error) {
	val := []byte{byte(DT_Float64)}
	return val, nil
}

type DTInt64 struct {
	dataType DL698DataType
	Value    float64
}

func (dt *DTInt64) Encode() ([]byte, error) {
	res := []byte{byte(DT_Int64)}
	return res, nil
}

type DTUint64 struct {
	dataType DL698DataType
	Value    float64
}

func (dt *DTUint64) Encode() ([]byte, error) {
	res := []byte{byte(DT_Uint64)}
	return res, nil
}
