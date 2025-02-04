package thumbnail

const (
	Always    = "Always"
	VideoOnly = "VideoOnly"
	AudioOnly = "AudioOnly"
	Never     = "Never"
)

func IsValid(s string) bool {
	return s == Always || s == VideoOnly || s == AudioOnly || s == Never
}
