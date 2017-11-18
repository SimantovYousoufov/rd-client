package gobrid

import "time"

const (
	TOKEN_ENV = "REAL_DEBRID_TOKEN"
	BASE_URL  = "https://api.real-debrid.com/rest/1.0"
	JSON      = "application/json"
)

const (
	GET_TORRENT  = "/torrents/info/"
	ADD_MAGNET   = "/torrents/addMagnet"
	SELECT_FILES = "/torrents/selectFiles/"

	UNRESTRICT_LINK = "/unrestrict/link"
)

const (
	POLLING_RATE = 2 * time.Second

	STATUS_DOWNLOADED       = "downloaded"
	WAITING_FILES_SELECTION = "waiting_files_selection"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	PATCH  = "PATCH"
	DELETE = "DELETE"
)
