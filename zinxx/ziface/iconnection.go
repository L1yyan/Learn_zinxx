package ziface

import "net"

//定义连接模块的抽象层
type Iconnection interface {
	//启动连接 让当前连接准备开始工作
	Start()

	//停止连接 结束当前连接的工作
	Stop()

	//获取当前连接所绑定的socket conn
	GetTCPConnection() *net.TCPConn

	//获取当前连接模块的连接ID
	GetConnID() uint32

	//获取远程客户端的 TCP状态 IP port
	RemoteAddr() net.Addr

	//发送数据 将数据发送给远程的客户端
	SendMsg(msgId uint32, data []byte) error
}

//定义一个处理连接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
