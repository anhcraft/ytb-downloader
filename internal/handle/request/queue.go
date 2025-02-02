package request

import (
	"sync"
)

type Queue struct {
	lock           sync.Mutex
	queue          []*Request
	updateCallback func(req *Request)
}

var updateCallbackOnce sync.Once
var queueInstance = &Queue{}

func GetQueue() *Queue {
	return queueInstance
}

func (rq *Queue) Poll() *Request {
	rq.lock.Lock()
	defer rq.lock.Unlock()
	if len(rq.queue) == 0 {
		return nil
	}
	req := rq.queue[0]
	rq.queue = rq.queue[1:]
	return req
}

func (rq *Queue) OfferBulk(reqs []*Request) {
	rq.lock.Lock()
	defer rq.lock.Unlock()
	rq.queue = append(rq.queue, reqs...)
}

func (rq *Queue) SetUpdateCallback(updateCallback func(req *Request)) {
	updateCallbackOnce.Do(func() {
		rq.updateCallback = updateCallback
	})
}

func (rq *Queue) OnUpdate(req *Request) {
	rq.updateCallback(req)
}
