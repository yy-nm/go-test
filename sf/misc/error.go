package misc

import "errors"

var ErrNilPointer = errors.New("ERROR: nil pointer")

var ErrPacketBodyLenTooBig = errors.New("ERROR: packet length too bigger")
var ErrPacketTailLength = errors.New("ERROR: packet length err")

var ErrConnRead = errors.New("ERROR: read error")
var ErrConnConnect = errors.New("ERROR: conn error")

var ErrConfigNotInit = errors.New("ERROR: Config not Read Config data")
var ErrConfigNotArr = errors.New("ERROR: current config is not a array")
var ErrConfigConvert = errors.New("ERROR: convert error")
var ErrConfigTypeNotMatch = errors.New("ERROR: config type is not match")

var ErrNetworkRecvConnCannotConn = errors.New("ERROR: Recv Conn cannot be connect")

var ErrCanNotConn = errors.New("ERROR: can not connect")

var ErrConManagerStop = errors.New("ERROR: ConnManager is Stop")
