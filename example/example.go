package main

import (
	"fmt"
	"github.com/sg3des/nmns"
)

func main() {
	dir := "data"

	if err := nmns.Init("./db.json", dir); err != nil {
		panic(err)
	}

	if err := nmns.Check("./db.json", dir); err != nil {
		panic(err)
	}

	db, err := nmns.Connect(dir)
	if err != nil {
		panic(err)
	}

	doc0 := map[string]string{"city": "Moscow", "country": "Russia"}
	doc1 := map[string]string{"city": "New York", "country": "USA"}
	doc2 := map[string]string{"city": "London", "country": "United Kingdom"}

	println("Write:")
	id0, err := db("world").Write(doc0)
	fmt.Println("	", id0, err)
	id1, err := db("world").Write(doc1)
	fmt.Println("	", id1, err)
	id2, err := db("world").Write(doc2)
	fmt.Println("	", id2, err)

	println("Read:")
	rdoc0, err := db("world").Read(id0)
	fmt.Println("	", rdoc0, err)
	rdoc1, err := db("world").Read(id1)
	fmt.Println("	", rdoc1, err)
	rdoc2, err := db("world").Read(id2)
	fmt.Println("	", rdoc2, err)

	println("Update:")
	udoc1 := map[string]string{"city": "Washington", "country": "USA"}
	db("world").Update(1, udoc1)
	rudoc1, err := db("world").Read(1)
	fmt.Println("	", rudoc1, err)

	println("All:")
	all, err := db("world").All()
	fmt.Println("	", all, err)

	println("Search:")
	filter0 := map[string]interface{}{"city": "London"}
	ids0, err := db("world").Search(filter0)
	fmt.Println("	    filter:", filter0, "id:", ids0, err)

	filter1 := map[string]interface{}{"@city": "^.o"}
	ids1, err := db("world").Search(filter1)
	fmt.Println("	reg filter:", filter1, "id:", ids1, err)

	filter2 := map[string]interface{}{"city": []string{"London", "Moscow"}}
	ids2, err := db("world").Search(filter2)
	fmt.Println("	arr filter:", filter2, "id:", ids2, err)

	println("Truncate field country:")
	db("world").Truncate("country")
	tall, err := db("world").All()
	fmt.Println("	", tall, err)

	println("All with limit 2:")
	all_l2, err := db("world").All(2)
	fmt.Println("	", all_l2, err)

	println("Search with limit 2:")
	all_sl2, err := db("world").Search(map[string]interface{}{"@city": ".*"}, 2)
	fmt.Println("	", all_sl2, err)
}
