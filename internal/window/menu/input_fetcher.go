package menu

import (
	"sync"
	"ytb-downloader/internal/handle/request"
)

var fetchingMutex sync.Mutex

func FetchInput(in string, callback func(req []*request.Request)) bool {
	if !fetchingMutex.TryLock() {
		return false
	}

	go func() {
		defer fetchingMutex.Unlock()
		req := request.ParseRequest(in)
		request.GetTable().AddBulk(req)
		callback(req)
	}()

	return true
}
