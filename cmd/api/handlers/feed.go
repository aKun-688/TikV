package handlers

import (
	"TikV/cmd/api/rpc"
	"TikV/dal/pack"
	"TikV/kitex_gen/feed"
	"TikV/pkg/errno"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"strconv"
)

// GetUserFeed 传递 获取用户视频流操作 的上下文至 Feed 服务的 RPC 客户端, 并获取相应的响应.
func GetUserFeed(ctx context.Context, c *app.RequestContext) {
	var feedVar FeedParam
	var laststTime int64
	var token string
	lastst_time := c.Query("latest_time")
	if len(lastst_time) != 0 {
		if latesttime, err := strconv.Atoi(lastst_time); err != nil {
			SendResponse(c, pack.BuildVideoResp(errno.ErrDecodingFailed))
			return
		} else {
			laststTime = int64(latesttime)
		}
	}

	feedVar.LatestTime = &laststTime

	token = c.Query("token")
	feedVar.Token = &token

	resp, err := rpc.GetUserFeed(ctx, &feed.DouyinFeedRequest{
		LatestTime: feedVar.LatestTime,
		Token:      feedVar.Token,
	})
	if err != nil {
		SendResponse(c, pack.BuildVideoResp(errno.ConvertErr(err)))
		return
	}
	SendResponse(c, resp)
}
