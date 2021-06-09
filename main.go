package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

type BlackHole struct {
	total int64
}

func (b *BlackHole) increment(size int64) {
	b.total += size
}

func (b BlackHole) print() {
	CallClear()
	fmt.Printf("Black Hole has eaten %d KB", b.total>>10)
}

func DirSize(path string) (int64, error) {

	var size int64

	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if !info.IsDir() {
			size += info.Size()
		}

		return err
	})

	return size, err
}

func RemoveContents(dir string) error {

	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	eatSize, err := DirSize(dir)
	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			continue
		}
	}

	blackHole.increment(eatSize)
	blackHole.print()

	return nil
}

var (
	blackHole = &BlackHole{}
	clear     map[string]func() //create a map for storing clear funcs
)

func main() {

	go func() {

		c := make(chan os.Signal, 1)

		signals := []os.Signal{
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGKILL,
			syscall.SIGTERM,
			syscall.SIGABRT,
			syscall.SIGQUIT,
		}

		signal.Notify(c, signals...)
		<-c

		os.Remove("hole")

		os.Exit(0)
	}()

	if _, err := os.Stat("hole"); os.IsNotExist(err) {
		os.Mkdir("hole", 0755)
	}

	ticker := time.NewTicker(500 * time.Millisecond)

	for range ticker.C {
		if err := RemoveContents("hole"); err != nil {
			panic(err)
		}
	}
}
