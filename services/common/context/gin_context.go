package context

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type UUID = string

type Context interface {
	GetContext() *gin.Context
	GetUserId() UUID
	GetDBTx() *gorm.DB
	// //RunTx(fn func(ctx Context) error, options ...*sql.TxOptions) error
	WithEntity(userId UUID, scopes []string) Context
	WithDb(db *gorm.DB) Context
}

type ctxWrap struct {
	*gin.Context
}

func (c *ctxWrap) GetContext() *gin.Context {
	return c.Context
}

// FromContext function provide method to wrapping base context to ContextX
func FromContext(ginCtx *gin.Context) Context {
	ctx := &ctxWrap{
		Context: ginCtx,
	}
	return ctx
}

func NewContext(ginCtx *gin.Context) Context {
	return &ctxWrap{ginCtx}
}

func (c *ctxWrap) GetUserId() UUID {
	return EntityFromContext(c.Context).GetId()
}

func (c *ctxWrap) GetDBTx() *gorm.DB {
	return DbTxFromContext(c.Context)
}

func (c *ctxWrap) WithEntity(userId UUID, scopes []string) Context {
	c.Context = EntityToContext(c.Context, userId, scopes)
	return c
}

func (c *ctxWrap) WithDb(db *gorm.DB) Context {
	c.Context = DbTxToContext(c.Context, db)
	return c
}
