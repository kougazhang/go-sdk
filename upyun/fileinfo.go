package upyun

import (
	"net/http"
	"os"
	"strings"
	"time"
)

type FileInfo struct {
	FileName    string
	FileSize    int64
	ContentType string
	IsFileDir   bool
	IsEmptyDir  bool
	MD5         string
	Time        time.Time

	Meta map[string]string

	/* image information */
	ImgType   string
	ImgWidth  int64
	ImgHeight int64
	ImgFrames int64
}

// the full path of file
func (i FileInfo) Name() string {
	return i.FileName
}

// length in bytes for regular files; system-dependent for others
func (i FileInfo) Size() int64 {
	return i.FileSize
}

// file mode bits
func (i FileInfo) Mode() os.FileMode {
	return 0
}

// modification time
func (i FileInfo) ModTime() time.Time {
	return i.Time
}

// abbreviation for Mode().IsFileDir()
func (i FileInfo) IsDir() bool {
	return i.IsFileDir
}

// underlying data source (can return nil)
func (i FileInfo) Sys() interface{} {
	return nil
}

/*
  Content-Type: image/gif
  ETag: "dc9ea7257aa6da18e74505259b04a946"
  x-upyun-file-type: GIF
  x-upyun-height: 379
  x-upyun-width: 500
  x-upyun-frames: 90
*/
func parseHeaderToFileInfo(header http.Header, getinfo bool) *FileInfo {
	fInfo := &FileInfo{}
	for k, v := range header {
		lk := strings.ToLower(k)
		if strings.HasPrefix(lk, "x-upyun-meta-") {
			if fInfo.Meta == nil {
				fInfo.Meta = make(map[string]string)
			}
			fInfo.Meta[lk] = v[0]
		}
	}

	if getinfo {
		// HTTP HEAD
		fInfo.FileSize = parseStrToInt(header.Get("x-upyun-file-size"))
		fInfo.IsFileDir = header.Get("x-upyun-file-type") == "folder"
		fInfo.Time = time.Unix(parseStrToInt(header.Get("x-upyun-file-date")), 0)
		fInfo.ContentType = header.Get("Content-Type")
		fInfo.MD5 = header.Get("Content-MD5")
	} else {
		fInfo.FileSize = parseStrToInt(header.Get("Content-Length"))
		fInfo.ContentType = header.Get("Content-Type")
		fInfo.MD5 = strings.ReplaceAll(header.Get("Content-Md5"), "\"", "")
		if fInfo.MD5 == "" {
			fInfo.MD5 = strings.ReplaceAll(header.Get("Etag"), "\"", "")
		}
		lastM := header.Get("Last-Modified")
		t, err := http.ParseTime(lastM)
		if err == nil {
			fInfo.Time = t
		}
		fInfo.ImgType = header.Get("x-upyun-file-type")
		fInfo.ImgWidth = parseStrToInt(header.Get("x-upyun-width"))
		fInfo.ImgHeight = parseStrToInt(header.Get("x-upyun-height"))
		fInfo.ImgFrames = parseStrToInt(header.Get("x-upyun-frames"))
	}
	return fInfo
}
