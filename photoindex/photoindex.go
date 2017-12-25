package main

// Create an index file from a root directory of photos
// Process all sub-directories
// Generate thumbnails for all images

import (
	"bytes"
	"fmt"
	"github.com/cozy/goexif2/exif"
	"github.com/nfnt/resize"
	"golang.org/x/image/tiff"
	"hash/fnv"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// Remember what directories we've seen to simplify thumbnail creation
var dirs map[string]bool

func hash(s string) string {
	h := fnv.New64a()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum64())
}

// TODO: We could be a bit smarter here- lots of duplicate code
func generateThumbnail(f *os.File, path string, ext string) ([]byte, error) {
	f.Seek(0, 0)
	var img image.Image
	var data []byte
	var err error
	if ext == ".JPG" {
		img, err = jpeg.Decode(f)
		if err != nil {
			fmt.Printf("Error seen decoding %v: %v\n", path, err)
			return data, err
		}
		m := resize.Resize(200, 0, img, resize.Lanczos3)
		buf := new(bytes.Buffer)
		err = jpeg.Encode(buf, m, nil)
		if err != nil {
			fmt.Printf("Error seen encoding %v: %v\n", path, err)
			return data, err
		}
		data = buf.Bytes()
	} else if ext == ".TIFF" || ext == ".TIF" {
		img, err = tiff.Decode(f)
		if err != nil {
			fmt.Printf("Error seen decoding %v: %v\n", path, err)
			return data, err
		}
		m := resize.Resize(200, 0, img, resize.Lanczos3)
		buf := new(bytes.Buffer)
		err = tiff.Encode(buf, m, nil)
		if err != nil {
			fmt.Printf("Error seen encoding %v: %v\n", path, err)
			return data, err
		}
		data = buf.Bytes()
	} else if ext == ".PNG" {
		img, err = png.Decode(f)
		if err != nil {
			fmt.Printf("Error seen decoding %v: %v\n", path, err)
			return data, err
		}
		m := resize.Resize(200, 0, img, resize.Lanczos3)
		buf := new(bytes.Buffer)
		err = png.Encode(buf, m)
		if err != nil {
			fmt.Printf("Error seen encoding %v: %v\n", path, err)
			return data, err
		}
		data = buf.Bytes()
	}
	return data, err
}

func main() {
	thumbsDirName := "thumbs"
	thumbSuffix := "_thumb"
	if len(os.Args) != 2 {
		fmt.Println("usage: photoindex root\nGenerates a photo index from root directory to stdout")
		os.Exit(1)
	}
	// TODO: We should stat() root to make sure it exists!
	dirs = make(map[string]bool)
	err := filepath.Walk(os.Args[1], func(path string, info os.FileInfo, walkerr error) error {
		if walkerr == nil {
			base := filepath.Base(path)
			if base == thumbsDirName { // We should completely skip the thumbnail directory
				return filepath.SkipDir
			}
			if info.IsDir() { // If we've walked onto a directory itself, just ignore it
				return nil
			}
			extension := filepath.Ext(path)
			ext := strings.ToUpper(extension)
			if ext != ".JPG" && ext != ".PNG" && ext != ".TIFF" && ext != ".TIF" {
				return nil
			}
			// Create the thumb directory if needed, add it to the map of directories
			dir := filepath.Dir(path)
			thumbDir := filepath.Join(dir, thumbsDirName)
			if !dirs[dir] {
				if err := os.Mkdir(thumbDir, 0777); err != nil && !os.IsExist(err) {
					fmt.Printf("Error creating thumb directory %v: %v\n", thumbDir, err)
					os.Exit(1)
				}
				dirs[dir] = true
			}

			thekey := hash(path)

			// Extract DateTime & get a thumbnail
			f, err := os.Open(path)
			if err != nil {
				fmt.Printf("Error opening %v; %v\n", path, err)
				os.Exit(1)
			}
			x, exiferr := exif.Decode(f)
			dtString := "NONE"
			var data []byte
			if exiferr == nil {
				// Pull out DateTime if there
				tag, err := x.Get(exif.DateTimeOriginal)
				if err == nil {
					dtString, err = tag.StringVal()
				}
			}
			// Get a thumbnail - we MUST have one
			if exiferr == nil {
				// We had EXIF, so the might already be a thumbnail...
				data, err = x.JpegThumbnail() // Use the one in the JPG
			}
			if exiferr != nil || err != nil {
				// We didn't get a JpegThumbnail for some reason so generate one.
				data, err = generateThumbnail(f, path, ext)
			}
			if err != nil {
				// Couldn't get a thumbnail- stop everything
				fmt.Printf("Cannot generate a thumbnail for %s\n", path)
				os.Exit(1)
			}

			// Write out the thumbnail
			thumbFilename := strings.TrimSuffix(base, extension) + thumbSuffix + extension
			thumbFilePath := filepath.Join(thumbDir, thumbFilename)
			tf, err := os.Create(thumbFilePath)
			if err != nil {
				fmt.Printf("Error writing thumbnail to %v: %v\n", thumbFilePath, err)
				os.Exit(1)
			}
			tf.Write(data)
			tf.Close()

			fmt.Printf("%s,%s,%s,%s\n", thekey, path, thumbFilePath, dtString)
			f.Close()
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}
