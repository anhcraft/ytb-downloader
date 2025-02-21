# YTB Downloader

![GitHub Release](https://img.shields.io/github/v/release/anhcraft/ytb-downloader)
![GitHub License](https://img.shields.io/github/license/anhcraft/ytb-downloader)


A simple GUI for yt-dlp. [Download here](https://github.com/anhcraft/ytb-downloader/releases)

![https://i.imgur.com/UtKVy6k.png](https://i.imgur.com/UtKVy6k.png)

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
- Multiple setting profiles

## Requirement
- [yt-dlp](https://github.com/yt-dlp/yt-dlp)
- [FFmpeg](https://ffmpeg.org/) including `ffmpeg` and `ffprobe`

## Installation
### Windows
- Download [FFmpeg](https://ffmpeg.org/) (including `ffmpeg` and `ffprobe`)
- Download [yt-dlp](https://github.com/yt-dlp/yt-dlp) binary file
- Configure using GUI or via file

## Linux
- Install FFmpeg: (including `ffmpeg` and `ffprobe`)
```bash
sudo apt install ffmpeg
```
- Locate FFmpeg binary file (e.g. `/usr/bin/ffmpeg`)
```bash
sudo dpkg -L ffmpeg
```
- Download [yt-dlp](https://github.com/yt-dlp/yt-dlp) binary file
- Configure using GUI or via file

## Scripting
- You can write [Tengo](https://tengolang.com/) script to control fetching request
- The script is loaded every time a batch of input is _fetched_. However, if it does not exist, the default handle is used
- The script is executed once for each request

### Input & Output variables
- Input: `_input` denotes a line of input (space-stripped guaranteed)
- Required output: `_action` and `_url`
  - `_action = skip`: skip this request
  - `_action = override`: override the input with `_url`, and continue with the default handle. Remember that the new input must be compatible to ytdlp
  - `_action = custom`: download the file from `_url` using custom downloader (not Yt-dlp)
    - You must also specify `_filepath` denoting the target file (subdirectories are automatically created). Prefix the path with `$DOWNLOAD_FOLDER/` to start at the download folder
  - other values: continue with the default handle
- Optional output:
  - `_title`: title of the request (defaults to the input)

### Modules
- All Tengo standard-library modules are enabled (including file and OS access)
- Additional modules (made by YTB-Downloader)
  - `url`: URL utilities ([View docs](https://github.com/anhcraft/ytb-downloader/blob/main/internal/scripting/module/url.go))

### Security
- **Warning**: Be careful when using the script, you must acknowledge what it does under the hood. Do not use script taken from untrusted sources.

### Example
- See [script.tengo](script.tengo) for an example

## FAQ

### 1. How do I find the installation folder?
- The installation folder must contain the setting file, as such you can open the settings window, click on "Locate setting file" to open the explorer at the installation folder (Windows-only)
- Another approach is to open the info window, then find the path to the executable file (which must be inside the installation folder)

### 2. I got `"ERROR: [WinError 32] The process cannot access the file because it is being used by another process:"`
- This might happen when you download the video into the folder currently opening in the IDE. To resolve, pick a different download location.

### 3. I am required to elevate privilege (administrator) to open the app
- The problem happens when you install the app for all users (under `Program Files (x86)`). To resolve, enable running as administrator (see: https://i.imgur.com/ALTWIM4.png) - this is one-time process and typically does not need to do it again when manually updating the app
- Alternatively, install the app for a specific user. The app will be installed to `C:\Users\User\AppData\Local\Programs\YTB Downloader`

### 4. I got `HTTP Error 403: Forbidden` when downloading from YouTube
- You might be playing YouTube from the browser while downloading from the GUI both at the same time

## Building
- To build the program, first install Fyne: 
```bash
go install fyne.io/fyne/v2/cmd/fyne@latest
````
- Build the program:

```bash
sudo ./build.sh
````
