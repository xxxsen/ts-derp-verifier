package validator

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
	"ts-derp-verifier/tailscale"

	"github.com/samber/lo"
	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

type TailscaleVerifier struct {
	client   *tailscale.Client
	interval time.Duration
	res      atomic.Value
}

func NewTailscaleVerifier(c *tailscale.Client, interval time.Duration) (*TailscaleVerifier, error) {
	if interval == 0 {
		interval = 10 * time.Minute
	}
	v := &TailscaleVerifier{client: c, interval: interval}
	if err := v.init(); err != nil {
		return nil, fmt.Errorf("init verifier failed, err:%w", err)
	}
	return v, nil
}

func (s *TailscaleVerifier) init() error {
	if err := s.readOnce(); err != nil {
		return fmt.Errorf("init first failed, err:%w", err)
	}
	go s.startLoop()
	return nil
}

func (s *TailscaleVerifier) startLoop() {
	ticker := time.NewTicker(s.interval)
	for range ticker.C {
		if err := s.readOnce(); err != nil {
			logutil.GetLogger(context.Background()).Error("refresh tailscale client failed", zap.Error(err))
			continue
		}
	}
}

func (s *TailscaleVerifier) readOnce() error {
	ctx := context.Background()
	res, err := s.client.ListDevices(ctx)
	if err != nil {
		return err
	}
	logutil.GetLogger(ctx).Debug("read devices success", zap.Int("device_count", len(res)))
	for _, item := range res {
		logutil.GetLogger(ctx).Debug("recv device", zap.Bool("authorized", item.Authorized), zap.String("node_public", item.NodeKey))
	}
	s.res.Store(res)
	return nil
}

func (s *TailscaleVerifier) Verify(node string) (bool, error) {
	res := s.res.Load()
	if res == nil {
		return false, fmt.Errorf("devices data not init")
	}
	devs := res.([]*tailscale.Device)
	return lo.ContainsBy(devs, func(dev *tailscale.Device) bool {
		return dev.Authorized == true && dev.NodeKey == node
	}), nil
}
