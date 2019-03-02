package manager

import (
	"github.com/kataras/iris/websocket"
	"sync"
)

type WSManager interface {
	GetWSConnection(string) *websocket.Connection
	AddMapRelationship(string, *websocket.Connection)
	DeleteMapRelationship(string)
}

type wsManager struct {
	Conns map[string]*websocket.Connection //存放用户和websocket连接之间的映射关系
	mu sync.RWMutex //更新映射关系时使用读写锁进行内存同步
}

var (
	wsmanager *wsManager
	wslock sync.Mutex
)

func WSInstance() *wsManager {
	if wsmanager != nil  {
		return wsmanager
	}
	wslock.Lock()
	defer wslock.Unlock()

	if wsmanager != nil {
		return wsmanager
	}
	wsmanager = &wsManager{
		Conns:make(map[string]*websocket.Connection),
	}
	return wsmanager
}

//获取用户对应的socket连接
func (m *wsManager) GetWSConnection(account string) *websocket.Connection {
	m.mu.RLock()
	defer m.mu.RUnlock()
	conn := (m.Conns)[account]
	return conn
}

//新增或更新映射关系
func (m *wsManager) AddMapRelationship(account string, conn *websocket.Connection) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Conns[account] = conn
}

//删除映射关系
func (m *wsManager) DeleteMapRelationship(account string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.Conns, account)
}