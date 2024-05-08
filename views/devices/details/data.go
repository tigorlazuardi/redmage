package details

import (
	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/models"
)

type Data struct {
	Error       string
	Device      *models.Device
	Images      models.ImageSlice
	TotalImages int64
	Params      api.ImageListParams
}

type splitBySubredditImages struct {
	Subreddit string
	Images    models.ImageSlice
}

func (data Data) splitImages() []*splitBySubredditImages {
	var out []*splitBySubredditImages

	for _, image := range data.Images {
		var found bool

	inner:
		for _, o := range out {
			if o.Subreddit == image.Subreddit {
				found = true
				o.Images = append(o.Images, image)
				break inner
			}
		}

		if !found {
			out = append(out, &splitBySubredditImages{
				Subreddit: image.Subreddit,
				Images:    models.ImageSlice{image},
			})
		}

	}

	return out
}
