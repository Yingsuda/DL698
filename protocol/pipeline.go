package protocol

import (
	"dev.magustek.com/bigdata/dass/iotdriver/OP2_DL_698/protocol/pmodel"
	"dev.magustek.com/bigdata/dass/iotdriver/OP2_DL_698/utils"
	"fmt"
	"gitee.com/iotdrive/tools/logs"
	"io"
	"net"
	"time"
)

type PipeLine struct {
	dial func() (net.Conn, error)
	conn net.Conn
	//Ctx        context.Context
	//Cancel     context.CancelFunc
	connStatus bool
	addr       string //sa
	ca         byte
	fh         *DL698FrameHandler //数据读取和写入 统一出入口，暂时没用
	pc         *DL698PacketCodec
	mp         *DL698MsgProcessor
}

func NewPipeLineWithDial(dial func() (net.Conn, error), SAddr string) *PipeLine {
	var p *PipeLine
	conn, err := dial()
	if err != nil {
		fmt.Println("dial err:", err)
		p = &PipeLine{
			dial:       dial,
			addr:       SAddr,
			connStatus: false,
			//fh:   NewJT808FrameHandler(conn),
			pc: GetJT808PacketCodec(),
			mp: GetJT808MsgProcessor(),
			ca: 0x33,
		}
	} else {
		fmt.Println("dial ok")
		p = &PipeLine{
			dial:       dial,
			addr:       SAddr,
			conn:       conn,
			connStatus: true,
			//fh:   NewJT808FrameHandler(conn),
			pc: GetJT808PacketCodec(),
			mp: GetJT808MsgProcessor(),
			ca: 0x33,
		}
	}
	return p
}

func NewPipeLineWithConn(conn net.Conn, SAddr string) *PipeLine {
	p := &PipeLine{
		addr:       SAddr,
		conn:       conn,
		connStatus: true,
		pc:         GetJT808PacketCodec(),
		mp:         GetJT808MsgProcessor(),
		ca:         0x33,
	}
	return p
}

func (p *PipeLine) ProcessConnWrite(pkd *utils.PacketData, pcd *utils.ProcessData) error {
	data, err := p.pc.Encode(pkd, pcd)
	if err != nil {
		logs.Error("ProcessConnWrite Encode err:", err)
		return err
	}
	logs.Debug("REQ:% 02X", data)
	_, err = p.conn.Write(data)
	if err != nil {
		logs.Error("conn write err:", err)
		return err
	}
	_ = p.conn.SetReadDeadline(time.Now().Add(time.Second * 1))
	res := make([]byte, 1024)
	n, err := p.conn.Read(res)
	if err != nil {
		logs.Error("conn read err:", err)
		return err
	}
	logs.Debug("RES:% 02X", res[:n])
	npkd, err := p.pc.Decode(res[:n])
	if err != nil {
		fmt.Println("Decode err:", err)
		return err
	}
	p.ca = npkd.Address.Ca

	npcd, err := p.mp.Process(npkd)
	if err != nil {
		logs.Error("Process err:", err)
		//return err
	}
	logs.Debug("New Process Data: %+v\n", npcd)
	return nil
}

// 数组
func (p *PipeLine) ProcessRead(data []byte) error {
	logs.Debug("RECV:% 02X\n", data)
	pkd, err := p.pc.Decode(data)
	if err != nil {
		logs.Error("Decode err:", err)
		return err
	}
	pcd, err := p.mp.Process(pkd)
	if err != nil {
		logs.Error("Process err:", err)
		return err
	}
	send, err := p.pc.Encode(pkd, pcd)
	if err != nil {
		logs.Error("PC Encode err:", err)
		return err
	}
	logs.Debug("SEND:% 02X\n", send)
	_, err = p.conn.Write(send)
	return err
}

func (p *PipeLine) Login() (err error) {
	if !p.connStatus {
		return fmt.Errorf("conn not ok")
	}
	pkd := &utils.PacketData{
		Control: utils.CTRLArea{
			FuncCode: 0x01,
			Dir:      0x01,
			Prm:      0x00,
			Framing:  0x00,
			Blur:     0x00,
		},
		Address: utils.AddressArea{
			Sa: utils.AddressSA{
				AddressType:  0x00,
				LogicAddress: 0x00,
				Address:      p.addr,
				AddressLen:   byte(len(p.addr) / 2),
			},
			Ca: p.ca,
		},
	}

	pcd := &utils.ProcessData{
		OutComing: &pmodel.LinkRequest{
			Piid_acd:    0,
			RequestType: 0, //0 登录  1：心跳  2：退出登录
			Heartbeat:   60,
		},
	}

	return p.ProcessConnWrite(pkd, pcd)
}

func (p *PipeLine) HeartBeat() (err error) {
	if !p.connStatus {
		return fmt.Errorf("conn not ok")
	}
	pkd := &utils.PacketData{
		Control: utils.CTRLArea{
			FuncCode: 0x01,
			Dir:      0x01,
			Prm:      0x00,
			Framing:  0x00,
			Blur:     0x00,
		},
		Address: utils.AddressArea{
			Sa: utils.AddressSA{
				AddressType:  0x00,
				LogicAddress: 0x00,
				Address:      p.addr,
				AddressLen:   byte(len(p.addr) / 2),
			},
			Ca: p.ca,
		},
	}

	pcd := &utils.ProcessData{
		OutComing: &pmodel.LinkRequest{
			Piid_acd:    0,
			RequestType: 1, //0 登录  1：心跳  2：退出登录
			Heartbeat:   60,
		},
	}

	return p.ProcessConnWrite(pkd, pcd)
}

func (p *PipeLine) Connect() {
	pkd := &utils.PacketData{
		Control: utils.CTRLArea{
			FuncCode: 0x03,
			Dir:      0x00,
			Prm:      0x01,
			Framing:  0x00,
			Blur:     0x00,
		},
		Address: utils.AddressArea{
			Sa: utils.AddressSA{
				AddressType:  0x00,
				LogicAddress: 0x00,
				Address:      p.addr,
				AddressLen:   byte(len(p.addr) / 2),
			},
			Ca: p.ca,
		},
	}

	pcd := &utils.ProcessData{
		OutComing: &pmodel.ConnectRequest{
			Piid:            0,
			ProtocolVersion: 22,
			MaxSend:         1024,
			MaxReceive:      1024,
			MaxSize:         1,
			MaxAPDULength:   7000,
			WantTimeout:     7200,
		},
	}

	p.ProcessConnWrite(pkd, pcd)
}

func (p *PipeLine) Start() {
	data := make([]byte, 1024)
	for {
		_ = p.conn.SetReadDeadline(time.Now().Add(time.Second * 1))
		n, err := p.conn.Read(data)
		if err != nil {
			if err == io.EOF {
				p.connStatus = false
				fmt.Println("tcp server close conn")
				_ = p.conn.Close()
				return
			}
			continue
			//其他错误：连接中断
		}

		err = p.ProcessRead(data[:n])
		if err != nil {
			fmt.Println("Process Read err:", err)
		}
	}
}

func (p *PipeLine) Status() bool {
	return p.connStatus
}
