package uid

import (
	"context"
	"errors"

	"github.com/airunny/filter/cache"
	filterContext "github.com/airunny/filter/context"
	"github.com/airunny/filter/variables"
)

const Name = "uid"

func init() {
	variables.Register(variables.NewSimpleVariable(&UID{}))
}

// UID 用户ID
type UID struct{}

func (s *UID) Name() string    { return Name }
func (s *UID) Cacheable() bool { return true }
func (s *UID) Value(ctx context.Context, _ interface{}, _ *cache.Cache) (interface{}, error) {
	uid, ok := filterContext.FromUserId(ctx)
	if !ok {
		return nil, errors.New("uid not found in context")
	}
	return uid, nil
}
