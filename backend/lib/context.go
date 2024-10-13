package lib

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type subKey struct{}

func SetSubKey(ctx context.Context, sub string) context.Context {
	return context.WithValue(ctx, subKey{}, sub)
}
func GetSubValue(ctx context.Context) (uuid.UUID, error) {
	val, ok := ctx.Value(subKey{}).(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("no set subject value")
	}
	sub, err := uuid.Parse(val)
	return sub, err
}
