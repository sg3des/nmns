package nmns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func (s *Nmns) Insert(table string, doc map[string]string) (id int, err error) {
	id = s.Index[table]
	for field, val := range doc {
		maxlen := s.Scheme[table][field]

		if len(val) > maxlen {
			val = val[0:maxlen]
		}

		var b []byte
		misslen := maxlen - len(val)
		b = append([]byte(val), make([]byte, misslen)...)
		_, err = s.Files[table][field].WriteAt(b, int64(id*maxlen))
	}

	s.Incrementindex(table)
	return
}

func (s *Nmns) Incrementindex(table string) {
	s.Index[table]++
	d, err := json.Marshal(s.Index)
	if err != nil {
		panic(err)
	}

	Index.Truncate(0)
	_, err = Index.WriteAt(d, 0)
	if err != nil {
		panic(err)
	}
}

func (s *Nmns) Read(table string, id int) (doc map[string]string, err error) {
	doc = make(map[string]string)
	for field, file := range s.Files[table] {
		vallen := s.Scheme[table][field]
		val := make([]byte, vallen)
		_, err = file.ReadAt(val, int64(id*vallen))

		if err != nil && err.Error() == "EOF" {
			err = nil
		}
		if err != nil {
			return
		}
		doc[field] = string(val)
	}
	return
}

func (s *Nmns) Search(table string, filter map[string]interface{}) (ids []int, err error) {

	for id := 0; id < s.Index[table]; id++ {
		add := false
		for sfield, val := range filter {
			field := strings.Trim(sfield, "@")

			vallen := s.Scheme[table][field]
			valread := make([]byte, vallen)
			_, end := s.Files[table][field].ReadAt(valread, int64(id*vallen))
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

			if !add {
				break
			}

		}
		if add {
			ids = append(ids, id)
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

func (s *Nmns) Delete(table string, id int) (err error) {
	if id > s.Index[table] {
		err = fmt.Errorf("%s", "id is missing")
		return
	}

	for field, size := range s.Scheme[table] {
		var n int

		r := make([]byte, size)
		n, err = s.Files[table][field].ReadAt(r, int64(id*size))
		fmt.Println("read", n, r)

		b := make([]byte, size)

		n, err = s.Files[table][field].WriteAt(b, int64(id*size))
		fmt.Println("writ", n, b)
		// s.Files[table][field].Seek(offset, whence)
	}
	return
}

func (s *Nmns) Update(table string, id int, doc map[string]string) (err error) {
	if id > s.Index[table] {
		err = fmt.Errorf("%s", "id is missing")
		return
	}

	for field, val := range doc {
		maxlen := s.Scheme[table][field]

		if len(val) > maxlen {
			val = val[0:maxlen]
		}

		var b []byte
		misslen := maxlen - len(val)
		b = append([]byte(val), make([]byte, misslen)...)
		_, err = s.Files[table][field].WriteAt(b, int64(id*maxlen))
	}

	return
}
