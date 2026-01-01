package main

import (
	"bytes"
	"fmt"
	"syscall"
	"unsafe"
)

type DWORD uint32
type WORD uint16

type file_time struct {
	low_date_time  DWORD
	high_date_time DWORD
}

type win32_find_data struct {
	file_attributes    DWORD
	creation_time      file_time
	last_access_time   file_time
	last_write_time    file_time
	file_size_high     DWORD
	file_size_low      DWORD
	filename           [260]byte
	alternate_filename [14]byte
	file_type          DWORD
	creator_type       DWORD
	finder_flags       WORD
}

type system_time struct {
	Year         WORD
	Month        WORD
	DayOfWeek    WORD
	Day          WORD
	Hour         WORD
	Minute       WORD
	Second       WORD
	Milliseconds WORD
}

func main() {
	lib := syscall.NewLazyDLL("kernel32.dll")
	findfirstfile := lib.NewProc("FindFirstFileA")
	findnextfile := lib.NewProc("FindNextFileA")
	filetime_to_systemtime := lib.NewProc("FileTimeToSystemTime")
	filename := "V:\\hivemind\\rpc\\*"
	byte_value, _ := syscall.BytePtrFromString(filename)
	var find_data win32_find_data
	handle, _, _ := findfirstfile.Call(
		uintptr(unsafe.Pointer(byte_value)),
		uintptr(unsafe.Pointer(&find_data)),
	)
	if handle == uintptr(syscall.InvalidHandle) {
		fmt.Println("FindFirstFile failed")
		return
	}
	var syst system_time
	filetime_to_systemtime.Call(
		uintptr(unsafe.Pointer(&find_data.creation_time)),
		uintptr(unsafe.Pointer(&syst)),
	)
	fmt.Println(string(bytes.Trim(find_data.filename[:], "\x00")))
	//fmt.Printf("%d %d %d (%d:%d)\n", syst.Year, syst.Month, syst.Day, syst.Hour, syst.Minute)

	//fmt.Println(unsafe.Pointer(handle))
	//fmt.Println(find_data)
	//fmt.Println(string(find_data.filename[:]))
	for i := 0; i < 10; i++ {
		var the_find win32_find_data
		findnextfile.Call(
			handle,
			uintptr(unsafe.Pointer(&the_find)),
		)
		//fmt.Println(ret)
		fmt.Println(string(bytes.Trim(the_find.filename[:], "\x00")))
	}
	//last_err, _, _ := lib.NewProc("GetLastError").Call()
	//fmt.Printf("last %d\n", last_err)
}
