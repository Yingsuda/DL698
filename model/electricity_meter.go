package model

import (
	"dev.magustek.com/bigdata/dass/iotdriver/OP2_DL_698/DLContorl"
	"dev.magustek.com/bigdata/dass/iotdriver/OP2_DL_698/protocol"
	"fmt"
	"gitee.com/iotdrive/tools/csapi"
	"gitee.com/iotdrive/tools/logs"
	"net"
	"reflect"
	"strconv"
	"time"
)

type ElectricityMeter struct {
	//电压
	ch       chan struct{}
	oadMap   map[string]*UploadPoint //获取数据
	pointMap map[int]*UploadPoint    //更新数据
	pl       *protocol.PipeLine
	cp       *DLContorl.ControlProxy //for control

	addr      string
	saAddress string
}

func (e *ElectricityMeter) dialTcp() (net.Conn, error) {
	return net.Dial("tcp", e.addr)
}

func NewElectricityMeter(TcpHost, SAAddress, ControlAddress, ControlUsername, ControlPass string) *ElectricityMeter {
	em := &ElectricityMeter{
		ch:        make(chan struct{}, 1),
		addr:      TcpHost,
		pointMap:  make(map[int]*UploadPoint),
		oadMap:    make(map[string]*UploadPoint),
		saAddress: SAAddress,
		cp:        DLContorl.NewControlProxy(ControlAddress, ControlUsername, ControlPass),
	}
	//em.pl = protocol.NewPipeLine(em.dialTcp, SAAddress)
	em.ch <- struct{}{}
	return em
}

func (e *ElectricityMeter) AddPoint(up *UploadPoint) (err error) {
	if oadInfo, ok := e.oadMap[up.Oad]; ok {
		if oadInfo.IsArray {
			if up.Dt == oadInfo.Dt {
				up.piMap[up.Uid] = up.index
				if int(up.index)+1 > len(up.values) {
					ex := make([]interface{}, int(up.index)+1-len(up.values))
					up.values = append(up.values, ex...)
				}
				up.values[up.index] = up.values
				e.pointMap[up.Uid] = oadInfo //一组的uid指向同一个uploadPoint
				e.RegisterOAD(oadInfo)
			} else {
				err = fmt.Errorf("oad array datatype need same")
			}
		} else {
			err = fmt.Errorf("not array oad is exist")
		}
	} else {
		up.piMap = make(map[int]byte)
		up.values = make([]interface{}, 0)
		up.piMap[up.Uid] = up.index
		if int(up.index)+1 > len(up.values) {
			ex := make([]interface{}, int(up.index)+1-len(up.values))
			up.values = append(up.values, ex...)
		}
		up.values[up.index] = up.values
		e.oadMap[up.Oad] = up
		//已经存在，且是数组，添加进去
		e.pointMap[up.Uid] = up //uid 不会重复，用来更新数据的
		e.RegisterOAD(up)
	}
	return err
}

func (e *ElectricityMeter) Reset() {
	e.pointMap = make(map[int]*UploadPoint)
	e.oadMap = make(map[string]*UploadPoint)

}

func (e *ElectricityMeter) Login(conn net.Conn) error {
	e.pl = protocol.NewPipeLineWithConn(conn, e.saAddress)
	err := e.pl.Login()
	if err != nil {
		return err
	}
	go func() {
		ticker := time.NewTicker(time.Second * 60)
		for {
			select {
			case <-ticker.C:
				if e.pl.Status() {
					err = e.pl.HeartBeat()
					if err != nil {
						logs.Error("Heart Beat err:", err)
					}
				} else {
					return
				}
			}
		}
	}()
	return nil
}

func (e *ElectricityMeter) Start() {

	e.pl.Start()
}

func (e *ElectricityMeter) ControlByOAD(oad string, value interface{}) error {
	<-e.ch
	defer func() {
		e.ch <- struct{}{}
	}()

	defer e.cp.Close()
	if up, ok := e.oadMap[oad]; ok {
		//control
		fmt.Printf("control UID %d ;Value %v ,Type is %T \n", up.Uid, value, value)
		//连接管理
		err := e.cp.CheckConn()
		if err != nil {
			return err
		}

		var av string
		switch value.(type) {
		case float64:
			av = strconv.FormatFloat(value.(float64), 'f', 2, 64)
		case string:
			av = value.(string)
		default:
			logs.Notice("Unknown type of value:", reflect.TypeOf(value))
			return fmt.Errorf("unknown type of value %v", reflect.TypeOf(value))
		}

		cmd := csapi.Command{
			UID:      uint32(up.Uid),
			GN:       up.GN,
			AV:       av,
			RT:       up.RT,
			Operator: 1,
		}

		//登录
		err = e.cp.Login()
		if err != nil {
			logs.Warning("登录控制代理失败！", err)
			return fmt.Errorf("登录控制代理失败！")
		}
		logs.Info("control message:UID:%v,GN:%v,AV:%v", cmd.UID, cmd.GN, cmd.AV)

		resp, err := e.cp.Write(cmd)
		if err != nil {
			logs.Error(err)
		}

		if resp.Status == false {
			logs.Error("control err.UID:%v,GN:%v,AV:%v,Msg:%v", resp.UID, cmd.GN, cmd.AV, resp.Msg)
			return fmt.Errorf("control err.UID:%v,GN:%v,AV:%v,Msg:%v", resp.UID, cmd.GN, cmd.AV, resp.Msg)
		}
		logs.Info("control success:UID:%v,Msg:%v,Status:%v,Type:%v", resp.UID, resp.Msg, resp.Status, resp.Type)
		return nil
	} else {
		fmt.Println("ControlByOAD OAD Not Exist")
		return fmt.Errorf("oad %s not exist,control err", oad)
	}
}
