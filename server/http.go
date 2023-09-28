package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"transferSrv/controller"
	"transferSrv/infra/config"
	"transferSrv/infra/library/log"
)

func StartHttpSrv() {

	cnf := config.GetCnf()

	router := controller.Route()
	// router.Run(cnf.WebCnf.Host + ":" + cnf.WebCnf.Port)

	srv := &http.Server{
		Addr:    cnf.WebCnf.Host + ":" + cnf.WebCnf.Port,
		Handler: router,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	//等待中断信号关闭服务器
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Info("Server exiting")

}
