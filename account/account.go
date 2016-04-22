package account

import ()

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
