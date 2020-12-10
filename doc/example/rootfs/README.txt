pre-install:
  All files in this directory will be copied to the rootfs at "/"
  before anything begins. These could be setup scripts, repository
  configurations etc.

post-install:
  All files in this directory will be copied over ther rootfs at "/"
  after everything is done. If there are same filenames already
  existing, they will be overwritten.

