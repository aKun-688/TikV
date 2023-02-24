package main

import (
	"TikV/cmd/user/service"
	"TikV/dal/pack"
	user "TikV/kitex_gen/user"
	"TikV/pkg/errno"
	"TikV/pkg/jwt"
	"context"
)

// UserSrvImpl implements the last service interface defined in the IDL.
type UserSrvImpl struct{}

// Register implements the UserSrvImpl interface.
func (s *UserSrvImpl) Register(ctx context.Context, req *user.DouyinUserRegisterRequest) (resp *user.DouyinUserRegisterResponse, err error) {
	// TODO: Your code here...
	if len(req.Username) == 0 || len(req.Password) == 0 {

		resp = pack.BuilduserRegisterResp(errno.ErrBind)
		return resp, nil
	}

	err = service.NewCreateUserService(ctx).CreateUser(req, Argon2Config)
	if err != nil {
		resp = pack.BuilduserRegisterResp(err)
		return resp, nil
	}

	// 新用户注册成功后直接登录
	uid, err := service.NewCheckUserService(ctx).CheckUser(req)
	if err != nil {
		resp = pack.BuilduserRegisterResp(err)
		return resp, nil
	}

	token, err := Jwt.CreateToken(jwt.CustomClaims{
		Id: int64(uid),
	})
	if err != nil {
		resp = pack.BuilduserRegisterResp(errno.ErrSignatureInvalid)
		return resp, nil
	}

	resp = pack.BuilduserRegisterResp(errno.Success)
	resp.UserId = uid
	resp.Token = token
	return resp, nil
}

// Login implements the UserSrvImpl interface.
func (s *UserSrvImpl) Login(ctx context.Context, req *user.DouyinUserRegisterRequest) (resp *user.DouyinUserRegisterResponse, err error) {
	if len(req.Username) == 0 || len(req.Password) == 0 {
		resp = pack.BuilduserRegisterResp(errno.ErrBind)
		return resp, nil
	}

	uid, err := service.NewCheckUserService(ctx).CheckUser(req)
	if err != nil {
		resp = pack.BuilduserRegisterResp(err)
		return resp, nil
	}

	token, err := Jwt.CreateToken(jwt.CustomClaims{
		Id: int64(uid),
	})
	if err != nil {
		resp = pack.BuilduserRegisterResp(errno.ErrSignatureInvalid)
		return resp, nil
	}

	resp = pack.BuilduserRegisterResp(errno.Success)
	resp.UserId = uid
	resp.Token = token
	return resp, nil
}

// GetUserById implements the UserSrvImpl interface.
func (s *UserSrvImpl) GetUserById(ctx context.Context, req *user.DouyinUserRequest) (resp *user.DouyinUserResponse, err error) {
	claim, err := Jwt.ParseToken(req.Token)
	if err != nil {
		resp = pack.BuilduserUserResp(errno.ErrTokenInvalid)
		return resp, nil
	}
	// else if claim.Id != int64(req.UserId) {
	// 	resp = pack.BuilduserUserResp(errno.ErrValidation)
	// 	return resp, nil
	// }

	if req.UserId < 0 {
		resp = pack.BuilduserUserResp(errno.ErrBind)
		return resp, nil
	}

	user, err := service.NewMGetUserService(ctx).MGetUser(req, claim.Id)
	if err != nil {
		resp = pack.BuilduserUserResp(err)
		return resp, nil
	}

	if claim.Id == req.UserId {
		user.IsFollow = true
	} else {
		// TODO 获取claim.id 是否已关注 req.userid
		user.IsFollow = false
	}

	resp = pack.BuilduserUserResp(errno.Success)
	resp.User = user
	return resp, nil
}
