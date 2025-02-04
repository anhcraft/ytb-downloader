package scripting

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleDownload(t *testing.T) {
	res, err := HandleDownload([]byte(`
url := import("url")

_action := "default"
_url := _input

process := func(input) {
    domain := url.extractDomain(input)
    query := url.extractQuery(input)

    if domain == "www.youtubetrimmer.com" {
        if query["v"] != undefined && len(query["v"]) > 0 {
            videoID := query["v"][0]
            newURL := "https://youtu.be/" + videoID
            _action = "override"
            _url = newURL
            return
        }
    }
}

process(_input)
`), "https://www.youtubetrimmer.com/view/?v=123")

	if err != nil {
		t.Fatal(err)
	} else {
		assert.Equal(t, "override", res.Action)
		assert.Equal(t, "https://youtu.be/123", res.Url)
	}
}
