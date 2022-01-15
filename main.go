package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Monster struct {
	Id        int
	Name      string
	Type      string
	Desc      string
	Atk       int
	Def       int
	Level     int
	Race      string
	Attribute string
	Archetype string
}

type Response struct {
	Data []Monster
}

func getMonsters(filename string) map[Monster]bool {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	monsters := make(map[Monster]bool)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "#extra" || text == "#side" {
			break
		}

		if len(text) == 8 {
			resp, err := http.Get("https://db.ygoprodeck.com/api/v7/cardinfo.php?id=" + text)
			if err != nil {
				log.Fatal(err)
			}

			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			var r Response
			err = json.Unmarshal(body, &r)
			if err != nil {
				log.Fatal(err)
			}

			if r.Data[0].Attribute != "" {
				monsters[r.Data[0]] = true
			}
		}
	}

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return monsters
}

func main() {

	monsters := getMonsters("/Users/bkalisetti658/desktop/dragonlink.ydk")

	for monster := range monsters {
		fmt.Println(monster.Name)
	}

}
