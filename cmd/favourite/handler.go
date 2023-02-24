package main

import (
	"TikV/cmd/favourite/service"
	"TikV/dal/pack"
	favorite "TikV/kitex_gen/favorite"
	"TikV/pkg/errno"
	"context"
)

// FavoriteSrvImpl implements the last service interface defined in the IDL.
type FavoriteSrvImpl struct{}

// FavoriteAction implements the FavoriteSrvImpl interface.
func (s *FavoriteSrvImpl) FavoriteAction(ctx context.Context, req *favorite.DouyinFavoriteActionRequest) (resp *favorite.DouyinFavoriteActionResponse, err error) {
	// TODO: Your code here...
	claim, err := Jwt.ParseToken(req.Token)
	if err != nil {
		resp = pack.BuildFavoriteActionResp(errno.ErrTokenInvalid)
		return resp, nil
	}

	if req.UserId == 0 || claim.Id != 0 {
		req.UserId = claim.Id
	}

	if req.ActionType != 1 && req.ActionType != 2 || req.UserId == 0 || req.VideoId == 0 {
		resp = pack.BuildFavoriteActionResp(errno.ErrBind)
		return resp, nil
	}

	err = service.NewFavoriteActionService(ctx).FavoriteAction(req)
	if err != nil {
		resp = pack.BuildFavoriteActionResp(err)
		return resp, nil
	}
	resp = pack.BuildFavoriteActionResp(err)
	return resp, nil
}

// FavoriteList implements the FavoriteSrvImpl interface.
func (s *FavoriteSrvImpl) FavoriteList(ctx context.Context, req *favorite.DouyinFavoriteListRequest) (resp *favorite.DouyinFavoriteListResponse, err error) {
	claim, err := Jwt.ParseToken(req.Token)
	if err != nil {
		resp = pack.BuildFavoriteListResp(errno.ErrTokenInvalid)
		return resp, nil
	}

	if req.UserId == 0 || claim.Id == 0 {
		resp = pack.BuildFavoriteListResp(errno.ErrBind)
		return resp, nil
	}

	videos, err := service.NewFavoriteListService(ctx).FavoriteList(req)
	if err != nil {
		resp = pack.BuildFavoriteListResp(err)
		return resp, nil
	}
	resp = pack.BuildFavoriteListResp(errno.Success)
	resp.VideoList = videos
	return resp, nil
}
