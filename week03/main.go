package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		// _, _ := w.Write([]byte("test"))
		_, _ = w.Write([]byte("pong"))
	})

	serverOut := make(chan struct{})

	mux.HandleFunc("shutdown", func(w http.ResponseWriter, r *http.Request) {
		serverOut <- struct{}{}
	})

	server := http.Server{
		Handler: mux,
		Addr:    "8080",
	}

	g.Go(func() error {
		err := server.ListenAndServe()
		if err != nil {
			log.Println("g1 error,will exit.", err.Error())
		}

		return err
	})

	g.Go(func() error {
		select {
		case <-ctx.Done():
			log.Println("g2 errgroup exit...")
		case <-serverOut:
			log.Println("g2, request `/shutdown`, server will out...")
		}

		timeoutCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

		defer cancel()

		err := server.Shutdown(timeoutCtx)
		log.Println("shutting down server...")
		return err
	})

	g.Go(func() error {
		// quit := make(chan os.Signal, 0)
		// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		quit := make(chan os.Signal, 0)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-ctx.Done():
			log.Println("g3, ctx execute cancel...")
			log.Println("g3 error,", ctx.Err().Error())
			// 当g2退出时，已经有错误了，此时的error 并不会覆盖到g中
			return ctx.Err()
		case sig := <-quit:
			return fmt.Errorf("g3 get os signal: %v", sig)
		}
	})

	// g.Wait 等待所有 go执行完毕后执行
	fmt.Printf("end, errgroup exiting, %+v\n", g.Wait())
}
