package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Gender string

const (
	Male   Gender = "MALE"
	Female Gender = "FEMALE"
)

var GetGender = map[string]Gender{
	"MALE":   Male,
	"FEMALE": Female,
}

type User struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Dob         time.Time `json:"date_of_birth"`
	Age         int
	Gender      Gender     `json:"gender"`
	Password    string     `json:"password"`
	CreatedTime time.Time  `json:"created_time"`
	UpdatedTime time.Time  `json:"updated_time"`
	Posts       []Post     `json:"posts"`
	Follows     []Follow   `json:"follows"`
	Likes       []Like     `json:"likes"`
	Reposted    []Reposted `json:"repost"`
}

type Follow struct {
	UserID string    `json:"userID"`
	Since  time.Time `json:"since"`
}

type Like struct {
	PostID string    `json:"postID"`
	At     time.Time `json:"at"`
}

type Reposted struct {
	PostID string    `json:"postID"`
	At     time.Time `json:"at"`
}

type Post struct {
	ID string    `json:"id"`
	At time.Time `json:"at"`
}

// bikin random 40 posts
func FakePostData() []Post {
	faker := gofakeit.New(0)

	var posts = make([]Post, 40)

	for i := 0; i < 40; i++ {
		posts[i] = Post{
			ID: faker.UUID(),
			At: faker.DateRange(time.Date(2023, 1, 1, 0, 0, 0, 0, &time.Location{}), time.Now()),
		}
	}
	return posts

}

func FakeUserData() []User {
	var users = make([]User, 500)
	faker := gofakeit.New(0)
	for i := 0; i < 500; i++ {
		// generate user data
		users[i] = User{
			ID:          faker.UUID(),
			Username:    faker.Name(),
			Email:       faker.Email(),
			Dob:         faker.DateRange(time.Date(2003, 1, 1, 0, 0, 0, 0, &time.Location{}), time.Now()),
			Password:    faker.Password(true, false, true, true, false, 8),
			CreatedTime: faker.DateRange(time.Date(2015, 1, 1, 0, 0, 0, 0, &time.Location{}), time.Now()),
			Age:         faker.Number(10, 50),
			Posts:       FakePostData(),
		}
	}

	// yang difollow user siapa aja
	for i, _ := range users {
		users[i].Follows = createRandomFollowe(users)
	}

	// user likes post apa aja
	for i, _ := range users {
		rand.Seed(time.Now().UnixNano())
		randomIndex := rand.Intn(len(users))
		pick := users[randomIndex]
		otherUserPosts := pick.Posts
		randomStart := faker.Number(0, 28)
		randomEnd := faker.Number(randomStart+1, 39)
		var likes []Like
		for _, otherPost := range otherUserPosts[randomStart:randomEnd] {
			like := Like{
				PostID: otherPost.ID,
				At:     faker.DateRange(time.Date(2023, 1, 1, 0, 0, 0, 0, &time.Location{}), time.Now()),
			}
			likes = append(likes, like)
		}

		users[i].Likes = likes
	}

	// user repost apa aja
	for i, _ := range users {
		rand.Seed(time.Now().UnixNano())
		randomIndex := rand.Intn(len(users))
		pick := users[randomIndex]
		otherUserPosts := pick.Posts
		randomStart := faker.Number(0, 28)
		randomEnd := faker.Number(randomStart+1, 39)
		var reposts []Reposted
		for _, otherPost := range otherUserPosts[randomStart:randomEnd] {
			repost := Reposted{
				PostID: otherPost.ID,
				At:     faker.DateRange(time.Date(2023, 1, 1, 0, 0, 0, 0, &time.Location{}), time.Now()),
			}
			reposts = append(reposts, repost)
		}
		users[i].Reposted = reposts
	}

	return users
}

// bikin 50 random followee
func createRandomFollowe(users []User) []Follow {
	faker := gofakeit.New(0)

	var friends = make([]Follow, 50)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 50; i++ {
		randomIndex := rand.Intn(len(users))
		pick := users[randomIndex]
		friends[i] = Follow{UserID: pick.ID, Since: faker.DateRange(time.Date(2015, 1, 1, 0, 0, 0, 0, &time.Location{}), time.Now())}
	}
	return friends

}


func main() {
	ctx := context.Background()
	users := FakeUserData()

	// // Connection to database
	dbUri := "bolt://localhost:7687"
	dbUser := "neo4j"
	dbPassword := "lintang-neo4j"
	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		panic(err)
	}
	defer driver.Close(ctx)
	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}

	var peoples []map[string]any

	for _, u := range users {
		var follows []map[string]any = make([]map[string]any, 50)
		// [{"id": "123", "since": 2023}]
		for i, follow := range u.Follows {
			follows[i] = make(map[string]any)
			follows[i]["user_id"] = follow.UserID
			follows[i]["since"] = follow.Since
		}

		var likes []map[string]any = make([]map[string]any, 40)
		for i, like := range u.Likes {
			likes[i] = make(map[string]any)
			likes[i]["post_id"] = like.PostID
			likes[i]["at"] = like.At
		}

		var reposts []map[string]any = make([]map[string]any, 40)
		for i, repost := range u.Reposted {
			reposts[i] = make(map[string]any)
			reposts[i]["post_id"] = repost.PostID
			reposts[i]["at"] = repost.At
		}

		var posts []map[string]any = make([]map[string]any, 40)
		for i, post := range u.Posts {
			posts[i] = make(map[string]any)
			posts[i]["id"] = post.ID
			posts[i]["at"] = post.At
		}

		peoples = append(peoples, map[string]any{
			"id":      u.ID,
			"follows": follows,
			"likes":   likes,
			"reposts": reposts,
			"posts":   posts,
		})
	}

	fmt.Println("harap sabar, lagi insert data ke neo4j")


	// membuat node person
	for _, person := range peoples {
		_, err := neo4j.ExecuteQuery(ctx, driver,
			"MERGE (p: Person {id: $person.id})",
			map[string]any{
				"person": person,
			}, neo4j.EagerResultTransformer,
			neo4j.ExecuteQueryWithDatabase("neo4j"))
		if err != nil {
			panic(err)
		}
	}

	// membuat node post
	for _, person := range peoples {
		_, err := neo4j.ExecuteQuery(ctx, driver,
			`
			MATCH (p: Person {id: $person.id})
			UNWIND $person.posts as post
			MERGE (newPost: Post {id: post.id})`, //, at: $post.at
			map[string]any{
				"person": person,
			}, neo4j.EagerResultTransformer,
			neo4j.ExecuteQueryWithDatabase("neo4j"))
		if err != nil {
			panic(err)
		}
	}

	// insert relasi follows
	for _, person := range peoples {
		if person["follows"] != "" {
			_, err = neo4j.ExecuteQuery(ctx, driver, `
				MATCH (p: Person {id: $person.id})
				UNWIND $person.follows as follow
				MATCH (otherUser: Person {id: follow.user_id})
				MERGE (p)-[r:FOLLOWS {since: follow.since}]->(otherUser)
				`, map[string]any{
				"person": person,
			}, neo4j.EagerResultTransformer,
				neo4j.ExecuteQueryWithDatabase("neo4j"))
		}
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}

	// membuat relasi user posted
	for _, person := range peoples {
		if person["posts"] != "" {
			_, err = neo4j.ExecuteQuery(ctx, driver, `
				MATCH (p: Person {id: $person.id})
				UNWIND $person.posts as post
				MATCH (userPost: Post {id: post.id})
				MERGE (p)-[r:POSTED {at: post.at}]->(userPost)
				`, map[string]any{
				"person": person,
			}, neo4j.EagerResultTransformer,
				neo4j.ExecuteQueryWithDatabase("neo4j"))
		}
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}

	// membuat relasi user like post
	for _, person := range peoples {
		if person["likes"] != "" {
			_, err = neo4j.ExecuteQuery(ctx, driver, `
				MATCH (p: Person {id: $person.id})
				UNWIND $person.likes as like
				MATCH (likedPost: Post {id: like.post_id})
				MERGE (p)-[r:LIKES {at: like.at}]->(likedPost)
				`, map[string]any{
				"person": person,
			}, neo4j.EagerResultTransformer,
				neo4j.ExecuteQueryWithDatabase("neo4j"))
		}
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}

	// membuat relasi user - repost
	for _, person := range peoples {
		if person["reposts"] != "" {
			_, err = neo4j.ExecuteQuery(ctx, driver, `
				MATCH (p: Person {id: $person.id})
				UNWIND $person.reposts as repost
				MATCH (repostedPost: Post {id: repost.post_id})
				MERGE (p)-[r:REPOSTS {at: repost.at}]->(repostedPost)
				`, map[string]any{
				"person": person,
			}, neo4j.EagerResultTransformer,
				neo4j.ExecuteQueryWithDatabase("neo4j"))
		}
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}


	for _, user := range users {
		fmt.Println("userID: " + user.ID)
	}

}
