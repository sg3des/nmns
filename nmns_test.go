package nmns

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

var (
	t_dir        = "testdata"
	t_schemefile = "db.json"
	t_dbname     = "data"

	t_scheme map[string]map[string]int

	t_s Nmns
	id  int
)

const dbjson = `{"Users":{"name":16,"age":3}}`

func init() {
	if err := os.MkdirAll(path.Join(t_dir, t_dbname), 0755); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(path.Join(t_dir, t_schemefile), []byte(dbjson), 0755); err != nil {
		panic(err)
	}

	if err := json.Unmarshal([]byte(dbjson), &t_scheme); err != nil {
		panic(err)
	}

}

func TestInit(t *testing.T) {
	if err := Init(path.Join(t_dir, t_schemefile), path.Join(t_dir, t_dbname)); err != nil {
		t.Error(err)
	}
	if _, err := os.Stat(path.Join(t_dir, t_dbname)); err != nil {
		t.Error("failed to create database directory")
	}

	for table, fields := range t_scheme {
		if _, err := os.Stat(path.Join(t_dir, t_dbname, table)); err != nil {
			t.Error("failed to create table directory (" + table + ")")
		}
		for field, _ := range fields {
			if _, err := os.Stat(path.Join(t_dir, t_dbname, table, field)); err != nil {
				t.Error("failed to create file by field (" + field + ")")
			}
		}
	}
}

func TestCheck(t *testing.T) {
	if err := Check(path.Join(t_dir, t_schemefile), path.Join(t_dir, t_dbname)); err != nil {
		t.Error(err)
	}
}

func TestConnect(t *testing.T) {
	t_s, err = Connect(path.Join(t_dir, t_dbname))
	if err != nil {
		t.Error(err)
	}

}

func TestInsert(t *testing.T) {
	doc := map[string]string{"name": "Valeriy", "age": "99"}
	id, err = t_s.Insert("Users", doc)
	if err != nil {
		t.Error(err)
	}
	t.Log(id)
}

func TestUpdate(t *testing.T) {
	doc := map[string]string{"name": "Valeriy", "age": "01"}
	err := t_s.Update("Users", id, doc)
	if err != nil {
		t.Error(err)
	}
}

func TestRead(t *testing.T) {
	doc, err := t_s.Read("Users", id)
	if err != nil {
		t.Error(err)
	}
	t.Log(doc)
}

func TestSearch(t *testing.T) {
	filter := map[string]interface{}{"name": "Valeriy"}
	ids, err := t_s.Search("Users", filter)
	if err != nil {
		t.Error(err)
	}
	t.Log(ids)
	// need to append
}

func TestDelete(t *testing.T) {
	if err := t_s.Delete("Users", id); err != nil {
		t.Error(err)
	}
	doc, err := t_s.Read("Users", id)
	if err != nil {
		t.Error(err)
	}
	t.Log(doc)
}

func TestAll(t *testing.T) {
	data, err := t_s.All("Users")
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}

func TestTruncate(t *testing.T) {
	err := t_s.Truncate("Users", "age")
	if err != nil {
		t.Error(err)
	}
	// need to append
}
