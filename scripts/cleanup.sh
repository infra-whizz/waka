#!/bin/sh

# order matters
for idx in "0" "1" "2" "3" "4" "5" "6" "7"; do
    echo "Cleaning up /dev/loop$idx device..."
    for dev in "loop${idx}p1" "loop${idx}p2" "loop${idx}p4" "loop${idx}p3"; do
	umount -f /dev/$dev 2>/dev/null
    done
    losetup -d /dev/loop$idx 2>/dev/null
done

echo "Remove build mounts"
rmdir /tmp/waka-build* 2>/dev/null

echo "Remove build RAW"
rm /tmp/build/* 2>/dev/null

echo "Done"
