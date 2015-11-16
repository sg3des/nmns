package nmns

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

//Index return index for table
func Index(dir, table string) *IndexStruct {
	var indexFile *os.File
	var err error

	indexFile, err = os.OpenFile(path.Join(dir, table+".id"), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
		// indexFile, err = os.Create(path.Join(dir, table+".id"))
	}
	return &IndexStruct{File: indexFile}
}

//IndexStruct is wrapper for os.File
type IndexStruct struct {
	File *os.File
}

//Write data on disk
func (i *IndexStruct) Write(id int) error {
	_, err := i.File.WriteAt([]byte(strconv.Itoa(id)), 0)
	return err
}

//Read data from disk
func (i *IndexStruct) Read() (int, error) {
	data, err := ioutil.ReadAll(i.File)
	if err != nil {
		return 0, err
	}

	id, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}

	return id, nil
}
