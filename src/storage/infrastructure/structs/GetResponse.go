package structs

type GetResponse struct {
	Id           uint   `json:"id"`
	FileName     string `json:"file_name"`
	ContentType  string `json:"content_type"`
	Size         int64  `json:"size"`
	URL          string `json:"url"`
	CollectionID string `json:"collection_id"`
}
