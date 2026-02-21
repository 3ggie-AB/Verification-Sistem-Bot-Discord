package loopers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"crypto-member/db"
	"crypto-member/models"

	"gorm.io/gorm/clause"
)

func summarizeWithGroq(apiKey string, articleURL string, context string) (string, error) {
	endpoint := "https://api.groq.com/openai/v1/chat/completions"

	prompt := fmt.Sprintf(`Anda adalah analis berita keuangan profesional.

Baca artikel dari URL berikut dan buat ringkasan dalam Bahasa Indonesia yang ringkas, jelas, dan objektif.

URL: %s

Tambahan Konteks dari saya :
%s

Instruksi:
- Tulis ringkasan utama (3-8 kalimat)
- Fokus pada poin penting
- Gunakan bahasa Indonesia yang mudah dipahami
- Jangan mencantumkan opini pribadi
- Jangan menyertakan tautan

Output hanya isi ringkasan.`, articleURL, context)

	payload := map[string]interface{}{
		"model": "openai/gpt-oss-120b",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.3,
	}

	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("groq error: %s", string(bodyBytes))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from Groq")
	}

	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}

func SaveArticleWithGroq(article models.CryptoNews, groqAPIKey string) error {

	summary, err := summarizeWithGroq(groqAPIKey, article.URL, article.Description)
	if err != nil {
		log.Println("Groq summarize failed:", err)
		return err
	}

	article.ContentIndo = summary

	return db.DB.
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&article).Error
}

func SaveCryptoNewsLoop(newsAPIKey string) {

	groqAPIKey := os.Getenv("GROQ_API_KEY")

	if groqAPIKey == "" {
		log.Fatal("GROQ_API_KEY not set")
	}

	interval := 2*time.Hour + 24*time.Minute

	for {

		articles, err := FetchCryptoNews(newsAPIKey)
		if err != nil {
			log.Println("Error fetching news:", err)
			time.Sleep(30 * time.Minute)
			continue
		}

		for _, a := range articles {

			news := models.CryptoNews{
				ArticleID:    a.ArticleID,
				SourceName:   a.Source.Name,
				SourceDomain: a.Source.Domain,
				Thumbnail:    a.Thumbnail,
				URL:          a.URL,
				Title:        a.Title,
				Description:  a.Description,
				Content:      a.Content,
				IsUpload:     false,
			}

			err := SaveArticleWithGroq(news, groqAPIKey)
			if err != nil {
				log.Fatal("Error saving article:", err)
				// Hapus kalau sudah ada di DB
				db.DB.
					Where("article_id = ?", news.ArticleID).
					Delete(&models.CryptoNews{})

				continue
			} else {
				log.Println("Saved:", a.Title)
			}
		}

		time.Sleep(interval)
	}
}