package rafts

type LogEntry struct {
	// 日志任期
	TermIndex int

	// 日志索引
	Index int

	// 日志内容
	Command interface{}
}
