#!/bin/bash
# Organize photos into YYYYMM folders
# Run daily via cron at 4:00 AM
#
# Crontab entry:
# 0 4 * * * /mnt/external/drive/photos/organize_photos.sh
#
# Handles patterns:
# - YYYYMMDD_* (e.g., 20210810_123456.jpg)
# - YYYY-MM-DD* (e.g., 2021-03-29-155232.jpg)
# - VID_YYYYMMDD* (e.g., VID_20210810_123456.mp4)
# - signal-YYYY-MM-DD* (e.g., signal-2025-12-15-20-41-05.jpg)

PHOTOS_DIR="/mnt/external/drive/photos"
LOG_FILE="/var/log/organize_photos.log"

cd "$PHOTOS_DIR" || exit 1

total_count=0

# Pattern 1: YYYYMMDD_* files
for file in 20[0-9][0-9][0-9][0-9][0-9][0-9]_*; do
  if [ -f "$file" ]; then
    folder="${file:0:6}"
    mkdir -p "$folder"
    mv "$file" "$folder/"
    ((total_count++))
  fi
done

# Pattern 2: YYYY-MM-DD* files -> rename to YYYYMMDD* and move
for file in 20[0-9][0-9]-[0-9][0-9]-[0-9][0-9]*; do
  if [ -f "$file" ]; then
    folder="${file:0:4}${file:5:2}"
    newname="${file:0:4}${file:5:2}${file:8:2}${file:10}"
    mkdir -p "$folder"
    mv "$file" "$folder/$newname"
    ((total_count++))
  fi
done

# Pattern 3: VID_YYYYMMDD* files
for file in VID_20[0-9][0-9][0-9][0-9][0-9][0-9]*; do
  if [ -f "$file" ]; then
    folder="${file:4:6}"
    mkdir -p "$folder"
    mv "$file" "$folder/"
    ((total_count++))
  fi
done

# Pattern 4: signal-YYYY-MM-DD* files -> rename to YYYYMMDD_*-signal and move
for file in signal-20[0-9][0-9]-[0-9][0-9]-[0-9][0-9]*; do
  if [ -f "$file" ]; then
    ext="${file##*.}"
    base="${file%.*}"
    # Extract YYYYMM for folder
    folder="${file:7:4}${file:12:2}"
    # Build new filename: YYYYMMDD_HHMMSSmmm-rest-signal.ext
    date="${file:7:4}${file:12:2}${file:15:2}"
    rest="${base:17}"
    # Remove dashes from time part
    newrest=$(echo "$rest" | sed "s/^\([0-9][0-9]\)-\([0-9][0-9]\)-\([0-9][0-9]\)-\([0-9]*\)/\1\2\3\4/")
    newname="${date}_${newrest}-signal.${ext}"
    mkdir -p "$folder"
    mv "$file" "$folder/$newname"
    ((total_count++))
  fi
done

if [ $total_count -gt 0 ]; then
  echo "$(date '+%Y-%m-%d %H:%M:%S') - Moved $total_count files" >> "$LOG_FILE"
fi
