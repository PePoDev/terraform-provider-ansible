package ansible

import (
	"embed"
	"log"
	"syscall"
	"unsafe"
)

//go:embed ansible
var ansibleFs embed.FS

func init() {
	ansibleFs.ReadFile("hello.txt")
}

func run() {
	fd, err := MemFsCreate("/file.bin")
	if err != nil {
		log.Fatal(err)
	}

	err = CopyToMem(fd, ansibleFs)
	if err != nil {
		log.Fatal(err)
	}

	err = ExecveAt(fd)
	if err != nil {
		log.Fatal(err)
	}
}

func ExecveAt(fd uintptr) (err error) {
	s, err := syscall.BytePtrFromString("")
	if err != nil {
		return err
	}
	ret, _, errno := syscall.Syscall6(322, fd, uintptr(unsafe.Pointer(s)), 0, 0, 0x1000, 0)
	if int(ret) == -1 {
		return errno
	}

	// never hit
	log.Println("should never hit")
	return err
}

func CopyToMem(fd uintptr, buf []byte) (err error) {
	_, err = syscall.Write(int(fd), buf)
	if err != nil {
		return err
	}

	return nil
}

func MemFsCreate(path string) (r1 uintptr, err error) {
	s, err := syscall.BytePtrFromString(path)
	if err != nil {
		return 0, err
	}

	r1, _, errno := syscall.Syscall(319, uintptr(unsafe.Pointer(s)), 0, 0)

	if int(r1) == -1 {
		return r1, errno
	}

	return r1, nil
}
