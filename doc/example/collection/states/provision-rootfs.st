id: provision-rootfs
description: Install rootfs packages, provisioning it
state:
  provision-rootfs:
    - ansible.packaging.os.zypper []rootfs_packages:

    # - ansible.packaging.os.zypper:
    #     packages:
    #       - aaa_base
    #       - emacs
    #     root: "/tmp/test"
    #     pipe_file: "/tmp/test.pipe.txt"
    #     stuff &pref:
