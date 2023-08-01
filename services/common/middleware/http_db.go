package middleware

import (
	"database/sql"
	"errors"
	"fmt"
	"smartkid/services/common/infra/db"
	"smartkid/services/common/infra/logger"

	appctx "smartkid/services/common/context"

	"github.com/gin-gonic/gin"
)

type (
	Handler        = func(c *gin.Context) error
	MiddlewareDbTx = func(handler Handler, rollback ...Handler) Handler
	MiddlewareDb   = func(handler Handler) Handler
)

func HttpDb(db db.SQL, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" {
			middleware := WithDb(db, log)
			handler := func(c *gin.Context) (err error) {
				c.Next()
				return nil
			}
			_ = middleware(handler)(c)

		} else {
			middleware := WithTx(db, log)
			handler := func(c *gin.Context) (err error) {
				log.Debug("handler")
				c.Next()
				return nil
			}
			rollback := func(c *gin.Context) error {
				log.Debug("rollback")
				return nil
			}
			_ = middleware(handler, rollback)(c)
		}
	}
}

func WithDb(db db.SQL, log logger.Logger, options ...*sql.TxOptions) MiddlewareDb {
	return func(handler Handler) Handler {
		return func(c *gin.Context) error {
			ctx := appctx.FromContext(c)
			ctx.WithDb(db.GetDB())
			return handler(c)
		}
	}
}

func WithTx(db db.SQL, log logger.Logger, options ...*sql.TxOptions) MiddlewareDbTx {
	return func(handler Handler, rollback ...Handler) Handler {
		return func(ginCtx *gin.Context) error {
			var err error
			var opt *sql.TxOptions
			if len(options) > 0 {
				opt = options[0]
			}
			var rb Handler
			if len(rollback) > 0 {
				rb = rollback[0]
			}

			ctx := appctx.FromContext(ginCtx)
			log.DebugCtx(ctx, "----------BEGIN TX")
			tx := db.GetDB().BeginTx(ctx.GetContext(), opt) // (ctx, opt)
			if tx.Error != nil {
				log.DebugCtx(ctx, fmt.Sprintf("----------BEGIN TX ERROR: %v", tx.Error))
				if rb != nil {
					return rb(ctx.GetContext())
				}

				return errors.New("TX error")
			}
			defer func() {
				if tx != nil {
					log.DebugCtx(ctx, "----------ROLLBACK RECOVERED TX")
					tx.Rollback()
					if rb != nil {
						_ = rb(ctx.GetContext())
					}
				}
			}()

			ctx = ctx.WithDb(tx)
			err = handler(ctx.GetContext())
			if err != nil {
				log.DebugCtx(ctx, fmt.Sprintf("----------ROLLBACK TX: %v", err))
				tx.Rollback()
			} else {
				log.DebugCtx(ctx, "----------COMMIT TX")
				tx.Commit()
			}

			tx = nil
			return err
		}
	}
}
