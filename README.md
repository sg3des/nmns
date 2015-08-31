# NMNS - NO MEMORY NO SQL

Tiny simple file database. Designed for small applications with small volume of stored data

## Decription
This is dir/file database, tables is the directory, field - files in them.

Для создания базы данных потребуется json файл описывающий структуру: 

	{
		"Table": {
			"name": size in bytes,
			"text": 128
		}
	}

Данные не хранятся в оперативной памяти, получение данных происходит при считывании необходимого количество байт из файла. Если длина значения 8 байт, то для получения 5-го ид - считываются данные из файла с 40 по 47 байт включительно.


## Functions

> Init - создаст базу данных, или перезапишет существующую.

> Check - проверяет актуальность структуры БД. В процессе разработки своего приложения, вы можете изменять json файл конфигурации - добавлять таблицы, поля или изменять размер полей, все изменения будут вноситься автоматически, без потери данных. 

> Connect - подключается к базе данных, и возвращает объект подключения (наприер "c")

> c.Insert - вставляет данные в базу данных и возвращает id

> c.Read - читает данные по конкретному id

> c.Search - поиск данных по условию, возвращает список найденных id, примеры фильтров:
	
	map[string]interface{}{"name":"Valeriy"} - полное совпадение по одному полю

	map[string]interface{}{"name":"Valeriy","age":"99"} - полное совпадение по обоим полям

	map[string]interface{}{"name":["Valeriy","Zarina"]} - полное совпадение хотябы по одному значению

	map[string]interface{}{"@name":"Val.*"} - префикс @ позволяет искать с помощью регулярных выражений 

	map[string]interface{}{"@name":"Val.*","age":"99"} - поиск по регулярному выражению по полям name и полное совпадение по полю age

> c.Update - обвновляет данные по заданному id

> c.Delete - удаляет данные по заданному id

## Benchmark

*Производительность  зависит от скорости жесткого диска*

###Speed
	1 second:
		Inserting new docs: 33400
		Read random docs: 123200

###Memory
	Read and write 100,000 times:
		Working memory: 3mb
		Total allocated memory: 70mb
