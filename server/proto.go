package main

type TProto struct {
	Cmd    int
	proto1 TProto1Req
	proto2 TProto1Rsp
	// proto2  TProto2
	proto3  TProto3
	proto4  TProto4
	proto5  TProto5
	proto6  TProto6
	proto7  TProto7
	proto8  TProto8
	proto9  TProto9
	proto10 TProto10
	proto11 TProto11
	proto12 TProto12
	proto13 TProto13
	proto14 TProto14
	proto15 TProto15
}

// ”我叫姚梦嘉， 我要来玩狼人杀“
type TProto1Req struct {
	Name string
}

// “姚梦嘉请坐2号座位”  / “2号姚梦嘉加入了游戏， 当前座位列表人员是 1,2，3,4,5,6,7,8”
type TProto1Rsp struct {
	Pos int
}

type TProto3 struct {
	Name1 string
	Name2 string
	Name3 string
	Name4 string
	Name5 string
	Name6 string
	Name7 string
	Name8 string

	EnterName string // 名字
	Pos       int    // 位置
	Status    int    // 1进 2出
}

// -----

// “游戏开始，  你的职业是 XXXX”
type TProto4 struct {
	zhieyey int // 1PM  2 JC  3SS
}

// ----------

// “天黑了”
type TProto5 struct {
}

// “警察你们有 2 5 8 三人，请决定查验人”
type TProto6 struct {
	Police []int
}

//    “”我们查2号姚梦嘉”
type TProto7 struct {
	Object int
}

// “杀手你们有 1 3 4 三人， 请决定要杀的人”
type TProto8 struct {
	Killer []int
}

//    “我们杀2号姚梦嘉”
type TProto9 struct {
	Object int
}

// ----------30秒过去 -----------------

// “2号身份。。。。。”
type TProto10 struct {
	SF int // 1 好 2 坏
}

// “天亮了， 昨晚2号死亡  2号请发遗言”
type TProto11 struct {
	Dead  int // 1 好 2 坏
	YiYan bool
	Turn  int
}

// -----  顺序发言 -------------
type TProto12 struct {
	Turn int
}

// 发言内容
type TProto13 struct {
	Chat string
}

// “大家投票”
type TProto14 struct {
	Object int
}

// “X号被投票出具， 游戏继续、结束XX胜利”
type TProto15 struct {
	Result   int
	GameOver bool
	Winner   int
}
