package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
)

var audioPath = "./audio"
var mp3Path = fmt.Sprintf("%s/%s", audioPath, "latest.mp3")
var dfpwmPath = fmt.Sprintf("%s/%s", audioPath, "latest.dfpwm")

type QuoteRequest struct {
	ID           string   `json:"_id"`
	Content      string   `json:"content"`
	Author       string   `json:"author"`
	Tags         []string `json:"tags"`
	AuthorSlug   string   `json:"authorSlug"`
	Length       int      `json:"length"`
	DateAdded    string   `json:"dateAdded"`
	DateModified string   `json:"dateModified"`
}

func main() {
	r := gin.Default()
	r.GET("/get", func(c *gin.Context) {
		quote, err := getRandomQuote()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		err = download(fmt.Sprintf("%s - %s", quote.Content, quote.Author))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		err = exec.Command("ffmpeg",
			"-y",
			"-i", mp3Path,
			"-ar", "44100",
			"-af", "volume=4",
			dfpwmPath,
		).Run()

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.FileAttachment(dfpwmPath, "latest.dfpwm")
	})
	r.Run()
}

// getRandomQuote obtains a random quote from the quotable API
func getRandomQuote() (QuoteRequest, error) {
	resp, err := http.Get("https://api.quotable.io/random")
	if err != nil {
		return QuoteRequest{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return QuoteRequest{}, err
	}

	var quote QuoteRequest

	err = json.Unmarshal(body, &quote)
	if err != nil {
		return QuoteRequest{}, err
	}

	return quote, nil
}

// download obtains the audio file from the Google TTS API
func download(text string) error {
	url := fmt.Sprintf("http://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s", url.QueryEscape(text), "en")
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	output, err := os.Create(mp3Path)
	if err != nil {
		return err
	}

	_, err = io.Copy(output, response.Body)
	if err != nil {
		return err
	}

	return nil
}
