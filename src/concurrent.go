package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type Actor struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Movies             []string `json:"movies"`
	AverageMoviesRating float64
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

	idChan := make(chan string)
	actorChan := make(chan Actor)
	resultChan := make(chan Actor)

	go getIds(idChan)
	go getActors(idChan, actorChan)
	go calculateAverageRatings(actorChan, resultChan)

	actors := collectResults(resultChan)
	sortByAverageMoviesRating(actors)

	for i, actor := range actors {
		if i >= 10 {
			break
		}
		fmt.Printf("%s: %.2f\n", actor.Name, actor.AverageMoviesRating)
	}

	duration := time.Since(start)
	fmt.Println("Duração da execução:", duration)
}

func getActors(idChan <-chan string, actorChan chan<- Actor) {
	for id := range idChan {
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

		actorChan <- actor
	}
	close(actorChan)
}

func calculateAverageRatings(actorChan <-chan Actor, resultChan chan<- Actor) {
	for actor := range actorChan {
		actor.AverageMoviesRating = calculateAverageRating(actor.Movies)
		resultChan <- actor
	}
	close(resultChan)
}

func collectResults(resultChan <-chan Actor) []Actor {
	var actors []Actor
	for actor := range resultChan {
		actors = append(actors, actor)
	}
	return actors
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

func getIds(idChan chan<- string) {
	file, err := os.Open("./actors.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		id := strings.Trim(line, `"`)
		idChan <- id
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	close(idChan)
}

func sortByAverageMoviesRating(actors []Actor) {
	sort.Slice(actors, func(i, j int) bool {
		return actors[i].AverageMoviesRating > actors[j].AverageMoviesRating
	})
}