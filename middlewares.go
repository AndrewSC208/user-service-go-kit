package users

import (
	"context"
	"time"

	kitlog "github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

func LoggingMiddleware(logger kitlog.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{next, logger}
	}
}

type loggingMiddleware struct {
	Service
	logger kitlog.Logger
}

func (mw loggingMiddleware) PostUser(ctx context.Context, u User) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostUser", "username", u.Username, "took", time.Since(begin), "err", err)
	}(time.Now())

	return mw.Service.PostUser(ctx, u)
}

func (mw loggingMiddleware) GetUser(ctx context.Context, username string) (u User, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetUser", "username", username, "took", time.Since(begin), "err", err)
	}(time.Now())

	return mw.Service.GetUser(ctx, username)
}

func (mw loggingMiddleware) PutUser(ctx context.Context, username string, u User) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PutProfile", "username", username, "took", time.Since(begin), "err", err)
	}(time.Now())

	return mw.Service.PutUser(ctx, username, u)
}

func (mw loggingMiddleware) PatchUser(ctx context.Context, username string, u User) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PatchProfile", "username", username, "took", time.Since(begin), "err", err)
	}(time.Now())

	return mw.Service.PatchUser(ctx, username, u)
}

func (mw loggingMiddleware) DeleteProfile(ctx context.Context, username string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "DeleteProfile", "username", username, "took", time.Since(begin), "err", err)
	}(time.Now())

	return mw.Service.DeleteUser(ctx, username)
}