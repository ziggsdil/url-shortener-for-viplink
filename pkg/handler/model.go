package handler

import (
	"fmt"
	"net/url"

	"git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/db"
)

type ShortLinkRequest struct {
	LongUrl string `json:"long_url"`
}

// my
type ShortVipLinkRequest struct {
	LongUrl        string `json:"long_url"`
	VipKey         string `json:"vip_key"`
	TimeToLive     int    `json:"ttl"`
	TimeToLiveUnit string `json:"ttl_unit"`
}

const (
	maxDays    = 2
	maxHours   = 48
	maxMinutes = 48 * 60
	maxSeconds = 48 * 60 * 60
)

// my
func (r ShortVipLinkRequest) Validate() error {
	if r.LongUrl == "" {
		return fmt.Errorf("invalid long url")
	}

	// TODO: написать проверку на существует ли уже такая vip ссылка
	if r.VipKey == "" {
		return fmt.Errorf("vip key is empty")
	}

	maxValues := map[string]int{
		"DAYS":    maxDays,
		"HOURS":   maxHours,
		"MINUTES": maxMinutes,
		"SECONDS": maxSeconds,
	}

	if r.TimeToLive > maxValues[r.TimeToLiveUnit] {
		return fmt.Errorf("date should be less than 2 days")
	}

	parsedUrl, err := url.Parse(r.LongUrl)
	if err != nil {
		return err
	}

	// можно еще проверить на http или https
	if parsedUrl.Scheme == "" {
		return fmt.Errorf("schema should be provided in long url")
	}

	/*	if (r.TimeToLiveUnit == "HOURS" && r.TimeToLive > 48) ||
		(r.TimeToLiveUnit == "DAYS" && r.TimeToLive > 2) ||
		(r.TimeToLiveUnit == "MINUTES" && r.TimeToLive > (48*60)) ||
		(r.TimeToLiveUnit == "SECONDS" && r.TimeToLive > (48*60*60)) {
		return fmt.Errorf("date should be less 2 days")
	}*/

	return nil
}

func (r ShortLinkRequest) Validate() error {
	if r.LongUrl == "" {
		return fmt.Errorf("invalid long url")
	}

	longUrl, err := url.Parse(r.LongUrl)
	if err != nil {
		return err
	}

	if longUrl.Scheme == "" {
		return fmt.Errorf("schema should be provided in long url")
	}

	return nil
}

type ShortLinkResponse struct {
	ShortUrl  string `json:"short_url"`
	SecretKey string `json:"secret_key"`
}

type InfoResponse struct {
	LongUrl  string `json:"long_url"`
	ShortUrl string `json:"short_url"`
	Clicks   int    `json:"clicks"`
}

func (r *InfoResponse) FromLink(link *db.Link, baseUrl string) {
	if link == nil {
		return
	}

	r.ShortUrl = shortLinkFunc(baseUrl, link.ShortSuffix)
	r.LongUrl = link.Link
	r.Clicks = link.Clicks
}
