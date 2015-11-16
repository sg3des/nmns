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
	tdir        = "testdata"
	tschemefile = "db.json"
	tdbname     = "data"

	tscheme map[string]map[string]int

	tdb  func(string) *TableStruct
	id   int
	terr error
	tdoc = map[string]string{"city": "Moscow", "country": "Russia"}
)

const dbjson = `{"world":{"city":32,"country":32}}`

func init() {
	if err := os.MkdirAll(path.Join(tdir, tdbname), 0755); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(path.Join(tdir, tschemefile), []byte(dbjson), 0755); err != nil {
		panic(err)
	}

	if err := json.Unmarshal([]byte(dbjson), &tscheme); err != nil {
		panic(err)
	}

}

func TestInit(t *testing.T) {
	if err := Init(path.Join(tdir, tschemefile), path.Join(tdir, tdbname)); err != nil {
		t.Error(err)
	}
	if _, err := os.Stat(path.Join(tdir, tdbname)); err != nil {
		t.Error("failed to create database directory")
	}

	for table, fields := range tscheme {
		if _, err := os.Stat(path.Join(tdir, tdbname, table)); err != nil {
			t.Error("failed to create table directory (" + table + ")")
		}
		for field := range fields {
			if _, err := os.Stat(path.Join(tdir, tdbname, table, field)); err != nil {
				t.Error("failed to create file by field (" + field + ")")
			}
		}
	}
}

func TestCheck(t *testing.T) {
	if err := Check(path.Join(tdir, tschemefile), path.Join(tdir, tdbname)); err != nil {
		t.Error(err)
	}
}

func TestConnect(t *testing.T) {
	tdb, err = Connect(path.Join(tdir, tdbname))
	if err != nil {
		t.Error(err)
	}

}

func TestWrite(t *testing.T) {
	id, err = tdb("world").Write(tdoc)
	if err != nil {
		t.Error(err)
	}
	t.Log(id)
}

func TestUpdate(t *testing.T) {
	doc := map[string]string{"city": "London", "country": "United Kingdom"}
	err := tdb("world").Update(id, doc)
	if err != nil {
		t.Error(err)
	}
}

func TestRead(t *testing.T) {
	doc, err := tdb("world").Read(id)
	if err != nil {
		t.Error(err)
	}
	t.Log(doc)
}

func TestSearch(t *testing.T) {
	filter := map[string]interface{}{"city": "Moscow"}
	ids, err := tdb("world").Search(filter)
	if err != nil {
		t.Error(err)
	}
	t.Log(ids)
	// need to append
}

func TestDelete(t *testing.T) {
	if err := tdb("world").Delete(id); err != nil {
		t.Error(err)
	}
	doc, err := tdb("world").Read(id)
	if err != nil {
		t.Error(err)
	}
	t.Log(doc)
}

func TestAll(t *testing.T) {
	data, err := tdb("world").All()
	if err != nil {
		t.Error(err)
	}
	t.Log(data)
}

func TestTruncate(t *testing.T) {
	err := tdb("world").Truncate("country")
	if err != nil {
		t.Error(err)
	}
	// need to append
}

func BenchmarkWrite(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tdb("world").Write(tdoc)
	}
}

func BenchmarkRead(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tdb("world").Read(rand.Intn(tdb("world").IndexNum))
	}
}
