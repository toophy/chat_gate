// easybuff
// 不要修改本文件, 每次消息有变动, 请手动生成本文件
// easybuff -s 描述文件目录 -o 目标文件目录 -l 语言(go,cpp)

package proto

import (
	. "github.com/toophy/toogo"
)

// 发送私聊信息
const C2C_chat_private_Id = 6

type C2C_chat_private struct {
	Target string // 聊天目标
	Data   string // 聊天信息
}

func (t *C2C_chat_private) Read(p *PacketReader) bool {
	defer RecoverRead(C2C_chat_private_Id)

	t.Target = p.ReadString()
	t.Data = p.ReadString()

	return true
}

func (t *C2C_chat_private) Write(p *PacketWriter) bool {
	defer RecoverWrite(C2C_chat_private_Id)

	p.WriteMsgId(C2C_chat_private_Id)
	p.WriteString(&t.Target)
	p.WriteString(&t.Data)
	p.WriteMsgOver()

	return true
}

// 创建角色
const C2G_createRole_Id = 3

type C2G_createRole struct {
	Name string // 角色名
	Sex  int8   // 性别
}

func (t *C2G_createRole) Read(p *PacketReader) bool {
	defer RecoverRead(C2G_createRole_Id)

	t.Name = p.ReadString()
	t.Sex = p.ReadInt8()

	return true
}

func (t *C2G_createRole) Write(p *PacketWriter) bool {
	defer RecoverWrite(C2G_createRole_Id)

	p.WriteMsgId(C2G_createRole_Id)
	p.WriteString(&t.Name)
	p.WriteInt8(t.Sex)
	p.WriteMsgOver()

	return true
}

// 登录聊天服务器
const C2G_login_Id = 1

type C2G_login struct {
	Account string // 帐号
	Time    int32  // 登录时间戳
	Sign    string // 验证码
}

func (t *C2G_login) Read(p *PacketReader) bool {
	defer RecoverRead(C2G_login_Id)

	t.Account = p.ReadString()
	t.Time = p.ReadInt32()
	t.Sign = p.ReadString()

	return true
}

func (t *C2G_login) Write(p *PacketWriter) bool {
	defer RecoverWrite(C2G_login_Id)

	p.WriteMsgId(C2G_login_Id)
	p.WriteString(&t.Account)
	p.WriteInt32(t.Time)
	p.WriteString(&t.Sign)
	p.WriteMsgOver()

	return true
}

// 发送聊天信息
const C2S_chat_Id = 5

type C2S_chat struct {
	Channel int32  // 频道
	Data    string // 聊天信息
}

func (t *C2S_chat) Read(p *PacketReader) bool {
	defer RecoverRead(C2S_chat_Id)

	t.Channel = p.ReadInt32()
	t.Data = p.ReadString()

	return true
}

func (t *C2S_chat) Write(p *PacketWriter) bool {
	defer RecoverWrite(C2S_chat_Id)

	p.WriteMsgId(C2S_chat_Id)
	p.WriteInt32(t.Channel)
	p.WriteString(&t.Data)
	p.WriteMsgOver()

	return true
}

// 响应创建角色
const G2C_createRole_ret_Id = 4

type G2C_createRole_ret struct {
	Ret int8   // 创建角色结果,
	Msg string // 创建角色失败描述
}

func (t *G2C_createRole_ret) Read(p *PacketReader) bool {
	defer RecoverRead(G2C_createRole_ret_Id)

	t.Ret = p.ReadInt8()
	t.Msg = p.ReadString()

	return true
}

func (t *G2C_createRole_ret) Write(p *PacketWriter) bool {
	defer RecoverWrite(G2C_createRole_ret_Id)

	p.WriteMsgId(G2C_createRole_ret_Id)
	p.WriteInt8(t.Ret)
	p.WriteString(&t.Msg)
	p.WriteMsgOver()

	return true
}

// 响应登录
const G2C_login_ret_Id = 2

type G2C_login_ret struct {
	Ret int8   // 登录结果,0:成功,其他为失败原因
	Msg string // 登录失败描述
}

func (t *G2C_login_ret) Read(p *PacketReader) bool {
	defer RecoverRead(G2C_login_ret_Id)

	t.Ret = p.ReadInt8()
	t.Msg = p.ReadString()

	return true
}

func (t *G2C_login_ret) Write(p *PacketWriter) bool {
	defer RecoverWrite(G2C_login_ret_Id)

	p.WriteMsgId(G2C_login_ret_Id)
	p.WriteInt8(t.Ret)
	p.WriteString(&t.Msg)
	p.WriteMsgOver()

	return true
}

// 返回聊天信息
const S2C_chat_Id = 7

type S2C_chat struct {
	Channel int32  // 频道
	Source  string // 发言人
	Data    string // 聊天信息
}

func (t *S2C_chat) Read(p *PacketReader) bool {
	defer RecoverRead(S2C_chat_Id)

	t.Channel = p.ReadInt32()
	t.Source = p.ReadString()
	t.Data = p.ReadString()

	return true
}

func (t *S2C_chat) Write(p *PacketWriter) bool {
	defer RecoverWrite(S2C_chat_Id)

	p.WriteMsgId(S2C_chat_Id)
	p.WriteInt32(t.Channel)
	p.WriteString(&t.Source)
	p.WriteString(&t.Data)
	p.WriteMsgOver()

	return true
}

// 返回私人聊天信息
const S2C_chat_private_Id = 8

type S2C_chat_private struct {
	Source string // 发言人
	Target string // 倾听者
	Data   string // 聊天信息
}

func (t *S2C_chat_private) Read(p *PacketReader) bool {
	defer RecoverRead(S2C_chat_private_Id)

	t.Source = p.ReadString()
	t.Target = p.ReadString()
	t.Data = p.ReadString()

	return true
}

func (t *S2C_chat_private) Write(p *PacketWriter) bool {
	defer RecoverWrite(S2C_chat_private_Id)

	p.WriteMsgId(S2C_chat_private_Id)
	p.WriteString(&t.Source)
	p.WriteString(&t.Target)
	p.WriteString(&t.Data)
	p.WriteMsgOver()

	return true
}
