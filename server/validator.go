package server

type VerifyFunc func(nodeKey string) (bool, error)
