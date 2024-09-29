package image

import (
	"bit-image/internal/postrges"
)

type ImageStore struct {
	DBHandler *postrges.ConnectionHandler
}

func NewImageStore(dbHandler *postrges.ConnectionHandler) *ImageStore {
	return &ImageStore{
		DBHandler: dbHandler,
	}
}
