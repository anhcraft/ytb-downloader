# YTB Downloader

A simple GUI for yt-dlp.

## Features
- Download the best format
- Download video-only or audio-only
- Continue downloading after interrupt
- Fetch videos from playlist
- Support other sites beside Youtube

## Requirement
- [yt-dlp](https://github.com/yt-dlp/yt-dlp)
- [FFmpeg](https://ffmpeg.org/) including `ffmpeg` and `ffprobe`

## Notes
- For Youtube, audio-only format is .m4a, otherwise, audio format is opus (default Windows player may not support opus)
- Tested sites: Youtube, Reddit, Tiktok, Twitter, Weibo