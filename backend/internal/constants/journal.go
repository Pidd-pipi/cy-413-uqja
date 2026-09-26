package constants

// 心情指数区间：书写提示按区间给出，低/中/高三档
const (
	JournalBandLow  = "low"
	JournalBandMid  = "mid"
	JournalBandHigh = "high"
)

// JournalPromptCount 每个心情区间固定返回的书写提示条数
const JournalPromptCount = 3

// JournalPrompts 按心情区间组织的书写提示，前端点击后带入日记正文
var JournalPrompts = map[string][]string{
	JournalBandLow: {
		"此刻最沉的那种感觉是什么？试着给它起个名字。",
		"如果今天只能照顾好一件小事，你最想先照顾哪一件？",
		"写一句你想听到别人对你说的话，先对自己说一遍。",
	},
	JournalBandMid: {
		"今天有哪三个瞬间值得被记下来，哪怕很普通？",
		"最近什么事情一直在占用你的心力？它值得你继续投入吗？",
		"此刻你的身体感觉如何？它想提醒你什么？",
	},
	JournalBandHigh: {
		"今天是什么让你心情不错？把细节写下来，方便以后回味。",
		"这份好心情里，有没有想感谢的人或事？",
		"如果把今天的状态存起来，下次低落时你想怎么使用它？",
	},
}
