package sf_misc

import "errors"

var ErrNilPointer = errors.New("ERROR: nil pointer")

var ErrPacketLengthTooBig = errors.New("ERROR: packet length too bigger")
var ErrPacketLength = errors.New("ERROR: packet length err")
var ErrPackerHeaderParse = errors.New("ERROR: msg header parse err")

var ErrConnRead = errors.New("ERROR: read error")

var ErrConfigNotInit = errors.New("ERROR: Config not Read Config data")
var ErrConfigNotArr = errors.New("ERROR: current config is not a array")
var ErrConfigConvert = errors.New("ERROR: convert error")
var ErrConfigTypeNotMatch = errors.New("ERROR: config type is not match")

var ErrNetworkRecvConnCannotConn = errors.New("ERROR: Recv Conn cannot be connect")
