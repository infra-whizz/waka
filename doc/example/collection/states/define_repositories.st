id: init-repos
description: Define repository provisioning
state:
  # Create basic directories for the repos
  # This also can be done in rootfs overlay,
  # but for the consistency sake, implemented
  # as a state rule.
  setup-repos-dirs:
    - ansible.files.file []add_directories:

  # Setup repositories for Zypper
  setup-repos-config:
    - ansible.packaging.os.zypper_repository []add_repositories:

