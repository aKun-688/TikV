package pack

import (
	"TikV/dal/db"
	"TikV/kitex_gen/user"
	"context"
	"errors"
	"gorm.io/gorm"
)

// FollowingList pack lists of following info.
func FollowingList(ctx context.Context, vs []*db.Relation, fromID int64) ([]*user.User, error) {
	users := make([]*db.User, 0)
	for _, v := range vs {
		user2, err := db.GetUerByID(ctx, int64(v.ToUserID))
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		users = append(users, user2)
	}

	return Users(ctx, users, fromID)
}

// FollowerList pack lists of follower info.
func FollowerList(ctx context.Context, vs []*db.Relation, fromID int64) ([]*user.User, error) {
	users := make([]*db.User, 0)
	for _, v := range vs {
		user2, err := db.GetUerByID(ctx, int64(v.UserID))
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		users = append(users, user2)
	}

	return Users(ctx, users, fromID)
}
