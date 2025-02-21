#!/bin/bash

os=$(uname -s)
version="1.0.8"

case "$os" in
    Darwin)
      fyne package -os darwin -icon "./assets/ytb.png" -name "YTB Downloader" -appVersion "$version"
        ;;
    Linux)
      fyne package -os linux -icon "./assets/ytb.png" -name "YTB Downloader" -appVersion "$version"
        ;;
    MINGW*) # Windows Git bash
      fyne package -os windows -icon "./assets/ytb.png" -name "YTB Downloader" -appVersion "$version" -executable "YTB Downloader.exe"
        ;;
    *)
        echo "Unsupported operating system."
        exit 1
        ;;
esac