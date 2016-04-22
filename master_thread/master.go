package master_thread

import (
	"fmt"
	. "github.com/toophy/chat_gate/account"
	"github.com/toophy/chat_gate/proto"
	"github.com/toophy/toogo"
	"sync"
)

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

// 首次运行
func (this *MasterThread) On_firstRun() {
	this.ChatSvrs = make(map[uint64]*toogo.Session, 100)
	this.PlatAccountsId = make(map[uint64]*PlatAccount, 10)
	this.PlatAccountsName = make(map[string]*PlatAccount, 10)
	this.PlatRoleId = make(map[uint64]*PlatRole, 10)
	this.LastAccountId = 1

	this.AddAccount("koko", "tyh", 1)
	this.AddAccount("toophy", "tyh", 1)
	this.RegistNetMsgDefault(this.On_defaultNetMsg)
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
	this.LogInfo("网络包错误,关闭会话连接 id=%d", sessionId)
	toogo.CloseSession(this.Get_thread_id(), sessionId)
}

// 注册消息
func (this *MasterThread) On_registNetMsg() {
	this.RegistNetMsg(proto.C2G_login_Id, this.on_c2g_login)
	this.RegistNetMsg(proto.S2G_registe_Id, this.on_s2g_registe)
}

func (this *MasterThread) On_defaultNetMsg(msg_id uint16, pack *toogo.PacketReader, sessionId uint64) bool {
	// 这里决定怎么转发

	switch msg_id {
	case proto.C2S_chat_Id:
		msg := proto.C2S_chat{}
		msg.Read(pack)

		svrSession := this.GetSession(toogo.Tgid_make_Sid(1, 1))
		if svrSession != nil {
			p := toogo.NewPacket(64, svrSession.SessionId)

			if p != nil {
				p.SetsubTgid(pack.LinkTgid)
				msg.Write(p)
				toogo.SendPacket(p)
			}
		}
	}

	return true
}

func (this *MasterThread) on_c2g_login(pack *toogo.PacketReader, sessionId uint64) bool {
	msg := proto.C2G_login{}
	msg.Read(pack)

	p := toogo.NewPacket(64, sessionId)

	toogo.SetSessionTgid(sessionId, toogo.Tgid_make_Rid(1, 1, 1))
	toogo.SetTgidSession(toogo.Tgid_make_Rid(1, 1, 1), sessionId)

	if p != nil {
		msgLoginRet := new(proto.G2C_login_ret)
		msgLoginRet.Ret = 0
		msgLoginRet.Msg = "ok"
		msgLoginRet.Write(p)

		toogo.SendPacket(p)
	}

	return true
}

func (this *MasterThread) on_s2g_registe(pack *toogo.PacketReader, sessionId uint64) bool {
	msg := proto.S2G_registe{}
	msg.Read(pack)

	this.LogInfo("s2g_registe %d", msg.Sid)

	toogo.SetSessionTgid(sessionId, msg.Sid)
	toogo.SetTgidSession(msg.Sid, sessionId)

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

func (this *MasterThread) AddAccount(name string, pkey string, skey uint64) bool {
	if len(name) == 0 || len(pkey) == 0 || skey == 0 {
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
		if _, okr := this.PlatRoleId[rid]; !okr {
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
	if len(name) == 0 || len(pkey) == 0 || skey == 0 {
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
