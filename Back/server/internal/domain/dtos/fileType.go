package dtos

type FileTypeRequest struct {
	Name string `json:"name"`
}

// dto to recieve an array of file types to create
type BulkCreateFileTypeRequest struct {
	Types []FileTypeRequest `json:"types"`
}

type FileTypeResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
