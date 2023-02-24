package main

import (
	"TikV/cmd/relation/service"
	"TikV/dal/pack"
	relation "TikV/kitex_gen/relation"
	"TikV/pkg/errno"
	"context"
)

// RelationSrvImpl implements the last service interface defined in the IDL.
type RelationSrvImpl struct{}

// RelationAction implements the RelationSrvImpl interface.
func (s *RelationSrvImpl) RelationAction(ctx context.Context, req *relation.DouyinRelationActionRequest) (resp *relation.DouyinRelationActionResponse, err error) {
	// TODO: Your code here...
	claim, err := Jwt.ParseToken(req.Token)
	if err != nil {
		resp = pack.BuildRelationActionResp(errno.ErrTokenInvalid)
		return resp, nil
	}

	if req.UserId == 0 || claim.Id != 0 {
		req.UserId = claim.Id
	}

	if req.ActionType < 1 || req.ActionType > 2 {
		resp = pack.BuildRelationActionResp(errno.ErrBind)
		return resp, nil
	}
	err = service.NewRelationActionService(ctx).RelationAction(req)
	if err != nil {
		resp = pack.BuildRelationActionResp(err)
		return resp, nil
	}
	resp = pack.BuildRelationActionResp(errno.Success)
	return resp, nil
}

// RelationFollowList implements the RelationSrvImpl interface.
func (s *RelationSrvImpl) RelationFollowList(ctx context.Context, req *relation.DouyinRelationFollowListRequest) (resp *relation.DouyinRelationFollowListResponse, err error) {
	claim, err := Jwt.ParseToken(req.Token)
	if err != nil {
		resp = pack.BuildFollowingListResp(errno.ErrTokenInvalid)
		return resp, nil
	}

	if req.UserId == 0 || claim.Id != 0 {
		req.UserId = claim.Id // 没有传入UserID，默认为自己
	}

	users, err := service.NewFollowingListService(ctx).FollowingList(req, claim.Id)
	if err != nil {
		resp = pack.BuildFollowingListResp(err)
		return resp, nil
	}

	resp = pack.BuildFollowingListResp(errno.Success)
	resp.UserList = users
	return resp, nil
}

// RelationFollowerList implements the RelationSrvImpl interface.
func (s *RelationSrvImpl) RelationFollowerList(ctx context.Context, req *relation.DouyinRelationFollowerListRequest) (resp *relation.DouyinRelationFollowerListResponse, err error) {
	claim, err := Jwt.ParseToken(req.Token)
	if err != nil {
		resp = pack.BuildFollowerListResp(errno.ErrTokenInvalid)
		return resp, nil
	}

	if req.UserId == 0 {
		req.UserId = claim.Id // 没有传入UserID，默认为自己
	}

	users, err := service.NewFollowerListService(ctx).FollowerList(req, claim.Id)
	if err != nil {
		resp = pack.BuildFollowerListResp(err)
		return resp, nil
	}

	resp = pack.BuildFollowerListResp(errno.Success)
	resp.UserList = users
	return resp, nil
}
