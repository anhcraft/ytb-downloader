package request

import (
	"net/url"
	"testing"
)

func TestRewriteYoutubeShortLink(t *testing.T) {
	u, e := url.Parse("https://youtu.be/dQw4w9WgXcQ?si=xyI9Ut-ovop31b")
	if e != nil {
		t.Fatal(e)
	}

	RewriteYoutubeShortLink(u)

	if u.String() != "https://youtube.com/watch?si=xyI9Ut-ovop31b&v=dQw4w9WgXcQ" {
		t.Fatal(u.String())
	}
}
