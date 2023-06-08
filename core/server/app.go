package server

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type AppInterface interface {
	rpcServer()
	httpServer()
	run() error
}

type AppEngine struct {
	gs     *grpc.Server
	hs     *http.Server
	gsAddr string
}

func New(gs *grpc.Server, hs *http.Server) *AppEngine {
	return &AppEngine{
		gs: gs,
		hs: hs,
	}
}

func (b *AppEngine) rpcServer() {
	reflection.Register(b.gs)
	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	}

	if err := b.gs.Serve(l); err != nil {
		log.Fatalf("listen: %s\n", err)
	}
}

func (b *AppEngine) httpServer() {
	if err := b.hs.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func (b *AppEngine) Run() error {
	//开启
	go b.httpServer()

	go b.rpcServer()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := b.hs.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
	return nil
}
