package nmns

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"testing"
)

var (
	t_dir        = "testdata"
	t_schemefile = "db.json"
	t_dbname     = "data"

	t_scheme map[string]map[string]int

	t_db  func(string) *TableStruct
	id    int
	t_err error
	t_doc = map[string]string{"city": "Moscow", "country": "Russia"}
)

const dbjson = `{"world":{"city":32,"country":32}}`

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
	t_db, t_err = Connect(path.Join(t_dir, t_dbname))
	if t_err != nil {
		t.Error(t_err)
	}

}

func TestWrite(t *testing.T) {
	id, err := t_db("world").Write(t_doc)
	if err != nil {
		t.Error(err)
	}
	t.Log(id)
}

func TestUpdate(t *testing.T) {
	doc := map[string]string{"city": "London", "country": "United Kingdom"}
	err := t_db("world").Update(id, doc)
	if err != nil {
		t.Error(err)
	}
}

func TestRead(t *testing.T) {
	doc, err := t_db("world").Read(id)
	if err != nil {
		t.Error(err)
	}
	t.Log(doc)
}

func TestSearch(t *testing.T) {
	filter := map[string]interface{}{"city": "Moscow"}
	ids, err := t_db("world").Search(filter)
	if err != nil {
		t.Error(err)
	}
	t.Log(ids)
	// need to append
}

func TestDelete(t *testing.T) {
	if err := t_db("world").Delete(id); err != nil {
		t.Error(err)
	}
	doc, err := t_db("world").Read(id)
	if err != nil {
		t.Error(err)
	}
	t.Log(doc)
}

func TestAll(t *testing.T) {
	data, err := t_db("world").All()
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}

func TestTruncate(t *testing.T) {
	err := t_db("world").Truncate("country")
	if err != nil {
		t.Error(err)
	}
	// need to append
}

func BenchmarkWrite(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t_db("world").Write(t_doc)
	}
}

func BenchmarkRead(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t_db("world").Read(rand.Intn(t_db("world").IndexNum))
	}
}
