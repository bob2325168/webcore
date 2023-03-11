package err

import "errors"

var (
	ErrRequestEmpty = errors.New("context: ctx.request empty")
)
