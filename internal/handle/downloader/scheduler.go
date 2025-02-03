package downloader

import (
	"sync"
	"sync/atomic"
	"ytb-downloader/internal/handle/request"
	"ytb-downloader/internal/settings"
)

var once sync.Once
var interrupt atomic.Bool
var workerCount atomic.Uint32

func InitDownloadScheduler() {
	once.Do(func() {
		go runDownloadScheduler()
	})
}

func TerminateDownloadScheduler() {
	interrupt.Store(true)
}

func runDownloadScheduler() {
	for !interrupt.Load() {
		if workerCount.Load() >= settings.Get().GetConcurrentDownloads() {
			continue
		}

		req := request.GetQueue().Poll()

		// Race condition is acceptable for "status", so here we re-check again
		if req == nil || req.Status() != request.StatusInQueue {
			continue
		}

		// Set status here instead of delegating to the worker
		req.SetStatus(request.StatusDownloading)
		request.GetQueue().OnUpdate(req)
		workerCount.Add(1)
		go download(req, func() {
			workerCount.Add(^uint32(0))
		})
	}
}
