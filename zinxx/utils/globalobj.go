package utils

import (
	"encoding/json"
	"os"

	"github.com/rocksun/hellogo/zinxx/ziface"
)

/*
	存储一切有关zinx框架的全局参数，供其他模块使用
	一些参数是可以通过zinx.json 由用户进行配置
*/

type GlobalObj struct {
	/*
		Server
	*/
	TcpServer ziface.Iserver //当前Zinx全局的Server对象
	Host      string         //当前主机监听的IP
	TcpPort   int            //当前服务器主机监听的端口号
	Name      string         //当前服务器的名称

	/*
		Zinx
	*/
	Version          string //当前zinx版本号
	MaxConn          int    //当前主机允许的最大连接数
	MaxPackageSzie   uint32 //当前zinx框架数据包的最大值
	WorkerPoolSize   uint32 //当前业务工作池的goroutine数量
	MaxWorkerTaskLen uint32 //zinx框架允许用户最多开辟多少个worker（限定条件）
}

//定义一个全局的对外GlobalObj

var GlobalObject *GlobalObj

// 从zinx.json 去加载用于自定义的参数
func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}

	//将json文件数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}

}

// 提供一个init方法，初始化当前的GlobalObject
func init() {
	GlobalObject = &GlobalObj{
		//如果配置文件没有加载，默认的值
		Name:             "ZinxServerApp",
		Version:          "V0.6",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSzie:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}

	//应该尝试从conf/zinx.json 去加载一些用户自定义的参数
	GlobalObject.Reload()
}
