package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Release struct {
	TagName string `json:"tag_name"`
}

func IsLatest(current string) (bool, string, error) {
	resp, err := http.Get("https://api.github.com/repos/anhcraft/ytb-downloader/releases/latest")
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, "", fmt.Errorf("failed to fetch latest release: %s", resp.Status)
	}

	var release Release
	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		return false, "", err
	}

	tagName := strings.TrimPrefix(release.TagName, "v")
	return compareVersions(current, tagName) >= 0, tagName, nil
}

func compareVersions(v1, v2 string) int {
	v1Parts := strings.Split(v1, ".")
	v2Parts := strings.Split(v2, ".")

	for i := 0; i < 3; i++ {
		v1Num, err := strconv.Atoi(v1Parts[i])
		if err != nil {
			return -1
		}
		v2Num, err := strconv.Atoi(v2Parts[i])
		if err != nil {
			return -1
		}

		if v1Num < v2Num {
			return -1
		}
		if v1Num > v2Num {
			return 1
		}
	}

	return 0
}
