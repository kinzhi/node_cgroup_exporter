package collector

import (
	"path/filepath"
	"strings"

	"github.com/prometheus/procfs"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	// The path of the proc filesystem.
	procPath     = kingpin.Flag("path.procfs", "procfs mountpoint.").Default(procfs.DefaultMountPoint).String()
	sysPath      = kingpin.Flag("path.sysfs", "sysfs mountpoint.").Default("/sys").String()
	cgroupPath   = kingpin.Flag("path.cgroupfs", "cgroupfs mountpoint.").Default("/sys/fs/cgroup").String()
	rootfsPath   = kingpin.Flag("path.rootfs", "rootfs mountpoint.").Default("/").String()
	udevDataPath = kingpin.Flag("path.udev.data", "udev data path.").Default("/run/udev/data").String()
)

func procFilePath(name string) string {
	return filepath.Join(*procPath, name)
}

func sysFilePath(name string) string {
	return filepath.Join(*sysPath, name)
}

func cgroupFilePath(name string) string {
	return filepath.Join(*cgroupPath, name)
}

func rootfsFilePath(name string) string {
	return filepath.Join(*rootfsPath, name)
}

func udevDataFilePath(name string) string {
	return filepath.Join(*udevDataPath, name)
}

func rootfsStripPrefix(path string) string {
	if *rootfsPath == "/" {
		return path
	}
	stripped := strings.TrimPrefix(path, *rootfsPath)
	if stripped == "" {
		return "/"
	}
	return stripped
}
