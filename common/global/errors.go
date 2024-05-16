package global

import (
	"errors"
)

var (
	ErrInvalidParms          = errors.New("Invalid parms")
	ErrChannelClosed         = errors.New("channel has been closed!")
	ErrInvalidResponseStruct = errors.New("invalid response structure.")
	ErrNoSuchDevice          = errors.New("No Such Device.")
	ErrNoSuchDeviceModel     = errors.New("No Such Device Model")
	ErrCoreNotReady          = errors.New("ICore is not ready")
	ErrEdgeIsNotOnline       = errors.New("edge is not online")
	ErrDeviceIsNotActive     = errors.New("device is not active")
	ErrDeviceIsOffline       = errors.New("device is offline")
	ErrUnknown               = errors.New("unknown error.")
)
