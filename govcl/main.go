package main

import (
	"encoding/json"
	"fmt"
	"log"

	conn "github.com/266game/goserver/Connection"
	tcpclient "github.com/266game/goserver/TCPClient"
	"github.com/ying32/govcl/vcl"
)

type TMainForm struct {
	*vcl.TForm
	Btn1    *vcl.TButton
	Edit1   *vcl.TEdit
	Memo1   *vcl.TMemo
	pClient *tcpclient.TTCPClient
}
type data struct {
	Cmd     int32  `json:"cmd"`
	Content string `json:"content"`
	Mark    int32  `json:"mark"`
}
type redata struct {
	Cmd     int32  `json:"cmd"`
	Content string `json:"content"`
	Mark    int32  `json:"mark"`
}

var (
	mainForm *TMainForm
)

func main() {
	vcl.Application.Initialize()
	vcl.Application.SetMainFormOnTaskBar(true)
	vcl.Application.CreateForm(&mainForm)
	vcl.Application.Run()

}

var dat = &data{}
var redat = &redata{}

// -- TMainForm

func (self *TMainForm) OnFormCreate(sender vcl.IObject) {
	self.SetCaption("梦嘉杀")
	self.Btn1 = vcl.NewButton(self)
	self.Btn1.SetParent(self)
	self.Btn1.SetBounds(10, 10, 88, 28)
	self.Btn1.SetCaption("发送(&S)")
	self.Btn1.SetOnClick(self.OnButtonClick)

	self.Edit1 = vcl.NewEdit(self)
	self.Edit1.SetParent(self)
	self.Edit1.SetBounds(10, 50, 500, 30)

	self.Memo1 = vcl.NewMemo(self)
	self.Memo1.SetParent(self)
	self.Memo1.SetBounds(10, 100, 500, 200)

	self.pClient = tcpclient.NewTCPClient()
	self.Memo1.Lines().Add(string("欢迎来到瞎比比狼人杀，请输入用户名："))

	// isone := 0
	self.pClient.OnRead = func(pData *conn.TData) {

		buf := pData.GetBuffer()
		// nLen := pData.GetLength()
		// log.Println("收到包了长度是", nLen, "\n", string(buf), "\n", buf)
		json.Unmarshal(buf, redat)
		log.Println("姚梦嘉", redat)
		switch redat.Cmd {
		case 1:
			self.Memo1.Lines().Add(string(redat.Content))
			// self.SendMsg(2, "")
			break
		case 2:
			self.Memo1.Lines().Add(string(redat.Content))
			// self.SendMsg(3, "")
			// self.Memo1.Lines().Add(string("玩家列表"))

		case 3:
			self.Memo1.Lines().Add(string(redat.Content))
		case 4:
			self.Memo1.Lines().Add(string(redat.Content))
		case 5:
			self.Memo1.Lines().Add(string(redat.Content))
			redat.Cmd = 6
		case 6:
			self.Memo1.Lines().Add(string(redat.Content))
		case 7:
			self.Memo1.Lines().Add(string(redat.Content))
		case 8:
			self.Memo1.Lines().Add(string(redat.Content))
			redat.Cmd = 9
		case 9:
			self.Memo1.Lines().Add(string(redat.Content))
		default:
			self.Memo1.Lines().Add("意外的协议")
		}
	}

	self.pClient.Connect("127.0.0.1:4567")
}
func (self *TMainForm) SendMsg(cmd int32, content string, nmark int32) {
	dat.Cmd = cmd
	dat.Content = content
	dat.Mark = nmark
	log.Println(dat, "准备发出数据")
	jsonStu, err := json.Marshal(dat)
	if err != nil {
		fmt.Println("生成json字符串错误")
	}
	self.pClient.WritePack([]byte(jsonStu))

}

func (self *TMainForm) OnButtonClick(sender vcl.IObject) {
	// 我叫姚梦嘉
	if redat.Cmd == 0 {
		redat.Cmd = 1
	}
	log.Println(redat.Mark, "redat.Markredat.Markredat.Markredat.Mark")
	self.SendMsg(redat.Cmd, self.Edit1.Text(), redat.Mark)
}
