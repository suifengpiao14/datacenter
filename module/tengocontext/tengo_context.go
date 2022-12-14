package tengocontext

import (
	context "context"

	"github.com/d5/tengo/v2"
)

type TengoContext struct {
	tengo.ObjectImpl
	Context context.Context
}

func (c *TengoContext) TypeName() string {
	return "context"
}
func (c *TengoContext) String() string {
	return "context"
}
