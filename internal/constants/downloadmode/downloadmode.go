package downloadmode

const (
	Default            = "Default"
	CustomDownloadOnly = "CustomDownloadOnly"
	YtdlpDownloadOnly  = "YtdlpDownloadOnly"
)

func IsValid(s string) bool {
	return s == Default || s == CustomDownloadOnly || s == YtdlpDownloadOnly
}

func HasCustomDownload(enum string) bool {
	return enum == Default || enum == CustomDownloadOnly
}

func HasYtdlpDownload(enum string) bool {
	return enum == Default || enum == YtdlpDownloadOnly
}
