package titime

import (
	"runtime"
	"syscall"
	"unsafe"
)

const _NRSchedSetaffinity = 203

func pinCPU(cpu uint) error {
	var mask [1024 / 8]uint8
	runtime.LockOSThread()
	mask[cpu/8] |= 1 << (cpu % 8)
	_, _, errno := syscall.RawSyscall(_NRSchedSetaffinity, uintptr(0), uintptr(len(mask)*8), uintptr(unsafe.Pointer(&mask)))
	if errno != 0 {
		return errno
	}
	return nil
}

func unpinCPU() error {
	var mask [1024 / 8]uint8
	for i := range mask {
		mask[i] = 0xff
	}
	_, _, errno := syscall.RawSyscall(_NRSchedSetaffinity, uintptr(0), uintptr(len(mask)*8), uintptr(unsafe.Pointer(&mask)))
	runtime.UnlockOSThread()
	if errno != 0 {
		return errno
	}
	return nil
}
