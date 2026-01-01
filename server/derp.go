package server

import (
	"context"
	"fmt"
	"net/http"

	"ts-derp-verifier/model"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/logutil"
	"github.com/xxxsen/common/webapi"
	"go.uber.org/zap"
)

type VerifyServer struct {
	c *config
}

func New(opts ...Option) (*VerifyServer, error) {
	c := &config{}
	for _, opt := range opts {
		opt(c)
	}
	if len(c.addr) == 0 {
		return nil, fmt.Errorf("no bind addr")
	}
	if c.verifyFn == nil {
		return nil, fmt.Errorf("no verify fn")
	}
	return &VerifyServer{c: c}, nil
}

func (s *VerifyServer) Run(ctx context.Context) error {
	engine, err := webapi.NewEngine("/", s.c.addr, webapi.WithRegister(s.register))
	if err != nil {
		return fmt.Errorf("init engine failed, err:%w", err)
	}
	return engine.Run()
}

func (s *VerifyServer) register(c *gin.RouterGroup) {
	c.POST("/derp/verify", s.verifyHandler)
}

func (s *VerifyServer) verifyHandler(c *gin.Context) {
	req := &model.DERPAdmitClientRequest{}
	if err := c.ShouldBindBodyWithJSON(req); err != nil {
		logutil.GetLogger(c).Error("bind body failed", zap.Error(err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	allow, err := s.c.verifyFn(req.NodePublic.String())
	if err != nil {
		logutil.GetLogger(c).Error("verify node public failed", zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	resp := &model.DERPAdmitClientResponse{Allow: allow}
	c.JSON(http.StatusOK, resp)
}
