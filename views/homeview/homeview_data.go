package homeview

import (
	"slices"
	"strings"

	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/models"
)

type Data struct {
	SubredditsList      api.ListSubredditsResult
	RecentlyAddedImages RecentlyAddedImages
	Error               error
}

type RecentlyAddedImages = []RecentlyAddedImage

type subredditMapValue struct {
	subreddit *models.Subreddit
	images    []*models.Image
}

type RecentlyAddedImage struct {
	Device     *models.Device
	Subreddits []Subreddit
}

type Subreddit struct {
	Subreddit *models.Subreddit
	Images    models.ImageSlice
}

func NewRecentlyAddedImages(images models.ImageSlice) RecentlyAddedImages {
	r := make(RecentlyAddedImages, 0, len(images))
	for _, image := range images {
		if image.R.Device == nil || image.R.Subreddit == nil {
			continue
		}
		var deviceFound bool
		for i, ra := range r {
			if ra.Device.ID == image.R.Device.ID {
				deviceFound = true
				var subredditFound bool
				for j, subreddit := range r[i].Subreddits {
					if subreddit.Subreddit.ID == image.R.Subreddit.ID {
						subredditFound = true
						r[i].Subreddits[j].Images = append(r[i].Subreddits[j].Images, image)
					}
				}
				if !subredditFound {
					r[i].Subreddits = append(r[i].Subreddits, Subreddit{
						Subreddit: image.R.Subreddit,
						Images:    models.ImageSlice{image},
					})
				}
			}
		}
		if !deviceFound {
			r = append(r, RecentlyAddedImage{
				Device: image.R.Device,
				Subreddits: []Subreddit{
					{
						Subreddit: image.R.Subreddit,
						Images:    models.ImageSlice{image},
					},
				},
			})
		}
	}

	for _, r := range r {
		slices.SortFunc(r.Subreddits, func(left, right Subreddit) int {
			leftName := strings.ToLower(left.Subreddit.Name)
			rightName := strings.ToLower(right.Subreddit.Name)
			return strings.Compare(leftName, rightName)
		})
	}

	slices.SortFunc(r, func(left, right RecentlyAddedImage) int {
		leftName := strings.ToLower(left.Device.Name)
		rightName := strings.ToLower(right.Device.Name)
		return strings.Compare(leftName, rightName)
	})

	return r
}
