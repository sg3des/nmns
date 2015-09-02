# NMNS - NO MEMORY NO SQL

Tiny simple file database. Designed for small applications with small volume of stored data

## Decription
This is dir/file database, tables is the directory, field - files in them.

To create the database required json file describing the structure:

	{
		"Table": {
			"name": size in bytes,
			"text": 128
		}
	}

Data not stored in memory, receiving data is performed by reading a predetermined number of bytes (size of a field) at a specific position (id * field size)

> If the length of the values of 8 bytes, 5 id data - read from 40 to 48 bytes


## Functions

> Init - It creates a database, or overwrite an existing one(*can be omitted*)

> Check - check the relevance of the database structure. 

	In the process of developing your application, you can modify the configuration json file - add tables, fields, or change the fields size - all the changes are made automatically without data loss.

> Connect - connected to a database, and returns a connection object (for example "c")

> c.Insert - inserts data in the database and returns the id

> c.Read - reades data on the given id

> c.Search - Search data by filter, returns a list of id, examples of filters:

	map[string]interface{}{"name":"Valeriy"} - full match by a single field

	map[string]interface{}{"name":"Valeriy","age":"99"} - full match on the two fields

	map[string]interface{}{"name":[]string{"Valeriy","Zarina"}} - at least full match one value

	map[string]interface{}{"@name":"Val.*"} - prefix @ allows you to search using regular expressions

	map[string]interface{}{"@name":"Val.*","age":"99"} - regular expression search by the field "name" and full match by field "age"

> c.Update - updates data on the given id

> c.Delete - deletes data on the given id

> c.Truncate - it clears all values:
	
	[]string{"name","age"} - delete all the data from these fields
	
	"name" - deletes all the data of the only field

	"" - complete cleaning table(all fields)

## Benchmark

* Performance will vary depending on the speed of your hard drive *

###Speed
	1 second:
		Inserting new docs: 33400
		Read random docs: 123200

###Memory
	Read and write 100,000 times:
		Working memory: 3mb
		Total allocated memory: 70mb
