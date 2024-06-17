package model

import (
	"context"
	"dev.magustek.com/bigdata/dass/iotdriver/OP2_DL_698/DLContorl"
	upload_points "gitee.com/iotdrive/leader.upload/model/points"
	"gitee.com/iotdrive/leader.upload/utils/global"
	"gitee.com/iotdrive/tools/logs"
	"github.com/astaxie/beego/config"
	"net"
	"strings"
)

var (
	ctx    context.Context
	Cancel context.CancelFunc
)

type Upload698 struct {
	e         *ElectricityMeter
	port      string
	saAddress string
	runModel  int //0:server  1:client

	controlAddress string
	controlUser    string
	controlPass    string
}

func NewUpload698() *Upload698 {
	ctx, Cancel = context.WithCancel(context.Background())
	return &Upload698{}
}

func (u *Upload698) SetupConf(conf config.Configer) {
	u.port = conf.DefaultString("DL698_PORT", "28080")
	if !strings.HasPrefix(u.port, ":") {
		u.port = ":" + u.port
	}
	logs.Info("Port:", u.port)
	u.saAddress = conf.DefaultString("saAddress", "12345678")
	logs.Info("saAddress:", u.saAddress)
	u.controlAddress = conf.DefaultString("controlAddress", "127.0.0.1:12000")
	logs.Info("controlAddress:", u.controlAddress)

	u.controlUser = conf.DefaultString("controlUser", "admin")
	u.controlPass = conf.DefaultString("controlPass", "admin")

	u.e = NewElectricityMeter("127.0.0.1:"+u.port, u.saAddress, u.controlAddress, u.controlUser, u.controlPass)
	go u.Run()
}

func (u *Upload698) SyncPoints(plist []*upload_points.Point) error {
	u.e.Reset()
	var pNum = 0
	for _, p := range plist {
		if p == nil {
			logs.Error(" PN = %v , is nil", p.GN)
			continue
		}
		logs.Debug("UID:", p.UID, "AD :", p.AD, " SR:", p.SR, " GN:", p.GN)
		//解析测点：生成UploadPoint

		up, err := convertUploadPoint(p.AD, p.SR)
		if err != nil {
			logs.Error("convertUploadPoint err ", err, " AD is :", p.AD, " SR is :", p.SR)
			continue
		}

		//测点是否支持控制
		DLContorl.RegisterControlHandle(up.Oad, u.e.ControlByOAD)

		up.Uid = int(p.UID)
		up.GN = p.GN
		up.RT = uint8(p.RT)
		//保存测点
		//根据OAD保存测点	 --用于获取数据
		//根据UID保存测点  --用于更新数据
		err = u.e.AddPoint(up)
		if err != nil {
			logs.Error("AddPoint err:", err, " point ad is ", p.AD, ";sr is ", p.SR)
			continue
		}
		pNum++
	}
	logs.Info("加载测点全点表，完成")
	logs.Info("测点总计数量: ", pNum)
	return nil
}

func (u *Upload698) PushData(item *upload_points.Item) {
	//数据变化，更新数据
	//fmt.Println("UID:", item.UID, " Value:", item.Value.IValue)
	u.e.UpdateElectricityInfo(int(item.UID), item.Value.IValue)
}

func (u *Upload698) SetHistoryHandle(f func(begin int64, end int64) error) {

}

func (u *Upload698) SetPushControlEventHandle(callback func(data []byte) ([]byte, error)) {

}

func (u *Upload698) Run() {
	defer global.Stop()
	if u.runModel == 0 {
		//监听：
		listener, err := net.Listen("tcp", u.port)
		if err != nil {
			logs.Error("Listen err:", err)
			return
		}
		for {
			conn, err := listener.Accept()
			if err != nil {
				logs.Error("accept err:", err)
				continue
			}
			err = u.e.Login(conn)
			if err != nil {
				logs.Error("Login error:", err)
				_ = conn.Close()
				continue
			}
			go u.e.Start()
		}
	} else if u.runModel == 1 {
		logs.Error("not client")
	} else {
		//监听：
		logs.Error("runModel err")
	}
}
