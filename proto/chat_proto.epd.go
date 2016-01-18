// easybuff
// 不要修改本文件, 每次消息有变动, 请手动生成本文件
// easybuff -s 描述文件目录 -o 目标文件目录 -l 语言(go,cpp)

package proto

import (
	. "github.com/toophy/toogo"
)

// ------ 枚举
const ()

// ------ 消息ID
const (
	C2G_login_Id          = 1  // 登录聊天服务器
	G2C_login_ret_Id      = 2  // 响应登录
	C2G_createRole_Id     = 3  // 创建角色
	G2C_createRole_ret_Id = 4  // 响应创建角色
	C2S_chat_Id           = 5  // 发送聊天信息
	C2C_chat_private_Id   = 6  // 发送私聊信息
	S2C_chat_Id           = 7  // 返回聊天信息
	S2C_chat_private_Id   = 8  // 返回私人聊天信息
	S2G_more_packet_Id    = 9  // Server发送给Gate消息包
	G2S_more_packet_Id    = 10 // Gate发送给GServer消息包
	S2G_registe_Id        = 12 // Server向Gate注册
	G2S_registe_Id        = 13 // Gate向Server返回注册结果
	S2C_monsterData_Id    = 14 // Server向Client返回怪物信息
)

// ------ 普通结构
// 位置数据
type TVec3 struct {
	X float32 // x坐标
	Y float32 // y坐标
	Z float32 // z坐标
}

func (t *TVec3) Read(p *PacketReader) bool {
	defer RecoverRead("TVec3")

	t.X = p.ReadFloat32()
	t.Y = p.ReadFloat32()
	t.Z = p.ReadFloat32()

	return true
}

func (t *TVec3) Write(p *PacketWriter) bool {
	defer RecoverWrite("TVec3")
	p.WriteFloat32(t.X)
	p.WriteFloat32(t.Y)
	p.WriteFloat32(t.Z)

	return true
}

// 怪物数据
type MonsterData struct {
	Name      string  // 怪物名
	CurrPos   TVec3   // 当前位置
	TargetPos TVec3   // 目标位置
	Speed     float32 // 速度
}

func (t *MonsterData) Read(p *PacketReader) bool {
	defer RecoverRead("MonsterData")

	t.Name = p.ReadString()
	t.CurrPos.Read(p)
	t.TargetPos.Read(p)
	t.Speed = p.ReadFloat32()

	return true
}

func (t *MonsterData) Write(p *PacketWriter) bool {
	defer RecoverWrite("MonsterData")
	p.WriteString(&t.Name)
	t.CurrPos.Write(p)
	t.TargetPos.Write(p)
	p.WriteFloat32(t.Speed)

	return true
}

// ------ 消息结构
// 登录聊天服务器
type C2G_login struct {
	Account string // 帐号
	Time    int32  // 登录时间戳
	Sign    string // 验证码
}

func (t *C2G_login) Read(p *PacketReader) bool {
	defer RecoverRead("C2G_login")

	t.Account = p.ReadString()
	t.Time = p.ReadInt32()
	t.Sign = p.ReadString()

	return true
}

func (t *C2G_login) Write(p *PacketWriter) bool {
	defer RecoverWrite("C2G_login")
	p.WriteMsgId(C2G_login_Id)
	p.WriteString(&t.Account)
	p.WriteInt32(t.Time)
	p.WriteString(&t.Sign)
	p.WriteMsgOver()

	return true
}

// 响应登录
type G2C_login_ret struct {
	Ret int8   // 登录结果,0:成功,其他为失败原因
	Msg string // 登录失败描述
}

func (t *G2C_login_ret) Read(p *PacketReader) bool {
	defer RecoverRead("G2C_login_ret")

	t.Ret = p.ReadInt8()
	t.Msg = p.ReadString()

	return true
}

func (t *G2C_login_ret) Write(p *PacketWriter) bool {
	defer RecoverWrite("G2C_login_ret")
	p.WriteMsgId(G2C_login_ret_Id)
	p.WriteInt8(t.Ret)
	p.WriteString(&t.Msg)
	p.WriteMsgOver()

	return true
}

// 创建角色
type C2G_createRole struct {
	Name string // 角色名
	Sex  int8   // 性别
}

func (t *C2G_createRole) Read(p *PacketReader) bool {
	defer RecoverRead("C2G_createRole")

	t.Name = p.ReadString()
	t.Sex = p.ReadInt8()

	return true
}

func (t *C2G_createRole) Write(p *PacketWriter) bool {
	defer RecoverWrite("C2G_createRole")
	p.WriteMsgId(C2G_createRole_Id)
	p.WriteString(&t.Name)
	p.WriteInt8(t.Sex)
	p.WriteMsgOver()

	return true
}

// 响应创建角色
type G2C_createRole_ret struct {
	Ret int8   // 创建角色结果,0:成功,其他为失败原因
	Msg string // 创建角色失败描述
}

func (t *G2C_createRole_ret) Read(p *PacketReader) bool {
	defer RecoverRead("G2C_createRole_ret")

	t.Ret = p.ReadInt8()
	t.Msg = p.ReadString()

	return true
}

func (t *G2C_createRole_ret) Write(p *PacketWriter) bool {
	defer RecoverWrite("G2C_createRole_ret")
	p.WriteMsgId(G2C_createRole_ret_Id)
	p.WriteInt8(t.Ret)
	p.WriteString(&t.Msg)
	p.WriteMsgOver()

	return true
}

// 发送聊天信息
type C2S_chat struct {
	Channel int32  // 频道
	Data    string // 聊天信息
}

func (t *C2S_chat) Read(p *PacketReader) bool {
	defer RecoverRead("C2S_chat")

	t.Channel = p.ReadInt32()
	t.Data = p.ReadString()

	return true
}

func (t *C2S_chat) Write(p *PacketWriter) bool {
	defer RecoverWrite("C2S_chat")
	p.WriteMsgId(C2S_chat_Id)
	p.WriteInt32(t.Channel)
	p.WriteString(&t.Data)
	p.WriteMsgOver()

	return true
}

// 发送私聊信息
type C2C_chat_private struct {
	Target string // 聊天目标
	Data   string // 聊天信息
}

func (t *C2C_chat_private) Read(p *PacketReader) bool {
	defer RecoverRead("C2C_chat_private")

	t.Target = p.ReadString()
	t.Data = p.ReadString()

	return true
}

func (t *C2C_chat_private) Write(p *PacketWriter) bool {
	defer RecoverWrite("C2C_chat_private")
	p.WriteMsgId(C2C_chat_private_Id)
	p.WriteString(&t.Target)
	p.WriteString(&t.Data)
	p.WriteMsgOver()

	return true
}

// 返回聊天信息
type S2C_chat struct {
	Channel int32  // 频道
	Source  string // 发言人
	Data    string // 聊天信息
}

func (t *S2C_chat) Read(p *PacketReader) bool {
	defer RecoverRead("S2C_chat")

	t.Channel = p.ReadInt32()
	t.Source = p.ReadString()
	t.Data = p.ReadString()

	return true
}

func (t *S2C_chat) Write(p *PacketWriter, tgid uint64) bool {
	defer RecoverWrite("S2C_chat")
	p.SetsubTgid(tgid)
	p.WriteMsgId(S2C_chat_Id)
	p.WriteInt32(t.Channel)
	p.WriteString(&t.Source)
	p.WriteString(&t.Data)
	p.WriteMsgOver()

	return true
}

// 返回私人聊天信息
type S2C_chat_private struct {
	Source string // 发言人
	Target string // 倾听者
	Data   string // 聊天信息
}

func (t *S2C_chat_private) Read(p *PacketReader) bool {
	defer RecoverRead("S2C_chat_private")

	t.Source = p.ReadString()
	t.Target = p.ReadString()
	t.Data = p.ReadString()

	return true
}

func (t *S2C_chat_private) Write(p *PacketWriter, tgid uint64) bool {
	defer RecoverWrite("S2C_chat_private")
	p.SetsubTgid(tgid)
	p.WriteMsgId(S2C_chat_private_Id)
	p.WriteString(&t.Source)
	p.WriteString(&t.Target)
	p.WriteString(&t.Data)
	p.WriteMsgOver()

	return true
}

// Server发送给Gate消息包
type S2G_more_packet struct {
}

func (t *S2G_more_packet) Read(p *PacketReader) bool {
	defer RecoverRead("S2G_more_packet")

	return true
}

func (t *S2G_more_packet) Write(p *PacketWriter, tgid uint64) bool {
	defer RecoverWrite("S2G_more_packet")
	p.SetsubTgid(tgid)
	p.WriteMsgId(S2G_more_packet_Id)
	p.WriteMsgOver()

	return true
}

// Gate发送给GServer消息包
type G2S_more_packet struct {
}

func (t *G2S_more_packet) Read(p *PacketReader) bool {
	defer RecoverRead("G2S_more_packet")

	return true
}

func (t *G2S_more_packet) Write(p *PacketWriter, tgid uint64) bool {
	defer RecoverWrite("G2S_more_packet")
	p.SetsubTgid(tgid)
	p.WriteMsgId(G2S_more_packet_Id)
	p.WriteMsgOver()

	return true
}

// Server向Gate注册
type S2G_registe struct {
	Sid uint64 // 小区编号
}

func (t *S2G_registe) Read(p *PacketReader) bool {
	defer RecoverRead("S2G_registe")

	t.Sid = p.ReadUint64()

	return true
}

func (t *S2G_registe) Write(p *PacketWriter, tgid uint64) bool {
	defer RecoverWrite("S2G_registe")
	p.SetsubTgid(tgid)
	p.WriteMsgId(S2G_registe_Id)
	p.WriteUint64(t.Sid)
	p.WriteMsgOver()

	return true
}

// Gate向Server返回注册结果
type G2S_registe struct {
	Ret uint8  // 返回结果,0:成功,其他失败
	Msg string // 返回失败原因
}

func (t *G2S_registe) Read(p *PacketReader) bool {
	defer RecoverRead("G2S_registe")

	t.Ret = p.ReadUint8()
	t.Msg = p.ReadString()

	return true
}

func (t *G2S_registe) Write(p *PacketWriter, tgid uint64) bool {
	defer RecoverWrite("G2S_registe")
	p.SetsubTgid(tgid)
	p.WriteMsgId(G2S_registe_Id)
	p.WriteUint8(t.Ret)
	p.WriteString(&t.Msg)
	p.WriteMsgOver()

	return true
}

// Server向Client返回怪物信息
type S2C_monsterData struct {
	Data MonsterData // 怪物数据
}

func (t *S2C_monsterData) Read(p *PacketReader) bool {
	defer RecoverRead("S2C_monsterData")

	t.Data.Read(p)

	return true
}

func (t *S2C_monsterData) Write(p *PacketWriter, tgid uint64) bool {
	defer RecoverWrite("S2C_monsterData")
	p.SetsubTgid(tgid)
	p.WriteMsgId(S2C_monsterData_Id)
	t.Data.Write(p)
	p.WriteMsgOver()

	return true
}
