package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	conn "github.com/266game/goserver/Connection"
	tcpserver "github.com/266game/goserver/TCPServer"
)

type data struct {
	Cmd     int32  `json:"cmd"`
	Content string `json:"content"`
	Mark    int    `json:"mark"`
}
type redata struct {
	Cmd     int32  `json:"cmd"`
	Content string `json:"content"`
	Mark    int    `json:"mark"`
}

// type targetlist struct {
// 	target   int
// 	operator int
// 	// Connect *conn.TData
// }

var numberList sync.Map
var identityList sync.Map
var connList sync.Map
var nameList sync.Map
var nMark int
var policeList []int
var killerList []int
var commonList []int
var dat = &data{}

//查验目标确认
var policetarget sync.Map
var lenpoliceget int
var npoliceget int
var killtarget sync.Map
var lenkillget int
var nkillget int

func timer() {
	timer := time.NewTimer(time.Second * 60)
	go func() {
		//等触发时的信号
		<-timer.C
		proto8()
	}()
}

//除了刚加入的玩家，其他玩家都会受到新玩家加入的信息
func proto2(content string) {
	redat := &redata{}

	redat.Cmd = 2
	redat.Content = content
	jsonStu, err := json.Marshal(redat)
	fmt.Println("json字符串", string(jsonStu))
	if err != nil {
		fmt.Println("生成json字符串错误")
	}
	// case 2: //广播xxx玩家进入了房间

	joinPlayer := func(k, v interface{}) bool {

		v.(*conn.TConnection).WritePack(jsonStu)
		return true
	}
	connList.Range(joinPlayer)
}

//给所有玩家广播玩家列表
func proto3() {
	redat := &redata{}

	joinPlayer := func(k, v interface{}) bool {
		redat.Cmd = 3
		// redat.Mark = k.(int) //告诉客户端，我发的是给另外一个客户端的
		fmt.Println("xxxxxxxxxxxxxxxxxxxxx")
		redat.Content = "玩家列表"
		jsonStu, err := json.Marshal(redat)
		if err != nil {
			fmt.Println("生成json字符串错误")
		}
		v.(*conn.TConnection).WritePack(jsonStu)
		for index := 1; index <= nMark; index++ {
			log.Println(index, "indexindexindexindex")

			if v, ok := numberList.Load(index); ok {
				if v2, ok2 := nameList.Load(index); ok2 {
					num := strconv.Itoa(v.(int))
					// num := v.(string)
					name := v2.(string)
					redat.Content = name + "[" + num + "号]"
				}
			}

			jsonStu, err := json.Marshal(redat)
			if err != nil {
				fmt.Println("生成json字符串错误")
			}
			// fmt.Println(v, "你是哪个客户端")

			v.(*conn.TConnection).WritePack(jsonStu)
		}
		return true
	}
	connList.Range(joinPlayer)

	if nMark == 8 {
		randomPlayer() //随机生成玩家身份
		proto4()       //开始游戏
	}
}

//随机生成玩家身份
func randomPlayer() {
	// var list = [8]int{2, 2, 3, 3}
	var list = [8]int{1, 1, 1, 1, 2, 2, 3, 3}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	nindex := 1
	strIdentity := "平民"
	for _, i := range r.Perm(len(list)) {
		val := list[i]
		switch val {
		case 1:
			strIdentity = "平民"
			commonList = append(commonList, nindex)
			break
		case 2:
			strIdentity = "警察"
			policeList = append(policeList, nindex)
			break
		case 3:
			strIdentity = "杀手"
			killerList = append(killerList, nindex)
			break
		}
		identityList.Store(nindex, strIdentity)
		nindex = nindex + 1
	}

}

//给所有玩家广播 人满开始游戏
func proto4() {
	//天黑请闭眼
	//你的身份是xx
	redat := &redata{}
	//新开启一个线程来处理触发后的事件
	//开启定时器,记录20s的时间
	timer()

	joinPlayer := func(k, v interface{}) bool {
		redat.Cmd = 4
		redat.Content = "=======================游戏开始========================="
		jsonStu, err := json.Marshal(redat)
		if err != nil {
			fmt.Println("生成json字符串错误")
		}
		v.(*conn.TConnection).WritePack(jsonStu)

		for index := 1; index <= nMark; index++ {
			if k == index {
				if v, ok := identityList.Load(index); ok {
					redat.Content = "天黑请闭眼，今夜你的身份是:" + v.(string)
				}

				jsonStu, err := json.Marshal(redat)
				if err != nil {
					fmt.Println("生成json字符串错误")
				}

				v.(*conn.TConnection).WritePack(jsonStu)
			}
		}
		return true
	}
	connList.Range(joinPlayer)
	proto5() //告诉队友，并且让玩家指定验人杀人
}

//你的队友是xx号
func proto5() {
	//警察,你的队友是xx
	//杀手,你的队友是xx
	//平民,你等别人动 你不能动
	redat := &redata{}
	redat.Cmd = 5
	var mList []int

	//处理警察
	for index := 0; index < 3; index++ {

		if index == 0 { //警察组
			mList = policeList
		} else if index == 1 { //杀手组
			mList = killerList
		} else if index == 2 { //平民组
			mList = commonList
		}
		for i := 0; i < len(mList); i++ {

			strTeam := ""
			//找到一个队伍里除了自己之外的队友
			if index != 2 {
				for j := 0; j < len(mList); j++ {
					if j != i {
						if v, ok := nameList.Load(mList[j]); ok {
							name := v.(string)
							strTeam = strTeam + "[" + strconv.Itoa(mList[j]) + "]号 " + name + ";"
						}
					}
				}
			}
			if v, ok := connList.Load(mList[i]); ok {
				redat.Mark = mList[i] //告诉客户端他是几号
				redat.Content = "你的队友:" + strTeam + "          请输入你们要查验的目标号码:"
				if index == 1 {
					redat.Content = "你的队友:" + strTeam + "          请输入你们要击杀的目标号码:"
				} else if index == 2 {
					redat.Content = "请等待其他睁眼玩家操作"
				}
				jsonStu, err := json.Marshal(redat)
				if err != nil {
					fmt.Println("生成json字符串错误")
				}
				v.(*conn.TConnection).WritePack(jsonStu)
			}

		}
	}
}

//强类型语言,返回值也要确认类型
func calNameandNum(num int, name string) string {
	strOrig := ""
	if num != 0 {
		if v, ok := nameList.Load(num); ok {
			strOrig = v.(string) + "[" + strconv.Itoa(num) + "]号"
		}
	} else {
		// if v, ok := nameList.Load(num); ok {
		// 	strOrig = v.(string) + "[" + strconv.Itoa(num) + "]号"
		// }
	}
	return strOrig
}

//xx号要查xx号
//xx号要杀xx号
func proto6() {

	redat := &redata{}
	redat.Cmd = 6

	var mList []int
	ndepart := 0 //所属组织
	ntarget, err := strconv.Atoi(dat.Content)
	if err != nil {
		fmt.Println("错误的string转int")
	}

	//辨别发送消息的人属于哪个组织,两个组织之间信息不共享
	//同组织之间成员信息共享
	for j := 0; j < len(policeList); j++ {
		if dat.Mark == policeList[j] {
			ndepart = 1
			mList = policeList
		}
	}
	for j := 0; j < len(killerList); j++ {
		if dat.Mark == killerList[j] {
			ndepart = 2
			mList = killerList
		}
	}
	//遍历某个职业
	for i := 0; i < len(mList); i++ {
		if v, ok := connList.Load(mList[i]); ok {
			redat.Mark = mList[i] //告诉客户端他是几号
			//进来两次
			if ndepart == 1 {
				if npoliceget == 0 {
					redat.Content = calNameandNum(dat.Mark, "") + "想查验" + calNameandNum(ntarget, "")
					//把查验目标放入查验map
					if lenpoliceget == 0 {
						policetarget.Store(dat.Mark, ntarget)
						//要保证每次只加一次
						lenpoliceget = lenpoliceget + 1
					}
					//每次都要往里面插入数据,key相同的进行替换
					getLen := func(k, v1 interface{}) bool {
						policetarget.Store(dat.Mark, ntarget)
						//这里是map计数的变量,绝不重复计数
						if (k.(int) != dat.Mark) && (i == 0) && (lenpoliceget < len(mList)) {
							lenpoliceget = lenpoliceget + 1
						}
						return true
					}
					policetarget.Range(getLen)

					jsonStu, err := json.Marshal(redat)
					if err != nil {
						fmt.Println("生成json字符串错误")
					}
					//告诉每一个小警察,他或者他队友要验谁
					v.(*conn.TConnection).WritePack(jsonStu)
					//检查查验目标是否一致
					v2 := 0
					joinPlayer := func(k, v1 interface{}) bool {
						//如果警察的选择内容map长度和警察的人员长度一样,说明每一个警察都已经做出了选择
						if lenpoliceget == len(mList) {
							//如果两个警察的选择都一样,就广播给警察,告诉他们已经确认目标,并且通过proto7查出目标身份
							if v2 == v1.(int) {

								redat.Content = "目标一致,警方已确认目标,查验" + calNameandNum(v1.(int), "")
								jsonStu, err := json.Marshal(redat)
								v.(*conn.TConnection).WritePack(jsonStu)
								if err != nil {
									fmt.Println("生成json字符串错误")
								}
								//把职业,查验目标,警察连接,这是第几次循环传到proto7
								proto7(ndepart, v1.(int), v.(*conn.TConnection), i)
							}

						}
						v2 = v1.(int)
						return true
					}
					policetarget.Range(joinPlayer)
				}
			} else if ndepart == 2 {
				if nkillget == 0 {

					redat.Content = calNameandNum(dat.Mark, "") + "想击杀" + calNameandNum(ntarget, "")
					//把查验目标放入查验map
					if lenkillget == 0 {
						killtarget.Store(dat.Mark, ntarget)
						//要保证每次只加一次
						lenkillget = lenkillget + 1
					}
					//每次都要往里面插入数据,key相同的进行替换
					getLen := func(k, v1 interface{}) bool {
						killtarget.Store(dat.Mark, ntarget)
						//这里是map计数的变量,绝不重复计数
						if (k.(int) != dat.Mark) && (i == 0) && (lenkillget < len(mList)) {
							lenkillget = lenkillget + 1
						}
						return true
					}
					killtarget.Range(getLen)
					jsonStu, err := json.Marshal(redat)
					if err != nil {
						fmt.Println("生成json字符串错误")
					}
					//告诉每一个杀手,他或者他队友要验谁
					v.(*conn.TConnection).WritePack(jsonStu)
					//检查查验目标是否一致
					v2 := 0
					joinPlayer := func(k, v1 interface{}) bool {
						//如果杀手的选择内容map长度和杀手的人员长度一样,说明每一个杀手都已经做出了选择
						if lenkillget == len(mList) {
							//如果两个警察的选择都一样,就广播给警察,告诉他们已经确认目标,并且通过proto7查出目标身份
							//第一个和第二个目标相同
							if v2 == v1.(int) {
								redat.Content = "目标一致,杀手已确认目标,击杀" + calNameandNum(v1.(int), "")
								jsonStu, err := json.Marshal(redat)
								v.(*conn.TConnection).WritePack(jsonStu)
								if err != nil {
									fmt.Println("生成json字符串错误")
								}
								//把职业,查验目标,警察连接,这是第几次循环传到proto7
								proto7(ndepart, v1.(int), v.(*conn.TConnection), i)
							}

						}
						//把第一个value复制给变量v2 等下可以跟第二个value比较
						v2 = v1.(int)
						return true
					}
					killtarget.Range(joinPlayer)
				}
			}

		}

	}

}

//xx身份是xx,
//xx已经死亡,从玩家列表中剔除
func proto7(ndepart int, ntarget int, conn *conn.TConnection, nindex int) {
	redat := &redata{}
	redat.Cmd = 7
	if ndepart == 1 {
		findout := func(k, v interface{}) bool {
			if k == ntarget {

				redat.Content = calNameandNum(ntarget, "") + "的身份是 " + v.(string)
				if nindex == 1 {
					npoliceget = ntarget //确认本夜查人目标并保存,之后就无法更改了
				}
				jsonStu, err := json.Marshal(redat)
				if err != nil {
					fmt.Println("生成json字符串错误")
				}
				conn.WritePack(jsonStu)
			}
			return true
		}
		identityList.Range(findout)
	} else {
		findout := func(k, v interface{}) bool {
			if k == ntarget {
				//index=1,说明已经遍历第二遍了,遍历到第二个杀手了,那就可以进行操作了
				if nindex == 1 {
					identityList.Delete(k) //要在第二次进入的时候才进行删除
					nkillget = ntarget     //确认本夜杀人目标并保存,之后就无法更改了
				}
			}
			return true
		}
		identityList.Range(findout)
	}
}

//定时器触发,天黑20s后,天亮了,xx死了,请他发表遗言
func proto8() {
	redat := &redata{}
	redat.Cmd = 8
	redat.Content = "昨晚死的人是" + calNameandNum(nkillget, "") + " ,请" + calNameandNum(nkillget, "") + "发表遗言:"
	jsonStu, err := json.Marshal(redat)
	if err != nil {
		fmt.Println("生成json字符串错误")
	}
	joinPlayer := func(k, v interface{}) bool {
		v.(*conn.TConnection).WritePack(jsonStu)
		return true
	}
	connList.Range(joinPlayer)
	npoliceget = 0 //初始化,晚上又可以当做标记保存当晚查验人
}

//接受xx发表的遗言,广播给所有人
func proto9() {
	redat := &redata{}
	redat.Cmd = 9
	log.Println(dat, "dat.Contentdat.Contentdat.Contentdat.Contentdat.Contentdat.Contentdat.Content")
	redat.Content = calNameandNum(nkillget, "") + "遗言:" + dat.Content
	nkillget = 0 //初始化,晚上又可以当做标记保存当晚杀人

	jsonStu, err := json.Marshal(redat)
	if err != nil {
		fmt.Println("生成json字符串错误")
	}
	joinPlayer := func(k, v interface{}) bool {
		v.(*conn.TConnection).WritePack(jsonStu)
		return true
	}
	connList.Range(joinPlayer)
}

//死者左手边 顺时针开始发言
func proto10() {
	// redat := &redata{}
	// redat.Cmd = 9
	// redat.Content = calNameandNum(nkillget, "") + "遗言:" + dat.Content
	// jsonStu, err := json.Marshal(redat)
	// if err != nil {
	// 	fmt.Println("生成json字符串错误")
	// }
	// joinPlayer := func(k, v interface{}) bool {
	// 	v.(*conn.TConnection).WritePack(jsonStu)
	// 	return true
	// }
	// connList.Range(joinPlayer)
}
func main() {
	pServer := tcpserver.NewTCPServer()
	nMark = 0
	pServer.OnRead = func(pData *conn.TData) {
		buf := pData.GetBuffer()
		// nLen := pData.GetLength()
		redat := &redata{}
		json.Unmarshal(buf, dat)

		log.Println(dat.Cmd)

		switch dat.Cmd {
		case 1: //接受到新加入的用户名
			nMark = nMark + 1
			redat.Cmd = 1
			redat.Content = "你的用户名是：" + dat.Content + "     座位号是[ " + strconv.Itoa(nMark) + " ]号"
			jsonStu, err := json.Marshal(redat)
			fmt.Println("json字符串", string(jsonStu))
			if err != nil {
				fmt.Println("生成json字符串错误")
			}
			// inf := info{
			// 	Number: nMark,
			// 	Name:   dat.Content,
			// 	// Connect: pData,
			// 	// Connect: pData.GetConnection(),
			// }
			pData.GetConnection().WritePack(jsonStu)

			//除了刚加入的玩家，其他玩家都会受到新玩家加入的信息
			strPro2Cont := "玩家" + dat.Content + "[ " + strconv.Itoa(nMark) + " ]号 加入了游戏"
			proto2(strPro2Cont)

			numberList.Store(nMark, nMark)
			nameList.Store(nMark, dat.Content)
			connList.Store(nMark, pData.GetConnection())

			//给所有玩家广播玩家列表
			proto3()
			break
		case 6:
			proto6()
			break
		case 9:
			proto9()
			break
			// case 2: //广播xxx玩家进入了房间

			// 	joinPlayer := func(k, v interface{}) bool {
			// 		redat.Cmd = 2
			// 		if v, ok := numberList.Load(nMark); ok {
			// 			if v2, ok2 := nameList.Load(nMark); ok2 {
			// 				num := strconv.Itoa(v.(int))
			// 				// num := v.(string)
			// 				name := v2.(string)
			// 				redat.Content = name + "[" + num + "号]"
			// 			}
			// 		}

			// 		jsonStu, err := json.Marshal(redat)
			// 		if err != nil {
			// 			fmt.Println("生成json字符串错误")
			// 		}
			// 		v.(*conn.TConnection).WritePack(jsonStu)
			// 		return true
			// 	}
			// 	connList.Range(joinPlayer)

			// 	break
			// case 3: //广播现在的玩家列表

			// 	joinPlayer := func(k, v interface{}) bool {
			// 		redat.Cmd = 3
			// 		redat.Mark = k.(int) //告诉客户端，我发的是给另外一个客户端的
			// 		fmt.Println("xxxxxxxxxxxxxxxxxxxxx")

			// 		for index := 1; index <= nMark; index++ {
			// 			log.Println(index, "indexindexindexindex")

			// 			if v, ok := numberList.Load(index); ok {
			// 				if v2, ok2 := nameList.Load(index); ok2 {
			// 					num := strconv.Itoa(v.(int))
			// 					// num := v.(string)
			// 					name := v2.(string)
			// 					redat.Content = name + "[" + num + "号]"
			// 				}
			// 			}

			// 			jsonStu, err := json.Marshal(redat)
			// 			if err != nil {
			// 				fmt.Println("生成json字符串错误")
			// 			}
			// 			fmt.Println(v, "你是哪个客户端")

			// 			v.(*conn.TConnection).WritePack(jsonStu)
			// 		}
			// 		return true
			// 	}
			// 	connList.Range(joinPlayer)

			// 	break
		}
		// if dat.c
		// fmt.Print("     00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F")
		// for i := 0; i < nLen; i++ {
		// 	if i%16 == 0 {
		// 		fmt.Printf("\n%04d", i/16)
		// 	}
		// 	fmt.Printf(" %02x", buf[i])

		// }
		// fmt.Println("\n", string(buf)) //打印出来

		// ymj := func(k, v interface{}) bool {
		// 	v.(*conn.TConnection).WritePack(buf)
		// 	return true
		// }
		// m.Range(ymj)
	}

	pServer.OnClientConnect = func(pConn *conn.TConnection) {
		//
		// log.Println("有客户端连接 姚梦嘉", pConn.GetTCPConn().RemoteAddr())
		// m.Store(pConn.GetTCPConn().RemoteAddr(), pConn)

	}

	pServer.Start(":4567")

	time.Sleep(time.Hour * 10)
}
