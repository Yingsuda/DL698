package model

import (
	"dev.magustek.com/bigdata/dass/iotdriver/OP2_DL_698/protocol"
	"fmt"
	"net"
	"time"
)

type ElectricityMeter struct {
	//电压
	ch       chan struct{}
	oadMap   map[string]*UploadPoint //获取数据
	pointMap map[int]*UploadPoint    //更新数据
	pl       *protocol.PipeLine

	addr      string
	saAddress string
}

func (e *ElectricityMeter) dialTcp() (net.Conn, error) {
	return net.Dial("tcp", e.addr)
}

func NewElectricityMeter(TcpHost string, SAAddress string) *ElectricityMeter {
	em := &ElectricityMeter{
		ch:        make(chan struct{}, 1),
		addr:      TcpHost,
		pointMap:  make(map[int]*UploadPoint),
		oadMap:    make(map[string]*UploadPoint),
		saAddress: SAAddress,
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
					_ = e.pl.HeartBeat()
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
