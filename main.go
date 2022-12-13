package main

import (
	"log"
	"os"
)

var ymlData = []byte(`
a:
  - b: a string from struct B
`)

var jsonData = []byte(`{
	"a": [
		{
			"b": "a string from struct B"
		}
	]
}`)

func main() {
	// Здесь у нас контейнер значений.
	// Из него мы можем получать значения по пути.
	// По сути это KV: Path[string]
	var doc Container

	// Путь, по которому мы хотим значение: "a.[0].b"
	path := NewPath(WithPath("a"), WithIndex(0), WithPath("b"))

	// Пример с ямлом
	doc, err := NewYamlContainer(ymlData)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
	}

	val, err := doc.GetValueAtPath(path)
	if err != nil {
		log.Fatalf("value error: %v", err)
	}

	log.Println("value:", val)

	// Пример с джисоном

	doc, err = NewJsonContainer(jsonData)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
	}

	val, err = doc.GetValueAtPath(path)
	if err != nil {
		log.Fatalf("value error: %v", err)
	}

	log.Println("value:", val)

	// Пример с Env'ом
	err = os.Setenv("SUPERAPP_a.[0].b", "fooo test value")
	if err != nil {
		log.Default().Println("Error set env:", err)
	}

	doc, err = NewEnvContainer("SUPERAPP_")
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
	}

	val, err = doc.GetValueAtPath(path)
	if err != nil {
		log.Fatalf("value error: %v", err)
	}

	log.Println("value:", val)

	// Теперь идём по схеме и по пути запрашиваем из контейнера значение
	// и парсим в string -> в целевой тип

}
