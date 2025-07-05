# Mount disk
## Find block devices
```bash
lsblk
```
## Find disk's UUID you want to mount,
e.g. `/dev/sdb1` 
```bash
blkid
#e.g. 
blkid /dev/sdb1

```
## Modify fstab
```bash
sudo vim /etc/fstab
```



## Create mount point
```bash
sudo mkdir /mnt/mydisk
```

## Mount disk
```bash
mount -a #there could be need to reload systemctl
```
