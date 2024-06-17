package DLContorl

import (
	"encoding/json"
	"errors"
	"gitee.com/iotdrive/tools/csapi"
	"gitee.com/iotdrive/tools/logs"
	"net"
	"time"
)

type ControlProxy struct {
	controlSystemAddress string //控制代理地址
	username             string //控制代理登录用户名
	password             string //控制代理登录密码
	conn                 *net.TCPConn
}

func NewControlProxy(controlSystemAddress string, username string, password string) *ControlProxy {
	return &ControlProxy{
		controlSystemAddress: controlSystemAddress,
		username:             username,
		password:             password,
	}
}

func (c *ControlProxy) Setup() (err error) {
	if c.controlSystemAddress == "" {
		err = errors.New("control system address can't be null")
		return
	}
	// //连接boxControlSystem
	// logs.Info("Connect address:", c.controlSystemAddress)
	// err = c.Conn()
	// if err != nil {
	// 	return
	// }
	// logs.Info("Connect to control system successfully")
	// //登录boxControlSystem
	// c.Login()
	// logs.Info("Login control system success")
	// go c.HeartBeat()
	return
}

func (c *ControlProxy) Conn() (err error) {
	c.conn, err = csapi.GetConn(c.controlSystemAddress)
	return
}

// func (c *controlProxy) HeartBeat() {
// 	for {
// 		logs.Debug("heart beat")
// 		ack := &csapi.Command{
// 			Operator: csapi.TYPE_IDLE,
// 		}
// 		c.Write(*ack)
// 		time.Sleep(time.Second * 10)
// 	}
// }

func (c *ControlProxy) Close() (err error) {
	if c.conn == nil {
		return errors.New("conn has been closed")
	}
	err = c.conn.Close()
	c.conn = nil
	return
}

func (c *ControlProxy) Login() (err error) {
	loginCmd := csapi.Command{
		GN:       c.username,
		AV:       c.password,
		Operator: csapi.TYPE_LOGIN,
	}
	_, err = c.Write(loginCmd)
	return
}

func (c *ControlProxy) Reconnect() (err error) {
	c.Close()
	for i := 0; i < 5; i++ {
		if c.Conn() == nil {
			return nil
		}
		time.Sleep(time.Second)
	}
	return errors.New("reconnect failed for 5 times")
}

// 建立连接与重连
func (c *ControlProxy) CheckConn() error {
	if c.conn == nil {
		err := c.Reconnect()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ControlProxy) Write(sendCmd csapi.Command) (resp csapi.Ack, err error) {
	data, err := json.Marshal(&sendCmd)
	if err != nil {
		return
	}
	err = csapi.Write(c.conn, data)
	if err != nil {
		logs.Error(err)
		return
	}
	b, err := csapi.Read(c.conn, csapi.NET_READ_DEFAULT_TIMEOUT)
	if err != nil {
		logs.Error(err)
		return
	}
	err = json.Unmarshal(b, &resp)
	return
}
