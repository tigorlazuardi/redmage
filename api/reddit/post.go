package reddit

import (
	"fmt"
	"html"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/tigorlazuardi/redmage/config"
	"github.com/tigorlazuardi/redmage/models"
)

type Listing struct {
	Kind string `json:"kind"`
	Data Data   `json:"data"`
}

func (l *Listing) GetPosts() []Post {
	return l.Data.Children
}

func (l *Listing) GetSubreddit() string {
	length := len(l.Data.Children)
	if length == 0 {
		return ""
	}
	return l.Data.Children[length-1].Data.Subreddit
}

// GetLastAfter returns the last post namee for pagination.
//
// Returns empty string if there is no more posts to look up.
func (l *Listing) GetLastAfter() string {
	posts := l.GetPosts()
	if len(posts) == 0 {
		return ""
	}
	return posts[len(posts)-1].GetName()
}

type (
	MediaEmbed       struct{}
	SecureMediaEmbed struct{}
	Gildings         struct{}
	Source           struct {
		URL    string `json:"url"`
		Width  int64  `json:"width"`
		Height int64  `json:"height"`
	}
)

type Resolutions struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
type (
	Variants struct{}
	Images   struct {
		Source      Source        `json:"source"`
		Resolutions []Resolutions `json:"resolutions"`
		Variants    Variants      `json:"variants"`
		ID          string        `json:"id"`
	}
)

type Preview struct {
	Images  []Images `json:"images"`
	Enabled bool     `json:"enabled"`
}
type LinkFlairRichtext struct {
	E string `json:"e"`
	T string `json:"t"`
}
type ThumbnailPreview struct {
	Y int    `json:"y"`
	X int    `json:"x"`
	U string `json:"u"`
}

type MediaMetadata struct {
	Status          string             `json:"status"`
	Kind            string             `json:"e"`
	Mimetype        string             `json:"m"`
	ExtraThumbnails []ThumbnailPreview `json:"p"`
	Thumbnail       ThumbnailPreview   `json:"s"`
	ID              string             `json:"id"`
}
type Items struct {
	OutboundURL string `json:"outbound_url,omitempty"`
	MediaID     string `json:"media_id"`
	ID          int    `json:"id"`
}
type GalleryData struct {
	Items []Items `json:"items"`
}
type AuthorFlairRichtext struct {
	E string `json:"e"`
	T string `json:"t"`
}
type PostData struct {
	ApprovedAtUtc              any                      `json:"approved_at_utc"`
	Subreddit                  string                   `json:"subreddit"`
	Selftext                   string                   `json:"selftext"`
	AuthorFullname             string                   `json:"author_fullname"`
	Saved                      bool                     `json:"saved"`
	ModReasonTitle             any                      `json:"mod_reason_title"`
	Gilded                     int                      `json:"gilded"`
	Clicked                    bool                     `json:"clicked"`
	IsGallery                  bool                     `json:"is_gallery"`
	Title                      string                   `json:"title"`
	LinkFlairRichtext          []LinkFlairRichtext      `json:"link_flair_richtext"`
	SubredditNamePrefixed      string                   `json:"subreddit_name_prefixed"`
	Hidden                     bool                     `json:"hidden"`
	Pwls                       int                      `json:"pwls"`
	LinkFlairCSSClass          string                   `json:"link_flair_css_class"`
	Downs                      int                      `json:"downs"`
	ThumbnailHeight            int64                    `json:"thumbnail_height"`
	TopAwardedType             any                      `json:"top_awarded_type"`
	HideScore                  bool                     `json:"hide_score"`
	MediaMetadata              map[string]MediaMetadata `json:"media_metadata"`
	Name                       string                   `json:"name"`
	Quarantine                 bool                     `json:"quarantine"`
	LinkFlairTextColor         any                      `json:"link_flair_text_color"`
	UpvoteRatio                float64                  `json:"upvote_ratio"`
	AuthorFlairBackgroundColor any                      `json:"author_flair_background_color"`
	Ups                        int                      `json:"ups"`
	Domain                     string                   `json:"domain"`
	MediaEmbed                 MediaEmbed               `json:"media_embed"`
	ThumbnailWidth             int64                    `json:"thumbnail_width"`
	AuthorFlairTemplateID      string                   `json:"author_flair_template_id"`
	IsOriginalContent          bool                     `json:"is_original_content"`
	UserReports                []any                    `json:"user_reports"`
	SecureMedia                any                      `json:"secure_media"`
	IsRedditMediaDomain        bool                     `json:"is_reddit_media_domain"`
	IsMeta                     bool                     `json:"is_meta"`
	Category                   any                      `json:"category"`
	SecureMediaEmbed           SecureMediaEmbed         `json:"secure_media_embed"`
	GalleryData                GalleryData              `json:"gallery_data"`
	LinkFlairText              string                   `json:"link_flair_text"`
	CanModPost                 bool                     `json:"can_mod_post"`
	Score                      int                      `json:"score"`
	ApprovedBy                 any                      `json:"approved_by"`
	IsCreatedFromAdsUI         bool                     `json:"is_created_from_ads_ui"`
	AuthorPremium              bool                     `json:"author_premium"`
	Thumbnail                  string                   `json:"thumbnail"`
	Edited                     any                      `json:"edited"`
	AuthorFlairCSSClass        string                   `json:"author_flair_css_class"`
	AuthorFlairRichtext        []AuthorFlairRichtext    `json:"author_flair_richtext"`
	Gildings                   Gildings                 `json:"gildings"`
	ContentCategories          any                      `json:"content_categories"`
	IsSelf                     bool                     `json:"is_self"`
	SubredditType              string                   `json:"subreddit_type"`
	Created                    float64                  `json:"created"`
	LinkFlairType              string                   `json:"link_flair_type"`
	Wls                        int                      `json:"wls"`
	RemovedByCategory          any                      `json:"removed_by_category"`
	BannedBy                   any                      `json:"banned_by"`
	AuthorFlairType            string                   `json:"author_flair_type"`
	TotalAwardsReceived        int                      `json:"total_awards_received"`
	AllowLiveComments          bool                     `json:"allow_live_comments"`
	SelftextHTML               any                      `json:"selftext_html"`
	Likes                      any                      `json:"likes"`
	SuggestedSort              any                      `json:"suggested_sort"`
	BannedAtUtc                any                      `json:"banned_at_utc"`
	URLOverriddenByDest        string                   `json:"url_overridden_by_dest"`
	ViewCount                  any                      `json:"view_count"`
	Archived                   bool                     `json:"archived"`
	NoFollow                   bool                     `json:"no_follow"`
	IsCrosspostable            bool                     `json:"is_crosspostable"`
	Pinned                     bool                     `json:"pinned"`
	Over18                     bool                     `json:"over_18"`
	AllAwardings               []any                    `json:"all_awardings"`
	Awarders                   []any                    `json:"awarders"`
	MediaOnly                  bool                     `json:"media_only"`
	CanGild                    bool                     `json:"can_gild"`
	Spoiler                    bool                     `json:"spoiler"`
	Locked                     bool                     `json:"locked"`
	AuthorFlairText            string                   `json:"author_flair_text"`
	TreatmentTags              []any                    `json:"treatment_tags"`
	Visited                    bool                     `json:"visited"`
	RemovedBy                  any                      `json:"removed_by"`
	ModNote                    any                      `json:"mod_note"`
	Distinguished              any                      `json:"distinguished"`
	SubredditID                string                   `json:"subreddit_id"`
	AuthorIsBlocked            bool                     `json:"author_is_blocked"`
	ModReasonBy                any                      `json:"mod_reason_by"`
	NumReports                 any                      `json:"num_reports"`
	RemovalReason              any                      `json:"removal_reason"`
	LinkFlairBackgroundColor   any                      `json:"link_flair_background_color"`
	ID                         string                   `json:"id"`
	IsRobotIndexable           bool                     `json:"is_robot_indexable"`
	ReportReasons              any                      `json:"report_reasons"`
	Author                     string                   `json:"author"`
	DiscussionType             any                      `json:"discussion_type"`
	NumComments                int                      `json:"num_comments"`
	SendReplies                bool                     `json:"send_replies"`
	WhitelistStatus            string                   `json:"whitelist_status"`
	ContestMode                bool                     `json:"contest_mode"`
	ModReports                 []any                    `json:"mod_reports"`
	AuthorPatreonFlair         bool                     `json:"author_patreon_flair"`
	AuthorFlairTextColor       string                   `json:"author_flair_text_color"`
	Permalink                  string                   `json:"permalink"`
	ParentWhitelistStatus      string                   `json:"parent_whitelist_status"`
	Stickied                   bool                     `json:"stickied"`
	URL                        string                   `json:"url"`
	SubredditSubscribers       int                      `json:"subreddit_subscribers"`
	CreatedUtc                 float64                  `json:"created_utc"`
	NumCrossposts              int                      `json:"num_crossposts"`
	Media                      any                      `json:"media"`
	IsVideo                    bool                     `json:"is_video"`
	PostHint                   string                   `json:"post_hint"`
	Preview                    Preview                  `json:"preview"`
}

type Post struct {
	Kind string   `json:"kind"`
	Data PostData `json:"data,omitempty"`
}

func (post *Post) IsNSFW() bool {
	return post.Data.Over18
}

func (post *Post) IsNSFWInt() int {
	if post.IsNSFW() {
		return 1
	}
	return 0
}

func (post *Post) IsImagePost() bool {
	return post.Data.PostHint == "image"
}

func (post *Post) GetImageURL() string {
	return post.Data.URL
}

func (post *Post) GetCreated() time.Time {
	return time.Unix(int64(post.Data.Created), 0)
}

func (post *Post) GetAuthor() string {
	return post.Data.Author
}

func (post *Post) GetTitle() string {
	return html.UnescapeString(post.Data.Title)
}

func (post *Post) GetAuthorURL() string {
	return fmt.Sprintf("https://www.reddit.com/user/%s", post.Data.Author)
}

func (post *Post) GetImageAspectRatio() float64 {
	width, height := post.GetImageSize()
	if height == 0 {
		return 0
	}
	return float64(width) / float64(height)
}

func (post *Post) GetName() string {
	return post.Data.Name
}

func (post *Post) GetImageTargetPath(cfg *config.Config, device *models.Device) string {
	if device.WindowsWallpaperMode == 1 {
		return post.GetWindowsWallpaperImageTargetPath(cfg, device)
	}
	baseDownloadDir := cfg.String("download.directory")
	p := path.Join(baseDownloadDir, device.Slug, post.GetSubreddit(), post.GetImageFilename())
	abs, _ := filepath.Abs(p)
	return abs
}

func (post *Post) GetImageTargetDir(cfg *config.Config, device *models.Device) string {
	if device.WindowsWallpaperMode == 1 {
		return post.GetWindowsWallpaperImageTargetDir(cfg, device)
	}
	baseDownloadDir := cfg.String("download.directory")
	p := path.Join(baseDownloadDir, device.Slug, post.GetSubreddit())
	abs, _ := filepath.Abs(p)
	return abs
}

func (post *Post) GetWindowsWallpaperImageTargetPath(cfg *config.Config, device *models.Device) string {
	baseDownloadDir := cfg.String("download.directory")
	filename := fmt.Sprintf("%s_%s", post.GetSubreddit(), post.GetImageFilename())
	p := path.Join(baseDownloadDir, device.Slug, filename)
	abs, _ := filepath.Abs(p)
	return abs
}

func (post *Post) GetWindowsWallpaperImageTargetDir(cfg *config.Config, device *models.Device) string {
	baseDownloadDir := cfg.String("download.directory")
	p := path.Join(baseDownloadDir, device.Slug)
	abs, _ := filepath.Abs(p)
	return abs
}

func (post *Post) GetThumbnailTargetPath(cfg *config.Config) string {
	baseDownloadDir := cfg.String("download.directory")
	p := path.Join(baseDownloadDir, "_thumbnails", post.GetSubreddit(), post.GetImageFilename())
	abs, _ := filepath.Abs(p)
	if before, found := strings.CutSuffix(abs, ".png"); found {
		return before + ".jpeg"
	}
	return abs
}

func (post *Post) GetThumbnailTargetDir(cfg *config.Config) string {
	baseDownloadDir := cfg.String("download.directory")
	p := path.Join(baseDownloadDir, "_thumbnails", post.GetSubreddit())
	abs, _ := filepath.Abs(p)
	return abs
}

func (post *Post) GetThumbnailRelativePath() string {
	p := path.Join("_thumbnails", post.GetSubreddit(), post.GetImageFilename())
	if before, found := strings.CutSuffix(p, ".png"); found {
		return before + ".jpeg"
	}
	return p
}

func (post *Post) GetImageRelativePath(device *models.Device) string {
	if device.WindowsWallpaperMode == 1 {
		return post.GetWindowsWallpaperImageRelativePath(device)
	}

	return path.Join(device.Slug, post.GetSubreddit(), post.GetImageFilename())
}

func (post *Post) GetWindowsWallpaperImageRelativePath(device *models.Device) string {
	filename := fmt.Sprintf("%s_%s", post.GetSubreddit(), post.GetImageFilename())
	return path.Join(device.Slug, filename)
}

func (post *Post) GetImageFilename() string {
	if !post.IsImagePost() {
		return ""
	}
	link := post.GetImageURL()
	u, _ := url.Parse(link)
	if u == nil {
		return ""
	}
	split := strings.Split(u.Path, "/")
	return split[len(split)-1]
}

func (post *Post) GetImageSize() (width, height int64) {
	if len(post.Data.Preview.Images) == 0 {
		return 0, 0
	}
	source := post.Data.Preview.Images[0].Source
	return source.Width, source.Height
}

func (post *Post) GetThumbnailURL() string {
	return post.Data.Thumbnail
}

func (post *Post) GetThumbnailSize() (width, height int64) {
	return post.Data.ThumbnailWidth, post.Data.ThumbnailHeight
}

func (post *Post) GetSubreddit() string {
	return post.Data.Subreddit
}

func (post *Post) GetPermalink() string {
	return post.Data.Permalink
}

func (post *Post) GetPostURL() string {
	return fmt.Sprintf("https://reddit.com%s", post.Data.Permalink)
}

func (post *Post) GetPostName() string {
	return post.Data.Name
}

func (post *Post) GetPostTitle() string {
	return post.Data.Title
}

func (post *Post) GetPostCreated() int64 {
	return int64(post.Data.Created)
}

func (post *Post) GetID() string {
	return post.Data.ID
}

type Data struct {
	After     string `json:"after"`
	Dist      int    `json:"dist"`
	Modhash   string `json:"modhash"`
	GeoFilter any    `json:"geo_filter"`
	Children  []Post `json:"children"`
	Before    any    `json:"before"`
}
