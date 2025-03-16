package znet

import (
	"fmt"
	"strconv"

	"github.com/rocksun/hellogo/zinxx/utils"
	"github.com/rocksun/hellogo/zinxx/ziface"
)

//消息处理模块实现

type MsgHandle struct {
	//存放每个id所对应的处理方法
	Apis map[uint32]ziface.IRouter
	//负责worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//worker池的数量
	WorkerPoolSize uint32
}

// 初始化
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
		TaskQueue: make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
		//从全局配置中获取也可以在配置文件中让用户设置
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
	}
}

// 调度对应的router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	//从request 找到id
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgid = ", request.GetMsgID(), "is NOT FOUND! Need Add!")
	}
	//根据id调度对应业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//判断当前msg绑定的API是否已经存在
	if _, judge := mh.Apis[msgID]; judge {
		//id已经注册
		panic("repeat api , msgid= " + strconv.Itoa(int(msgID)))
	}
	//添加绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID= ", msgID, " success!")
}

//启动一个worker工作池 开启工作池的动作只能发生一次 一个框架只有一个池子
func (mh* MsgHandle) StartWorkerPool() {
	//根据workerpoolsize分别开启worker 每个worker用一个go来承载
	for  i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//1 给当前的worker对应的channel消息队列开辟空间 
		mh.TaskQueue[i] = make (chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//2 启动当前worker，阻塞等待消息从channel进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}
//启动后一个worker工作流程
func (mh* MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("WorkerID = ",workerID," is started . .")
	//不断的阻塞等待对应消息队列的消息
	for{
		select {
		//如果有消息过来，出列的就是一个客户端的Request，执行当前re绑定的业务	
		case request := <- taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

//将消息交给taskqueue 又worker处理
func (mh* MsgHandle) SendMsgToTaskQueue (request ziface.IRequest) {
	//1 平均分配给不同的worker
	//根据客户端简历的connID进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID= ",request.GetConnection().GetConnID()," requst MsgID= ",request.GetMsgID()," to workerID= ",workerID)
	//2将消息发送给对应的worker的Task'queue
	mh.TaskQueue[workerID] <- request
}