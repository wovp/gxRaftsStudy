package rafts

type AppendEntriesArg struct {
	// 当前消息的任期和索引
	Term     int
	LeaderId int

	// 前一条日志的索引和任期
	PreIndex int
	PreTerm  int

	// 要复制的日志条目
	Entries []LogEntry

	// 领导者已提交的日志记录索引
	LeaderCommitedIndex int
}

type AppendEntriesResponse struct {
	// 当前消息的任期
	Term    int
	Success bool
}

func (rf *Raft) AppendEntries(arg *AppendEntriesArg, res *AppendEntriesResponse) {
	res.Term = rf.currentTerm
	res.Success = false
	// 检查任期
	if arg.Term < rf.currentTerm {
		return
	}

	lastIndex := len(rf.logs) - 1
	RfLastIndex := rf.logs[lastIndex].Index
	RfLastTerm := rf.logs[lastIndex].TermIndex
	// 检查日志
	// 当前一个日志的任期或索引不匹配的花，就说明是非法的
	if RfLastTerm != arg.PreTerm || RfLastIndex != arg.PreIndex {
		return
	}

	// 复制日志
	// 先解决重复日志的问题
	index := arg.PreIndex
	for i, x := range arg.Entries {
		index += 1
		// 已经处理过这个日志了
		if len(rf.logs) > index {
			// 如果任期相同 就说明是重复日志 从下一个日志开始处理
			// 外层if已经保证索引相同
			if rf.logs[index].TermIndex == x.TermIndex {
				continue
			}
			// 如果不相同 说明之前index后接收到的日志是错误的 需要覆盖
			// 先保存index之前的日志
			rf.logs = rf.logs[:index]
		}
		// 在rf日志后拼接新的日志
		// 在第一个if中已经去过重了
		rf.logs = append(rf.logs, arg.Entries[i:]...)
		break
	}

	// 提交日志
	if rf.commitIndex < arg.LeaderCommitedIndex {
		// 数组的下标
		lastLogArrayIndex := len(rf.logs) - 1

		// 当领导者的提交日志索引大于追随者的最大日志，说明新的日志还没有传过来
		if arg.LeaderCommitedIndex > rf.logs[lastLogArrayIndex].Index {
			rf.commitIndex = rf.logs[lastLogArrayIndex].Index
		} else {
			// 如果追随者的日志索引大于参数的索引，说明领导者还没有提交后面的日志
			rf.commitIndex = arg.LeaderCommitedIndex
		}
		rf.apply()
	}

	return
}
