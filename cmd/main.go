package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/xxxsen/common/logger"
	"github.com/xxxsen/ts-derp-verifier/config"
	"github.com/xxxsen/ts-derp-verifier/server"
	"github.com/xxxsen/ts-derp-verifier/tailscale"
	"github.com/xxxsen/ts-derp-verifier/validator"
	"go.uber.org/zap"
)

var (
	cfg = flag.String("config", "./config.json", "config file")
)

func main() {
	flag.Parse()
	c, err := config.Load(*cfg)
	if err != nil {
		log.Fatalf("read config failed, file:%s, err:%v", *cfg, err)
	}

	logkit := logger.Init(c.Log.File, c.Log.Level, int(c.Log.FileCount), int(c.Log.FileSize), int(c.Log.KeepDays), c.Log.Console)

	ts := tailscale.NewClient(c.Tailnet, c.APIKey)
	verifier, err := validator.NewTailscaleVerifier(ts, time.Duration(c.RefreshInterval)*time.Second)
	if err != nil {
		logkit.Fatal("init verifier failed", zap.Error(err))
	}

	svc, err := server.New(server.WithBind(c.Listen), server.WithVerifier(verifier.Verify))
	if err != nil {
		logkit.Fatal("init verify server failed", zap.Error(err))
	}
	logkit.Info("start verify server", zap.String("bind", c.Listen))
	if err := svc.Run(context.Background()); err != nil {
		logkit.Fatal("run server failed", zap.Error(err))
	}
}
