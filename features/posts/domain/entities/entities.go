package entities

import (
	"github.com/k0marov/go-socnet/core/abstract/ownable_likeable/contexters"
	"github.com/k0marov/go-socnet/core/general/static_store"
	"github.com/k0marov/go-socnet/features/posts/domain/models"
	"github.com/k0marov/go-socnet/features/posts/domain/values"
	profile_entities "github.com/k0marov/go-socnet/features/profiles/domain/entities"
)

type Post struct {
	models.PostModel
	Images []values.PostImage
	Likes  int
}

type ContextedPost struct {
	Post
	contexters.OwnLikeContext
	Author profile_entities.ContextedProfile
}

func ImagePathsToUrls(models []models.PostImageModel) (images []values.PostImage) {
	for _, model := range models {
		images = append(images, values.PostImage{
			URL:   static_store.PathToURL(model.Path),
			Index: model.Index,
		})
	}
	return
}
