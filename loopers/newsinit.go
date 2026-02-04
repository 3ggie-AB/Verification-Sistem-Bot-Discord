package loopers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"crypto-member/db"
	"crypto-member/models"

	"github.com/PuerkitoBio/goquery"
	"gorm.io/gorm/clause"
)

// translateText pake MyMemory
func translateText(text string) (string, error) {
	if text == "" {
		return "", nil
	}

	encodedText := url.QueryEscape(text)
	apiURL := fmt.Sprintf("https://api.mymemory.translated.net/get?q=%s&langpair=en|id", encodedText)

	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	var result struct {
		ResponseData struct {
			TranslatedText string `json:"translatedText"`
		} `json:"responseData"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", err
	}

	return result.ResponseData.TranslatedText, nil
}

// fetchFullContent scrape full article content
func fetchFullContent(articleURL string) (string, error) {
	res, err := http.Get(articleURL)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch URL: %s, status: %d", articleURL, res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	var paragraphs []string
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			paragraphs = append(paragraphs, text)
		}
	})

	content := strings.Join(paragraphs, "\n\n")
	return content, nil
}

// Translate batch supaya gak kena limit MyMemory
func translateFullContent(text string) string {
	if text == "" {
		return ""
	}

	maxLen := 400
	var final strings.Builder

	for i := 0; i < len(text); i += maxLen {
		end := i + maxLen
		if end > len(text) {
			end = len(text)
		}
		part := text[i:end]
		translatedPart, err := translateText(part)
		if err != nil {
			log.Println("Translate failed for part:", err)
			translatedPart = "" // skip kalau gagal
		}
		final.WriteString(translatedPart)
		final.WriteString("\n\n")
	}

	return final.String()
}

// SaveArticleFullScrape simpan artikel + scrape full content + translate
func SaveArticleFullScrape(article models.CryptoNews) error {
	// Scrape full content
	fullContent, err := fetchFullContent(article.URL)
	if err != nil {
		log.Println("Failed to fetch full content, fallback to API content:", err)
		fullContent = article.Content
	}
	article.Content = fullContent

	// Translate full content batch
	article.ContentIndo = translateFullContent(fullContent)

	// Simpan ke DB, jangan duplikat
	return db.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&article).Error
}


// SaveCryptoNewsLoop ambil semua article dari endpoint, scrape, translate, simpan
func SaveCryptoNewsLoop(apiKey string) {
	interval := 2*time.Hour + 24*time.Minute // 10x sehari
	for {
		articles, err := FetchCryptoNews(apiKey) // ambil endpoint
		if err != nil {
			fmt.Println("Error fetching news:", err)
			time.Sleep(30 * time.Minute)
			continue
		}

		for _, a := range articles {
			// Map endpoint ke models.CryptoNews
			news := models.CryptoNews{
				ArticleID:    a.ArticleID,
				SourceName:   a.Source.Name,
				SourceDomain: a.Source.Domain,
				Thumbnail:    a.Thumbnail,
				URL:          a.URL,
				Title:        a.Title,
				Description:  a.Description,
				Content:      a.Content, // sementara, nanti di-overwrite
				IsUpload:     false,
			}

			err := SaveArticleFullScrape(news)
			if err != nil {
				fmt.Println("Error saving article:", err)
			} else {
				fmt.Println("Saved full article:", a.Title)
			}
		}

		time.Sleep(interval)
	}
}
