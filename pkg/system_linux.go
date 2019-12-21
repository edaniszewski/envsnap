package pkg

import (
	"runtime"

	log "github.com/sirupsen/logrus"
)

// LoadSystemInfo loads the system info (os, arch, etc) into a SysInfo struct.
func LoadSystemInfo() (SysInfo, error) {
	stdout, stderr, err := runCommand("uname", "-srpo")
	if err != nil {
		errString := stderr.String()
		if errString == "" {
			errString = "<no output>"
		}
		log.Debugf("command error: %v", errString)
		return SysInfo{}, err
	}
	info := toSlice(stdout.Bytes())

	return SysInfo{
		OS:            runtime.GOOS,
		Kernel:        info[0],
		KernelVersion: info[1],
		Processor:     info[2],
		Arch:          info[3],
	}, nil
}
