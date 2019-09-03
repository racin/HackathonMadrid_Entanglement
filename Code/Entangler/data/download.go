package data

type DownloadRequest struct {
	Key string
}

type DownloadResponse struct {
	*DownloadRequest
	Value []byte
}
