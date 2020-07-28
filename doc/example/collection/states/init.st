id: init
description: Bootstrap generic environment
state:
  # bind-mount /proc, /sys etc
  #~linux-tree/bind-mount:

  ~init-repos:
  ~provision-rootfs:

  # bind-umount /proc, sys etc
  #~linux-tree/bind-umount:

  #+overlay/additional:

