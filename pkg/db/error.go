package db

import "fmt"

var (
	LinkNotFoundError = fmt.Errorf("link doesn't exist")

	SuffixNotFoundError = fmt.Errorf("short suffix doesn't exist")

	SecretKeyNotFoundError = fmt.Errorf("secret key doesn't exist")

	RowsToDeleteNotFoundError = fmt.Errorf("rows to delete don't exist")

	NothingToUpdateError = fmt.Errorf("nothing to update, all dates is valid")
)
