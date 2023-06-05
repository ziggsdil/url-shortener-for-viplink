package handler

import (
	"net/http"
)

func (h *Handler) Healthcheck(w http.ResponseWriter, _ *http.Request) {
	h.renderer.RenderOK(w)
}
