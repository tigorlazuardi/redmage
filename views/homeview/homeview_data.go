package homeview

import (
	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/models"
)

type Data struct {
	SubredditsList      api.ListSubredditsResult
	RecentlyAddedImages RecentlyAddedImages
	Error               error
}

type RecentlyAddedImages = map[int32]deviceMapValue

type deviceMapValue struct {
	device     *models.Device
	subreddits map[int32]subredditMapValue
}

type subredditMapValue struct {
	subreddit *models.Subreddit
	images    []*models.Image
}

func NewRecentlyAddedImages(images models.ImageSlice) RecentlyAddedImages {
	r := make(RecentlyAddedImages)
	for _, image := range images {
		if image.R.Device == nil || image.R.Subreddit == nil {
			continue
		}
		if _, ok := r[image.R.Device.ID]; !ok {
			r[image.R.Device.ID] = deviceMapValue{
				device:     image.R.Device,
				subreddits: make(map[int32]subredditMapValue),
			}
		}
		if _, ok := r[image.R.Device.ID].subreddits[image.R.Subreddit.ID]; !ok {
			r[image.R.Device.ID].subreddits[image.R.Subreddit.ID] = subredditMapValue{}
		}
		images := append(r[image.R.Device.ID].subreddits[image.R.Subreddit.ID].images, image)
		r[image.R.Device.ID].subreddits[image.R.Subreddit.ID] = subredditMapValue{
			subreddit: image.R.Subreddit,
			images:    images,
		}
	}
	return r
}
