package anime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// AnimeResult defines the structure for unmarshaling the JSON response
// from the Jikan API for anime search results.
type AnimeResult struct {
	Data []struct {
		MalID    int    `json:"mal_id"`
		Title    string `json:"title"`
		Images   struct {
			JPG struct {
				ImageURL string `json:"image_url"`
			} `json:"jpg"`
		} `json:"images"`
		Score    float64 `json:"score"`
		Episodes int     `json:"episodes"`
		Synopsis string  `json:"synopsis"`
		Url      string  `json:"url"` // Note: JSON field is 'url', Go field is 'Url' (capitalized for export)
	} `json:"data"`
}

// SearchAnime performs a search query against the Jikan API for anime
// and returns a formatted string with the anime details or an error.
func SearchAnime(query string) (string, error) {
	// Base URL for the Jikan anime API.
	baseUrl := "https://api.jikan.moe/v4/anime"

	// Construct the request URL, ensuring the query is URL-encoded.
	// We limit the results to 1 as we only need the top match.
	reqUrl := fmt.Sprintf("%s?q=%s&limit=1", baseUrl, url.QueryEscape(query))

	// Perform the HTTP GET request.
	resp, err := http.Get(reqUrl)
	if err != nil {
		return "", fmt.Errorf("failed to make API request: %w", err)
	}
	defer resp.Body.Close() // Ensure the response body is closed after reading.

	// Check if the API request was successful (HTTP status 200 OK).
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: received status %s", resp.Status)
	}

	// Decode the JSON response into the AnimeResult struct.
	var result AnimeResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("failed to decode API response: %w", err)
	}

	// Check if any anime data was found in the response.
	if len(result.Data) == 0 {
		return "Anime tidak ditemukan.", nil // Return a friendly message if no results.
	}

	// Take the first (most relevant) anime from the results.
	anime := result.Data[0]

	// Format the anime details into a readable message string.
	// We trim whitespace from the synopsis for cleaner presentation.
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
