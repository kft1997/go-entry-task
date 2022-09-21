package util

import "errors"

var (
	ErrTokenNotExist   = errors.New("token doesn't exist")
	ErrInvalidPassword = errors.New("passwd incorrect")
	ErrUpdateDB        = errors.New("update database error")
	ErrDelRedis        = errors.New("delete redis error")
	ErrRedisNotExist   = errors.New("info doesn't exist in redis")
)
