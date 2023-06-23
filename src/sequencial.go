package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"bufio"
	"strings"
	"encoding/json"
	"sort"
	"time"
)

type Actor struct {
	ID     				string   `json:"id"`
	Name   				string   `json:"name"`
	Movies 				[]string `json:"movies"`
	AverageMoviesRating  float64
}

type Movie struct {
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	AverageRating   float64  `json:"averageRating"`
	NumberOfVotes   int      `json:"numberOfVotes"`
	StartYear       int      `json:"startYear"`
	LengthInMinutes int      `json:"lengthInMinutes"`
	Genres          []string `json:"genres"`
}

const (
	MOVIES_URL = "http://150.165.15.91:8001/movies/"
	ACTORS_URL = "http://150.165.15.91:8001/actors/"
)


func main() {
	
	start := time.Now()
	ids := getIds()
	actors := getActors(ids)
	for i := range actors {
		actors[i].AverageMoviesRating = calculateAverageRating(actors[i].Movies)
	}
	sortByAverageMoviesRating(actors)
	for i, actor := range actors {
		if i >= 10 {
			break
		}
		fmt.Printf("%d. %s: %.2f\n", i+1, actor.Name, actor.AverageMoviesRating)
	}
	duration := time.Since(start)
	fmt.Println("Duração da execução:", duration)
	
}

func getActors(ids []string) []Actor {
	var actors []Actor
	var cont = 0

	for _, id := range ids {
		url := ACTORS_URL + id
		response, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer response.Body.Close()

		content, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var actor Actor
		err = json.Unmarshal(content, &actor)
		if err != nil {
			fmt.Println(err)
			continue
		}

		actors = append(actors, actor)
		cont++
	}
	return actors
}

func calculateAverageRating(movieIDs []string) float64 {
	if len(movieIDs) == 0 {
		return 0.0
	}

	str := movieIDs[0]
	movieIDs = strings.Split(str, " ")

	total := 0.0
	numMovies := len(movieIDs)

	for _, movieID := range movieIDs {
		movie := getMovie(movieID)
		total += movie.AverageRating
	}

	averageRating := total / float64(numMovies)
	return averageRating
}

func getMovie(id string) Movie {
	var movie Movie

	url := MOVIES_URL + id
	response, _ := http.Get(url)
	defer response.Body.Close()

	content, _ := ioutil.ReadAll(response.Body)

	_ = json.Unmarshal(content, &movie)

	return movie
}

func getIds() []string {
	ids := []string{}
	file, err := os.Open("./actors.txt")
	if err != nil {
    	fmt.Println(err)
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		id := strings.Trim(line, `"`)
		ids = append(ids, id)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	return ids
}

func sortByAverageMoviesRating(actors []Actor) {
	sort.Slice(actors, func(i, j int) bool {
		return actors[i].AverageMoviesRating > actors[j].AverageMoviesRating
	})
}
