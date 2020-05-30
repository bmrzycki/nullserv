package main

import "time"

type NullFile struct {
	data    []byte
	content string
}

// NotFoundFiles to return HTTP 404, not found
var NotFoundFiles = map[string]bool{
	"7z":   true,
	"avi":  true,
	"bin":  true,
	"bz2":  true,
	"doc":  true,
	"docx": true,
	"exe":  true,
	"gz":   true,
	"iso":  true,
	"java": true,
	"mov":  true,
	"mp3":  true,
	"mp4":  true,
	"pdf":  true,
	"ppt":  true,
	"pptx": true,
	"lz":   true,
	"lzma": true,
	"sh":   true,
	"swf":  true,
	"rar":  true,
	"tar":  true,
	"tb2":  true,
	"tbz":  true,
	"tbz2": true,
	"tgz":  true,
	"txz":  true,
	"webp": true,
	"webm": true,
	"xls":  true,
	"xlsx": true,
	"xz":   true,
	"zip":  true,
}

// AltSuffix maps alternative extensions to other valid NullFiles
var AltSuffix = map[string]string{
	"asp":  "html",
	"aspx": "html",
	"cgi":  "html",
	"dll":  "html",
	"do":   "html",
	"htm":  "html",
	"jpeg": "jpg",
	"m4a":  "mp4",
	"m4b":  "mp4",
	"m4p":  "mp4",
	"m4r":  "mp4",
	"m4v":  "mp4",
	"res":  "reset",
	"stat": "stats",
	"tif":  "tiff",
	"ver":  "version",
	"vers": "version",
	"xht":  "xhtml",
}

// NullFiles maps extensions to a minimal, valid, bytestream. The special case
// where stream == nil sends no data because an empty file is valid. Regardless
// of data the associated mime type is always sent.
var NullFiles = map[string]NullFile{
	"": NullFile{nil, "text/plain"}, // Catch-all empty suffix
	"bmp": NullFile{[]byte{
		'\x42', '\x4d', '\x1e', '\x00', '\x00', '\x00', '\x00', '\x00',
		'\x00', '\x00', '\x1a', '\x00', '\x00', '\x00', '\x0c', '\x00',
		'\x00', '\x00', '\x01', '\x00', '\x01', '\x00', '\x01', '\x00',
		'\x18', '\x00', '\x00', '\x00', '\xff', '\x00',
	}, "image/bmp"},
	"css": NullFile{nil, "text/css"},
	"csv": NullFile{nil, "text/csv"},
	"gif": NullFile{[]byte{
		'\x47', '\x49', '\x46', '\x38', '\x39', '\x61', '\x01', '\x00',
		'\x01', '\x00', '\x80', '\x00', '\x00', '\xff', '\xff', '\xff',
		'\x00', '\x00', '\x00', '\x21', '\xf9', '\x04', '\x01', '\x00',
		'\x00', '\x00', '\x00', '\x2c', '\x00', '\x00', '\x00', '\x00',
		'\x01', '\x00', '\x01', '\x00', '\x00', '\x02', '\x02', '\x44',
		'\x01', '\x00', '\x3b',
	}, "image/gif"},
	"html": NullFile{[]byte("<!DOCTYPE html><title>x</title>"),
		"text/html"},
	"ico": NullFile{[]byte{
		'\x00', '\x00', '\x01', '\x00', '\x01', '\x00', '\x01', '\x01',
		'\x00', '\x00', '\x01', '\x00', '\x18', '\x00', '\x30', '\x00',
		'\x00', '\x00', '\x16', '\x00', '\x00', '\x00', '\x28', '\x00',
		'\x00', '\x00', '\x01', '\x00', '\x00', '\x00', '\x02', '\x00',
		'\x00', '\x00', '\x01', '\x00', '\x18', '\x00', '\x00', '\x00',
		'\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00',
		'\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00',
		'\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00',
		'\xff', '\x00', '\x00', '\x00', '\x00', '\x00',
	}, "image/x-icon"},
	"jpg": NullFile{[]byte{
		'\xff', '\xd8', '\xff', '\xe0', '\x00', '\x10', '\x4a', '\x46',
		'\x49', '\x46', '\x00', '\x01', '\x01', '\x01', '\x00', '\x01',
		'\x00', '\x01', '\x00', '\x00', '\xff', '\xdb', '\x00', '\x43',
		'\x00', '\xff', '\xff', '\xff', '\xff', '\xff', '\xff', '\xff',
		'\xff', '\xff', '\xff', '\xff', '\xff', '\xff', '\xff', '\xff',
		'\xff', '\xff', '\xff', '\xff', '\xff', '\xff', '\xff', '\xff',
		'\xff', '\xff', '\xff', '\xff', '\xff', '\xff', '\xff', '\xff',
		'\xff', '\xff', '\xff', '\xff', '\xff', '\xff', '\xff', '\xff',
		'\xff', '\xff', '\xff', '\xff', '\xff', '\xff', '\xff', '\xff',
		'\xff', '\xff', '\xff', '\xff', '\xff', '\xff', '\xff', '\xff',
		'\xff', '\xff', '\xff', '\xff', '\xff', '\xff', '\xff', '\xff',
		'\xff', '\xff', '\xc0', '\x00', '\x0b', '\x08', '\x00', '\x01',
		'\x00', '\x01', '\x01', '\x01', '\x11', '\x00', '\xff', '\xc4',
		'\x00', '\x14', '\x00', '\x01', '\x00', '\x00', '\x00', '\x00',
		'\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00',
		'\x00', '\x00', '\x00', '\x03', '\xff', '\xc4', '\x00', '\x14',
		'\x10', '\x01', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00',
		'\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00',
		'\x00', '\x00', '\xff', '\xda', '\x00', '\x08', '\x01', '\x01',
		'\x00', '\x00', '\x3f', '\x00', '\x47', '\xff', '\xd9',
	}, "image/jpg"},
	"js":   NullFile{nil, "application/javascript"},
	"json": NullFile{[]byte("{}"), "application/json"},
	"png": NullFile{[]byte{
		'\x89', '\x50', '\x4e', '\x47', '\x0d', '\x0a', '\x1a', '\x0a',
		'\x00', '\x00', '\x00', '\x0d', '\x49', '\x48', '\x44', '\x52',
		'\x00', '\x00', '\x00', '\x01', '\x00', '\x00', '\x00', '\x01',
		'\x08', '\x06', '\x00', '\x00', '\x00', '\x1f', '\x15', '\xc4',
		'\x89', '\x00', '\x00', '\x00', '\x0a', '\x49', '\x44', '\x41',
		'\x54', '\x78', '\x9c', '\x63', '\x00', '\x01', '\x00', '\x00',
		'\x05', '\x00', '\x01', '\x0d', '\x0a', '\x2d', '\xb4', '\x00',
		'\x00', '\x00', '\x00', '\x49', '\x45', '\x4e', '\x44', '\xae',
		'\x42', '\x60', '\x82',
	}, "image/png"},
	"tiff": NullFile{[]byte{
		'\x4d', '\x4d', '\x00', '\x2a', '\x00', '\x00', '\x00', '\x08',
		'\x00', '\x07', '\x01', '\x00', '\x00', '\x03', '\x00', '\x00',
		'\x00', '\x01', '\x00', '\x01', '\x00', '\x00', '\x01', '\x01',
		'\x00', '\x03', '\x00', '\x00', '\x00', '\x01', '\x00', '\x01',
		'\x00', '\x00', '\x01', '\x06', '\x00', '\x03', '\x00', '\x00',
		'\x00', '\x01', '\x00', '\x00', '\x00', '\x00', '\x01', '\x11',
		'\x00', '\x03', '\x00', '\x00', '\x00', '\x01', '\x00', '\x00',
		'\x00', '\x00', '\x01', '\x17', '\x00', '\x03', '\x00', '\x00',
		'\x00', '\x01', '\x00', '\x01', '\x00', '\x00', '\x01', '\x1a',
		'\x00', '\x05', '\x00', '\x00', '\x00', '\x01', '\x00', '\x00',
		'\x00', '\x64', '\x01', '\x1b', '\x00', '\x05', '\x00', '\x00',
		'\x00', '\x01', '\x00', '\x00', '\x00', '\x64', '\x00', '\x00',
		'\x00', '\x00', '\x20', '\x20', '\x20', '\x62', '\x79', '\x20',
		'\x61', '\x6c', '\x6f', '\x6b',
	}, "image/tiff"},
	"stats": NullFile{[]byte("{ \"ok\": 0 }"), "application/json"},
	"svg": NullFile{[]byte{
		'\x3c', '\x73', '\x76', '\x67', '\x20', '\x78', '\x6d', '\x6c',
		'\x6e', '\x73', '\x3d', '\x22', '\x68', '\x74', '\x74', '\x70',
		'\x3a', '\x2f', '\x2f', '\x77', '\x77', '\x77', '\x2e', '\x77',
		'\x33', '\x2e', '\x6f', '\x72', '\x67', '\x2f', '\x32', '\x30',
		'\x30', '\x30', '\x2f', '\x73', '\x76', '\x67', '\x22', '\x2f',
		'\x3e',
	}, "image/svg+xml"},
	"wasm": NullFile{[]byte{
		'\x00', '\x61', '\x73', '\x6d', '\x01', '\x00', '\x00', '\x00',
	}, "application/wasm"},
	"xhtml": NullFile{[]byte{
		'\x3c', '\x68', '\x74', '\x6d', '\x6c', '\x20', '\x78', '\x6d',
		'\x6c', '\x6e', '\x73', '\x3d', '\x22', '\x68', '\x74', '\x74',
		'\x70', '\x3a', '\x2f', '\x2f', '\x77', '\x77', '\x77', '\x2e',
		'\x77', '\x33', '\x2e', '\x6f', '\x72', '\x67', '\x2f', '\x31',
		'\x39', '\x39', '\x39', '\x2f', '\x78', '\x68', '\x74', '\x6d',
		'\x6c', '\x22', '\x2f', '\x3e',
	}, "application/xhtml+xml"},
	"xml": NullFile{[]byte{
		'\x3c', '\x21', '\x44', '\x4f', '\x43', '\x54', '\x59', '\x50',
		'\x45', '\x20', '\x5f', '\x5b', '\x3c', '\x21', '\x45', '\x4c',
		'\x45', '\x4d', '\x45', '\x4e', '\x54', '\x20', '\x5f', '\x20',
		'\x45', '\x4d', '\x50', '\x54', '\x59', '\x3e', '\x5d', '\x3e',
		'\x3c', '\x5f', '\x2f', '\x3e',
	}, "application/xml"},
}

func GenVersion() {
	NullFiles["version"] = NullFile{[]byte("{\n" +
		"  \"build_date\": \"" + BuildInfo["build_date"] + "\",\n" +
		"  \"commit_date\": \"" + BuildInfo["commit_date"] + "\",\n" +
		"  \"reset_date\": \"" +
		time.Now().Format(time.RFC1123Z) + "\",\n" +
		"  \"sha\": \"" + BuildInfo["sha"] + "\",\n" +
		"  \"version\": \"" + BuildInfo["version"] + "\"\n" +
		"}",
	), "application/json"}
}
