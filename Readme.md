Pictures

Generic tool to build a simple metadata database from a directory of pictures (possibly nested)

Database Format

Index - maps a unique id to each picture we find to
    * Path to picture
    * Path to thumbnail
    * JSON of EXIF data

Labels - maps a label to a set of qualifying images.  We use labels to represent Years and Months.

Goal is to be simple & understandable but machine readable.

https://github.com/deckarep/golang-set
https://github.com/cozy/goexif2
https://github.com/nfnt/resize
https://github.com/segmentio/ksuid
https://github.com/nanobox-io/golang-scribble


{
    "title": "My photo archive",
    "basePath" : "/Volumes/bills_files/photos export",
    "created": "2018:01:07",
    "index": "photo_index.txt",
    "labels": {
        "2003": [ "12345", "82732", "33829" ],
        "2004": [ "882923", "094039", "8208208" ],
        "January": [ "12345", "838202" ]
    }
}

Single executable- ```pictures``` which does the following
* Can be used to initialize/update the database from the index
    * First time creation
    * Generate all default labels (replace any existing)
    * Preserve any existing non-default labels
* Can be used to host the database & provide a REST API so it can be interrogated & new labels added
* Eventually it should also be able to generate the photo index I suppose...?

```
	pictures

	usage:
		pictures build index myindex.txt --from=/mypix
		pictures build database mypix.json --from=myindex.txt
		pictures host mypix.json --port=8080
```

Labels can be JSON based...

But Index should be a simple CSV flat file

```
INDEX,FILEPATH,THUMBPATH,DATETIME
12345,/2003/DSC00023.jpg,/2003/thumbs/DSC0023_thumb.jpg,2003:01:05 23:33:02
45678,/2005/DSC23023.jpg,/2005/thumbs/DSC23023_thumb.jpg,2005:03:12 11:38:02
```

photoindex - executable that generates/re-creates the index file from a directory of pictures.

Each ID should be a hash of the FILEPATH- so that we can avoid re-indexing files that we've already seen.  Runnning the photoindex executable on same directory twice should produce identical output.
Use a SHA1 hash- https://gobyexample.com/sha1-hashes
SHA1 might be too slow- can also try the fnv library in standard golang: https://golang.org/pkg/hash/fnv/
    the hash.New64a() looks promising...

Run a quick test- fetch all pathnames & see if we get any collisions!  Try different hashing algorithms.


```

process() {
    For each file {
        create thumb directory if it doesn't exist
        generate id via SHA1 hash of path
        if *.jpg {
            Get EXIF
            Create thumbnail
         } else if *.tiff or *.png {
             Create a thumbnail (using go image library)
         }
         Create entry in index with id, path, thumb (if there) and exif (if there)
         if EXIF exists {
             Pull year and month from DateTime
             Label this id for year and month tags
         } else {
             Label this id with "nodate" tag
         }
    }
    For each directory {
        process()
    }
}

```