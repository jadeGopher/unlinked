package pagination

func GetPageCount(totalCount, pageSize int64) (pageCount int64) {
	pageCount = totalCount / pageSize
	if totalCount%pageSize != 0 {
		pageCount++
	}
	return pageCount
}
