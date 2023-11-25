package handle

type Status uint

const (
	Queued      Status = 0
	Downloading Status = 1
	Done        Status = 2
	Error       Status = 3
)

func (o Status) String() string {
	switch o {
	case 0:
		return "Queued"
	case 1:
		return "Downloading"
	case 2:
		return "Done"
	case 3:
		return "Error"
	default:
		return "Unknown"
	}
}
