# YTB Downloader

![GitHub Release](https://img.shields.io/github/v/release/anhcraft/ytb-downloader)
![GitHub License](https://img.shields.io/github/license/anhcraft/ytb-downloader)


A simple GUI for yt-dlp. [Download here](https://github.com/anhcraft/ytb-downloader/releases)

![https://i.imgur.com/vJOuWBb.png](https://i.imgur.com/vJOuWBb.png)

## Features
- Support various websites (as long as yt-dlp supports)
  - Tested sites: Youtube, Reddit, Tiktok, Twitter/X, Weibo
- Fetch videos from playlist
- Download the best format
- Download video-only or audio-only
- Embed thumbnail
- Concurrent downloading
- Continue downloading after interrupt
- Logging to file and console
- CJK font support

## Requirement
- [yt-dlp](https://github.com/yt-dlp/yt-dlp)
- [FFmpeg](https://ffmpeg.org/) including `ffmpeg` and `ffprobe`

## Installation
### Windows
- Download [FFmpeg](https://ffmpeg.org/) (including `ffmpeg` and `ffprobe`)
- Download [yt-dlp](https://github.com/yt-dlp/yt-dlp) binary file
- Configure using GUI or via file (see below)

## Linux
- Install FFmpeg: `sudo apt install ffmpeg` (including `ffmpeg` and `ffprobe`)
- Locate FFmpeg binary file: `sudo dpkg -L ffmpeg` (usually `/usr/bin/ffmpeg`)
- Download [yt-dlp](https://github.com/yt-dlp/yt-dlp) binary file
- Configure using GUI or via file (see below)

## Configuration
Open the file `settings.json` via an editor or command: `sudo nano settings.json`.

Example configuration for Linux environment:
```json
{
  "format":"Default",
  "embedThumbnail":"AudioOnly",
  "ytdlpPath":"/yt-dlp_linux",
  "ffmpegPath":"/usr/bin/ffmpeg",
  "concurrentDownloads":1,
  "concurrentFragments":3,
  "logPath": "./log.txt"
}
```

## Notes
- If you use the Windows installer, the default installation path is `C:\Program Files (x86)\YTB Downloader`
- `"ERROR: [WinError 32] The process cannot access the file because it is being used by another process:"`: might happen when you download the video into the folder currently opening in the IDE. To fix, pick a different download location.
- Cannot run on Windows? Enable running as administrator (see: https://i.imgur.com/ALTWIM4.png)

## Building
- To build the program, first install Fyne: `go install fyne.io/fyne/v2/cmd/fyne@latest`
- Build the program: `sudo ./build.sh`
