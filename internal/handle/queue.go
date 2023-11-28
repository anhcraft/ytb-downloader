package handle

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"ytb-downloader/internal/settings"
)

var processes = make([]*Process, 0)
var inputted = map[string]int{}

func ClearProcesses() {
	if !isDownloading {
		processes = make([]*Process, 0)
		inputted = map[string]int{}
	}
}

func CountProcess() int {
	return len(processes)
}

func GetProcess(i int) *Process {
	return processes[i]
}

func SubmitUrl(link string, format string, onUpdate func()) {
	u, err := url.Parse(link)
	if err != nil {
		return
	}
	// Validate & normalize URLs
	if u.Path == "/playlist" {
		// async handling
		// TODO lock & disable until finishing
		go submitPlaylistUrl(link, format, onUpdate)
	} else if u.Path == "/watch" {
		link = "https://www.youtube.com/watch?v=" + u.Query().Get("v")
		submitVideoUrl(link, "", format, onUpdate)
	} else if u.Host == "youtu.be" {
		link = "https://www.youtube.com/watch?v=" + u.Path[1:]
		submitVideoUrl(link, "", format, onUpdate)
	} else {
		submitVideoUrl(link, "", format, onUpdate)
	}
}

func submitVideoUrl(link string, name string, format string, onUpdate func()) {
	if _, ok := inputted[link]; ok {
		return
	}
	inputted[link] = len(processes)
	p := &Process{
		Name:   name,
		URL:    link,
		Format: format,
		Status: Queued,
	}
	processes = append(processes, p)
	onUpdate()
	log.Printf("new video link: %s\n", link)
	if p.Name == "" {
		p.Name = p.URL
		go fetchVideoName(p, onUpdate)
	}
}

func fetchVideoName(p *Process, onUpdate func()) {
	// somehow printing into the console does not support UTF8
	// so the workaround is using a temporary file

	temp, err := os.CreateTemp("", hash(p.URL))
	defer func(temp *os.File) {
		err := temp.Close()
		if err != nil {
			log.Println("error closing temp file:", err)
		}
	}(temp)
	if err != nil {
		log.Println("error creating temp file:", err)
		return
	}
	tempPath := temp.Name()

	cmd := exec.Command(settings.Get().GetYTdlpPath(), "--skip-download", "--ignore-errors", "--no-warnings", "--print-to-file", "title", tempPath, p.URL)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	log.Printf("Executing command %s\n", cmd.String())
	if err := cmd.Run(); err != nil {
		log.Println("error running command:", err)
		return
	}

	bytes, err := os.ReadFile(tempPath)
	if err != nil {
		log.Println("error creating temp file:", err)
		return
	}

	p.Name = string(bytes)
	onUpdate()
}

func hash(link string) string {
	hash := sha256.New()
	hash.Write([]byte(link))
	return hex.EncodeToString(hash.Sum(nil))
}

func submitPlaylistUrl(link string, format string, onUpdate func()) {
	log.Printf("new playlist link: %s\n", link)

	// ./yt-dlp.exe --flat-playlist --ignore-errors --no-warnings --print-to-file "title,url" "temp.txt" ""
	// somehow printing into the console does not support UTF8
	// so the workaround is using a temporary file

	temp, err := os.CreateTemp("", hash(link))
	defer func(temp *os.File) {
		err := temp.Close()
		if err != nil {
			log.Println("error closing temp file:", err)
		}
	}(temp)
	if err != nil {
		log.Println("error creating temp file:", err)
		return
	}
	tempPath := temp.Name()

	cmd := exec.Command(settings.Get().GetYTdlpPath(), "--skip-download", "--flat-playlist", "--ignore-errors", "--no-warnings", "--print-to-file", "url,title", tempPath, link)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	log.Printf("Executing command %s\n", cmd.String())
	if err := cmd.Run(); err != nil {
		log.Println("error running command:", err)
		return
	}

	bytes, err := os.ReadFile(tempPath)
	if err != nil {
		log.Println("error creating temp file:", err)
		return
	}
	lines := strings.Split(string(bytes), "\n")

	log.Printf("found %d videos in the playlist\n", len(lines)>>1)
	for i := 0; i+1 < len(lines); i += 2 {
		// TODO better way to check private videos
		log.Println(lines[i+1])
		if strings.Contains(lines[i+1], "[Private video]") {
			continue
		}
		submitVideoUrl(lines[i], lines[i+1], format, onUpdate)
	}
}
