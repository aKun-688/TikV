package handlers

import (
	"TikV/cmd/api/rpc"
	"TikV/dal/pack"
	"TikV/kitex_gen/comment"
	"TikV/pkg/errno"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"strconv"
)

// 传递 评论操作 的上下文至 Comment 服务的 RPC 客户端, 并获取相应的响应.
func CommentAction(ctx context.Context, c *app.RequestContext) {
	var paramVar CommentActionParam
	token := c.Query("token")
	video_id := c.Query("video_id")
	action_type := c.Query("action_type")

	vid, err := strconv.Atoi(video_id)
	if err != nil {
		SendResponse(c, pack.BuildCommentActionResp(errno.ErrBind))
		return
	}
	act, err := strconv.Atoi(action_type)
	if err != nil {
		SendResponse(c, pack.BuildCommentActionResp(errno.ErrBind))
		return
	}

	paramVar.Token = token
	paramVar.VideoId = int64(vid)
	paramVar.ActionType = int32(act)

	rpcReq := comment.DouyinCommentActionRequest{
		VideoId:    paramVar.VideoId,
		Token:      paramVar.Token,
		ActionType: paramVar.ActionType,
	}

	if act == 1 {
		comment_text := c.Query("comment_text")
		rpcReq.CommentText = &comment_text
	} else {
		comment_id := c.Query("comment_id")
		cid, err := strconv.Atoi(comment_id)
		if err != nil {
			SendResponse(c, pack.BuildCommentActionResp(errno.ErrBind))
			return
		}
		cid64 := int64(cid)
		rpcReq.CommentId = &cid64
	}

	resp, err := rpc.CommentAction(ctx, &rpcReq)
	if err != nil {
		SendResponse(c, pack.BuildCommentActionResp(errno.ConvertErr(err)))
		return
	}
	SendResponse(c, resp)
}

// 传递 获取评论列表操作 的上下文至 Comment 服务的 RPC 客户端, 并获取相应的响应.
func CommentList(ctx context.Context, c *app.RequestContext) {
	var paramVar CommentListParam
	videoid, err := strconv.Atoi(c.Query("video_id"))
	if err != nil {
		SendResponse(c, pack.BuildCommentListResp(errno.ErrBind))
		return
	}
	paramVar.VideoId = int64(videoid)
	paramVar.Token = c.Query("token")

	if len(paramVar.Token) == 0 || paramVar.VideoId < 0 {
		SendResponse(c, pack.BuildCommentListResp(errno.ErrBind))
		return
	}

	resp, err := rpc.CommentList(ctx, &comment.DouyinCommentListRequest{
		VideoId: paramVar.VideoId,
		Token:   paramVar.Token,
	})
	if err != nil {
		SendResponse(c, pack.BuildCommentListResp(errno.ConvertErr(err)))
		return
	}
	SendResponse(c, resp)
}
