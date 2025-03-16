package znet

import (
	"errors"
	"fmt"
	"io"
	"net"

	// "github.com/rocksun/hellogo/zinxx/utils"
	"github.com/rocksun/hellogo/zinxx/utils"
	"github.com/rocksun/hellogo/zinxx/ziface"
)

// 连接模块
type Connection struct {
	//当前conn隶属于哪个server
	TcpServer ziface.Iserver
	//当前连接的socket TCP套接字
	Conn *net.TCPConn

	//连接的ID
	ConnID uint32

	//当前的连接状态
	isClosed bool

	//告知当前连接已经推出/停止 channel(由reader告知writer退出)
	ExitChan chan bool

	//用于读缓冲的管道，用于读写goroutine之间的通信
	msgChan chan []byte
	//消息的管理msgid和对应的业务处理api
	MsgHandler ziface.IMsgHandle
}

// 初始化连接模块的方法
func NewConnection(server ziface.Iserver, conn *net.TCPConn, connID uint32, MsgHandle ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: MsgHandle,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
	}

	//将conn加入到connManager中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

// 实现方法
// 连接的读业务
func (c *Connection) startReader() {
	fmt.Println("Reader Gorountine is running")
	defer fmt.Println("connID = ", c.ConnID, "[Reader is exit!],remote addr is", c.RemoteAddr().String())
	defer c.Stop()
	for {
		//读取客户端的数据到buf中，最大字节由全局配置定义
		// buf := make([]byte,utils.GlobalObject.MaxPackageSzie)
		// _, err := c.Conn.Read(buf)
		// if err != nil {
		// 	fmt.Println("recv buf err",err)
		// 	continue
		// }

		//创建一个拆包解包对象
		dp := NewDataPack()

		//读取客户端的msghead的二进制流 8个字节，
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("Read msg head error :", err)
			break
		}

		//拆包 得到msgid，和msgdatalen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
		}
		//根据datalen 再次读取data 放在 msg.data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				break
			}
		}
		msg.SetData(data)
		// 得到当前conn数据的Request请求数据
		req := &Request{
			conn: c,
			msg:  msg,
		}

		//判断是否开启工作池
		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经开启工作池机制，将消息发送给worker工作池处理即可
			c.MsgHandler.SendMsgToTaskQueue(req)
		} else {
			//从路由中，找到注册绑定的Conn对应的Router调用
			//根据绑定好的MsgID找到对应处理api业务执行
			go c.MsgHandler.DoMsgHandler(req)
		}

	}
}

// 写消息goroutine，专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running...]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	//不断阻塞等待channel消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			//有数据写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data error: ", err)
				return
			}
		case <-c.ExitChan:
			//代表reader已经退出，此时writer也要退出
			return
		}

	}
}

// 启动连接 让当前连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start ... ConnID = ", c.ConnID)
	//启动从当前连接的读数据的业务
	go c.startReader()
	//TODO 启动从当前连接写数据的业务
	go c.StartWriter()
}

// 停止连接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop().. ConnID = ", c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true

	//关闭socket连接
	c.Conn.Close()

	//告知Writer关闭
	c.ExitChan <- true

	//将当前连接从connMgr中摘除掉
	c.TcpServer.GetConnMgr().Remove(c)
	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

// 获取当前连接所绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前连接模块的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的 TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 提供一个sendmsg方法，把发送给客户端的数据先封包再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}

	//将data封包 msgdaatalen msgid msgdata
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg")
	}
	//将数据发送给客户端
	c.msgChan <- binaryMsg
	return nil
}
