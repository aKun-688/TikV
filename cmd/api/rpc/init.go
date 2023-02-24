package rpc

import "TikV/pkg/ttviper"

// InitRPC init rpc client
func InitRPC(Config *ttviper.Config) {
	UserConfig := ttviper.ConfigInit("TIKV_USER", "userConfig")
	initUserRpc(&UserConfig)

	FeedConfig := ttviper.ConfigInit("TIKV_FEED", "feedConfig")
	initFeedRpc(&FeedConfig)

	PublishConfig := ttviper.ConfigInit("TIKV_PUBLISH", "publishConfig")
	initPublishRpc(&PublishConfig)

	FavoriteConfig := ttviper.ConfigInit("TIKV_FAVORITE", "favoriteConfig")
	initFavoriteRpc(&FavoriteConfig)

	CommentConfig := ttviper.ConfigInit("TIKV_COMMENT", "commentConfig")
	initCommentRpc(&CommentConfig)

	RelationConfig := ttviper.ConfigInit("TIKV_RELATION", "relationConfig")
	initRelationRpc(&RelationConfig)
}
