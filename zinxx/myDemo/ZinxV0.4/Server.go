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

// Test PreHandle
func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle .. ")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping ... \n"))
	if err != nil {
		fmt.Println("Call Back before Ping Error!")

	}
}

// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle .. ")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping ping ping ... \n"))
	if err != nil {
		fmt.Println("Call Back ping...Ping...ping Error!")

	}
}

// Test PostHandle
func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle .. ")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping ... \n"))
	if err != nil {
		fmt.Println("Call Back after ping Error!")

	}
}

func main() {
	//1 创建一个Server句柄，使用Zinx的api
	s := znet.NewServer("[zinxV0.4]")

	//2 给当前znx框架添加一个自定义Router
	s.AddRouter(&PingRouter{})

	//3 启动Server
	s.Serve()
}
