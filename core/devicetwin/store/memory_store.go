package store

import (
	"sync"

	"github.com/edgehook/ithings/common/types/v1"
)

var (
	defaultMaxCachedTwins = 20
)

/*
* store device twin into simple memory cache.
 */
type MemoryStore struct {
	ReportedTwinsCache *sync.Map
	DesiredTwinsCache  *sync.Map
}

func NewMemoryStore() *MemoryStore {
	var reportedMap sync.Map
	var desiredMap sync.Map
	return &MemoryStore{
		ReportedTwinsCache: &reportedMap,
		DesiredTwinsCache:  &desiredMap,
	}
}

func (m *MemoryStore) Initialize() error {
	return nil
}

func (m *MemoryStore) StoreReportedTwins(deviceID string, twinProperties []*v1.TwinProperty) error {
	var cacheMap *sync.Map

	cacheMap = m.ReportedTwinsCache
	v, exist := cacheMap.Load(deviceID)
	if !exist {
		v = NewTwinsCache()
	}

	tc, isThisType := v.(*TwinsCache)
	if !isThisType {
		tc = NewTwinsCache()
	}

	data := &TwinsData{
		Twins: twinProperties,
	}

	tc.Put(data)

	cacheMap.Store(deviceID, tc)

	return nil
}

func (m *MemoryStore) LoadReportedTwins(deviceID string, count int) []*TwinsData {
	var cacheMap *sync.Map

	cacheMap = m.ReportedTwinsCache
	v, exist := cacheMap.Load(deviceID)
	if !exist {
		return nil
	}
	tc, _ := v.(*TwinsCache)
	if tc == nil {
		return nil
	}

	data := tc.Get(count)
	if data == nil {
		return nil
	}
	return data
}

func (m *MemoryStore) UpdateDesiredTwins(deviceID string, desiredTwins []*v1.TwinProperty) error {
	var cacheMap *sync.Map

	cacheMap = m.DesiredTwinsCache
	v, exist := cacheMap.Load(deviceID)
	if !exist {
		v = NewTwinsCache()
	}

	tc, isThisType := v.(*TwinsCache)
	if !isThisType {
		tc = NewTwinsCache()
	}

	//update the head twins since the desired twin
	// just has on cache(TwinsCache.Size is always 1) .
	tc.UpdateHeadTwins(desiredTwins)

	cacheMap.Store(deviceID, tc)
	return nil
}

func (m *MemoryStore) GetDesiredTwins(deviceID string) *TwinsData {
	var cacheMap *sync.Map

	cacheMap = m.DesiredTwinsCache
	v, exist := cacheMap.Load(deviceID)
	if !exist {
		return nil
	}

	tc, _ := v.(*TwinsCache)
	if tc == nil {
		return nil
	}

	data := tc.Get(1)
	if data == nil || len(data) == 0 {
		return nil
	}

	return data[0]
}

type TwinsCache struct {
	Cache []*TwinsData
	Head  int
	Size  int
	mutex *sync.Mutex
}

func NewTwinsCache() *TwinsCache {
	var mutex sync.Mutex

	return &TwinsCache{
		Cache: make([]*TwinsData, defaultMaxCachedTwins),
		Size:  0,
		Head:  0,
		mutex: &mutex,
	}
}

func (t *TwinsCache) Put(data *TwinsData) {
	if data == nil {
		return
	}
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.Size < defaultMaxCachedTwins {
		t.Cache[t.Size] = data
		t.Size++
	} else {
		t.Cache[t.Head] = data
		t.Head = (t.Head + 1) % defaultMaxCachedTwins
	}
}

func (t *TwinsCache) Get(count int) []*TwinsData {
	if count == 0 {
		return nil
	}

	twinsData := make([]*TwinsData, count)
	tempHead := t.Head
	for i := 0; i < count; i++ {
		twinsData[i] = t.Cache[tempHead]
		tempHead = (tempHead + 1) % defaultMaxCachedTwins
	}
	return twinsData
}

func (t *TwinsCache) updateTwins(twinP *v1.TwinProperty) {
	if twinP == nil {
		return
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.Cache[t.Head] == nil {
		t.Cache[t.Head] = &TwinsData{
			Twins: make([]*v1.TwinProperty, 0),
		}
	}
	twinsData := t.Cache[t.Head]

	svc, prop := twinP.Service, twinP.PropertyName

	for i, twin := range twinsData.Twins {
		if twin == nil {
			continue
		}

		if svc != twin.Service {
			continue
		}

		if prop == twin.PropertyName {
			twinsData.Twins[i] = twinP
			return
		}
	}

	twinsData.Twins = append(twinsData.Twins, twinP)
}

func (t *TwinsCache) UpdateHeadTwins(twins []*v1.TwinProperty) {
	if twins == nil {
		return
	}

	//update all desired twins
	for _, twinProperty := range twins {
		if twinProperty == nil {
			continue
		}

		t.updateTwins(twinProperty)
	}
}
