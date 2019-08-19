package main

import (
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/icrowley/fake"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type People struct {
	Id int `json:"id"`
	Name      string `json:"name"`
	Title      string `json:"title"`
	Likes     int    `json:"likes"`
	CreatedAt string `json:"created_at"`
	Avatar    string `json:"avatar"`
}

type Comment struct {
	Id int
	Message  string  `json:"message"`
	PersonId int     `json:"person_id"`
	From     *People `json:"from"`
}

type PeopleDetail struct {
	Id int
	Name      string `json:"name"`
	Title      string `json:"title"`
	Likes     int    `json:"likes"`
	CreatedAt string `json:"created_at"`
	Avatar    string `json:"avatar"`
	About    string `json:"about"`
	Address  string `json:"address"`
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	r := gin.New()

	MaxCount := rand.Intn(300) + 100
	MaxComments := rand.Intn(300) + 100

	people := make([]People, MaxCount)
	peopleDetails := make([]PeopleDetail, MaxCount)
	comments := make([]Comment, MaxComments)

	for i := 0; i < MaxCount; i++ {
		people[i] = People{
			Id: i,
			Name:      fake.FullName(),
			Title:      fake.JobTitle(),
			CreatedAt: time.Now().Add(time.Minute * time.Duration(rand.Intn(2000))).Format(time.RFC3339),
			Likes:     rand.Intn(20),
			Avatar:    fmt.Sprintf("/static/%v.jpg", rand.Intn(19)+1),
		}
		peopleDetails[i] = PeopleDetail{
			Id: i,
			About:   fake.Paragraph(),
			Address: fmt.Sprintf("%s, %s, %s", fake.Zip(), fake.StreetAddress(), fake.Country()),
			Name:      fake.FullName(),
			Title:      fake.JobTitle(),
			CreatedAt: time.Now().Add(time.Minute * time.Duration(rand.Intn(2000))).Format(time.RFC3339),
			Likes:     rand.Intn(20),
			Avatar:    fmt.Sprintf("/static/%v.jpg", rand.Intn(19)+1),
		}
	}
	for i := 0; i < MaxComments; i++ {
		comments[i] = Comment{
			Id: i,
			PersonId: rand.Intn(MaxCount),
			From:     &people[rand.Intn(MaxCount)],
			Message:  fake.Sentence(),
		}
	}
	r.Use(static.Serve("/static", static.LocalFile("./static", false)))

	r.GET("/api/people", func(c *gin.Context) {
		count := len(people)
		take, err := strconv.Atoi(c.Query("take"))
		if err != nil {
			take = 10
		}
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 0
		}

		from := page * take
		to := (page + 1) * take

		if from > count {
			c.JSON(http.StatusOK, gin.H{
				"count": count,
				"page":  page,
				"from":  from,
				"to":    to,
				"take":  take,
				"length": 0,
				"data":  []int{},
			})
		} else {

			if to > count {
				to = count
			}

			tmp := people[from:to]
			c.JSON(http.StatusOK, gin.H{
				"count": count,
				"page":  page,
				"from":  from,
				"to":    to,
				"take":  take,
				"length": len(tmp),
				"data":  tmp,
			})
		}

	})

	r.GET("/api/people/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		c.JSON(http.StatusOK, peopleDetails[id])
	})

	r.GET("/api/people/:id/comments", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		filtered := make([]Comment, 0)
		for i := 0; i < MaxComments; i++ {
			if comments[i].PersonId == id {
				filtered = append(filtered, comments[i])
			}
		}
		c.JSON(http.StatusOK, filtered)

	})

	log.Println("Homing frontend assignment api v1.0")
	log.Println("Listening on port " + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}

}
