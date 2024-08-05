package appctx

import (
	"context"

	"github.com/google/uuid"

	"github.com/sirupsen/logrus"
)

const identifier = "req_id"

type Context struct {
	Context context.Context
	Logger  logrus.FieldLogger
}

func NewContext(ctx context.Context) *Context {
	logger := logrus.StandardLogger()
	u := uuid.New()

	return &Context{
		Context: ctx,
		Logger:  logger.WithField(identifier, u.String()),
	}
}
