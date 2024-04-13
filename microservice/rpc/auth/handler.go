package main

import (
	auth "rpc.xjco2913/kitex_gen/rpc/xjco2913/auth"
	"context"
)

// LoginServiceImpl implements the last service interface defined in the IDL.
type LoginServiceImpl struct{}

// Login implements the LoginServiceImpl interface.
func (s *LoginServiceImpl) Login(ctx context.Context, req *auth.LoginReq) (resp *auth.LoginResp, err error) {
	// TODO: Your code here...
	return
}
