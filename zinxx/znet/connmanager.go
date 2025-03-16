package znet

import (
	"errors"
	"fmt"
	"sync"

	"github.com/rocksun/hellogo/zinxx/ziface"
)

/*连接管理模块*/

type ConnManager struct {
	connections map[uint32]ziface.Iconnection //管理的连接集合
	connLock    sync.RWMutex                  //保护连接集合的读写锁

}

// 创建当前连接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.Iconnection),
	}
}

// 添加连接
func (connMgr *ConnManager) Add(conn ziface.Iconnection) {
	//保护共享资源map， 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.UnLock()

	//将conn加入到ConnManager中
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("connID = ", conn.GetConnID(), "Connection add to ConnManager successfully conn num = ", connMgr.Len)
}

// 删除连接
func (connMgr *ConnManager) Remove(conn ziface.Iconnection) {
	//保护共享资源map， 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.UnLock()

	//删除连接信息
	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("Connection delete from ConnManager successfully conn num = ", connMgr.Len)
}

// 根据connID获取连接
func (connMgr *ConnManager) Get(connID uint32) (ziface.Iconnection, error) {
	//保护共享资源map， 加写锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.UnLock()

	if conn, ok := connMgr.connections[connID]; ok {
		//存在ID对应连接
		return conn, nil
	} else {
		return nil, errors.New("connection not Found!")
	}
}

// 得到当前连接总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

// 清除并终止所有连接
func (connMgr *ConnManager) ClearConn() {
	//保护共享资源map， 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除conn并停止conn的工作
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()

		//删除
		delete(connMgr.connections, connID)

		fmt.Println("Clear All connections succ ! conn num = ", connMgr.Len())
	}
}
