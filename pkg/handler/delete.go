package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	database "git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/db"
	apierrors "git.yandex-academy.ru/school/2023-06/backend/go/homeworks/intro_lecture/ya-url-shortener-for-viplink/pkg/errors"
)

const (
	timeToUpdate = 12
)

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	secretKey := chi.URLParam(r, "secretKey")

	err := h.db.DeleteBySecretKey(ctx, secretKey)
	switch {
	case errors.Is(err, database.SecretKeyNotFoundError):
		fmt.Printf("failed to found secret key \"%s\"\n", secretKey)
		h.renderer.RenderError(w, apierrors.NotFoundError{})
		return
	case err != nil:
		fmt.Printf("failed to delete secret key\n")
		h.renderer.RenderError(w, apierrors.NotFoundError{})
		return
	}

	h.renderer.RenderOK(w)
}

// DeleteInvalidRows delete all invalid rows one time
func (h *Handler) DeleteInvalidRows() {
	fmt.Println("go routine is started")
	ctx := context.Background()
	// 2 times in day check database and delete all "died" urls
	ticker := time.NewTicker(timeToUpdate * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		h.CheckBaseOnExpirationDate(ctx)
		err := h.db.DeleteInvalidRows(ctx)
		switch {
		case errors.Is(err, database.RowsToDeleteNotFoundError):
			fmt.Printf("failed to found rows to delete\n")
		case err != nil:
			fmt.Printf("failed to delete rows\n")
		default:
			fmt.Printf("all rows is_deleted=true was deleted successfully\n")
		}
	}
}

func (h *Handler) CheckBaseOnExpirationDate(ctx context.Context) {
	err := h.db.IsLinkExpired(ctx)
	switch {
	case errors.Is(err, database.NothingToUpdateError):
		fmt.Printf("failed to found rows to update\n")
	case err != nil:
		fmt.Printf("failed to update rows\n")
	default:
		fmt.Printf("update was success\n")
	}
}
