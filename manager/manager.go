package manager

import (
	"downloader/models"
	"fmt"
	"path"
	"strings"
)

type Manager struct {
	Downloads    []*models.Download
	BasePath     string
	OnChangeChan chan *models.Download
}

func NewManager(basePath string) *Manager {
	m := &Manager{
		Downloads:    make([]*models.Download, 0),
		BasePath:     basePath,
		OnChangeChan: make(chan *models.Download),
	}

	go m.listenToChanges()

	return m
}

func (m *Manager) AddDownload(url string) {
	fileName := m.getFileNameFromUrl(url)
	filePath := path.Join(m.BasePath, fileName)
	download := models.NewDownload(url, filePath, m.OnChangeChan)
	m.Downloads = append(m.Downloads, download)
	download.Init()
}

func (m *Manager) getFileNameFromUrl(url string) string {
	urlSegments := strings.Split(url, "/")
	lastUrlSegment := urlSegments[len(urlSegments)-1]
	// remove query strings
	nameSegments := strings.Split(lastUrlSegment, "?")
	return nameSegments[0]
}

func (m *Manager) listenToChanges() {
	for {
		d := <-m.OnChangeChan
		fmt.Printf("%s \n", d.StringRepresentation())
	}
}
