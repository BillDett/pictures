package db

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

// IdSet is a set of picture ids
type IdSet map[string]bool

// Database is the pictures metadata database
type Database struct {
	Title    string           `json:"title"`
	Created  string           `json:"created"`
	Filepath string           `json:"filepath"`
	Index    string           `json:"index"`
	Ids      IdSet            `json:"ids"`
	Labels   map[string]IdSet `json:"labels"`
}

// InitDatabaseFrom ensures that we have a properly initialized database even if filename does not exist
func InitDatabaseFrom(filename string, indexfilename string) (*Database, error) {
	db := NewDatabase("My Photo Archive", filename, indexfilename)
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File " + filename + " does not exist...creating")
		} else { // File permissions/corrupted, etc...
			return nil, err
		}
	} else {
		fmt.Printf("Opened %s, reading...\n", filename)
		//finfo, _ := f.Stat()
		dbBytes, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading in database file %v: %v\n", filename, err)
			return nil, err
		}
		fmt.Printf("We read in %d bytes\n", len(dbBytes))
		err = json.Unmarshal(dbBytes, db)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error unmarshalling database file %v: %v\n", filename, err)
			return nil, err
		}
	}
	return db, nil
}

// NewDatabase creates an instance of Database
func NewDatabase(title string, filepath string, index string) *Database {
	d := new(Database)
	d.Title = title
	d.Index = index
	d.Filepath = filepath
	d.Ids = make(IdSet)
	d.Labels = make(map[string]IdSet)
	d.Created = time.Now().Format("Mon Jan _2 2006 15:04:05")
	return d
}

var undated string

func init() {
	undated = "undated"
}

// ProcessIndexRecord creates/updates the corresponding entry in the database
func (db *Database) ProcessIndexRecord(record []string) {
	// [53884f829ad4e0cc /Users/bill/temp/pix/2004/CIMG0024.JPG /Users/bill/temp/pix/2004/thumbs/CIMG0024_thumb.JPG 2004:04:25 02:20:24]
	key := record[0]
	//path := record[1]
	//thumbPath := record[2]
	dateTime := record[3]

	// Add to the list of all ids
	db.Ids[key] = true

	// Set the date labels if we can
	if dateTime != "NONE" {
		t, err := time.Parse("2006:01:02 15:04:05", dateTime)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding dateTime from index: %v: %v\n", dateTime, err)
			return
		}
		year := strconv.FormatInt(int64(t.Year()), 10)
		ids := db.Labels[year]
		if ids == nil {
			set := make(IdSet)
			db.Labels[year] = set
		}
		db.Labels[year][key] = true
		month := t.Month().String()
		ids = db.Labels[month]
		if ids == nil {
			set := make(IdSet)
			db.Labels[month] = set
		}
		db.Labels[month][key] = true

	} else {
		if db.Labels[undated] == nil {
			db.Labels[undated] = make(IdSet)
		}
		db.Labels[undated][key] = true
	}
}

// Union calculates the union of two IdSets
func (s1 IdSet) Union(s2 IdSet) IdSet {
	result := make(IdSet)
	for k := range s1 {
		result[k] = true
	}
	for k := range s2 {
		result[k] = true
	}
	return result
}

// Intersection calculates the difference between to IdSets
func (s1 IdSet) Intersection(s2 IdSet) IdSet {
	result := make(IdSet)
	// Loop over the smaller set
	if len(s1) < len(s2) {
		for k := range s1 {
			if s2[k] {
				result[k] = true
			}
		}
	} else {
		for k := range s2 {
			if s1[k] {
				result[k] = true
			}
		}
	}
	return result
}

/*
	Custom methods to marshal/unmarshal the IdSet as a slice of values

*/
func (ids IdSet) MarshalJSON() ([]byte, error) {
	// We only want to save the keys as a slice
	keys := make([]string, len(ids))
	i := 0
	for k := range ids {
		keys[i] = k
		i++
	}
	return json.Marshal(keys)
}
func (ids *IdSet) UnmarshalJSON(data []byte) error {
	result := make(IdSet)
	var keys []string
	if err := json.Unmarshal(data, &keys); err != nil {
		return err
	}
	for _, k := range keys {
		result[k] = true
	}
	*ids = result
	return nil
}

// Dump will print out the database in human readable format
func (db *Database) Dump() {
	json, err := json.MarshalIndent(db, "", "  ")
	//json, err := json.Marshal(db)
	if err == nil {
		fmt.Printf("%s\n", json)
	} else {
		fmt.Fprintf(os.Stderr, "Error marshalling: %v\n", err.Error())
	}
}

// Save will write out the database in human readable format
func (db *Database) Save() error {
	json, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling: %v\n", err.Error())
		return err
	}
	err = ioutil.WriteFile(db.Filepath, json, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating database file %v: %v\n", db.Filepath, err)
		return err
	}
	return nil
}

// Photodatabase manages the database at the name given
func Photodatabase(filename string, indexfilename string) error {
	/*
		Create or update a photo metadata database using the index file.
	*/

	db, err := InitDatabaseFrom(filename, indexfilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing database from %v: %v\n", filename, err)
		return err
	}

	// Process the Index
	idxfile, err := os.Open(indexfilename)
	if err != nil {
		return err
	}
	r := csv.NewReader(idxfile)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		db.ProcessIndexRecord(record)
	}
	idxfile.Close()

	db.Save()

	return nil
}
