package facebook

import (
	"fmt"

	"github.com/iawia002/annie/downloader"
	"github.com/iawia002/annie/request"
	"github.com/iawia002/annie/utils"
)

// Extract is the main function for extracting data
func Extract(url string) ([]downloader.Data, error) {
	var err error
	html, err := request.Get(url, url, nil)
	if err != nil {
		return downloader.EmptyList, err
	}
	title := utils.MatchOneOf(html, `<title id="pageTitle">(.+)</title>`)[1]

	streams := map[string]downloader.Stream{}
	for _, quality := range []string{"sd", "hd"} {
		srcElement := utils.MatchOneOf(
			html, fmt.Sprintf(`%s_src:"(.+?)"`, quality),
		)
		if srcElement == nil {
			continue
		}
		u := srcElement[1]
		size, err := request.Size(u, url)
		if err != nil {
			return downloader.EmptyList, err
		}
		urlData := downloader.URL{
			URL:  u,
			Size: size,
			Ext:  "mp4",
		}
		streams[quality] = downloader.Stream{
			URLs:    []downloader.URL{urlData},
			Size:    size,
			Quality: quality,
		}
	}

	return []downloader.Data{
		{
			Site:    "Facebook facebook.com",
			Title:   title,
			Type:    "video",
			Streams: streams,
			URL:     url,
		},
	}, nil
}
