#!/bin/ash

./home_app >> $NFS_MOUNT/$(date +%Y-%m-%d.log) 2>&1
