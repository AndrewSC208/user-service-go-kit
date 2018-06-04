package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"

	svc "github.com/AndrewSC208/user-service-go-kit"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
		storeUrl = flag.String("db.url", "postgresql://root@localhost:26257/bank?sslmode=disable", "STORE db url")
	)
	flag.Parse()

	db, err := gorm.Open("postgres", storeUrl)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&svc.UserModel{})

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var s svc.Service
	{
		// create new service, and pass store in
		s = svc.NewService(*db)

		// Setup logging
		s = svc.LoggingMiddleware(logger)(s)
	}

	var h http.Handler
	{
		h = svc.MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}