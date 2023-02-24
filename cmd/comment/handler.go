package main

import (
	"TikV/cmd/comment/service"
	"TikV/dal/pack"
	comment "TikV/kitex_gen/comment"
	"TikV/pkg/errno"
	"context"
)

// CommentSrvImpl implements the last service interface defined in the IDL.
type CommentSrvImpl struct{}

// CommentAction implements the CommentSrvImpl interface.
func (s *CommentSrvImpl) CommentAction(ctx context.Context, req *comment.DouyinCommentActionRequest) (resp *comment.DouyinCommentActionResponse, err error) {
	// TODO: Your code here...
	claim, err := Jwt.ParseToken(req.Token)
	if err != nil {
		resp = pack.BuildCommentActionResp(errno.ErrTokenInvalid)
		return resp, nil
	}

	if req.UserId == 0 || claim.Id != 0 {
		req.UserId = claim.Id
	}

	if req.ActionType != 1 && req.ActionType != 2 || req.UserId == 0 || req.VideoId == 0 {
		resp = pack.BuildCommentActionResp(errno.ErrBind)
		return resp, nil
	}

	err = service.NewCommentActionService(ctx).CommentAction(req)
	if err != nil {
		resp = pack.BuildCommentActionResp(err)
		return resp, nil
	}
	resp = pack.BuildCommentActionResp(err)
	return resp, nil
}

// CommentList implements the CommentSrvImpl interface.
func (s *CommentSrvImpl) CommentList(ctx context.Context, req *comment.DouyinCommentListRequest) (resp *comment.DouyinCommentListResponse, err error) {
	claim, err := Jwt.ParseToken(req.Token)
	if err != nil {
		resp = pack.BuildCommentListResp(errno.ErrTokenInvalid)
		return resp, nil
	}

	if req.VideoId == 0 || claim.Id == 0 {
		resp = pack.BuildCommentListResp(errno.ErrBind)
		return resp, nil
	}

	comments, err := service.NewCommentListService(ctx).CommentList(req, claim.Id)
	if err != nil {
		resp = pack.BuildCommentListResp(err)
		return resp, nil
	}
	resp = pack.BuildCommentListResp(errno.Success)
	resp.CommentList = comments
	return resp, nil
}
