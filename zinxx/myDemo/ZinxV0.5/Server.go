package main

import (
	"fmt"

	"github.com/rocksun/hellogo/zinxx/ziface"
	"github.com/rocksun/hellogo/zinxx/znet"
)

//ping test 自定义路由

type PingRouter struct {
	znet.BaseRouter
}



// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle .. ")
	//先读取客户端的数据，再回写ping
	fmt.Println("recv from client : msgID = ",request.GetMsgID(),",data = ",string(request.GetData()))
	err := request.GetConnection().SendMsg(1,[]byte("ping...ping...ping..."));if err !=nil {
		fmt.Println(err)
	}
}



func main() {
	//1 创建一个Server句柄，使用Zinx的api
	s := znet.NewServer("[zinxV0.5]")

	//2 给当前znx框架添加一个自定义Router
	s.AddRouter(&PingRouter{})

	//3 启动Server
	s.Serve()
}
