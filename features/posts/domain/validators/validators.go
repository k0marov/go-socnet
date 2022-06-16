package validators

import (
	"github.com/k0marov/socnet/core/client_errors"
	"github.com/k0marov/socnet/features/posts/domain/values"
)

type PostValidator func(newPost values.NewPostData) (client_errors.ClientError, bool)
