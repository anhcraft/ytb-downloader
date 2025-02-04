package request

import (
	"slices"
	"strconv"
	"ytb-downloader/internal/constants/downloadmode"
	"ytb-downloader/internal/constants/format"
	"ytb-downloader/internal/constants/thumbnail"
	"ytb-downloader/internal/settings"
)

func SupplyQueue(req []*Request) {
	commonArgs := append(settings.Get().ExtraYtdlpOptionsAsArray(), "--ignore-errors", "--no-warnings",
		"--progress", "--newline",
		"--progress-template", "[[PROGRESS]] %(progress._percent_str)s,%(progress._downloaded_bytes_str)s,%(progress._total_bytes_str)s,%(progress._speed_str)s,%(progress._eta_str)s",
		"--concurrent-fragments", strconv.FormatUint(uint64(settings.Get().GetConcurrentFragments()), 10),
		"--abort-on-unavailable-fragments",
		"-P", settings.Get().GetDownloadFolder())

	if downloadmode.HasYtdlpDownload(settings.Get().GetDisallowOverwrite()) {
		commonArgs = append(commonArgs, "--no-overwrites")
	}

	if fp := settings.Get().GetFfmpegPath(); len(fp) > 0 {
		commonArgs = append(commonArgs, "--ffmpeg-location", fp)
	}

	fmt := settings.Get().GetFormat()
	embedThumbnail := settings.Get().GetEmbedThumbnail()
	shouldEmbedThumbnail := embedThumbnail != thumbnail.Never &&
		(embedThumbnail == thumbnail.Always || fmt == embedThumbnail)

	if shouldEmbedThumbnail {
		commonArgs = append(commonArgs, "--embed-thumbnail")
	}

	// Choose the best quality format
	// Remux the video to mp4 or audio to m4a to support thumbnail embedding
	if fmt == format.VideoOnly {
		commonArgs = append(commonArgs, "-f", "bestvideo")
		if shouldEmbedThumbnail {
			commonArgs = append(commonArgs, "--remux-video", "mp4")
		}
	} else if fmt == format.AudioOnly {
		commonArgs = append(commonArgs, "-f", "bestaudio")
		if shouldEmbedThumbnail {
			commonArgs = append(commonArgs, "-x", "--audio-quality", "0", "--audio-format", "m4a")
		}
	} else if shouldEmbedThumbnail {
		commonArgs = append(commonArgs, "--merge-output-format", "mp4")
	}

	for _, r := range req {
		if r.Custom() {
			r.SetFormat(format.Default)
			continue
		}

		r.SetFormat(fmt)

		args := slices.Clone(commonArgs)
		args = append(args, r.rawUrl)
		r.SetDownloadCmdArgs(args)
	}

	GetQueue().OfferBulk(req)
}
