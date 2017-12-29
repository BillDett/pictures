# Pictures

Generic tool to build a simple metadata database from a directory of pictures.  The goal is an extremely simple data structure that is easily machine readable and stays usable regardless of language, operating system or software changes.

## Creating an index and database 

The first thing we need is an index of all pictures located in the given directory:

```
$ ./pictures build index --from=/Users/bill/Pictures > index.txt
```

Generates an index file ```index.txt``` which contains a record for each picture encountered (recursively) in the given directory.  Each picture has a thumbnail image created and stored in a ```/thumbs``` directory created inside each directory seen.

The index file format is very simple:

PICTURE_ID,FILEPATH,THUMBPATH,DATETIME

e.g.
```
21d5f22733d7931d,/Users/bill/Pictures/Andree.jpg,/Users/bill/Pictures/thumbs/Andree_thumb.jpg,2003:08:05 22:55:15
```

If there is no embedded EXIF data containing DATETIME, then the picture will have ```NONE``` in that column.  Supported image formats are JPEG, TIFF, and PNG.

The index represents all pictures, but there is no organization.  To organize the pictures we need a database of labels.  We can generate a default database which will label the pictures by timestamp as follows:

```
$ ./pictures build database --index=index.txt --from=database.json
```

If ```database.json``` does not exist, it will be created.  If it already exists, then any existing non-default labels will not be affected.

The database file is a basic JSON document with labels (tags) as object members pointing to lists of PICTURE_ID values.  Each list of PICTURE_ID acts like a set of values.  The default labels are generated for Year and Month.

It is safe to re-run a ```pictures build database``` on an existing database file even if the index file has changed.

## Hosting an index and database

Once an index and database have been created, the ```pictures``` executable can be used to host them interactively with a small web site that allows viewing and management of labels.

```
$ ./pictures host --index=index.txt --database=database.json --port=8080
```

Opening ```http://localhost:8080``` in a browser will show the main page of images showing the latest available month's worth of pictures.

