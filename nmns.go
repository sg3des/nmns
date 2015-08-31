package nmns

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
)

var (
	dir = "./data/"

	err error

	Index *os.File
)

type Nmns struct {
	Scheme map[string]map[string]int
	Files  map[string]map[string]*os.File
	Index  map[string]int
}

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

func writeIndex(dir string, newIndex map[string]int) error {
	fmt.Println(dir, newIndex)
	d, _ := json.Marshal(newIndex)

	b, err := ioutil.ReadFile(path.Join(dir, "index.json"))
	if err != nil || len(b) == 0 {
		err := ioutil.WriteFile(path.Join(dir, "index.json"), d, 0755)
		return err
	}

	var curIndex map[string]int
	err = json.Unmarshal(b, &curIndex)
	if err != nil {
		return err
	}

	for table, size := range newIndex {
		curIndex[table] = size
	}

	dataIndex, err := json.Marshal(curIndex)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path.Join(dir, "index.json"), dataIndex, 0755)
}

func Check(schemefile, dir string) error {
	newScheme, err := readScheme(schemefile)
	if err != nil {
		return err
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

func cp(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}

	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
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

	return writeIndex(dir, map[string]int{table: 0})
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

func Connect(dir string) (Nmns, error) {
	var s Nmns
	s.Scheme, err = readScheme(path.Join(dir, "scheme.json"))
	if err != nil {
		return s, err
	}

	s.Index, err = readIndex(path.Join(dir, "index.json"))
	if err != nil {
		return s, err
	}

	s.Files = make(map[string]map[string]*os.File)
	for table, values := range s.Scheme {
		s.Files[table] = make(map[string]*os.File)
		for field, _ := range values {
			f, err := os.OpenFile(path.Join(dir, table, field), os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
			if err != nil {
				return s, err
			}
			s.Files[table][field] = f
		}
	}
	return s, nil
}

func readScheme(file string) (scheme map[string]map[string]int, err error) {
	readfile, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	err = json.Unmarshal(readfile, &scheme)
	return
}

func readIndex(file string) (index map[string]int, err error) {
	Index, err = os.OpenFile(file, os.O_RDWR, 0755)
	if err != nil {
		return
	}

	readfile, err := ioutil.ReadAll(Index)
	if err != nil {
		return
	}

	err = json.Unmarshal(readfile, &index)
	return
}

func sortmap(m map[string]interface{}) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return
}

func keys(m map[string]map[string]int) []string {
	var keys []string
	for key, _ := range m {
		keys = append(keys, key)
	}
	return keys
}

func difference(slice1, slice2 []string) []string {
	var diff []string

	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}
