package context

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

var entityKey = "entity"

type EntityModel interface {
	GetId() UUID
	GetScopes() []string
}

type Entity struct {
	Id     UUID     `json:"id,omitempty"`
	Scopes []string `json:"scopes,omitempty"`
}

func NewEntityModel(id UUID, scopes []string) EntityModel {
	return Entity{
		Id:     id,
		Scopes: scopes,
	}
}

func (e Entity) GetId() UUID {
	return e.Id
}

func (e Entity) GetScopes() []string {
	return e.Scopes
}

// EntityFromContext finds the tracer information from the context. REQUIRES Middleware to have run.
func EntityFromContext(ctx *gin.Context) EntityModel {
	raw, _ := ctx.Value(entityKey).(Entity)
	return raw
}

// EntityToContext inject entity id into context.
func EntityToContext(ctx *gin.Context, id UUID, scopes []string) *gin.Context {
	entity := NewEntityModel(id, scopes)
	if currentEntity, ok := ctx.Value(entityKey).(Entity); !ok || !reflect.DeepEqual(currentEntity, entity) {
		ctx.Set(entityKey, entity)
	}
	return ctx
}
