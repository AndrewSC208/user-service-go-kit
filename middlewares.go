package users

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) PostProfile(ctx context.Context, u User) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostProfile", "id", u.Username, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostUser(ctx, u)
}

func (mw loggingMiddleware) GetProfile(ctx context.Context, username string) (u User, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetProfile", "username", username, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetUser(ctx, username)
}

func (mw loggingMiddleware) PutProfile(ctx context.Context, username string, u User) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PutProfile", "username", username, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PutUser(ctx, username, u)
}

func (mw loggingMiddleware) PatchProfile(ctx context.Context, username string, u User) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PatchProfile", "username", username, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PatchUser(ctx, username, u)
}

func (mw loggingMiddleware) DeleteProfile(ctx context.Context, username string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "DeleteProfile", "username", username, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.DeleteUser(ctx, username)
}