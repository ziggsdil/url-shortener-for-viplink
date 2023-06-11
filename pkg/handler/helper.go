package handler

import (
	"encoding/json"
	"fmt"
	"io"
)

func (h *Handler) parseBody(from io.ReadCloser, to interface{}) error {
	body, err := io.ReadAll(from)
	if err != nil || len(body) == 0 {
		return fmt.Errorf("empty body")
	}
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, to)
	if err != nil {
		return err
	}

	return nil
}
