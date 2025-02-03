package component

import (
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"sync"
	"time"
)

const delay = time.Millisecond * 500

func NewAutoSaveInput(reader func() string, writer func(val string), validator func(string) error) *widget.Entry {
	bind := binding.NewString()
	_ = bind.Set(reader())

	var timer *time.Timer
	var mutex sync.Mutex

	bind.AddListener(binding.NewDataListener(func() {
		mutex.Lock()
		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(delay, func() {
			val, err := bind.Get()
			if err != nil {
				return
			}
			if validator != nil {
				if err = validator(val); err != nil {
					return
				}
			}
			writer(val)
		})
		mutex.Unlock()
	}))

	entry := widget.NewEntryWithData(bind)
	entry.Validator = validator
	return entry
}
