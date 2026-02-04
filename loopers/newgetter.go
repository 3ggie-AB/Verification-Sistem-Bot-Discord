package loopers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type NewsArticle struct {
	ArticleID   string `json:"article_id"`
	Source      struct {
		Name   string `json:"name"`
		Domain string `json:"domain"`
	} `json:"source"`
	Thumbnail   string `json:"thumbnail"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

type ApiResponse struct {
	Data struct {
		Results []NewsArticle `json:"results"`
	} `json:"data"`
}

func FetchCryptoNews(apiKey string) ([]NewsArticle, error) {
	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Get(fmt.Sprintf("https://api.thenewsapi.net/crypto?within=24h&sentiments=negative,neutral,positive&categories=price-analysis,markets,altcoin,policy,technology,blockchain,security,AI,web3,stablecoins,news&langs=en,id&page=1&size=10&apikey=%s", apiKey))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var resp struct {
		Data struct {
			Results []NewsArticle `json:"results"`
		} `json:"data"`
	}

	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return resp.Data.Results, nil
}