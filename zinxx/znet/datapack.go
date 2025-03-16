package znet

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/rocksun/hellogo/zinxx/utils"
	"github.com/rocksun/hellogo/zinxx/ziface"
)

//封包 拆包 具体模块

type DataPack struct {}

//拆包封包实例的一个初始化方法

func NewDataPack() *DataPack{
	return &DataPack{}
}

//获取包的头的长度方法
func(dp *DataPack)	GetHeadLen() uint32 {
	//datalen uint32(4字节) +Id uint32(4字节)
	return 8
}

//封包方法
//datalen|msgid|data
func(dp *DataPack)	Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	//将Datalen写进databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil,err
	}
	//将Msgid写进databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil,err
	}
	//将data数据写进databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil,err
	}
	return dataBuff.Bytes(),nil
}
//拆包方法(将包的head信息都读出来) 之后再根据head信息里的data长度，再进行一次读
func(dp  *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)
	
	//只解压head信息，得到datalen和msgid
	msg := &Message{}
	//读长度
	if err := binary.Read(dataBuff,binary.LittleEndian, &msg.DataLen);err != nil {
		return nil,err
	}

	//读id
	if err := binary.Read(dataBuff,binary.LittleEndian, &msg.Id);err != nil {
		return nil,err
	}

	//判断datalen是否已经长处我们允许的最大长度
	if utils.GlobalObject.MaxPackageSzie >0 && msg.DataLen>utils.GlobalObject.MaxPackageSzie {
		return nil , errors.New("too Large msg data recv")
	}
	return msg, nil
}