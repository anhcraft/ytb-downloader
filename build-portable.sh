#!/bin/bash

zipfile="YTB-Downloader.Windows.Portable.zip"

7z a "$zipfile" "assets" "YTB Downloader.exe"

if [ $? -eq 0 ]; then
  echo "Successfully created $zipfile"
else
  echo "An error occurred while creating the zip file."
fi
