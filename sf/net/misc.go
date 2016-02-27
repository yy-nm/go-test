package net

import (
	"fmt"
	"test/sf/config"
)

type ServiceType string
type ServiceId int

func Get_service_type(conf config.Config) ServiceType {
	t, e := conf.String()
	if e != nil {
		fmt.Println(conf.Type())
		panic("service type must be string")
	}
	return ServiceType(t)
}

func Get_service_id(conf config.Config) ServiceId {
	i, e := conf.Float()
	if e != nil {
		panic("service id must be number")
	}

	return ServiceId(i)
}

func Config_get_string(conf config.Config) string {
	s, e := conf.String()
	if e != nil {
		panic("Config_get_string func must get string")
	}

	return s
}
