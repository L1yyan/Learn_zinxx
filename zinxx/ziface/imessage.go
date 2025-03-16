package ziface

/*
	将请求的消息封装到一个message中，定义抽象接口
*/

type IMessage interface {
	//获取消息id
	GetMsgId() uint32
	//获取消息长度
	GetMsgLen() uint32
	//获取消息内容
	GetData() []byte

	//设置消息id
	SetMsgId(uint32)
	//设置消息内容
	SetData([]byte)
	//设置消息长度
	SetMsgLen(uint32)
}
