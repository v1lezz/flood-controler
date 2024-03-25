package main

import (
	"context"
	"github.com/caarlos0/env/v10"
	"golang.org/x/sync/errgroup"
	"log"
	"os/signal"
	"syscall"
	"task/internal/app"
	"task/internal/config"
	http_port "task/internal/ports/http"
	"task/internal/redis_checker"
)

func main() {
	srvCFG := config.ServerConfig{}
	if err := env.Parse(&srvCFG); err != nil {
		log.Fatal("err parse server config")
	}
	log.Printf("SERVER CONFIG:%+v\n", srvCFG)
	redisCFG := config.RedisConfig{}
	if err := env.Parse(&redisCFG); err != nil {
		log.Fatal("err parse redis config")
	}
	log.Printf("REDIS CONFIG:%+v\n", redisCFG)
	chCFG := config.CheckerConfig{}
	if err := env.Parse(&chCFG); err != nil {
		log.Fatal("err parse checker config")
	}
	log.Printf("CHECKER CONFIG:%+v\n", chCFG)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	eg, gCtx := errgroup.WithContext(ctx)
	httpSrv := http_port.NewServer(srvCFG.Port)
	rch := redis_checker.NewRedisChecker(redisCFG, chCFG, eg)
	a, err := app.NewApp(httpSrv, rch)
	if err != nil {
		log.Fatal(err)
	}
	http_port.AppRouter(a, gCtx)
	if err = a.Run(eg, gCtx); err != nil {
		log.Fatal(err)
	}
}
