//Package nmns an lighweight file relational database.
//example:
//
//db, err := nmns.Connect(dir)
//
//doc := map[string]string{"city":"Moscow","country":"Russia"}
//id, err := db("world").Write(doc)
//
//doc, err := db("world").Read(id) - where id is interge
//
//filter := map[string]interface{}{"country":"Russia"}
//ids, err := db("world").Search(filter)
//
package nmns

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

//Creates a database dir with file,scheme and other, or overwrite an existing - can be omitted
func Init(schemefile, dir string) error {
	scheme, err := readScheme(schemefile)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	d, _ := json.Marshal(scheme)
	if err = ioutil.WriteFile(path.Join(dir, "scheme.json"), d, 0755); err != nil {
		return err
	}

	for table, values := range scheme {
		err = createTable(dir, table, values)
		if err != nil {
			return err
		}
	}

	return nil
}

//Check the relevance of the database structure. In the process of developing your application, you can modify the configuration json file - add tables, fields, or change the fields size - all the changes are made automatically without data loss.
func Check(schemefile, dir string) error {
	newScheme, err := readScheme(schemefile)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path.Join(dir, "scheme.json")); err != nil {
		if err := Init(schemefile, dir); err != nil {
			return err
		}
		return nil
	}

	curScheme, err := readScheme(path.Join(dir, "scheme.json"))
	if err != nil {
		return err
	}

	for table, newValues := range newScheme {
		if curValues, ok := curScheme[table]; ok {
			for newv, news := range newValues {
				if curs, ok := curValues[newv]; ok {
					if news != curs {
						//if different sizes
						err = reindex(dir, table, newv, curs, news)
						if err != nil {
							return err
						}
					}
				} else {
					//if cell in table does not esist
					err = createCell(dir, table, newv)
					if err != nil {
						return err
					}
				}
			}
		} else {
			//if table does not exist
			err = createTable(dir, table, newValues)
			if err != nil {
				return err
			}
		}
	}

	data, err := json.Marshal(newScheme)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(dir, "scheme.json"), data, 0755)
	if err != nil {
		return err
	}

	return nil
}

func createTable(dir, table string, values map[string]int) error {
	err := os.MkdirAll(path.Join(dir, table), 0755)
	if err != nil {
		return err
	}

	for field, _ := range values {
		err := createCell(dir, table, field)
		if err != nil {
			return err
		}
	}

	return Index(dir, table).Write(0)
}

func createCell(dir, table, name string) error {
	f, err := os.Create(path.Join(dir, table, name))
	if err != nil {
		return err
	}
	return f.Close()
}

func reindex(dir, table, name string, curs, news int) error {
	curPath := path.Join(dir, table, name)
	tmpPath := path.Join(dir, table, ".tmp")

	openfile, err := os.Open(curPath)
	if err != nil {
		return err
	}

	tmpfile, err := os.Create(tmpPath)
	if err != nil {
		return err
	}

	read := make([]byte, curs)
	for i := 0; err != nil; i++ {
		_, err = openfile.ReadAt(read, int64(curs*i))
		if curs > news {
			read = read[0:news]
		} else {
			misslen := news - curs
			read = append(read, make([]byte, misslen)...)
		}
		_, err := tmpfile.Write(read)
		if err != nil {
			return err
		}
	}

	if err := openfile.Close(); err != nil {
		return err
	}
	if err := tmpfile.Close(); err != nil {
		return err
	}
	if err := os.Remove(curPath); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, curPath); err != nil {
		return err
	}
	return nil
}

func readScheme(file string) (scheme map[string]map[string]int, err error) {
	readfile, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	err = json.Unmarshal(readfile, &scheme)
	return
}

//connected to a database, and returns a connection fuction
func Connect(dir string) (func(string) *TableStruct, error) {
	var s Nmns
	var err error
	s.Scheme, err = readScheme(path.Join(dir, "scheme.json"))
	if err != nil {
		return s.Table, err
	}

	s.Tables = make(map[string]*TableStruct)
	for table, values := range s.Scheme {

		t := &TableStruct{Name: table}
		t.Fields = make(map[string]*FieldStruct)
		t.IndexFile = Index(dir, table)
		id, err := t.IndexFile.Read()
		if err != nil {
			return s.Table, err
		}
		t.IndexNum = id

		for field, size := range values {
			f, err := os.OpenFile(path.Join(dir, table, field), os.O_RDWR, 0755)
			if err != nil {
				return s.Table, err
			}
			t.Fields[field] = &FieldStruct{Name: field, Size: size, File: f}
		}
		s.Tables = map[string]*TableStruct{table: t}
	}
	return s.Table, nil
}
