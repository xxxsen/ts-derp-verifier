package server

type config struct {
	addr     string
	verifyFn VerifyFunc
}

type Option func(c *config)

func WithBind(addr string) Option {
	return func(c *config) {
		c.addr = addr
	}
}

func WithVerifier(v VerifyFunc) Option {
	return func(c *config) {
		c.verifyFn = v
	}
}
