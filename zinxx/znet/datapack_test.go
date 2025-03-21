package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 知识负责测试datapack拆包封包 的单元测试
func TestDataPack(t *testing.T) {
	//模拟的服务器
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error", err)
			}

			go func(conn net.Conn) {
				dp := NewDataPack()
				for {
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error")
						break
					}
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err:", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						
						msg := msgHead.(*Message)
						msg.Data = make([]byte,msg.GetMsgLen())

						_, err := io.ReadFull(conn,msg.Data)
						if err != nil {
							fmt.Println("server unpack data err:", err)
							return 
						}
						fmt.Println("|--> recv msgid: ",msg.Id,"  |datalen = ",msg.DataLen,"  |Data = ",string(msg.Data),"  |")
					}
				}

			}(conn)
		}
	}()



	//模拟的客户端
	conn, err := net.Dial("tcp","127.0.0.1:7777")
	if err !=nil {
		fmt.Println("client dial err : ", err)
		return 
	}

	dp := NewDataPack()

	msg1 := &Message{
		Id : 1,
		DataLen: 4,
		Data: []byte {'z','i','n','x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("clent pack msg1 error",err)
		return 
	}

	msg2 := &Message{
		Id : 1,
		DataLen: 7,
		Data: []byte {'n','i','h','a','o','!','!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("clent pack msg1 error",err)
		return 
	}

	sendData1 = append (sendData1,sendData2...)
	conn.Write(sendData1)

	//客户端阻塞
	select{}
}
