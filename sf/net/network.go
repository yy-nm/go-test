package net

// 理论上处理以包为单位的逻辑层是不需要知道链接是否建立,
// 但由于一些退出机制的需求而有连接断开的需求(这个也不是刚性需求, 可以通过其他机制进行实现)
type NetworkCallback interface {
	MsgRecv(m Msg)
	ConnBroken(id ConnId)
	//ConnBuild(id ConnId)
}

type Network interface {
	SendMsg(m Msg)
	CloseConn(m Msg)
	Register(ncb NetworkCallback) (old NetworkCallback)
	Start()
	Stop()
}

type network struct {
	cm     ConnMgr
	fltIn  Filter // filter for income message
	fltOut Filter // filter for send message
	tfm    Transform
}
