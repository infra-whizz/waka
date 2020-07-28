id: linux-tree
description: Binding linux kernel FS
state:
  # bind-mount /proc, /sys etc
  bind-mount:
    - ansible.system.mount []mount_kernel_dirs:

  # bint-umount /proc, /sys etc
  bind-umount:
    - ansible.system.mount []umount_kernel_dirs:
