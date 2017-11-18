package main

import (
	"github.com/simantovyousoufov/rd-client/pkg/gobrid"
	"time"
	"errors"
	"os"
	"io"
	"net/http"
	"fmt"
	"math"
	"io/ioutil"
	"strings"
)

type TorrentDownload struct {
	statuses []string
	done       bool
	err        error
	progress   int
	link       string
	torrent    *gobrid.Torrent
}

func (t *TorrentDownload) AddStatus(status string) {
	t.statuses = append(t.statuses, status)
}

type MagnetCommand struct {
	links     []string
	c         *gobrid.RealDebrid
	downloads []*TorrentDownload
}

type MagnetCommandConfig struct {
	client *gobrid.RealDebrid
	//links  []string
}

func NewMagnetCommand(opts MagnetCommandConfig) *MagnetCommand {
	m := &MagnetCommand{
		c:         opts.client,
		links:     []string{},
		downloads: []*TorrentDownload{},
	}

	return m
}

func (m *MagnetCommand) AddLink(link string) {
	m.links = append(m.links, link)
}

func (m *MagnetCommand) AddLinksFromFile(path string) error {
	f, err := os.Open(path)

	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(f)

	if err != nil {
		return err
	}

	links := strings.Split(string(data), "\n")

	for _, l := range links {
		m.AddLink(l)
	}

	return nil
}

func (m *MagnetCommand) Download() {
	for _, link := range m.links {
		d := &TorrentDownload{
			done:     false,
			progress: 0,
			link:     link,
		}

		m.downloads = append(m.downloads, d)

		go func() {
			d.err = m.retrieve(d)
		}()
	}

	ticker := time.NewTicker(2 * time.Second)

	for range ticker.C {
		allComplete := true

		for _, td := range m.downloads {
			if td.err != nil {
				td.AddStatus(fmt.Sprintf("Received an error for torrent %s: %s.\n", td.link[:10], td.err))
				continue
			}

			if ! td.done {
				allComplete = false
				break
			}
		}

		if allComplete {
			ticker.Stop()
			break
		}

		ClearScreen()
		m.PrintTorrentStates()
	}

	fmt.Printf("Done!\n")
}

func (m *MagnetCommand) PrintTorrentStates() {
	for i, td := range m.downloads {
		fmt.Printf("==============================%d==============================\n", i)

		if td.torrent != nil {
			fmt.Printf("Name: %s\n", td.torrent.OriginalFilename)
		} else {
			fmt.Printf("Link: %s\n", td.link[:10])
		}

		fmt.Printf("Download Progress: %d/100\n", td.progress)
		fmt.Printf("Is complete: %t\n", td.done)

		if len(td.statuses) > 0 {
			fmt.Printf("Last Status: %s\n", td.statuses[len(td.statuses) - 1])
		}

		fmt.Println("=============================================================")
	}
}

func (m *MagnetCommand) retrieve(td *TorrentDownload) error {
	t, err := m.AddTorrent(td.link)

	if err != nil {
		return err
	}

	td.torrent = t
	td.AddStatus(fmt.Sprintf("Succesfully created torrent: %s.\n", t.OriginalFilename))

	td.AddStatus(fmt.Sprintf("Awaiting waiting_files_selector or downloaded."))
	td.torrent, err = m.AwaitTorrentStatus(td.torrent, []string{gobrid.WAITING_FILES_SELECTION, gobrid.STATUS_DOWNLOADED}, 10)

	td.AddStatus(fmt.Sprintf("Ready to select all files."))

	if err != nil {
		return err
	}

	err = m.SelectAllFilesForTorrent(td.torrent)

	if err != nil {
		return err
	}

	td.AddStatus(fmt.Sprintf("Selected all files."))
	td.AddStatus(fmt.Sprintf("Awaiting downloaded."))

	td.torrent, err = m.AwaitTorrentStatus(td.torrent, []string{gobrid.STATUS_DOWNLOADED}, 2)

	if err != nil {
		return err
	}

	td.AddStatus(fmt.Sprintf("Torrent downloaded, ready to download to local."))

	download, err := m.c.Unrestrict.UnrestrictLink(td.torrent.Links[0])

	out, err := os.Create(download.Filename)

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(POLLING_RATE)

		for range ticker.C {
			td.progress, err = m.GetDownloadProgress(download, out)

			if td.progress == 100 {
				break
			}

			if err != nil {
				td.AddStatus(fmt.Sprintf("Encountered an error: %s.", err))
			}

			td.AddStatus(fmt.Sprintf("%d%% downloaded.", td.progress))
		}
	}()

	m.DownloadTorrent(download, out)

	td.done = true

	return nil
}

func (m *MagnetCommand) AddTorrent(magnetUrl string) (*gobrid.Torrent, error) {
	magnet, err := m.c.Torrents.AddTorrentMagnet(&gobrid.AddMagnetLinkRequest{
		Magnet: magnetUrl,
	})

	if err != nil {
		return nil, err
	}

	torrent, err := m.c.Torrents.GetTorrent(magnet.Id)

	if err != nil {
		return nil, err
	}

	return torrent, err
}

// @todo status enum?
func (m *MagnetCommand) AwaitTorrentStatus(t *gobrid.Torrent, statuses []string, maxAttempts int) (*gobrid.Torrent, error) {
	var err error = nil
	ticker := time.NewTicker(2 * time.Second)
	attempts := 0

	for range ticker.C {
		if attempts > maxAttempts {
			ticker.Stop()

			return nil, errors.New("Timed out waiting for status")
		}

		t, err := m.c.Torrents.GetTorrent(t.Id)

		if err != nil {
			ticker.Stop()
			break
		}

		for _, s := range statuses {
			if t.Status == s {
				ticker.Stop()

				return t, nil
			}
		}

		attempts++
	}

	return nil, err
}

func (m *MagnetCommand) SelectAllFilesForTorrent(t *gobrid.Torrent) error {
	// @todo await for `waiting_files_selection` status here?

	return m.c.Torrents.SelectFiles(t, []int{}, true)
}

func (m *MagnetCommand) DownloadTorrent(download *gobrid.Download, dst *os.File) error {
	res, err := http.Get(download.Download)
	HandleError(err)

	defer res.Body.Close()
	defer dst.Close()

	_, err = io.Copy(dst, res.Body)

	return err
}

func (m *MagnetCommand) GetDownloadProgress(d *gobrid.Download, f *os.File) (int, error) {
	info, err := f.Stat()

	if err != nil {
		return 0, err
	}

	ratio := float64(info.Size()) / float64(d.Filesize)

	return int(math.Floor(ratio * 100)), nil
}
