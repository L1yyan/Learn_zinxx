package znet

import (

	"fmt"
	"net"

	"github.com/rocksun/hellogo/zinxx/utils"
	"github.com/rocksun/hellogo/zinxx/ziface"
)

// iSever的接口实现，定义一个Server的服务器类
type Server struct {
	//服务器名称
	Name string
	//服务器绑定的IP版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int
	//当前server消息管理模块，用来绑定msgID和对应的处理业务API
	MsgHandle ziface.IMsgHandle
	//该server的连接管理器
	ConnMgr ziface.IconnManager
}



// 启动服务器
func (s *Server) Start() {
	//TODO 可以把下面这些日志文件打到txt文件里
	fmt.Printf("[Zinx] Server Name: %s, Listenner at IP: %s,Port: %d ,is starting\n",utils.GlobalObject.Name,utils.GlobalObject.Host,utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",utils.GlobalObject.Version,utils.GlobalObject.MaxConn,utils.GlobalObject.MaxPackageSzie)
	
	
	go func() {
		//0 开启消息队列和工作池
		s.MsgHandle.StartWorkerPool()
		
		// 1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error :", err)
			return
		}

		//2 监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("Listen", s.IPVersion, "err", err)
			return
		}
		fmt.Println("start zinx server succ", s.Name, "succ, Listening..")
		
		var cid uint32 = 0

		//3 阻塞的等待客户端连接，处理客户端连接业务（读写）
		for {
			//如果偶客户端连接过来，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			//设置最大连接个数的判断，如果超过最大连接个数，则关闭新连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				//TODO 给客户端响应一个超出最大连接的错误包
				fmt.Println("Too Many Connection! MaxConn= ",utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}
			//将处理新连接的业务方法和Conn 进行绑定 得到我们的连接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandle)
			cid ++
			
			//启动当前的业务连接处理
			go dealConn.Start()
		}
	}()

}

// 停止服务器
func (s *Server) Stop() {
	//TODO 将服务器一些资源·状态或者一些已经开辟的连接信息进行停止或回收
	fmt.Println("[STOP] Zinx server name: ",s.Name)
	s.ConnMgr.ClearConn()
}

// 运行服务器
func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	//TODO 做一些启动服务器之后的额外业务

	//阻塞状态
	select {}
}


//路由功能，给当前的服务注册一个路由方法，供客户端的连接处理使用
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandle.AddRouter(msgID, router)
	fmt.Println("Add Router Succ!!")
}

func (s *Server) GetConnMgr() ziface.IconnManager {
	return s.ConnMgr
}
func NewServer(name string) ziface.Iserver {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		MsgHandle: 	NewMsgHandle(),
		ConnMgr: NewConnManager(),
	}
	return s
}
