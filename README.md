# NMNS - NO MEMORY NO SQL

Tiny simple file database. Designed for small applications with small volume of stored data

## Decription
This is dir/file database, tables is the directory, field - files in them.

To create the database is required json file describing the structure:

	{
		"Table": {
			"name": size in bytes,
			"text": 128
		}
	}

Data does not store in memory, receiving data is performed by reading a predetermined number of bytes (field size) at the specific position (id * field size)

> If the length of the values is 8 bytes, 5 id data is read from 40 to 48 bytes


## Functions

- `Init(dir,scheme.json)` - It creates a database, or overwrite an existing one(*can be omitted*)

- `Check(dir,scheme.json)` - check the relevance of the database structure. 

 > In the process of developing your application, you can modify the json scheme file - add tables, fields, or change the fields size - all the changes are made automatically without data loss.

- `Connect(dir)` - makes a connection to database, and returns a connection function (for example db(string)). `db("name").ACTION` - makes ready to use the specified table

    - `db("name").Write(doc)` - writes data in the database and returns the id, doc example: `doc := map[string]string{"city":"Moscow","country":"Russia"}`

    - `db("name").Read(id)` - reades data on the given id, id is `int`, example: `id := 1`

    - `db("name").Search(filter,[limit])` - Searches data by filter, returns a list of id, examples of filters:

        `map[string]interface{}{"name":"Valeriy"}` - full match by a single field

        `map[string]interface{}{"name":"Valeriy","age":"99"}` - full match on the two fields

        `map[string]interface{}{"name":[]string{"Valeriy","Zarina"}}` - at least full match on one value

        `map[string]interface{}{"@name":"Val.*"}` - prefix @ allows you to search using regular expressions

        `map[string]interface{}{"@name":"Val.*","age":"99"}` - for data searching use regular expression  by the field "name" and full match by the field "age"

    `[limit]` - is optional integer parameter


    - `db("name").Update(id,doc)` - updates data on the given id

    - `db("name").Delete(id)` - deletes data on the given id

    - `db("name").All([limit])` - gets all the data

    - `db("name").Truncate([fields])` - clear values, ex:

        `db("name").Truncate("city","country")` - deletes all the data from these fields

        `db("name").Truncate("city")`  - deletes all the data of the only field

        `db("name").Truncate()` - full cleaning of the table(all fields)


## Benchmark

***Performance will be vary depending on the speed of your hard drive***

###Speed
		Write: 300000 docs per 2.788s
		Read docs: 500000 docs per 3.344s