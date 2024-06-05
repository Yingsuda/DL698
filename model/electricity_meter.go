package model

import (
	"dl698/protocol"
	"fmt"
	"net"
)

type ElectricityMeter struct {
	//电压
	ch       chan struct{}
	oadMap   map[string]*UploadPoint //获取数据
	pointMap map[int]*UploadPoint    //更新数据
	pl       *protocol.PipeLine
	addr     string
}

func (e *ElectricityMeter) dialTcp() (net.Conn, error) {
	return net.Dial("tcp", e.addr)
}

func NewElectricityMeter(TcpHost string, SAAddress string) *ElectricityMeter {
	em := &ElectricityMeter{
		ch:       make(chan struct{}, 1),
		addr:     TcpHost,
		pointMap: make(map[int]*UploadPoint),
		oadMap:   make(map[string]*UploadPoint),
	}

	em.pl = protocol.NewPipeLine(em.dialTcp, SAAddress)
	em.ch <- struct{}{}
	return em
}

func (e *ElectricityMeter) AddPoint(up *UploadPoint) {
	if oadInfo, ok := e.oadMap[up.Oad]; ok {
		if oadInfo.IsArray {
			if up.Dt == oadInfo.Dt {
				oadInfo.values = append(oadInfo.values, up.Value)
				oadInfo.piMap[up.Uid] = byte(len(oadInfo.values) - 1)
				e.pointMap[up.Uid] = oadInfo //一组的uid指向同一个uploadPoint
				e.RegisterOAD(oadInfo)
			} else {
				fmt.Println("oad array datatype need same")
			}
		} else {
			fmt.Println("not array oad is exist")
		}
	} else {
		//fmt.Println("New:")
		up.values = make([]interface{}, 0)
		up.values = append(up.values, up.Value)
		up.piMap = make(map[int]byte)
		up.piMap[up.Uid] = byte(len(up.values) - 1)
		e.oadMap[up.Oad] = up
		//已经存在，且是数组，添加进去
		e.pointMap[up.Uid] = up //uid 不会重复，用来更新数据的
		e.RegisterOAD(up)
	}

}

func (e *ElectricityMeter) Login() error {
	return e.pl.Login()
}

func (e *ElectricityMeter) Start() {
	//for s, oad := range e.oadMap {
	//	fmt.Printf("oad %s STR:%+v \n", s, *oad)
	//}
	//fmt.Println("pointMap :", e.pointMap)
	e.pl.Start()
}
