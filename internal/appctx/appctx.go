package appctx

import (
	"context"

	"github.com/sirupsen/logrus"
)

const Identifier = "uuid"

type Context struct {
	Context context.Context
	Logger  logrus.FieldLogger
}

func NewContext(ctx context.Context) *Context {
	logger := logrus.StandardLogger()
	fields := logrus.Fields{}

	i := ctx.Value(Identifier)
	if i != nil {
		fields = logrus.Fields{
			Identifier: i,
		}
	}

	return &Context{
		Context: ctx,
		Logger:  logger.WithFields(fields),
	}
}
