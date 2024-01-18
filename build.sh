#!/bin/bash

os=$(uname -s)

case "$os" in
    Darwin)
      fyne package -os darwin -icon "./assets/ytb.png"
        ;;
    Linux)
      fyne package -os linux -icon "./assets/ytb.png"
        ;;
    MINGW*) # Windows Git bash
      fyne package -os windows -icon "./assets/ytb.png"
        ;;
    *)
        echo "Unsupported operating system."
        exit 1
        ;;
esac