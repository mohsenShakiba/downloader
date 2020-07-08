package models

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type DownloadStatus string

type Download struct {
	Url               string
	Status            DownloadStatus
	Progress          *Progress
	Error             *string
	HttpResponse      *http.Response
	File              *os.File
	FilePath          string
	OnChangeChan      chan<- *Download
	MeanDownloadSpeed int64
}

const (
	Pending    = "PENDING"
	InProgress = "DOWNLOADING"
	Failed     = "FAILED"
	Stopped    = "STOPPED"
	Finished   = "Finished"
)

func NewDownload(url string, filePath string, onChangeChan chan<- *Download) *Download {
	return &Download{
		Url:          url,
		Status:       Pending,
		FilePath:     filePath,
		OnChangeChan: onChangeChan,
	}
}

func (d *Download) MarkAsFinished() {
	d.Status = Finished
	d.Progress = nil
	d.Error = nil
	d.HttpResponse = nil
	d.notifyChange()
}

func (d *Download) Stop() {
	d.Status = Stopped
	d.notifyChange()
}

func (d *Download) Failed(err string) {
	d.Status = Failed
	d.Error = &err
	d.Progress = nil
	d.HttpResponse = nil
	d.notifyChange()
}

func (d *Download) Init() {
	f, err := os.OpenFile(d.FilePath, os.O_APPEND|os.O_CREATE, 0666)

	if err != nil {
		d.Failed(err.Error())
		return
	}

	d.File = f
	d.Status = InProgress
	d.Error = nil

	d.resetMeanDownloadSpeed()
	d.initDownload(d.Url)
}

func (d *Download) AddProgress(inc int) {
	if d.Status != InProgress {
		return
	}

	d.Progress.DownloadedLength += int64(inc)
	d.MeanDownloadSpeed += int64(inc)
	d.notifyChange()
}

func (d *Download) notifyChange() {
	d.OnChangeChan <- d
}

func (d *Download) initDownload(url string) {
	resp, err := http.Get(url)

	if err != nil {
		d.Failed(err.Error())
		return
	}

	defer resp.Body.Close()

	d.Progress = &Progress{
		TotalLength:      resp.ContentLength,
		DownloadedLength: 0,
	}

	if resp.StatusCode != http.StatusOK {
		err := fmt.Sprintf("unexpected status code %d", resp.StatusCode)
		d.Failed(err)
		return
	}

	var l sync.Mutex

	l.Lock()

	go func() {
		for {
			s := d.Progress.RemainingSegmentSize()

			if s <= 0 {
				l.Unlock()
				break
			}

			b := make([]byte, s)
			_, _ = io.ReadFull(resp.Body, b)
			_, _ = d.File.Write(b)
			d.AddProgress(len(b))
		}
	}()

	l.Lock()

}

func (d *Download) resetMeanDownloadSpeed() {
	ticker := time.NewTicker(time.Second)

	go func() {
		for {
			<-ticker.C
			d.MeanDownloadSpeed = 0
		}
	}()
}

func (d *Download) StringRepresentation() string {

	var progress float32

	switch d.Status {
	case Finished:
		progress = 100
	case InProgress:
		progress = d.Progress.Progress()
	case Stopped:
		progress = d.Progress.Progress()
	}

	return fmt.Sprintf("%s    %s    %s    %f %d", d.Url, d.FilePath, d.Status, progress, d.MeanDownloadSpeed)
}
