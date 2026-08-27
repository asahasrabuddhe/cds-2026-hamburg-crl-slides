cmd := newChild(modeChild)
cmd.SysProcAttr = &syscall.SysProcAttr{
	Cloneflags: syscall.CLONE_NEWUSER | syscall.CLONE_NEWNS,
	UidMappings: []syscall.SysProcIDMap{
		{ContainerID: 0, HostID: uid, Size: 1},
	},
	GidMappings: []syscall.SysProcIDMap{
		{ContainerID: 0, HostID: gid, Size: 1},
	},
	GidMappingsEnableSetgroups: false,
}
