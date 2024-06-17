package model

import (
	"dev.magustek.com/bigdata/dass/iotdriver/OP2_DL_698/dataExchange"
	"dev.magustek.com/bigdata/dass/iotdriver/OP2_DL_698/utils"
	"errors"
	"fmt"
	"gitee.com/iotdrive/tools/logs"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// 数组怎么支持
type UploadPoint struct {
	Oad string //oad name
	Uid int    //for update value
	GN  string //for control ; array not support control
	RT  uint8
	Dt  utils.DL698DataType //数据类型
	//数组类型：
	Value   interface{} //单个点的值
	Unit    byte        //单位 代号：0-255
	IsArray bool        //
	index   byte
	//
	piMap  map[int]byte  //保存数组Index，用于查找
	values []interface{} //校验规则
}

func convertUploadPoint(ad, sr string) (*UploadPoint, error) {
	up := new(UploadPoint)
	if ad == "" {
		return nil, errors.New("ad is empty")
	}
	if sr == "" {
		return nil, errors.New("sr is empty")
	}

	dt, err := utils.GetDLDataType(sr)
	if err != nil {
		return nil, err
	}

	if dt == utils.DT_Scaler_Uint {
		//单位类型获取；Unit
		ss := strings.Split(ad, ";")
		if len(ss) == 2 {
			unitNum, err := strconv.Atoi(ss[1])
			if err != nil {
				return nil, fmt.Errorf("unit num exchange err:%s", err.Error())
			}

			if unitNum >= 0 && unitNum <= 255 {
				up.Unit = byte(unitNum)
			} else {
				return nil, fmt.Errorf("unit num is %d ;out of range (0-255)", unitNum)
			}
		} else {
			return nil, fmt.Errorf("SR is Scaler_Uint ,AD need unit;ex: oad;unit")
		}
		ad = ss[0]
	}

	//数组
	ss := strings.Split(ad, "[")
	if len(ss) == 2 {
		up.Oad = ss[0]
		up.IsArray = true
		ad = ss[1]
		ss = strings.Split(ad, "]")
		if len(ss) == 2 {
			indexNum, err := strconv.Atoi(ss[0])
			if err != nil {
				return nil, fmt.Errorf("index num exchange err:%s", err.Error())
			}
			if indexNum >= 0 && indexNum <= 255 {
				up.index = byte(indexNum)
			} else {
				return nil, fmt.Errorf("index num is %d ;out of range (0-255)", indexNum)
			}
		} else {
			return nil, fmt.Errorf("this point is array ,AD need array index ;ad ex as: oad[index]")
		}
	} else if len(ss) == 1 {
		up.Oad = ss[0]
	} else {
		return nil, fmt.Errorf("ad is error")
	}
	up.Dt = dt
	if len(up.Oad) != 8 {
		return up, fmt.Errorf("oad length must be 8,now oad %s length  is %d ", up.Oad, len(up.Oad))
	}
	return up, nil
}

func (e *ElectricityMeter) UpdateElectricityInfo(uid int, val interface{}) {
	<-e.ch
	defer func() {
		e.ch <- struct{}{}
	}()

	//更新数据
	if up, ok := e.pointMap[uid]; ok {
		if up.IsArray {
			//fmt.Println("pimap:", up.piMap)
			if index, ok := up.piMap[uid]; ok {
				up.values[index] = val
				//fmt.Println("values:", up.values)
			} else {
				logs.Error("uid ", uid, " not exist ,update value error")
			}
		} else {
			up.Value = val
		}
	}
}

func (e *ElectricityMeter) GetElectricityInfo(oad string) (*UploadPoint, error) {
	//fmt.Println("len:", len(e.ch))
	<-e.ch
	defer func() {
		e.ch <- struct{}{}
	}()

	//获取数据
	if up, ok := e.oadMap[oad]; ok {
		return up, nil
	}
	return nil, fmt.Errorf("oad not exit")
}

func (e *ElectricityMeter) Getdl698DataByOAD(oad string) (utils.DL698Data, byte) {
	//fmt.Println("Get Data OAD is :", oad)
	up, err := e.GetElectricityInfo(oad)
	if err != nil {
		logs.Error("GetElectricityInfo ", oad, " error:", err)
		return nil, 0xff
	}
	var val float64
	if up.IsArray {
		dlArray := &utils.DTArray{}
		for _, value := range up.values {
			//所有Go类型数据转为Float64
			if up.Dt != utils.DT_OCTET_STR {
				val, err = getFloat64(value)
				if err != nil {
					logs.Error("Get float64 err:", err)
					return nil, 0xff
				}
			}

			switch up.Dt {
			case utils.DT_Int32:
				dlArray.Value = append(dlArray.Value, &utils.DTInt32{
					Value: val,
				})
			case utils.DT_Uint32:
				dlArray.Value = append(dlArray.Value, &utils.DTUint32{
					Value: val,
				})
			case utils.DT_OCTET_STR:
				switch value.(type) {
				case string:
					dlArray.Value = append(dlArray.Value, &utils.DTOctetString{
						Value: value.(string),
					})
				default:
					return nil, 0xff
				}
			case utils.DT_Int8:
				dlArray.Value = append(dlArray.Value, &utils.DTUint8{
					Value: val,
				})
			case utils.DT_Int16:
				dlArray.Value = append(dlArray.Value, &utils.DTInt16{
					Value: val,
				})
			case utils.DT_Uint8:
				dlArray.Value = append(dlArray.Value, &utils.DTUint8{
					Value: val,
				})
			case utils.DT_Uint16:
				dlArray.Value = append(dlArray.Value, &utils.DTUint16{
					Value: val,
				})
			case utils.DT_Int64:
				dlArray.Value = append(dlArray.Value, &utils.DTInt64{
					Value: val,
				})

			case utils.DT_Uint64:
				dlArray.Value = append(dlArray.Value, &utils.DTUint64{
					Value: val,
				})
			case utils.DT_Float32:
				dlArray.Value = append(dlArray.Value, &utils.DTFloat32{
					Value: val,
				})
			case utils.DT_Float64:
				dlArray.Value = append(dlArray.Value, &utils.DTFloat64{
					Value: val,
				})
			case utils.DT_DateTime_S:
				dlArray.Value = append(dlArray.Value, &utils.DTDateTimeS{
					Value: time.Now(),
				})
			default:
				logs.Error("unSupport DTType!")
				return nil, 0xff
			}
		}
		return dlArray, 0
	} else {
		//fmt.Println("Value:", up.Value)
		//所有Go类型数据转为Float64
		if up.Dt != utils.DT_OCTET_STR {
			val, err = getFloat64(up.Value)
			if err != nil {
				logs.Error("Get float64 err:", err)
				return nil, 0xff
			}
		}
		switch up.Dt {
		case utils.DT_Scaler_Uint:
			return &utils.DTScalerUint{
				Value: val,
				Unit:  up.Unit,
			}, 0
		case utils.DT_Int32:
			return &utils.DTInt32{
				Value: val,
			}, 0
		case utils.DT_Uint32:
			return &utils.DTUint32{
				Value: val,
			}, 0
		case utils.DT_OCTET_STR:
			switch up.Value.(type) {
			case string:
				return &utils.DTOctetString{
					Value: up.Value.(string),
				}, 0
			default:
				logs.Error("Now DTTYpe is DT_OCTET_STR,Value need be String")
				return nil, 0xff
			}
		case utils.DT_Int8:
			return &utils.DTInt8{
				Value: val,
			}, 0
		case utils.DT_Int16:
			return &utils.DTInt16{
				Value: val,
			}, 0
		case utils.DT_Uint8:
			return &utils.DTUint8{
				Value: val,
			}, 0
		case utils.DT_Uint16:
			return &utils.DTUint16{
				Value: val,
			}, 0
		case utils.DT_Int64:
			return &utils.DTInt64{
				Value: val,
			}, 0
		case utils.DT_Uint64:
			return &utils.DTUint64{
				Value: val,
			}, 0
		case utils.DT_Float32:
			return &utils.DTFloat32{
				Value: val,
			}, 0
		case utils.DT_Float64:
			return &utils.DTFloat64{
				Value: val,
			}, 0
		case utils.DT_DateTime_S:
			return &utils.DTDateTimeS{
				Value: time.Now(),
			}, 0
		default:
			logs.Error("unSupport DTType!")
			return nil, 0xff
		}
	}
}

func (e *ElectricityMeter) RegisterOAD(up *UploadPoint) {
	_ = dataExchange.RegisterOAD(up.Oad, e.Getdl698DataByOAD)
}

func getFloat64(ptrval interface{}) (val float64, err error) {
	if ptrval == nil {
		return 0, fmt.Errorf("ptrVal is nil")
	}
	switch ptrval.(type) {
	case uint8:
		val = float64(ptrval.(uint8))
	case int8:
		val = float64(ptrval.(int8))
	case int16:
		val = float64(ptrval.(int16))
	case uint16:
		val = float64(ptrval.(uint16))
	case uint32:
		val = float64(ptrval.(uint32))
	case int32:
		val = float64(ptrval.(int32))
	case uint64:
		val = float64(ptrval.(uint64))
	case int64:
		val = float64(ptrval.(int64))
	case int:
		val = float64(ptrval.(int))
	case uint:
		val = float64(ptrval.(uint))
	case float32:
		val = float64(ptrval.(float32))
	case float64:
		val = ptrval.(float64)
		if math.IsNaN(val) || math.IsInf(val, 0) {
			return 0, errors.New("浮点数编码错误")
		}
	case bool:
		if ptrval.(bool) {
			val = 1
		} else {
			val = 0
		}
	case string:
		val, err = strconv.ParseFloat(ptrval.(string), 64)
		if err != nil {
			return 0, err
		}
	default:
		err = errors.New("float64不支持的数据类型：" + reflect.TypeOf(ptrval).String())
	}
	return val, err
}
