package ziface

/*连接管理模块抽象层*/

type IconnManager interface {
	//添加连接
	Add (conn Iconnection)
	//删除连接
	Remove(conn Iconnection)
	//根据connID获取连接
	Get(connID uint32) (Iconnection, error)
	//得到当前连接总数
	Len() int
	//清除并终止所有连接
	ClearConn()
}