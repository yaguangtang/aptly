package http

import (
	"context"
	"os"
	"path/filepath"

	"github.com/aptly-dev/aptly/aptly"
	"github.com/aptly-dev/aptly/utils"
)

// DownloadTemp starts new download to temporary file and returns File
//
// Temporary file would be already removed, so no need to cleanup
func DownloadTemp(ctx context.Context, headers string, downloader aptly.Downloader, url string) (*os.File, error) {
	return DownloadTempWithChecksum(ctx, headers, downloader, url, nil, false)
}

// DownloadTempWithChecksum is a DownloadTemp with checksum verification
//
// Temporary file would be already removed, so no need to cleanup
func DownloadTempWithChecksum(ctx context.Context, headers string, downloader aptly.Downloader, url string, expected *utils.ChecksumInfo, ignoreMismatch bool) (*os.File, error) {
	tempdir, err := os.MkdirTemp(os.TempDir(), "aptly")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempdir)

	tempfile := filepath.Join(tempdir, "buffer")

	if expected != nil && downloader.GetProgress() != nil {
		downloader.GetProgress().InitBar(expected.Size, true, aptly.BarMirrorUpdateDownloadIndexes)
		defer downloader.GetProgress().ShutdownBar()
	}

	err = downloader.DownloadWithChecksum_withheader(ctx, headers, url, tempfile, expected, ignoreMismatch)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(tempfile)
	if err != nil {
		return nil, err
	}

	return file, nil
}
