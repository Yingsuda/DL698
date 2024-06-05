package protocol

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
)

var (
	ErrFrameReadEmpty = errors.New("Read empty frame")
)

type FrameHandler interface {
	Receive() ([]byte, error) // data -> frame，并写入io.Writer
	Send([]byte) error        // 从io.Reader中提取frame payload，并返回给上层
}

type DL698FrameHandler struct {
	rBuf   *bufio.Reader
	writer io.Writer
}

func NewDL698FrameHandler(conn net.Conn) *DL698FrameHandler {
	return &DL698FrameHandler{
		rBuf:   bufio.NewReader(conn),
		writer: conn,
	}
}

func (fh *DL698FrameHandler) Receive() ([]byte, error) {
	buf := make([]byte, 1024)
	_, err := fh.rBuf.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("fail to read stream to framePayload: %v", err)
	}
	// 移除末尾多余的0
	buf = bytes.TrimRight(buf, "\x00")
	//fmt.Printf("R: % 02X\n", buf)
	if len(buf) == 0 {
		return nil, ErrFrameReadEmpty
	}

	return buf, nil
}

func (fh *DL698FrameHandler) Send(payload []byte) (err error) {
	var p = payload
	if len(p) == 0 {
		err = fmt.Errorf("the payload is empty when sending, skip")
		return
	}
	for {
		var n int
		n, err = fh.writer.Write([]byte(p))
		if err != nil {
			return
		}
		if n >= len(p) {
			break
		}
		if n < len(p) {
			p = p[n:] // 没写完所有数据，再写一次
		}
	}

	//fmt.Printf("S: % 02X\n", p)
	return
}
