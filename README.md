# YTB Downloader

![GitHub Release](https://img.shields.io/github/v/release/anhcraft/ytb-downloader)
![GitHub License](https://img.shields.io/github/license/anhcraft/ytb-downloader)


A simple GUI for yt-dlp. [Download here](https://github.com/anhcraft/ytb-downloader/releases)

![https://i.imgur.com/zpvzS7y.png](https://i.imgur.com/zpvzS7y.png)

## Features
- Support various websites (as long as yt-dlp supports)
  - Manually-tested sites: YouTube, Facebook, Reddit, TikTok, Twitter/X, Weibo, BiliBili, Soundcloud, Vimeo
  - View the complete list: https://github.com/yt-dlp/yt-dlp/blob/master/supportedsites.md
- Fetch videos from playlist
  - **Note**: YouTube playlist is automatically flattened, then you can remove or keep each video from the table
- Download the best format
- Download video-only or audio-only
- Embed thumbnail
- Concurrent downloading
- Continue downloading after interrupt
- Logging to file and console
- CJK font support
- Scripting support using [Tengo](https://tengolang.com/) (see below)

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

## Scripting
- You can write [Tengo](https://tengolang.com/) script to control fetching request
- The script is loaded every time a batch of input is _fetched_. However, if it does not exist, the default handle is used
- The script is executed once for each request
- Input: `_input` denotes a line of input (space-stripped guaranteed)
- Output: `_action` and `_value`
  - `_action = skip`: skip this request
  - `_action = override`: override the input with `_value`, and continue with the default handle. Remember that the new input must be compatible to ytdlp
  - other values: continue with the default handle
- All Tengo stdlib modules are enabled (including file and OS access)
- Additional modules (made by YTB-Downloader)
  - `url`: URL utilities ([View docs](https://github.com/anhcraft/ytb-downloader/blob/main/internal/scripting/module/url.go))
- **Warning**: Be careful when using the script, you must acknowledge what it does under the hood. Do not use script taken from untrusted sources.
- See [script.tengo](script.tengo) for an example

## Notes
- If you use the Windows installer, the default installation path is `C:\Program Files (x86)\YTB Downloader`
- `"ERROR: [WinError 32] The process cannot access the file because it is being used by another process:"`: might happen when you download the video into the folder currently opening in the IDE. To fix, pick a different download location.
- Cannot run on Windows? Enable running as administrator (see: https://i.imgur.com/ALTWIM4.png)
- Got `HTTP Error 403: Forbidden`? You might be playing YouTube from the browser while downloading from the GUI both at the same time

## Building
- To build the program, first install Fyne: `go install fyne.io/fyne/v2/cmd/fyne@latest`
- Build the program: `sudo ./build.sh`
