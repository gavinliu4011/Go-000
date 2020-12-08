package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"week03/pkg/sync/errgroup"
)

type Server struct {
	srv *http.Server
}

func NewServer(addr string, router *http.ServeMux) *Server {
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}
	return &Server{srv: srv}
}

func (s *Server) Start() error {
	log.Println("开始监听端口。。。。。。")
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("开始关闭连接。。。。。。")
	return s.srv.Shutdown(ctx)
}

func main() {
	// 路由
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("进来了。。。。。。。")
		w.WriteHeader(http.StatusOK)
	})
	svr := NewServer(":8080", router)

	done := make(chan struct{}, 1)
	quit := make(chan os.Signal, 1)
	// 监听到指定信号就给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	g := errgroup.WithCancel(context.Background())

	g.Go(func(ctx context.Context) error {
		for {
			log.Println("等待linux signal 信号。。。。。。")
			select {
			case <-quit:
				log.Println("键入了一个linux signal 信号。。。。")
				return errors.New("什么导致了退出😱")
			case <-ctx.Done():
				// 这里也需要ctx.Done()不然其他goroutine报错时无法退出
				return ctx.Err()
			}
		}
	})

	g.Go(func(ctx context.Context) error {
		// return errors.New("假设")
		return svr.Start()
	})

	g.Go(func(ctx context.Context) error {
		<-ctx.Done()
		log.Println("ctx 被取消了。。。。。")
		ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := svr.Shutdown(ctx2); err != nil {
			log.Println(err)
		}
		done <- struct{}{}
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Println("捕捉到一枚错误", err)
	}
	<-done
	log.Println("服务完成.......")
}
