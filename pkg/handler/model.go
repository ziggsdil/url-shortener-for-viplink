package handler

import (
	"fmt"
	apierrors "git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/errors"
	"net/url"

	"git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/db"
)

type ShortLinkRequest struct {
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

const (
	ttlUnitDays    = "DAYS"
	ttlUnitHours   = "HOURS"
	ttlUnitMinutes = "MINUTES"
	ttlUnitSeconds = "SECONDS"
)

const (
	defaultTtl     = 10
	defaultTtlUnit = "HOURS"
)

func (r *ShortLinkRequest) Validate() error {
	if r.LongUrl == "" {
		return fmt.Errorf("invalid long url")
	}

	if r.TimeToLive < 0 {
		return apierrors.BadRequest{}
	}

	// значение 0 является валидным, оно отрабатывается в случае, если ничего не посылается
	if r.TimeToLive == 0 {
		r.TimeToLive = defaultTtl
	}

	if r.TimeToLiveUnit == "" {
		r.TimeToLiveUnit = defaultTtlUnit
	}

	maxValues := map[string]int{
		ttlUnitDays:    maxDays,
		ttlUnitHours:   maxHours,
		ttlUnitMinutes: maxMinutes,
		ttlUnitSeconds: maxSeconds,
	}

	// TODO: написать проверку на отрицательные числа
	// написать проверку на пустые поля ttl и ttl_unit
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
