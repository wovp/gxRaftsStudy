package rafts

import (
	"sync"
	"time"
)

type Raft struct {
	// 状态
	state int

	// 互斥锁
	mu sync.Mutex

	// 当前最新任期
	currentTerm int

	// 当前任期内投票的机器ID
	votedFor int

	// 心跳事件
	heartPopTime time.Time

	// 日志事件
	logs []LogEntry

	// 已经提交的日志索引
	commitIndex int
}

// 提交日志
func (rf *Raft) apply() {

}
