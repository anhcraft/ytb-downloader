package request

import "net/url"

// RewriteYoutubeShortLink rewrite https://youtu.be/X?Y -> https://youtube.com/watch?v=&Y
func RewriteYoutubeShortLink(u *url.URL) {
	if u.Host == "youtu.be" {
		// we want to keep the original query parameters
		query := u.Query()
		query.Set("v", u.Path[1:])
		u.RawQuery = query.Encode()
		u.Host = "youtube.com"
		u.Path = "/watch"
	}
}
