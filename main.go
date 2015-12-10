package main

import (
	"fmt"
	"github.com/toophy/chat_gate/proto"
	"github.com/toophy/toogo"
)

// ScreenSvr
type ScreenSvr struct {
	Skey uint16
}

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
func (this *MasterThread) On_packetError(sessionId uint64) {
	toogo.CloseSession(this.Get_thread_id(), sessionId)
}

// 注册消息
func (this *MasterThread) On_registNetMsg() {
	this.RegistNetMsg(proto.C2G_login_Id, this.on_c2g_login)
	this.RegistNetMsg(proto.C2S_chat_Id, this.on_c2s_chat)
	this.RegistNetMsg(proto.S2G_more_packet_Id, this.on_s2g_more_packet)
}

func (this *MasterThread) on_c2g_login(pack *toogo.PacketReader, sessionId uint64) bool {
	msg := proto.C2G_login{}
	msg.Read(pack)

	p := toogo.NewPacket(64, sessionId)

	if p != nil {
		msgLoginRet := new(proto.G2C_login_ret)
		msgLoginRet.Ret = 0
		msgLoginRet.Msg = "ok"
		msgLoginRet.Write(p)

		toogo.SendPacket(p)
	}

	return true
}

func (this *MasterThread) on_c2s_chat(pack *toogo.PacketReader, sessionId uint64) bool {
	msg := proto.C2S_chat{}
	msg.Read(pack)

	this.LogInfo("Say : %s", msg.Data)
	// 封包一层

	// pM := toogo.NewPacket(256, m.SessionId)
	// if pM != nil {
	// 	// defer RecoverWrite(S2G_more_packet_Id)
	// 	pM.WriteMsgId(proto.S2G_more_packet_Id)

	// 	pC := toogo.NewPacket(128, m.SessionId)
	// 	if pC != nil {
	// 		msgChat := new(proto.C2S_chat)
	// 		msgChat.Channel = 1
	// 		msgChat.Data = "我是OKOK"
	// 		msgChat.Write(pC)
	// 		pC.PacketWriteOver()

	// 		pM.WriteUint16(1)
	// 		pM.WriteDataEx(pC.GetData(), pC.GetPos())
	// 		pM.WriteMsgOver()

	// 		toogo.SendPacket(pM)
	// 	}
	// }

	return true
}

func (this *MasterThread) on_s2g_more_packet(pack *toogo.PacketReader, sessionId uint64) bool {
	defer toogo.RecoverRead(proto.S2G_more_packet_Id)

	// 整包, 多少个消息? 还是一个消息
	// 消息长度, 去掉消息头, 消息总长度
	subPackCount := pack.ReadUint16()
	for i := uint16(0); i < subPackCount; i++ {
		if !this.ProcSubNetPacket(pack, sessionId, proto.S2G_more_packet_Id) {
			return false
		}
	}

	return true
}

func main() {
	main_thread := new(MasterThread)
	main_thread.Init_thread(main_thread, toogo.Tid_master, "master", 1000, 100, 10000)
	toogo.Run(main_thread)
}
