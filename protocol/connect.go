package protocol

type Link_Request struct {
	piid_acd    byte
	requestType byte   //0:登录  1：心跳  2：退出登录
	heartbeat   uint16 //心跳周期
	requestDate string
}

func (l *Link_Request) GetType() APDUType {
	return LINK_Request
}

func (l *Link_Request) Encode(bytes []byte) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (l *Link_Request) Decode() {
	//TODO implement me
	panic("implement me")
}

type Link_Response struct {
	piid         byte
	result       byte   //结果 bit7:始终可信标识 1：可信  0：不可信   bit2-bit0：结果  0：成功  1：地址重复  2：非法设备  3：容量不足
	requestDate  string //请求时间
	receiveDate  string //收到的时间
	responseDate string //回复时间
}

func (l *Link_Response) GetType() APDUType {
	//TODO implement me
	panic("implement me")
}

func (l *Link_Response) Encode(bytes []byte) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (l *Link_Response) Decode() {
	//TODO implement me
	panic("implement me")
}
