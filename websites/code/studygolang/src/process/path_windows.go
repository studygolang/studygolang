package process

import (
	"path/filepath"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

func ExecutableDir() (string, error) {
	var (
		dir, p string
		err    error
	)
	p, err = getWindowsProcessBinaryPath()
	if err != nil {
		return "", err
	}
	dir = filepath.Dir(p)
	dir = strings.Replace(dir, "\\", "/", -1)
	return dir, nil
}

func getWindowsProcessBinaryPath() (string, error) {
	b := make([]uint16, 300)
	n, e := getModuleFileName(uint32(len(b)), &b[0])
	if e != nil {
		return "", e
	}
	return string(utf16.Decode(b[0:n])), nil
}

func getModuleFileName(buflen uint32, buf *uint16) (n uint32, err error) {
	h, err := syscall.LoadLibrary("kernel32.dll")
	if err != nil {
		return 0, err
	}
	defer syscall.FreeLibrary(h)
	addr, err := syscall.GetProcAddress(h, "GetModuleFileNameW")
	if err != nil {
		return 0, err
	}
	r0, _, e1 := syscall.Syscall(addr, 3, uintptr(0), uintptr(unsafe.Pointer(buf)), uintptr(buflen))
	n = uint32(r0)
	if n == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
