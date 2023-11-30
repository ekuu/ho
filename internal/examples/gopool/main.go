package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ekuu/ho/gopool"
)

func main() {
	runHttpServer()
}

func runHttpServer() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		time.Sleep(3 * time.Second)
		writer.Write([]byte("hello world"))
	})

	// 启动http server
	s := http.Server{Addr: "localhost:8080"}
	gopool.GoErr(func(ctx context.Context) error {
		err := s.ListenAndServe()
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	})

	// 监听是否调用了cancel()，如果调用了则意味着收到程序终止信号，则调用Shutdown，关闭http服务。
	gopool.GoErr(func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				// 等待所有连接完成传输后，优雅关闭
				return s.Shutdown(context.Background())
			}
		}
	})

	// 作为client访问server
	gopool.GoErr(func(ctx context.Context) error {
		// 等待ListenAndServe
		time.Sleep(100 * time.Millisecond)

		resp, err := http.Get("http://localhost:8080/")
		if err != nil {
			return err
		}
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", b)
		fmt.Println("client request end.")
		return nil
	})

	//gopool.GoErr(func(ctx context.Context) error {
	//	return errors.New("test error")
	//})

	gopool.WaitSignal()
}
