package net

import (
	"sync/atomic"
)

type ConnId int32

var invalidConnId ConnId = 0

var g_conn_id uint32 = 1

func getNewId() uint32 {
	return atomic.AddUint32(&g_conn_id, 1)
}

func GenUniConnId() ConnId {
	var i uint32
	for {
		i = getNewId()
		if CheckConnIdInvalid(ConnId(i)) {
			break
		}
	}

	return ConnId(i)
}

func GetInvalidConnId() ConnId {
	return invalidConnId
}

func CheckConnIdInvalid(id ConnId) bool {
	return id == invalidConnId
}
