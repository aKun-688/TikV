package handlers

import (
	"TikV/cmd/api/rpc"
	"TikV/dal/pack"
	"TikV/kitex_gen/relation"
	"TikV/pkg/errno"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"strconv"
)

// 传递 关注操作 的上下文至 Relation 服务的 RPC 客户端, 并获取相应的响应.
func RelationAction(ctx context.Context, c *app.RequestContext) {
	var paramVar RelationActionParam
	token := c.Query("token")
	to_user_id := c.Query("to_user_id")
	action_type := c.Query("action_type")

	tid, err := strconv.Atoi(to_user_id)
	if err != nil {
		SendResponse(c, pack.BuildRelationActionResp(errno.ErrBind))
		return
	}

	act, err := strconv.Atoi(action_type)
	if err != nil {
		SendResponse(c, pack.BuildRelationActionResp(errno.ErrBind))
		return
	}

	paramVar.Token = token
	paramVar.ToUserId = int64(tid)
	paramVar.ActionType = int32(act)

	rpcReq := relation.DouyinRelationActionRequest{
		ToUserId:   paramVar.ToUserId,
		Token:      paramVar.Token,
		ActionType: paramVar.ActionType,
	}

	resp, err := rpc.RelationAction(ctx, &rpcReq)
	if err != nil {
		SendResponse(c, pack.BuildRelationActionResp(errno.ConvertErr(err)))
		return
	}
	SendResponse(c, resp)
}

// 传递 获取正在关注列表操作 的上下文至 Relation 服务的 RPC 客户端, 并获取相应的响应.
func RelationFollowList(ctx context.Context, c *app.RequestContext) {
	var paramVar UserParam
	uid, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		SendResponse(c, pack.BuildFollowingListResp(errno.ErrBind))
		return
	}
	paramVar.UserId = int64(uid)
	paramVar.Token = c.Query("token")

	if len(paramVar.Token) == 0 || paramVar.UserId < 0 {
		SendResponse(c, pack.BuildFollowingListResp(errno.ErrBind))
		return
	}

	resp, err := rpc.RelationFollowList(ctx, &relation.DouyinRelationFollowListRequest{
		UserId: paramVar.UserId,
		Token:  paramVar.Token,
	})
	if err != nil {
		SendResponse(c, pack.BuildFollowingListResp(errno.ConvertErr(err)))
		return
	}
	SendResponse(c, resp)
}

// 传递 获取粉丝列表操作 的上下文至 Relation 服务的 RPC 客户端, 并获取相应的响应.
func RelationFollowerList(ctx context.Context, c *app.RequestContext) {
	var paramVar UserParam
	uid, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		SendResponse(c, pack.BuildFollowerListResp(errno.ErrBind))
		return
	}
	paramVar.UserId = int64(uid)
	paramVar.Token = c.Query("token")

	if len(paramVar.Token) == 0 || paramVar.UserId < 0 {
		SendResponse(c, pack.BuildFollowerListResp(errno.ErrBind))
		return
	}

	resp, err := rpc.RelationFollowerList(ctx, &relation.DouyinRelationFollowerListRequest{
		UserId: paramVar.UserId,
		Token:  paramVar.Token,
	})
	if err != nil {
		SendResponse(c, pack.BuildFollowerListResp(errno.ConvertErr(err)))
		return
	}
	SendResponse(c, resp)
}
