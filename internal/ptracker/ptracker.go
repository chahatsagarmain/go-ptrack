package ptracker

import (
	"fmt"
	"os"
)

// ptracker will check and output specific directories for process info
// /proc/<pid - if directory does not exist then the process is killed and return the error
// /proc/<pid>/mem - memory information
// /proc/<pid>/cwd - get current working directory

type Tracker interface {
	GetProcPath(string, int) string
	GetStatus(int) (int, error)
	GetMem(int) (string, error)
	GetCwd(int) (string, error)
	GetExe(int) (string, error)
	GetCommandLine(int) (string, error)
	GetIO(int) (string, error)
	GetSysCall(int) (string, error)
	GetFD(int) (string,int, error)
}

func GetProcPath(filename string, pid int) string {
	return fmt.Sprintf("/proc/%d/%s", pid, filename)
}

func GetStatus(pid int) (int, error) {
	procPath := fmt.Sprintf("/proc/%d", pid)
	if _, err := os.Stat(procPath); os.IsNotExist(err) {
		return 0, fmt.Errorf("process not found \n")
	}
	return 1, nil
}

func GetCommandLine(pid int) (string, error) {
	if _, err := GetStatus(pid); err != nil {
		return "", err
	}
	procPath := GetProcPath("cmdline", pid)
	var data []byte
	var err error
	if data, err = os.ReadFile(procPath); err != nil {
		return "", err
	}
	return string(data), nil
}

func GetCwd(pid int) (string, error) {
	if _, err := GetStatus(pid); err != nil {
		return "", err
	}
	cwdPath := fmt.Sprintf("/proc/%d/cwd", pid)
	target, err := os.Readlink(cwdPath)
	if err != nil {
		return "", err
	}
	return target, nil
}

func GetExe(pid int) (string, error) {
	if _, err := GetStatus(pid); err != nil {
		return "", err
	}
	cwdPath := fmt.Sprintf("/proc/%d/exe", pid)
	target, err := os.Readlink(cwdPath)
	if err != nil {
		return "", err
	}
	return target, nil
}

func GetIO(pid int) (string, error) {
	if _, err := GetStatus(pid); err != nil {
		return "", err
	}
	cwdPath := fmt.Sprintf("/proc/%d/io", pid)
	target, err := os.ReadFile(cwdPath)
	if err != nil {
		return "", err
	}
	return string(target), nil
}

func GetSysCall(pid int) (string, error) {
	if _, err := GetStatus(pid); err != nil {
		return "", err
	}
	cwdPath := fmt.Sprintf("/proc/%d/syscall", pid)
	target, err := os.ReadFile(cwdPath)
	if err != nil {
		return "", err
	}
	return string(target), nil
}

func GetMem(pid int) (string, error) {
	if _, err := GetStatus(pid); err != nil {
		return "", err
	}
	cwdPath := fmt.Sprintf("/proc/%d/statm", pid)
	target, err := os.ReadFile(cwdPath)
	if err != nil {
		return "", err
	}
	return string(target), nil
}

func GetFD(pid int) (string, int, error) {
	if _, err := GetStatus(pid); err != nil {
		return "",0, err
	}
	fdPath := fmt.Sprintf("/proc/%d/fd", pid)
	entries, err := os.ReadDir(fdPath)
	if err != nil {
		return "",0, err
	}

	result := ""
	for _, entry := range entries {
		fd := entry.Name()
		linkPath := fmt.Sprintf("%s/%s", fdPath, fd)
		target, err := os.Readlink(linkPath)
		if err != nil {
			target = fmt.Sprintf("unreadable: %v", err)
		}
		result += fmt.Sprintf("fd %s -> %s\n", fd, target)
	}
	return result,len(entries), nil
}
