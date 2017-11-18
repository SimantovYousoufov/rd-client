package gobrid

import "fmt"

type ApiError struct {
	ErrorName string `json:"error"`
	ErrorDetails string `json:"error_details"`
	ErrorCode int `json:"error_code"`
}

func (err *ApiError) Error() string {
	return fmt.Sprintf("%s: %s", err.ErrorName, err.ErrorDetails)
}

type AddMagnetLinkRequest struct {
	Magnet string `json:"magnet"`
}

type AddMagnetLinkResponse struct {
	Id string `json:"id"`
	Uri string `json:"uri"`
}

type Torrent struct {
	Id string `json:"id"`
	Filename string `json:"filename"`
	OriginalFilename string `json:"original_filename"`
	Hash string `json:"hash"`
	Bytes int64 `json:"bytes"`
	OriginalBytes int64 `json:"original_bytes"`
	Host string `json:"host"`
	Split int `json:"split"`
	Progress int `json:"progress"`
	Status string `json:"status"` // @todo enum
	Added string `json:"added"`
	Files []File `json:"files"`
	Links []string `json:"links"`
	Ended string `json:"ended"`
	Speed int `json:"speed"`
	Seeders int `json:"seeders"`
}

func (t *Torrent) IsWaitingForFileSelection() bool {
	return t.Status == WAITING_FILES_SELECTION
}

func (t *Torrent) IsDownloaded() bool {
	return t.Status == STATUS_DOWNLOADED
}

type File struct {
	Id int `json:"id"`
	Path string `json:"path"`
	Bytes int `json:"bytes"`
	Selected int `json:"selected"`
}

type Download struct {
	Id string `json:"id"`
	Filename string `json:"filename"`
	MimeType string `json:"mime_type"`
	Filesize int64 `json:"filesize"`
	Link string `json:"link"`
	Host string `json:"host"`
	Chunks int `json:"chunks"`
	Crc int `json:"crc"`
	Download string `json:"download"`
	Streamable int `json:"streamable"`
}