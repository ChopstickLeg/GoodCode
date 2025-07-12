package repository

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func GetRepoId(r *http.Request) (int64, error) {
	repoIdStr := r.URL.Query().Get("repoId")
	if repoIdStr == "" {
		return 0, errors.New("repoId query parameter is required")
	}

	repoId, err := strconv.ParseInt(repoIdStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid repoId: %w", err)
	}

	return repoId, nil
}
