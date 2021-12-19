package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

/*
基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，
要保证能够一个退出，全部注销退出
 */


var signals = []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT}

func Run(srv *http.Server, ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	eg, errCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return Start(srv)
	})
	eg.Go(func() error {
		<-errCtx.Done()
		return Stop(srv, errCtx)
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)
	eg.Go(func() error {
		for {
			select {
			case <-errCtx.Done():
				return errCtx.Err()
			case <-c:
				cancel()
			}
		}
	})
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}


func Start(srv *http.Server) error {
	http.HandleFunc("/hello", hello)
	err := srv.ListenAndServe()
	return err

}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Golang!")
}


func Stop(srv *http.Server, ctx context.Context) error {
	return srv.Shutdown(ctx)
}

func main() {
	ctx := context.Background()
	server := &http.Server{Addr: ":9090"}
	err := Run(server, ctx)
	if err != nil {
		fmt.Printf("got some errors : %v", err)
	}
}