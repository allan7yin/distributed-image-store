package mappers

import (
	"bit-image/pkg/common"
	"bit-image/pkg/common/entities"
	"fmt"
)

func (image entities.Image) image_to_file_mapper(image entities.Image) common.File {
	file := common.File{
		id: mapToFileID(image.Base.id, image.Base.id),
		// Extract the hash from the Image metadata
		hash: image.ImageMetaData.Hash,
	}

	return file
}

func mapToFileID(userID, imageID string) string {
	// Assuming fileIDPrefix is a format string like "%s-%s" or similar
	return fmt.Sprintf(fileIDPrefix, userID, imageID)
}
