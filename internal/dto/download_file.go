package dto

type DownloadFileRequest struct {
	File_id  string
	Owner_id string
}

type DownloadFileResponse struct {
	Path      string
	Mime_type string
	Name      string
}
