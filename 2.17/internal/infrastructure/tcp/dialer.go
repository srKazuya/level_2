// Package tcp provides...
package tcp

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"time"
)

type Dialer interface {
	Connect(timeout int, host, port string)
	Read(ctx context.Context, errChan chan error, conn net.Conn)
	Write(ctx context.Context, cancel context.CancelFunc, errChan chan error, conn net.Conn)
}

type dialer struct {
	log *slog.Logger
}

func NewDialer(log *slog.Logger) Dialer {
	return &dialer{log: log}
}

var (
	errFailedDial    = errors.New("ошибка подключения: ")
	errReqFailed     = errors.New("ошибка запроса: ")
	errCopyConnStdin = errors.New("ошибка копирования из соединнения: ")
	errSendMsg       = errors.New("ошибка отправки сообщения: ")
)

// Connect implements Dialer.
func (d *dialer) Connect(t int, h string, p string) {
	op := "inf.TCP.Dialer.Run"
	log := d.log.With(slog.String("op", op))
	address := h + ":" + p

	var conn net.Conn
	var err error

	if t <= 0 {
		conn, err = net.Dial("tcp", address)
	} else {
		conn, err = net.DialTimeout("tcp", address, time.Duration(t)*time.Second)
	}

	if err != nil {
		log.Error(errFailedDial.Error(), slog.Any("error", err))
		return
	}
	defer conn.Close()

	log.Info("соединение установлено", slog.Any("addr:", address))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errReadCh := make(chan error, 1)
	errWriteCh := make(chan error, 1)

	// Параллельно читаем и пишем
	go d.Read(ctx, errReadCh, conn)
	go d.Write(ctx, cancel, errWriteCh, conn)

	select {
	case err := <-errReadCh:
		if err != nil {
			log.Error("ошибка чтения", slog.Any("err", err))
		}
	case err := <-errWriteCh:
		if err != nil {
			log.Error("ошибка записи", slog.Any("err", err))
		}
	case <-ctx.Done():
		log.Info("соединение закрыто", slog.Any("addr:", address))
	}
}

func (d *dialer) Read(ctx context.Context, errCh chan error, c net.Conn) {
	_, err := io.Copy(os.Stdout, c)
	if err != nil {
		errCh <- fmt.Errorf("%w, %v", errCopyConnStdin, err)
	} else {
		errCh <- nil
	}
}

func (d *dialer) Write(ctx context.Context, cancel context.CancelFunc, errCh chan error, c net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			fmt.Println("EOF Ctrl+D")
			cancel()
			c.Close()
			return
		}
		if err != nil {
			errCh <- fmt.Errorf("%w, %v", errCopyConnStdin, err)
			return
		}

		_, err = c.Write([]byte(line))
		if err != nil {
			errCh <- fmt.Errorf("%w, %v", errSendMsg, err)
			return
		}
	}
}

