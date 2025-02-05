package znet

import (
	"fmt"
	"net"

	"github.com/rocksun/hellogo/zinxx/utils"
	"github.com/rocksun/hellogo/zinxx/ziface"
)

// 连接模块
type Connection struct {
	//当前连接的socket TCP套接字
	Conn *net.TCPConn

	//连接的ID
	ConnID uint32

	//当前的连接状态
	isClosed bool


	//告知当前连接已经推出/停止 channel
	ExitChan chan bool

	//该连接处理的方法Router
	Router ziface.IRouter
}

//初始化连接模块的方法
func NewConnection(conn *net.TCPConn,connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn: conn,
		ConnID: connID,
		Router : router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}
	return c
}

//实现方法
//连接的读业务
func (c *Connection) startReader() {
	fmt.Println("Reader Gorountine is running")
	defer fmt.Println("connID = ",c.ConnID,"Reader is exit, remote addr is",c.RemoteAddr().String())
	defer c.Stop()
	for {
		//读取客户端的数据到buf中，最大字节由全局配置定义
		buf := make([]byte,utils.GlobalObject.MaxPackageSzie)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err",err)
			continue
		}

		// 得到当前conn数据的Request请求数据
		req := &Request {
			conn : c,
			data : buf,
		}

		//执行注册的路由方法
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(req)
		//从路由中，找到注册绑定的Conn对应的Router调用
	}
}

//启动连接 让当前连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start ... ConnID = ",c.ConnID)
	//启动从当前连接的读数据的业务
	go c.startReader()
	//TODO 启动从当前连接写数据的业务



}

//停止连接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop().. ConnID = ",c.ConnID)
	if c.isClosed {
		return 
	}
	c.isClosed = true

	//关闭socket连接
	c.Conn.Close()

	//回收资源
	close(c.ExitChan)

}

//获取当前连接所绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前连接模块的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的 TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//发送数据 将数据发送给远程的客户端
func (c *Connection) Send(data []byte) error {
	return nil
}

