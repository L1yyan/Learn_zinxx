package ziface
/*
消息管理抽象层
*/

type IMsgHandle interface {
	//执行对应的router消息处理方法
	DoMsgHandler(request IRequest)
	//为消息添加具体的处理
	AddRouter(msgID uint32,router IRouter)
	//启动worker工作池
	StartWorkerPool()
	//将消息发送给消息任务队列处理
	SendMsgToTaskQueue (request IRequest)
	
}