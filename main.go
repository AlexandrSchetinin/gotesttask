package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
)

func Config() (string, string, int) {
	fContent, err := ioutil.ReadFile("config.ini")
	if err != nil {
		panic(err)
	}
	config := string(fContent)
	var i int
	address := ""
	for i = 8; '\r' != config[i]; i++ {
		address += string(config[i])
	}

	port := ""
	for i += 7; '\r' != config[i]; i++ {
		port += string(config[i])
	}

	PeakCache_S := ""
	for i += 12; i < len(config); i++ {
		PeakCache_S += string(config[i])
	}
	PeakCache, err := strconv.Atoi(PeakCache_S)
	if err != nil {
		log.Fatal(err)
	}

	return address, port, PeakCache
}

func MakeRequest(address string) []byte {
	resp, err := http.Get("http://" + address)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body
}

type Cache struct {
	Id      int
	Address string
}

func GetCache() []Cache {
	filename, err := os.Open("cache.json")
	if err != nil {
		log.Fatal(err)
	}

	defer filename.Close()

	data, err := ioutil.ReadAll(filename)

	if err != nil {
		log.Fatal(err)
	}

	var result []Cache

	jsonErr := json.Unmarshal(data, &result)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return result
}

var name string = ".html"

func main() {
	cache := GetCache()
	address, port, PeakCache := Config()
	i := 0
	for i = 0; i < len(cache); i++ {
		if address == cache[i].Address {
			name = strconv.Itoa(cache[i].Id) + name
			break
		}
	}
	if i == len(cache) {

		if len(cache) < PeakCache {
			var Newcache Cache
			Newcache.Id = cache[len(cache)-1].Id + 1
			Newcache.Address = address
			cache = append(cache, Newcache)
			name = strconv.Itoa(cache[len(cache)-1].Id) + name

			json_data, err := json.Marshal(cache)

			if err != nil {
				log.Fatal(err)
			}

			err = ioutil.WriteFile("cache.json", json_data, 0777)
			if err != nil {
				// Если произошла ошибка выводим ее в консоль
				fmt.Println(err)
			}
			f, err := os.Create("cache/address" + name)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			_, err = f.WriteString(string(MakeRequest(address)))
			if err != nil {
				panic(err)
			}
		} else {
			for i := 0; i < len(cache)-1; i++ {
				cache[i].Id = cache[i+1].Id
				cache[i].Address = cache[i+1].Address
			}
			cache[len(cache)-1].Id = 0
			cache[len(cache)-1].Address = ""
			cache = cache[:len(cache)-1]
			var Newcache Cache
			Newcache.Id = cache[len(cache)-1].Id + 1
			Newcache.Address = address
			cache = append(cache, Newcache)
			name = strconv.Itoa(cache[len(cache)-1].Id) + name

			json_data, err := json.Marshal(cache)

			if err != nil {
				log.Fatal(err)
			}

			err = ioutil.WriteFile("cache.json", json_data, 0777)
			if err != nil {
				// Если произошла ошибка выводим ее в консоль
				fmt.Println(err)
			}
			f, err := os.Create("cache/address" + name)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			_, err = f.WriteString(string(MakeRequest(address)))
			if err != nil {
				panic(err)
			}

		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, mux)

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var tpl = template.Must(template.ParseFiles("cache/address" + name))
	tpl.Execute(w, nil)

}
