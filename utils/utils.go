package utils

import (
	"user/internal/entities"
)

func MetaDataInfo(metaData *entities.MetaData) *entities.MetaData {
	if metaData.Total < 1 {
		return nil
	}
	if (metaData.CurrentPage)*(metaData.PerPage) < metaData.Total {
		metaData.Next = metaData.CurrentPage + 1
	}
	if metaData.CurrentPage > 1 {
		metaData.Prev = metaData.CurrentPage - 1
	}

	return metaData
}
