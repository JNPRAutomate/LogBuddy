package logbuddy

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"os"
	"time"
	"io/ioutil"
	"path"
	"path/filepath"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindata_file_info struct {
	name string
	size int64
	mode os.FileMode
	modTime time.Time
}

func (fi bindata_file_info) Name() string {
	return fi.name
}
func (fi bindata_file_info) Size() int64 {
	return fi.size
}
func (fi bindata_file_info) Mode() os.FileMode {
	return fi.mode
}
func (fi bindata_file_info) ModTime() time.Time {
	return fi.modTime
}
func (fi bindata_file_info) IsDir() bool {
	return false
}
func (fi bindata_file_info) Sys() interface{} {
	return nil
}

var _static_js_messages_js = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xb4\x55\x4d\x4f\xdb\x4c\x10\x3e\xdb\xbf\x62\xe4\x0b\x09\x6f\x64\x87\xb7\xe2\x92\x28\x87\x34\xa0\x96\xb6\x40\x9b\x44\xea\x11\x19\x7b\x62\xb6\x98\x5d\x77\x77\x0d\x45\x90\xff\xde\xd9\x0f\x1b\x1b\xaa\x16\x54\x7a\xb2\x3d\x3b\xf3\x3c\xcf\x7c\xec\x38\x49\x16\x82\x2b\x9d\x72\xad\x60\x23\x24\x5c\xa1\x52\x69\x81\xa0\x6f\x2b\x54\xe1\x75\x2a\xe1\xe8\xe4\x68\x7d\x76\xbc\x7a\x07\x33\x18\x4f\xad\xe5\x60\xbe\x9e\x7b\xcb\x9e\xb3\x2c\x0f\xbf\x78\xc3\xff\xce\x70\xb8\x5c\x7a\xc3\x1b\x67\x58\xad\xe7\xcb\x06\x66\x6f\xec\x81\xe6\x8b\x8f\x67\xfd\x83\xbd\xc6\xfb\xf4\x73\x03\xb8\xbf\x3f\x0d\xc3\x64\x37\x0c\x3e\xac\x4e\x4f\x18\xcf\xa0\x10\x65\xca\x0b\x28\xd9\x25\x82\x38\xff\x86\x99\x06\xa5\x65\x9d\xe9\x5a\x62\x18\x06\x8b\x92\x21\xd7\xc7\x3e\x8f\x09\xdc\x85\x41\xb0\xa6\x6c\x26\xf4\x5c\x5c\xa4\x9c\x63\x69\x5e\x57\x28\xaf\x51\x52\xf2\x1b\x56\xd0\xf7\x36\xdc\x4d\x88\x26\xe9\x47\x17\xc8\x51\xa6\x1a\x15\xa4\xc0\xf1\x06\x32\x7b\xea\x6b\x64\x95\x9e\xe0\x8d\x0f\x51\x05\xa9\xdd\xd4\x3c\xd3\x4c\xf0\x81\xa9\xdf\x08\x32\xc7\x37\x02\x65\xd9\x32\xcb\x36\x34\x92\xd8\x06\x06\xfe\x14\x66\x33\xe0\x75\x59\xc2\xfd\x3d\x74\x4c\x35\xcf\x71\xc3\x38\xe6\xd6\x3f\x68\x4f\x60\x6c\xd4\x06\x12\x29\x5f\x0e\xa6\x28\x31\xa5\xcf\x78\xc1\x36\xb7\x83\xbb\xc8\x10\x47\x13\x70\xfc\x91\x8f\x22\x43\x2b\x25\xea\x6a\xa1\x83\xee\xe7\x76\x38\x0d\xb7\x61\x93\x98\x2b\xd1\xd3\xac\x72\xa5\x59\x65\x1f\x95\x90\xda\xca\xf3\x6a\x1a\x7a\xeb\x17\xb1\x2a\x9a\x38\xdf\xc8\x38\xda\x0f\xf3\xb2\x6d\x38\xda\x1e\xf0\x2e\xc9\x90\xf0\xf4\x05\x53\x71\xe6\x0e\x4c\xe1\xbf\xe2\xf9\x4a\x64\x97\xa8\x07\xd1\x8d\x9a\x24\x49\x04\xff\x41\x2e\xb2\xfa\x8a\x2a\x1f\x97\x22\x4b\x4d\x64\x7c\x21\x94\xa6\x83\x28\x29\x45\xa1\x22\xca\x25\x48\x12\x4a\x3b\x2f\x91\x3a\x27\x14\x76\x70\x63\xc1\xad\xa9\x4b\x8c\xd7\x2e\x99\xa0\x45\x2e\x50\x1f\x96\x68\x5e\xdf\xde\x1e\xe5\x83\x88\x70\x0f\x52\x9d\x46\xc3\x98\x51\x31\xe5\xfb\xf5\xf1\x27\x02\xd8\x31\x19\xa0\xc5\x80\x85\x01\xcd\x77\xa6\xb6\x47\x2d\xbb\x9f\x18\xd5\x17\xd0\xdc\xb5\x5f\x49\x20\x17\x25\x4a\xa4\xdc\x0a\x63\x8c\x73\xa2\x1d\x92\xfd\xca\x8e\x99\xed\x7a\x95\x4a\x85\x0f\x87\x44\x69\xa7\x8a\x3c\x62\x53\x7f\x9a\xa1\x59\x7b\x55\x1d\x68\x0f\xd5\xf8\x79\x05\x06\xf8\xa5\x49\x77\xc2\x0d\xf3\x16\xb0\xa4\x6a\x3e\x11\xd0\xbb\xe2\x5e\x45\x92\x48\x2c\x98\xd2\x34\x5b\xb4\x79\xa4\xc6\xdc\x0f\xe1\xbf\x93\xe1\xf7\x91\x17\xf0\x97\x1c\xfd\xde\xa2\x94\x42\x3e\xea\xac\xb5\xbd\xc6\x68\x35\xdd\x7d\x34\x4f\xa2\x42\xde\x67\x34\x96\xd7\x9e\xe5\x53\xc2\x74\x93\xbc\x35\xab\x71\x65\x7a\xd5\x6c\xc2\xd2\xf4\x8f\xd3\xd6\x69\x5a\xf7\x70\x97\xe3\x4a\x0a\x2d\x4c\xe9\x63\x1b\xf2\xcc\x35\xd2\xbd\xf5\xb1\x42\x9e\x0f\xba\xbb\x75\xd0\x0e\xd1\x68\x3c\x6a\x77\x93\x43\x72\x40\x0d\xce\xd0\xed\x30\xa3\x57\x54\x24\x97\x03\xfe\x20\xb1\x7f\x96\x2a\xaa\xa7\x4a\xfd\xd6\xb4\xa5\xfc\xbd\x3a\xf7\xc3\x1a\x35\x01\x8d\x88\x25\x7e\xaf\x51\xe9\xe7\xeb\xa0\x80\x45\xbb\xea\x5f\xac\xc3\xff\x88\xfb\x32\x7e\x06\x00\x00\xff\xff\xd9\x10\x1e\x1c\xe4\x07\x00\x00")

func static_js_messages_js_bytes() ([]byte, error) {
	return bindata_read(
		_static_js_messages_js,
		"static/js/messages.js",
	)
}

func static_js_messages_js() (*asset, error) {
	bytes, err := static_js_messages_js_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "static/js/messages.js", size: 2020, mode: os.FileMode(436), modTime: time.Unix(1420702737, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"static/js/messages.js": static_js_messages_js,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() (*asset, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"static": &_bintree_t{nil, map[string]*_bintree_t{
		"js": &_bintree_t{nil, map[string]*_bintree_t{
			"messages.js": &_bintree_t{static_js_messages_js, map[string]*_bintree_t{
			}},
		}},
	}},
}}

// Restore an asset under the given directory
func RestoreAsset(dir, name string) error {
        data, err := Asset(name)
        if err != nil {
                return err
        }
        info, err := AssetInfo(name)
        if err != nil {
                return err
        }
        err = os.MkdirAll(_filePath(dir, path.Dir(name)), os.FileMode(0755))
        if err != nil {
                return err
        }
        err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
        if err != nil {
                return err
        }
        err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
        if err != nil {
                return err
        }
        return nil
}

// Restore assets under the given directory recursively
func RestoreAssets(dir, name string) error {
        children, err := AssetDir(name)
        if err != nil { // File
                return RestoreAsset(dir, name)
        } else { // Dir
                for _, child := range children {
                        err = RestoreAssets(dir, path.Join(name, child))
                        if err != nil {
                                return err
                        }
                }
        }
        return nil
}

func _filePath(dir, name string) string {
        cannonicalName := strings.Replace(name, "\\", "/", -1)
        return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

