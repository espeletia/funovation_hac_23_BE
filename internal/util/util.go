package util

import "strings"

func GetYoutubeIDFromURL(url string) string {
	// Split the URL by the "v=" parameter
	splitURL := strings.Split(url, "v=")

	// If the URL has the "v=" parameter, return the ID after it
	if len(splitURL) > 1 {
		return splitURL[1]
	}

	// If the URL doesn't have the "v=" parameter, split the URL by the last "/"
	splitURL = strings.Split(url, "/")

	// Return the last element of the split URL
	return splitURL[len(splitURL)-1]
}
