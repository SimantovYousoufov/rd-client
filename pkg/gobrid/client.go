package gobrid

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"net/url"
	"strings"
)

type RealDebrid struct {
	Torrents   TorrentService
	Unrestrict UnrestrictService
}

type Client struct {
	token  string
	client *http.Client
}

func NewClient(token string) *RealDebrid {
	client := &Client{
		token:  token,
		client: &http.Client{},
	}

	return &RealDebrid{
		TorrentService{
			Client: client,
		},
		UnrestrictService{
			Client: client,
		},
	}
}

func (c *Client) makeRequest(method string, uri string, data url.Values) (*http.Request, error) {
	req, err := http.NewRequest(method, BASE_URL+uri, strings.NewReader(data.Encode()))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	return req, nil
}

func (c *Client) readResponse(res *http.Response, dest interface{}) error {
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		e := &ApiError{}

		json.Unmarshal(body, &e)

		return e
	}

	if dest == nil {
		return nil
	}

	json.Unmarshal(body, dest)

	return nil
}

type UnrestrictService struct {
	*Client
}

func (u *UnrestrictService) UnrestrictLink(link string) (*Download, error) {
	b := url.Values{}
	b.Add("link", link)

	req, err := u.makeRequest(POST, UNRESTRICT_LINK, b)

	if err != nil {
		return nil, err
	}

	response := &Download{}

	r, err := u.client.Do(req)

	if err != nil {
		return nil, err
	}

	err = u.readResponse(r, response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

type TorrentService struct {
	*Client
}

func (t *TorrentService) GetTorrent(id string) (*Torrent, error) {
	req, err := t.makeRequest(GET, GET_TORRENT+id, nil)

	if err != nil {
		return nil, err
	}

	r, err := t.client.Do(req)
	response := &Torrent{}

	if err != nil {
		return nil, err
	}

	err = t.readResponse(r, response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (t *TorrentService) AddTorrentMagnet(request *AddMagnetLinkRequest) (*AddMagnetLinkResponse, error) {
	b := url.Values{}
	b.Add("magnet", request.Magnet)

	req, err := t.makeRequest(POST, ADD_MAGNET, b)

	if err != nil {
		return nil, err
	}

	response := &AddMagnetLinkResponse{}

	r, err := t.client.Do(req)

	if err != nil {
		return nil, err
	}

	err = t.readResponse(r, response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (t *TorrentService) SelectFiles(tor *Torrent, ids [] int, all bool) error {
	fileIds := []string{}

	if ! all {
		for _, id := range ids {
			fileIds = append(fileIds, string(id))
		}
	} else {
		fileIds = append(fileIds, "all")
	}

	joined := strings.Join(fileIds, ",")

	b := url.Values{}
	b.Add("files", joined)

	req, err := t.makeRequest(POST, SELECT_FILES+tor.Id, b)

	if err != nil {
		return err
	}

	r, err := t.client.Do(req)

	if err != nil {
		return err
	}

	err = t.readResponse(r, nil)

	return err
}
