package shellquote

import (
	"github.com/sergeymakinen/go-quote/unix"
	"github.com/sergeymakinen/go-quote/windows"
	"runtime"
	"strings"
)

func Join(args []string) string {
	switch runtime.GOOS {
	case "windows":
		return strings.Join(Map(args, windows.Argv.Quote), " ")
	}

	return strings.Join(Map(args, unix.ANSIC.Quote), " ")
}

func Map(args []string, f func(string) string) []string {
	var out []string
	for _, arg := range args {
		out = append(out, f(arg))
	}
	return out
}
