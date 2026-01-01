package model

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/netip"
	"strings"
)

type NodePublic struct {
	raw string
}

func (p *NodePublic) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		p.raw = strings.TrimSpace(s)
		return nil
	}
	var obj map[string]any
	if err := json.Unmarshal(b, &obj); err != nil {
		return err
	}
	raw, err := parseNodePublicObject(obj)
	if err != nil {
		return err
	}
	p.raw = raw
	return nil
}

func (p NodePublic) String() string {
	return p.raw
}

func parseNodePublicObject(obj map[string]any) (string, error) {
	v, ok := obj["k"]
	if !ok {
		v, ok = obj["K"]
	}
	if !ok {
		return "", errors.New("node public key missing")
	}
	arr, ok := v.([]any)
	if !ok || len(arr) != 32 {
		return "", errors.New("node public key must be 32 bytes")
	}
	buf := make([]byte, 32)
	for i, item := range arr {
		num, ok := item.(float64)
		if !ok || num < 0 || num > 255 {
			return "", errors.New("node public key byte out of range")
		}
		buf[i] = byte(num)
	}
	return "nodekey:" + hex.EncodeToString(buf), nil
}

type DERPAdmitClientRequest struct {
	NodePublic NodePublic
	Source     netip.Addr
}

type DERPAdmitClientResponse struct {
	Allow bool
}
