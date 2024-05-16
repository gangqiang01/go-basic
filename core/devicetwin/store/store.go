package store

import (
	"github.com/edgehook/ithings/common/types/v1"
)

// //////////////////////////
// twin data.
// /////////////////////////
type TwinsData struct {
	Twins []*v1.TwinProperty
}

/*
* Data store interface for device hot data.
 */
type DBStore interface {
	Initialize() error
	StoreReportedTwins(deviceID string, twinProperties []*v1.TwinProperty) error
	LoadReportedTwins(deviceID string, count int) []*TwinsData
	UpdateDesiredTwins(deviceID string, desiredTwins []*v1.TwinProperty) error
	GetDesiredTwins(deviceID string) *TwinsData
}

// new database store.
func NewDBStore(backend string) DBStore {
	switch backend {
	case "influxdb":
		return NewMemoryStore()
	case "postgresql":
		return NewMemoryStore()
	}

	return NewMemoryStore()
}
