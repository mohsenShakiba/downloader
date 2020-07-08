package models

const (
	defaultDownloadSegmentSize = 1024
)

type Progress struct {
	TotalLength      int64
	DownloadedLength int64
}

func (p *Progress) RemainingSegmentSize() int {
	remainingSize := p.TotalLength - p.DownloadedLength
	if remainingSize > defaultDownloadSegmentSize {
		return defaultDownloadSegmentSize
	}
	return int(remainingSize)
}

func (p *Progress) Progress() float32 {
	return float32(float64(p.DownloadedLength) / float64(p.TotalLength))
}
