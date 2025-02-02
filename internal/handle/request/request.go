package request

import (
	"net/url"
	"sync/atomic"
)

const (
	StatusInQueue uint32 = iota // including title download
	StatusDownloading
	StatusCompleted
	StatusFailed
	StatusTerminated
)

type Request struct {
	// [*] following fields are expected to be read-only
	url    *url.URL
	rawUrl string // cache

	// [*] following fields use relaxed consistency, no lock required
	// init on "Fetch" button
	title             string
	titleFetchCmdArgs []string
	titleFetched      bool

	// init on "Download" button
	format          string
	downloadCmdArgs []string

	// [*] update during downloading
	downloadProgress  string
	downloadedSize    string
	downloadTotalSize string
	downloadSpeed     string
	downloadEta       string
	downloadError     error

	// [*] following fields are expected to be atomic
	status atomic.Uint32
}

func NewRequest(url *url.URL) *Request {
	return &Request{
		url:    url,
		rawUrl: url.String(),
		title:  url.String(),
	}
}

func (req *Request) Url() *url.URL {
	return req.url
}

func (req *Request) RawUrl() string {
	return req.rawUrl
}

func (req *Request) Title() string {
	return req.title
}

func (req *Request) SetTitle(title string) {
	req.title = title
}

func (req *Request) TitleFetched() bool {
	return req.titleFetched
}

func (req *Request) SetTitleFetched(titleFetched bool) {
	req.titleFetched = titleFetched
}

func (req *Request) Format() string {
	return req.format
}

func (req *Request) SetFormat(format string) {
	req.format = format
}

func (req *Request) TitleFetchCmdArgs() []string {
	return req.titleFetchCmdArgs
}

func (req *Request) SetTitleFetchCmdArgs(titleFetchCmdArgs []string) {
	req.titleFetchCmdArgs = titleFetchCmdArgs
}

func (req *Request) DownloadCmdArgs() []string {
	return req.downloadCmdArgs
}

func (req *Request) SetDownloadCmdArgs(downloadCmdArgs []string) {
	req.downloadCmdArgs = downloadCmdArgs
}

func (req *Request) DownloadProgress() string {
	return req.downloadProgress
}

func (req *Request) SetDownloadProgress(downloadProgress string) {
	req.downloadProgress = downloadProgress
}

func (req *Request) DownloadedSize() string {
	return req.downloadedSize
}

func (req *Request) SetDownloadedSize(downloadedSize string) {
	req.downloadedSize = downloadedSize
}

func (req *Request) DownloadTotalSize() string {
	return req.downloadTotalSize
}

func (req *Request) SetDownloadTotalSize(downloadTotalSize string) {
	req.downloadTotalSize = downloadTotalSize
}

func (req *Request) DownloadSpeed() string {
	return req.downloadSpeed
}

func (req *Request) SetDownloadSpeed(downloadSpeed string) {
	req.downloadSpeed = downloadSpeed
}

func (req *Request) DownloadEta() string {
	return req.downloadEta
}

func (req *Request) SetDownloadEta(downloadEta string) {
	req.downloadEta = downloadEta
}

func (req *Request) DownloadError() error {
	return req.downloadError
}

func (req *Request) SetDownloadError(downloadError error) {
	req.downloadError = downloadError
}

func (req *Request) Status() uint32 {
	return req.status.Load()
}

func (req *Request) SetStatus(status uint32) {
	switch status {
	case StatusDownloading:
		req.status.CompareAndSwap(StatusInQueue, StatusDownloading)
	case StatusCompleted:
		req.status.CompareAndSwap(StatusDownloading, StatusCompleted)
	case StatusFailed:
		req.status.CompareAndSwap(StatusDownloading, StatusFailed)
	case StatusTerminated:
		req.status.CompareAndSwap(StatusInQueue, StatusTerminated)
		req.status.CompareAndSwap(StatusDownloading, StatusTerminated)
	case StatusInQueue:
	}
}

func (req *Request) DescribeStatus() string {
	switch req.Status() {
	case StatusInQueue:
		return "In Queue"
	case StatusDownloading:
		return "Downloading"
	case StatusCompleted:
		return "Completed"
	case StatusFailed:
		return "Failed"
	case StatusTerminated:
		return "Terminated"
	}
	return ""
}
