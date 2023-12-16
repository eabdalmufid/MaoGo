package api

import (
	"fmt"
	"net/http"
	"strings"
)

func GetTiktokVideo(url string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()
	location := resp.Header.Get("location")
	data := strings.Split(location, "/")
	data = data[5:]

	return fmt.Sprintf("https://www.tikwm.com/video/media/play/%s.mp4", data[0][:19]), nil
}
