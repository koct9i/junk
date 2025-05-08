package main

import (
	"bufio"
	"context"
	"os"
	"io"
	"os/signal"

	"github.com/urfave/cli/v3"

	"golang.org/x/sys/unix"
)

func status() error {
	file, err := os.Open("/proc/self/status")
	if err != nil {
		return err
	}
	defer file.Close()
	file.ReadAt()

	size := 64 << 10
	var reader *bufio.Reader
	for {
		reader := bufio.NewReaderSize(file, size)
		_, err := reader.Peek(size)
		if err != nil {
			if err == bufio.ErrBufferFull {
				size *= 2
				_, err = file.Seek(0, io.SeekStart)
				if err != nil {
					return err
				}
				continue
			}
			return err
		}
		
	}

}

func main() {
	command := &cli.Command{}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := command.Run(ctx, os.Args); err != nil {
		os.Exit(1)
	}
}
