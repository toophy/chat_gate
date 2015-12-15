package main

import (
	"fmt"
	"github.com/toophy/chat_gate/proto"
	"github.com/toophy/toogo"
	"sync"
)

// 虚拟帐号
type PlatAccount struct {
	Id   uint64 // 帐号唯一Id
	Name string // 帐号名(唯一)
	Pkey string // 帐号来源平台
}

// 角色
type PlatRole struct {
	Id        uint64 // 角色唯一Id
	Name      string // 可以为空,表示没有创建角色, 但是登录过
	Skey      uint64 // 创建角色所在的小区
	IpAddress string // 最近登录的IP地址
}

// 主线程
type MasterThread struct {
	toogo.Thread

	ChatSvrsLock sync.RWMutex
	ChatSvrs     map[uint64]*toogo.Session

	LastAccountId    uint64
	PlatAccountsId   map[uint64]*PlatAccount
	PlatAccountsName map[string]*PlatAccount
	PlatRoleId       map[uint64]*PlatRole
}

func GetPkeyId(pkey string) uint64 {
	switch pkey {
	case "tyh":
		return 1
	case "6998":
		return 2
	case "qqzone":
		return 3
	}
	return 0
}

func MakeAid(aid uint64, pid uint64) uint64 {
	return aid<<10 | pid
}

func MakeRid(aid uint64, pid uint64, sid uint64) uint64 {
	return sid<<42 | aid<<10 | pid
}

func MakeRidEx(apid uint64, sid uint64) uint64 {
	return sid<<42 | apid
}

func (this *MasterThread) AddAccount(name string, pkey string, skey uint64) bool {
	if len(name) == 0 || len(pkey) == 0 || len(skey) == 0 {
		return false
	}

	fix_acc_name := pkey + "_" + name
	if v, ok := this.PlatAccountsName[fix_acc_name]; !ok {
		pa := new(PlatAccount)
		pa.Id = MakeAid(this.LastAccountId, GetPkeyId(pkey))
		this.LastAccountId++

		pa.Name = fix_acc_name
		pa.Pkey = pkey
		this.PlatAccountsName[fix_acc_name] = pa
		this.PlatAccountsId[pa.Id] = pa

		pr := new(PlatRole)
		pr.Id = MakeRidEx(pa.Id, skey)
		pr.Skey = skey
		this.PlatRoleId[pr.Id] = pr
		return true
	} else {

		rid := MakeRidEx(v.Id, skey)
		if vr, okr := this.PlatRoleId[rid]; !okr {
			pr := new(PlatRole)
			pr.Id = rid
			pr.Skey = skey
			this.PlatRoleId[pr.Id] = pr
		}

		return true
	}

	return false
}

func (this *MasterThread) GetRole(name string, pkey string, skey uint64) *PlatRole {
	if len(name) == 0 || len(pkey) == 0 || len(skey) == 0 {
		return nil
	}

	fix_acc_name := pkey + "_" + name
	if v, ok := this.PlatAccountsName[fix_acc_name]; ok {

		rid := MakeRidEx(v.Id, skey)
		if v, ok := this.PlatRoleId[rid]; ok {
			return v
		}
	}

	return nil
}

// 首次运行
func (this *MasterThread) On_firstRun() {
	this.ChatSvrs = make(map[uint64]*toogo.Session, 100)
	this.PlatAccountsId = make(map[uint64]*PlatAccount, 10)
	this.PlatAccountsName = make(map[string]*PlatAccount, 10)
	this.PlatRoleId = make(map[uint64]*PlatRole, 10)
	this.LastAccountId = 1

	this.AddAccount("koko", "tyh", "1")
	this.AddAccount("bububu", "qqzone", "1")
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
	this.RegistNetMsg(proto.S2G_registe_Id, this.on_s2g_registe)
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

	// 封包一层
	targetSession := this.GetSession(1)
	if targetSession != nil {

		pM := toogo.NewPacket(256, targetSession.SessionId)
		if pM != nil {
			// defer RecoverWrite(G2S_more_packet_Id)
			pM.WriteMsgId(proto.G2S_more_packet_Id)

			pC := toogo.NewPacket(128, targetSession.SessionId)
			if pC != nil {
				_, msg_len, start_pos := pack.GetReadMsg()

				fmt.Println(pack.GetData()[start_pos : start_pos+uint64(msg_len)])
				pC.CopyMsg(pack.GetData()[start_pos:start_pos+uint64(msg_len)], uint64(msg_len))
				pC.PacketWriteOver()

				pM.WriteUint16(1)
				pM.WriteDataEx(pC.GetData(), pC.GetPos())
				pM.WriteMsgOver()

				toogo.SendPacket(pM)
			}
		}
	}

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

func (this *MasterThread) on_s2g_registe(pack *toogo.PacketReader, sessionId uint64) bool {
	msg := proto.S2G_registe{}
	msg.Read(pack)

	this.ChatSvrsLock.Lock()
	defer this.ChatSvrsLock.Unlock()

	this.ChatSvrs[msg.Sid] = toogo.GetSessionById(sessionId)

	return true
}

func (this *MasterThread) GetSession(sid uint64) *toogo.Session {
	this.ChatSvrsLock.RLock()
	defer this.ChatSvrsLock.RUnlock()

	if v, ok := this.ChatSvrs[sid]; ok {
		return v
	}

	return nil
}

func main() {
	main_thread := new(MasterThread)
	main_thread.Init_thread(main_thread, toogo.Tid_master, "master", 1000, 100, 10000)
	toogo.Run(main_thread)
}
