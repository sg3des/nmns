package nmns

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
)

//Nmns is helper object - library
type Nmns struct {
	Scheme map[string]map[string]int
	Tables map[string]*TableStruct
}

//TableStruct is table object containing the data on it, the name of the table, the last index entry fields and facilities
type TableStruct struct {
	Name      string
	Fields    map[string]*FieldStruct
	IndexNum  int
	IndexFile *IndexStruct
}

//FieldStruct is field object
type FieldStruct struct {
	Name string
	File *os.File
	Size int
}

//Table is the function of preparation of the table to work with it
func (s *Nmns) Table(table string) *TableStruct {
	if t, ok := s.Tables[table]; ok {
		return t
	}
	return nil
}

//Write is writes data in the database and returns the id
func (t *TableStruct) Write(doc map[string]string) (int, error) {
	if t == nil {
		err := fmt.Errorf("%s", "Requested table does not exist")
		return 0, err
	}

	id := t.IndexNum
	for field, val := range doc {
		if _, ok := t.Fields[field]; !ok {
			//if field not exist - just skip it
			continue
		}
		maxlen := t.Fields[field].Size

		if len(val) > maxlen {
			val = val[0:maxlen]
		}

		var b []byte

		var off = int64(id * maxlen)
		misslen := maxlen - len(val)
		b = append([]byte(val), make([]byte, misslen)...)
		_, err := t.Fields[field].File.WriteAt(b, off)
		if err != nil {
			return 0, err
		}
	}
	t.IndexNum++
	err := t.IndexFile.Write(t.IndexNum)
	return id, err
}

//Read - Reades data on the given id
func (t *TableStruct) Read(id int) (doc map[string]string, err error) {
	if t == nil {
		err = fmt.Errorf("%s", "Requested table does not exist")
		return
	}

	doc = make(map[string]string)
	for name, field := range t.Fields {
		val := make([]byte, field.Size)
		_, err = field.File.ReadAt(val, int64(id*field.Size))

		if err != nil && err.Error() == "EOF" {
			err = nil
		}
		if err != nil {
			return
		}
		doc[name] = strings.Trim(string(val), "\x00")
	}
	// if empty(doc) {
	// 	err = fmt.Errorf("%s %d %s", "item by id", id, "is null")
	// }
	return
}

//Search data by filter, returns a list of id
func (t *TableStruct) Search(filter map[string]interface{}, limit ...int) (data map[int]map[string]string, err error) {
	if t == nil {
		err = fmt.Errorf("%s", "Requested table does not exist")
		return
	}

	if t.IndexNum == 0 {
		return
	}

	l := t.IndexNum
	if len(limit) != 0 && limit[0] <= t.IndexNum {
		l = limit[0]
	}

	var ids []int
	for id := 0; id < l; id++ {
		add := false
		for sfield, val := range filter {
			field := strings.Trim(sfield, "@")

			vallen := t.Fields[field].Size
			valread := make([]byte, vallen)

			_, end := t.Fields[field].File.ReadAt(valread, int64(id*vallen))

			if end != nil && end.Error() != "EOF" {
				return
			}
			valread = bytes.Trim(valread, "\x00")

			switch val.(type) {
			case string:
				add, err = compare(sfield[0:1], val.(string), valread)
			case []string:
				tmpadd := false
				for _, v := range val.([]string) {
					tmpadd, err = compare(sfield[0:1], v, valread)
					if tmpadd {
						add = tmpadd
						break
					}
				}
			}

			if err != nil {
				return
			}

			if !add {
				break
			}

		}
		if add {
			ids = append(ids, id)
		}
	}

	data = make(map[int]map[string]string)
	var doc map[string]string
	for _, id := range ids {
		doc, err = t.Read(id)
		if err != nil {
			return
		}
		if !empty(doc) {
			data[id] = doc
		}
	}

	return
}

func compare(a string, expr string, valread []byte) (add bool, err error) {
	switch a {
	case "@":
		add, err = match(expr, valread)
	default:
		add = eq(expr, string(valread))
	}
	return
}

func eq(a, b string) bool {
	if a == b {
		return true
	}
	return false
}

func match(expr string, b []byte) (m bool, err error) {
	reg, err := regexp.Compile(expr)
	if err != nil {
		return
	}
	m = reg.Match(b)
	return
}

//Delete - deletes data on the given id
func (t *TableStruct) Delete(id int) (err error) {
	if id > t.IndexNum {
		err = fmt.Errorf("%s", "id is missing")
		return
	}

	for _, field := range t.Fields {
		b := make([]byte, field.Size)
		_, err = field.File.WriteAt(b, int64(id*field.Size))
		if err != nil {
			return err
		}
	}
	return
}

//Update - updates data on the given id
func (t *TableStruct) Update(id int, doc map[string]string) (err error) {
	if t == nil {
		err = fmt.Errorf("%s", "Requested table does not exist")
		return
	}

	if id > t.IndexNum {
		err = fmt.Errorf("%s", "id is missing")
		return
	}

	for name, val := range doc {
		if _, ok := t.Fields[name]; !ok {
			//if field not exist - just skip it
			continue
		}
		field := t.Fields[name]

		if len(val) > field.Size {
			val = val[0:field.Size]
		}

		var b []byte
		misslen := field.Size - len(val)
		b = append([]byte(val), make([]byte, misslen)...)
		_, err = field.File.WriteAt(b, int64(id*field.Size))
	}

	return
}

//All - Return all data
func (t *TableStruct) All(limit ...int) (data map[int]map[string]string, err error) {
	if t == nil {
		err = fmt.Errorf("%s", "Requested table does not exist")
		return
	}

	var doc map[string]string
	data = make(map[int]map[string]string)
	l := t.IndexNum
	if len(limit) != 0 && limit[0] <= t.IndexNum {
		l = limit[0]
	}

	for id := 0; id < l; id++ {
		doc, err = t.Read(id)
		if empty(doc) {
			continue
		}
		data[id] = doc
	}
	return
}

//Truncate - clear values
func (t *TableStruct) Truncate(fields ...string) error {
	if t == nil {
		err := fmt.Errorf("%s", "Requested table does not exist")
		return err
	}

	switch len(fields) {
	case 0:
		for _, field := range t.Fields {
			err := field.File.Truncate(0)
			if err != nil {
				return err
			}
		}
	case 1:
		if err := t.Fields[fields[0]].File.Truncate(0); err != nil {
			return err
		}
	default:
		for _, field := range fields {
			err := t.Fields[field].File.Truncate(0)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func empty(m map[string]string) bool {
	for _, v := range m {
		if len(v) > 0 {
			return false
		}
	}
	return true
}
