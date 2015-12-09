package main

import (
	"fmt"
	"github.com/toophy/chat_gate/proto"
	"github.com/toophy/toogo"
)

// 主线程
type MasterThread struct {
	toogo.Thread
}

// 首次运行
func (this *MasterThread) On_firstRun() {
}

// 响应线程最先运行
func (this *MasterThread) On_preRun() {
	// 处理各种最先处理的问题
}

// 响应线程运行
func (this *MasterThread) On_run() {
}

// 响应线程退出
func (this *MasterThread) On_end() {
}

// 响应网络事件
func (this *MasterThread) On_netEvent(m *toogo.Tmsg_net) bool {

	name_fix := m.Name
	if len(name_fix) == 0 {
		name_fix = fmt.Sprintf("Conn[%d]", m.SessionId)
	}

	switch m.Msg {
	case "listen failed":
		this.LogFatal("%s : Listen failed[%s]", name_fix, m.Info)

	case "listen ok":
		this.LogInfo("%s : Listen(%s) ok.", name_fix, toogo.GetSessionById(m.SessionId).GetIPAddress())

	case "accept failed":
		this.LogFatal(m.Info)
		return false

	case "accept ok":
		this.LogDebug("%s : Accept ok", name_fix)

	case "connect failed":
		this.LogError("%s : Connect failed[%s]", name_fix, m.Info)

	case "connect ok":
		this.LogDebug("%s : Connect ok", name_fix)

	case "read failed":
		this.LogError("%s : Connect read[%s]", name_fix, m.Info)

	case "pre close":
		this.LogDebug("%s : Connect pre close", name_fix)

	case "close failed":
		this.LogError("%s : Connect close failed[%s]", name_fix, m.Info)

	case "close ok":
		this.LogDebug("%s : Connect close ok.", name_fix)
	}

	return true
}

// -- 当网络消息包解析出现问题, 如何处理?
func (this *MasterThread) On_packetError(sessionId uint32) {
	toogo.CloseSession(this.Get_thread_id(), sessionId)
}

// 注册消息
func (this *MasterThread) On_registNetMsg() {
	this.RegistNetMsg(proto.C2G_login_Id, this.on_c2g_login)
	this.RegistNetMsg(proto.C2S_chat_Id, this.on_c2s_chat)
}

func (this *MasterThread) on_c2g_login(pack *toogo.PacketReader, sessionId uint32) bool {
	msg := proto.C2G_login{}
	msg.Read(pack)

	p := new(toogo.PacketWriter)
	d := make([]byte, 64)
	p.InitWriter(d)

	this.LogInfo("on_c2g_login")

	msgLoginRet := new(proto.G2C_login_ret)
	msgLoginRet.Ret = 0
	msgLoginRet.Msg = "ok"
	msgLoginRet.Write(p)

	p.PacketWriteOver()
	session := toogo.GetSessionById(sessionId)

	m := new(toogo.Tmsg_packet)
	m.Data = p.GetData()
	m.Len = uint32(p.GetPos())
	m.Count = uint16(p.Count)

	toogo.PostThreadMsg(session.MailId, m)
	return true
}

func (this *MasterThread) on_c2s_chat(pack *toogo.PacketReader, sessionId uint32) bool {
	msg := proto.C2S_chat{}
	msg.Read(pack)

	this.LogInfo("Say : %s", msg.Data)
	// 封包一层

	return true
}

func main() {
	main_thread := new(MasterThread)
	main_thread.Init_thread(main_thread, toogo.Tid_master, "master", 1000, 100, 10000)
	toogo.Run(main_thread)
}

// 封包处理(大批量,或者一个网络口(代理))
// 简化发送消息的代码
// Lua支持消息包
// Lua支持线程间消息
//
// Session上面做标记, 大包, 小包 等
// 每种包会有不同的解包机制, 包头大小
// 如何在accecpt时就知道这个包属于哪种包头?
// 是一个什么消息么?
// 比如, 一开始默认都是小包
// 验证成功后, 根据对方资质, 变成大包
// 或者
// 有一个Listen侦听到的都是服务器连接, 都是大包头?
// 其实只有gate服连接才是大包头, 其他都是小包头
//
// 服务器之间都是大包头, 也即是包长度4字节, 包消息数量2字节,
// 只有gate和客户端之间才是小包头, 包长度2字节, 包消息数量1字节,
//
// 只有服务器和gate服的连接才会包中有包, 其他都是简单的消息处理
//
// session需要分清楚
// 1. 这是什么连接(服务器,客户端)
// 2. 服务器的连接中, 是否有gate服
//    a. CG连接 小包
//    b. SS连接 大包
//    c. GS连接 混合大包
//
