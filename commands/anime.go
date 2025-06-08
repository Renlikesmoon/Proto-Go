package anime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type AnimeResult struct {
	Data []struct {
		MalID int `json:"mal_id"`
		Title string `json:"title"`
		Images struct {
			JPG struct {
				ImageURL string `json:"image_url"`
			} `json:"jpg"`
		} `json:"images"`
		Score float64 `json:"score"`
		Episodes int `json:"episodes"`
		Synopsis string `json:"synopsis"`
		Url string `json:"url"`
	} `json:"data"`
}

func SearchAnime(query string) (string, error) {
	baseUrl := "https://api.jikan.moe/v4/anime"
	reqUrl := fmt.Sprintf("%s?q=%s&limit=1", baseUrl, url.QueryEscape(query))

	resp, err := http.Get(reqUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %s", resp.Status)
	}

	var result AnimeResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	if len(result.Data) == 0 {
		return "Anime tidak ditemukan.", nil
	}

	anime := result.Data[0]
	message := fmt.Sprintf(
		"🎬 *%s*\n⭐ Score: %.2f | 🧩 Episodes: %d\n\n_%s_\n\n🔗 %s",
		anime.Title,
		anime.Score,
		anime.Episodes,
		strings.TrimSpace(anime.Synopsis),
		anime.Url,
	)

	return message, nil
}
