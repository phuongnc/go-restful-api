package context

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// dbTxContextKey private key for context, this is important
var dbTxCtxKey = "dbTx"

// DbTxFromContext get the database transaction from context. REQUIRES Middleware to have run.
func DbTxFromContext(ctx *gin.Context) *gorm.DB {
	raw, _ := ctx.Value(dbTxCtxKey).(*gorm.DB)
	return raw
}

// DbTxToContext inject database transaction into provided context.
func DbTxToContext(ctx *gin.Context, db *gorm.DB) *gin.Context {
	ctx.Set(dbTxCtxKey, db)
	return ctx
}
