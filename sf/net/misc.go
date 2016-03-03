package net

import (
	"fmt"
	"test/sf/config"
)

type Service struct {
	Stype ServiceType
	Sid   ServiceId
}
type ServiceType int
type ServiceId int

const (
	_ ServiceType = iota
	SERVICE_LOGIN
	SERVICE_ACCOUNT
	SERVICE_DATACENTER
	SERVICE_GATE
	SERVICE_U_WATCHER
)

const Services map[string]ServiceType = map[string]ServiceType{
	"login":      SERVICE_LOGIN,
	"account":    SERVICE_ACCOUNT,
	"datacenter": SERVICE_DATACENTER,
	"gate":       SERVICE_GATE,
	"u_watcher":  SERVICE_U_WATCHER,
}

func GetServiceType(conf config.Config) ServiceType {
	t, e := conf.String()
	if e != nil {
		fmt.Println(conf.Type())
		panic("service type must be string")
	}
	st, ok := Services[t]
	if !ok {
		panic("service type not find in Services")
	}

	return st
}

func GetServiceId(conf config.Config) ServiceId {
	i, e := conf.Float()
	if e != nil {
		panic("service id must be number")
	}

	return ServiceId(i)
}

func ConfigGetString(conf config.Config) string {
	s, e := conf.String()
	if e != nil {
		panic("Config_get_string func must get string")
	}

	return s
}
