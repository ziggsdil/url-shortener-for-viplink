package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	database "git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/db"
	apierrors "git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/errors"
)

const (
	shortLinkLen = 5
	secretKeyLen = 8
)

var shortLinkFunc = func(baseUrl, suffix string) string { return fmt.Sprintf("http://%s/%s", baseUrl, suffix) }

// my
func (h *Handler) ShortVIP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request ShortLinkRequest
	err := h.parseBody(r.Body, &request)
	if err != nil {
		fmt.Printf("error when parse body: %v\n", err)
		h.renderer.RenderError(w, apierrors.BadRequest{})
		return
	}

	if err := request.Validate(); err != nil {
		fmt.Printf("error when validate request: %v\n", err)
		h.renderer.RenderError(w, apierrors.BadRequest{})
		return
	}

	var shortSuffix string
	if request.VipKey == "" {
		for {
			shortSuffix, err = h.generateShortSuffix(ctx)
			if err == nil {
				break
			}

			if !errors.Is(err, apierrors.SuffixAlreadyExistsError{}) {
				fmt.Printf("error when generate short link: %v\n", err)
				h.renderer.RenderError(w, apierrors.InternalError{})
				return
			}
		}
	}
	// проверка не занята ли уже короткая ссылка
	// TODO: нужно проверить что содержится в vipkey
	// не получается пройти тест из за return на 63 строке, он возвращает сразу потому что встречает, что уже такая ссылка есть, возможно стоит убрать эту проверку.
	fmt.Println(request.VipKey)
	shortLink, err := h.db.SelectBySuffix(ctx, request.VipKey)
	fmt.Println("1123------------------", shortLink, err)
	switch {
	case err == nil:
		fmt.Printf("vip url \"%s\" already exist", request.VipKey)
		h.renderer.RenderError(w, apierrors.BadRequest{})
		return
	case errors.Is(err, database.SuffixNotFoundError):
	default:
		fmt.Printf("error when select vip url: %v\n", err)
		h.renderer.RenderError(w, apierrors.InternalError{})
		return
	}

	// проверка существует ли уже короткая ссылка на длинный url
	link, err := h.db.SelectByLink(ctx, request.LongUrl)
	switch {
	case err == nil:
		fmt.Printf("long url \"%s\" has been already shorten with suffix %s\n", link.Link, link.ShortSuffix)
		h.renderer.RenderJSON(w, ShortLinkResponse{ShortUrl: shortLinkFunc(h.url, link.ShortSuffix)})
		return
	case errors.Is(err, database.LinkNotFoundError):
	default:
		fmt.Printf("error when select long url: %v\n", err)
		h.renderer.RenderError(w, apierrors.InternalError{})
		return
	}

	var secretKey string
	for {
		secretKey, err = h.generateSecretKey(ctx)
		if err == nil {
			break
		}

		if !errors.Is(err, apierrors.SecretKeyAlreadyExistsError{}) {
			fmt.Printf("error when generate secret key: %v\n", err)
			h.renderer.RenderError(w, apierrors.InternalError{})
			return
		}
	}

	// преобразование единицы измерения временно интервала
	// TODO: возможно стоит отрефакторить
	var duration time.Duration
	switch request.TimeToLiveUnit {
	case "SECONDS":
		duration = time.Duration(request.TimeToLive) * time.Second
	case "MINUTES":
		duration = time.Duration(request.TimeToLive) * time.Minute
	case "HOURS":
		duration = time.Duration(request.TimeToLive) * time.Hour
	case "DAYS":
		duration = time.Duration(request.TimeToLive) * time.Hour * 24
	default:
		fmt.Printf("Incorrect time type")
	}

	expirationDate := time.Now().UTC().Add(duration) // приводим к типу UTC для сравнения вне зависимости от временной зоны

	// TODO: refactor
	if request.VipKey == "" {
		err = h.db.Save(ctx, shortSuffix, request.LongUrl, secretKey, expirationDate, false)
	} else {
		err = h.db.Save(ctx, request.VipKey, request.LongUrl, secretKey, expirationDate, false)

	}
	if err != nil {
		fmt.Printf("error when saving short link: %v\n", err)
		h.renderer.RenderError(w, apierrors.InternalError{})
		return
	}

	fmt.Printf("short link \"%s\" with suffix \"%s\" has been successfully saved\n", request.LongUrl, request.VipKey)
	// TODO: refactor
	if request.VipKey == "" {
		h.renderer.RenderJSON(w, ShortLinkResponse{ShortUrl: shortLinkFunc(h.url, shortSuffix), SecretKey: secretKey})
	} else {
		h.renderer.RenderJSON(w, ShortLinkResponse{ShortUrl: shortLinkFunc(h.url, request.VipKey), SecretKey: secretKey})
	}
}

func (h *Handler) generateShortSuffix(ctx context.Context) (string, error) {
	shortSuffix, err := h.generate(shortLinkLen)
	if err != nil {
		return "", err
	}

	// check if short suffix has already been used
	_, err = h.db.SelectBySuffix(ctx, shortSuffix)
	switch {
	case err == nil:
		return "", apierrors.SuffixAlreadyExistsError{}
	case errors.Is(err, database.SuffixNotFoundError):
		return shortSuffix, nil
	default:
		return "", err
	}
}

func (h *Handler) generateSecretKey(ctx context.Context) (string, error) {
	secretKey, err := h.generate(secretKeyLen)
	if err != nil {
		return "", err
	}

	// check if secret key has already been used
	_, err = h.db.SelectBySecretKey(ctx, secretKey)
	switch {
	case err == nil:
		return "", apierrors.SecretKeyAlreadyExistsError{}
	case errors.Is(err, database.SecretKeyNotFoundError):
		return secretKey, nil
	default:
		return "", err
	}
}

func (h *Handler) generate(length int) (string, error) {
	// generate random bytes
	randomBytes := make([]byte, length)
	encodedBytes := make([]byte, hex.EncodedLen(length))
	n, err := rand.Read(randomBytes)
	if n != length {
		fmt.Printf("invalid bytes generated: expected %d, got %d\n", length, n)
		return "", apierrors.InternalError{}
	}
	if err != nil {
		return "", err
	}

	// encode to make human-readable string
	hex.Encode(encodedBytes, randomBytes)
	return string(encodedBytes), nil
}
