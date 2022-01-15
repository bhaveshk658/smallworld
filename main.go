package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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

func compare(m1 Monster, m2 Monster) bool {
	if m1.Race == m2.Race && m1.Attribute != m2.Attribute && m1.Level != m2.Level && m1.Atk != m2.Atk && m1.Def != m2.Def {
		return true
	}

	if m1.Race != m2.Race && m1.Attribute == m2.Attribute && m1.Level != m2.Level && m1.Atk != m2.Atk && m1.Def != m2.Def {
		return true
	}

	if m1.Race != m2.Race && m1.Attribute != m2.Attribute && m1.Level == m2.Level && m1.Atk != m2.Atk && m1.Def != m2.Def {
		return true
	}

	if m1.Race != m2.Race && m1.Attribute != m2.Attribute && m1.Level != m2.Level && m1.Atk == m2.Atk && m1.Def != m2.Def {
		return true
	}

	if m1.Race != m2.Race && m1.Attribute != m2.Attribute && m1.Level != m2.Level && m1.Atk != m2.Atk && m1.Def == m2.Def {
		return true
	}

	return false

}

func getConnection(m1 Monster, m2 Monster) string {
	if m1.Race == m2.Race {
		return "Type: " + m1.Race
	}
	if m1.Attribute == m2.Attribute {
		return "Attribute: " + m1.Attribute
	}
	if m1.Level == m2.Level {
		return "Level: " + strconv.Itoa(m1.Level)
	}
	if m1.Atk == m2.Atk {
		return "Attack: " + strconv.Itoa(m1.Atk)
	}
	if m1.Def == m2.Def {
		return "Defense: " + strconv.Itoa(m1.Def)
	}
	return "ERROR"
}

func generateCombinations(monsters map[Monster]bool) [][]Monster {
	var combinations [][]Monster
	for m1 := range monsters {
		for m2 := range monsters {
			for m3 := range monsters {
				if compare(m1, m2) && compare(m2, m3) && m1.Name != m3.Name {
					combination := []Monster{m1, m2, m3}
					combinations = append(combinations, combination)
				}
			}
		}
	}

	return combinations
}

func generateCombinationsFromStarter(monsters map[Monster]bool, starter string) [][]Monster {
	var combinations [][]Monster
	for m1 := range monsters {
		if m1.Name == starter {
			for m2 := range monsters {
				for m3 := range monsters {
					if compare(m1, m2) && compare(m2, m3) && m1.Name != m3.Name {
						combination := []Monster{m1, m2, m3}
						combinations = append(combinations, combination)
					}
				}
			}
			break
		}
	}
	return combinations

}

func printCombinations(combinations [][]Monster) {
	for _, combination := range combinations {
		m1 := combination[0]
		m2 := combination[1]
		m3 := combination[2]
		fmt.Printf("%s --(%s)--> %s --(%s)--> %s\n", m1.Name, getConnection(m1, m2), m2.Name,
			getConnection(m2, m3), m3.Name)
	}
}

func main() {

	monsters := getMonsters("/Users/bkalisetti658/desktop/dragonlink.ydk")

	combinations := generateCombinationsFromStarter(monsters, "Effect Veiler")
	printCombinations(combinations)

}
