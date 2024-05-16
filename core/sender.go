package core

import (
	"time"

	"github.com/edgehook/ithings/common/global"
	"github.com/edgehook/ithings/common/types"
)

/*
* SendSyncRequestToEdge
* send the request and wait response.
* timeOut = 0, we will wait it forver.
* this function is thread safe.
 */
func SendSyncRequestToEdge(req *types.Request, timeOut time.Duration) (*types.Response, error) {
	if defaultICore == nil {
		return nil, global.ErrCoreNotReady
	}

	return defaultICore.iCore.SendSyncRequestToEdge(req, timeOut)
}
