package format

const (
	Default   = "Default"
	VideoOnly = "VideoOnly"
	AudioOnly = "AudioOnly"
)

func IsValid(s string) bool {
	return s == Default || s == VideoOnly || s == AudioOnly
}
