package main

import (
	"fmt"
	"github.com/sg3des/nmns"
	"runtime"
)

func main() {

	fmt.Println("start")
	dir := "data"
	if err := nmns.Init("db.json", dir); err != nil {
		panic(err)
	}

	if err := nmns.Check("db.json", dir); err != nil {
		panic(err)
	}

	s, err := nmns.Connect(dir)
	if err != nil {
		panic(err)
	}

	fmt.Println("")
	doc1 := map[string]string{"url": "yandex.ru", "cms": "wordpress"}
	id1, err := s.Insert("Urls", doc1)
	fmt.Println("insert id1", id1, doc1, err)

	doc2 := map[string]string{"url": "google.com"}
	id2, err := s.Insert("Urls", doc2)
	fmt.Println("insert id2", id2, doc2, err)
	fmt.Println("")

	docr1, err := s.Read("Urls", id1)
	fmt.Println("read id1", id1, docr1, err)
	docr2, err := s.Read("Urls", id2)
	fmt.Println("read id2", id2, docr2, err)
	fmt.Println("")

	ids, err := s.Search("Urls", map[string]interface{}{"@url": "google", "@version": ".*"})
	fmt.Println("search ids", ids, err)
	fmt.Println("")

	fmt.Println("delete id 1", s.Delete("Urls", 1))
	docd1, err := s.Read("Urls", 1)
	fmt.Println("read deleted id 1", docd1, err)
	fmt.Println("")

	docmissing, err := s.Read("Urls", 2)
	fmt.Println("read missing id 2", docmissing, err)
	fmt.Println("")

	docupdate := map[string]string{"url": "yandex.ru", "cms": "joomla", "version": "1.1"}
	fmt.Println("update id 0", docupdate, s.Update("Urls", 0, docupdate))

	docur, err := s.Read("Urls", 0)
	fmt.Println("read updated id 0", docur, err)

}

func bench(s *nmns.Nmns) {
	var mem runtime.MemStats

	doc := map[string]string{"url": "google.com", "cms": "wp", "version": "123"}
	for i := 0; i < 100000; i++ {
		id, _ := s.Insert("Urls", doc)
		s.Read("Urls", id)
	}

	runtime.ReadMemStats(&mem)
	fmt.Println("Alloc(bytes):     ", mem.Alloc)
	fmt.Println("TotalAlloc(bytes):", mem.TotalAlloc)
}
