# Hibernation Setup Guide

## Problem Summary

The machine was failing to hibernate due to insufficient swap space. With 15GB RAM and only 979MB swap, the kernel couldn't write the memory image to disk.

**Error in logs:**
```
PM: Cannot get swap writer
```

## Solution Overview

Created a 16GB hibernation file to store the memory image during hibernation. This file is used in conjunction with `rtcwake` to wake the machine at a specific time.

## Prerequisites

- Root access
- Filesystem with 16GB+ free space
- System with swap support enabled in kernel

## Step-by-Step Configuration

### 1. Create Hibernation File

```bash
# Create 16GB hibernation file (slightly larger than RAM)
sudo fallocate -l 16G /hibernatefile
sudo chmod 600 /hibernatefile
sudo mkswap /hibernatefile
```

### 2. Get Resume Parameters

Get the UUID of the filesystem containing the hibernation file:
```bash
findmnt -no UUID -T /hibernatefile
# Result: 8ce92d7a-9caf-4cf4-a4e5-943959f644a5
```

Get the physical offset of the file:
```bash
sudo filefrag -v /hibernatefile
# Look for "physical_offset" in the first line
# Result: 2766848
```

### 3. Configure GRUB

Edit `/etc/default/grub`:
```bash
GRUB_CMDLINE_LINUX_DEFAULT="quiet resume=UUID=8ce92d7a-9caf-4cf4-a4e5-943959f644a5 resume_offset=2766848"
```

Apply changes:
```bash
sudo update-grub
sudo update-initramfs -u -k all
```

### 4. Reboot

```bash
sudo reboot
```

### 5. Verify Kernel Parameters

After reboot, check that kernel loaded the resume parameters:
```bash
cat /proc/cmdline
# Should contain:
# resume=UUID=8ce92d7a-9caf-4cf4-a4e5-943959f644a5 resume_offset=2766848
```

## Updated nightsleep.sh Script

```bash
#!/bin/bash
WAKE_TIME="09:00"
CURRENT=$(date +%H:%M)

# CRITICAL: Activate hibernation swap file
swapon /hibernatefile 2>/dev/null || true

# Calculate seconds until wake time
if [[ "$CURRENT" < "$WAKE_TIME" ]]; then
    # Wake time is later today (e.g., 1 AM -> 9 AM)
    TARGET=$(date -d "$WAKE_TIME" +%s)
else
    # Wake time is tomorrow (e.g., 10 PM -> 9 AM next day)
    TARGET=$(date -d "tomorrow $WAKE_TIME" +%s)
fi

NOW=$(date +%s)
SLEEP_SECS=$((TARGET - NOW))

logger -t sleeper "Hibernating for $SLEEP_SECS seconds, wake at $WAKE_TIME"
/usr/sbin/rtcwake -m disk -s $SLEEP_SECS
systemctl hibernate
```

**Installation:**
```bash
chmod +x /root/nightsleep.sh
```

**Cron setup:**
```bash
# Edit crontab
sudo crontab -e

# Add line to run at 01:20 daily
20 1 * * * /root/nightsleep.sh
```

## Verification Checklist

After reboot, verify everything is configured correctly:

```bash
# 1. Check kernel command line
cat /proc/cmdline | grep resume
# Expected: resume=UUID=8ce92d7a-9caf-4cf4-a4e5-943959f644a5 resume_offset=2766848

# 2. Check hibernation file exists and is valid swap
file /hibernatefile
# Expected: Linux swap file, 4k page size, little endian, version 1

# 3. Check swap is active
swapon --show
# Expected: Should show /hibernatefile with 16G size

# 4. Check hibernation resume config
cat /sys/power/resume
cat /sys/power/resume_offset
```

## Troubleshooting

### Hibernation fails with "Cannot get swap writer"
- **Cause:** Swap space smaller than RAM
- **Fix:** Ensure `/hibernatefile` is active with `swapon /hibernatefile`

### Machine wakes immediately instead of at scheduled time
- **Cause:** `rtcwake` failed or hardware doesn't support wake alarm
- **Fix:** Check `rtcwake` is installed: `which rtcwake`
- **Check:** Test manually: `sudo rtcwake -m show -s 60`

### System hangs during hibernation
- **Cause:** Incompatible hardware or drivers
- **Fix:** Check kernel logs: `dmesg | grep -i hibernation`
- **Alternative:** Use `mem` sleep mode instead of `disk`

### Hibernate file not auto-activating on boot
- **Solution:** Add to `/etc/fstab`:
  ```
  /hibernatefile none swap sw 0 0
  ```

## Script
```bash
#!/bin/bash
WAKE_TIME="09:00"
CURRENT=$(date +%H:%M)
swapon /hibernatefile 2>/dev/null || true                                                                                                                          MCP

# Calculate seconds until wake time
if [[ "$CURRENT" < "$WAKE_TIME" ]]; then
    # Wake time is later today (1 AM case)
    TARGET=$(date -d "$WAKE_TIME" +%s)
else
    # Wake time is tomorrow (already past 7 AM)
    TARGET=$(date -d "tomorrow $WAKE_TIME" +%s)
fi
NOW=$(date +%s)
SLEEP_SECS=$((TARGET - NOW))
logger -t sleeper "Hibernating for $SLEEP_SECS seconds, wake at $WAKE_TIME"
/usr/sbin/rtcwake -m disk -s $SLEEP_SECS

```

## References

- [Arch Wiki - Suspend and hibernate](https://wiki.archlinux.org/title/Power_management/Suspend_and_hibernate)
- [Kernel documentation - swsusp](https://www.kernel.org/doc/Documentation/power/swsusp.txt)
- `man rtcwake`
- `man systemd-sleep`

