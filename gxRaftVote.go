package rafts

import "time"

type requestVoteArgs struct {
	Term      int
	Candidate int
}

type responseVoteArgs struct {
	Term    int
	Success bool
}

// rf 要 投票给别的机器
func (rf *Raft) requestVote(arg *requestVoteArgs, res *responseVoteArgs) {
	// 加锁
	rf.mu.Lock()
	defer rf.mu.Unlock()

	// 初始化返回值
	res.Term = rf.currentTerm
	res.Success = false

	// 候选人的任期更小
	if arg.Term < rf.currentTerm {
		return
	}

	// 候选人任期更大
	if arg.Term > rf.currentTerm {
		// 更新投票者的任期
		rf.currentTerm = arg.Term

		// 重置状态
		rf.state = Follower

		// 重置投票
		rf.votedFor = -1
	}

	// 过滤初始化 和 重复投票的问题
	if rf.votedFor == -1 || rf.votedFor == arg.Candidate {
		// 开始投票
		rf.votedFor = arg.Candidate
		res.Success = true
		rf.heartPopTime = time.Now()
	}

	return
}
