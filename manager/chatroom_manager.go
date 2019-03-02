package manager

import (
	"sync"
)

type CRManager interface {
	GetIdleRoom() int
	UpdateRoomStatus(int, bool)
}

type crManager struct {
	Rooms map[int]bool //存放空闲房间号
	mu sync.RWMutex //更新房间状态时使用读写锁进行内存同步
}

var (
	crmanager *crManager
	crlock sync.Mutex
	roomCnt = 100
)

func CRInstance() *crManager {
	if crmanager != nil  {
		return crmanager
	}
	crlock.Lock()
	defer crlock.Unlock()

	if crmanager != nil {
		return crmanager
	}
	var rooms = make(map[int]bool)
	//初始化房间状态
	for index := 0; index < roomCnt; index++ {
		rooms[index] = false
	}
	crmanager = &crManager{
		Rooms:rooms,
	}
	return crmanager
}

//返回空闲房间号
func (m *crManager) GetIdleRoom() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for index := 0; index < roomCnt; index++ {
		if m.Rooms[index] {
			return index
		}
	}
	return -1
}

//更新房间状态
func (m *crManager) UpdateRoomStatus(roomID int, status bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Rooms[roomID] = status
}