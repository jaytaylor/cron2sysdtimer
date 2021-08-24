package systemd

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
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

var _templates_service_tmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x44\x8d\x31\x0e\xc2\x30\x0c\x45\xf7\x9c\xc2\x27\xa8\xb8\x40\x07\x04\xac\x2c\x6d\xa7\xaa\x43\x15\x7e\x85\x87\x24\x55\xe2\x22\x90\x95\xbb\xa3\x26\x54\x4c\xb6\xbf\xff\xd3\x1b\x07\xcf\x32\x99\x2b\x92\x8d\xbc\x0a\x07\xdf\xaa\x52\x73\x9f\x1d\x28\x67\x4a\x88\x2f\xb6\xa0\xcd\xb3\xa8\x12\x2f\xd4\x9c\x17\x41\xa4\x9c\x4d\x59\x4a\xfb\x88\x54\x09\xfe\xb1\xff\xcc\xd8\x55\x72\x32\x3d\x3b\x84\x4d\x3a\x99\xa3\x74\xb0\xed\xc9\xdc\xde\xb0\xe5\x2c\xf0\x25\x38\x37\x57\xaa\xff\xac\x68\x83\x47\x7a\x86\xc3\x36\xa4\x2a\xdb\x67\xa9\xff\x82\xbf\xea\x1b\x00\x00\xff\xff\xc5\x00\xf6\x4c\xc2\x00\x00\x00")

func templates_service_tmpl() ([]byte, error) {
	return bindata_read(
		_templates_service_tmpl,
		"templates/service.tmpl",
	)
}

var _templates_timer_tmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8a\x0e\xcd\xcb\x2c\x89\xe5\x72\x49\x2d\x4e\x2e\xca\x2c\x28\xc9\xcc\xcf\xb3\xad\xae\x56\xd0\xf3\x4b\xcc\x4d\x55\xa8\xad\x55\x28\xc9\xcc\x4d\x2d\x52\x28\xcd\xcb\x2c\xe1\xe2\x8a\x0e\x01\x71\x62\xb9\xfc\xf3\x9c\x13\x73\x52\xf3\x52\x12\x8b\xc0\x4a\x9d\x8b\xf2\xf3\x8a\x0b\x52\x93\x15\x6a\x6b\xb9\x00\x01\x00\x00\xff\xff\x6d\x06\x2f\x26\x4e\x00\x00\x00")

func templates_timer_tmpl() ([]byte, error) {
	return bindata_read(
		_templates_timer_tmpl,
		"templates/timer.tmpl",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
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
var _bindata = map[string]func() ([]byte, error){
	"templates/service.tmpl": templates_service_tmpl,
	"templates/timer.tmpl": templates_timer_tmpl,
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
	Func func() ([]byte, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"templates": &_bintree_t{nil, map[string]*_bintree_t{
		"service.tmpl": &_bintree_t{templates_service_tmpl, map[string]*_bintree_t{
		}},
		"timer.tmpl": &_bintree_t{templates_timer_tmpl, map[string]*_bintree_t{
		}},
	}},
}}
