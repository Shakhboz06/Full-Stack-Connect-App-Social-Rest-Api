package db

import (
	"context"
	"database/sql"
	"fmt"
	"go-project/internal/store"
	"log"
	"math/rand"

	"github.com/brianvoe/gofakeit/v7"
)

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(1000)
	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("Error to seed the users")
			return
		}
	}

	posts := generatePosts(10000, users)

	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error to seed the posts")
			return
		}
	}

	comments := generateComments(10000, posts, users)

	for _, comment := range comments{
		if err := store.Comments.Create(ctx,  comment); err != nil {
			log.Println("Error to seed the comments")
			return
		}

	}
	
	followers := generateFollowers(1000, users)

	for _, follower := range followers{
		if err := store.Followers.Follow(ctx, follower.FollowerID, follower.UserID); err != nil {
			log.Println("Error to seed the followers")
			return
		}

	}
	
}


func generateUsers(num int) []*store.Users {

	gofakeit.Seed(8675309)

	users := make([]*store.Users, num)
	password := store.Password{}
	password.Set(gofakeit.Password(true, true, true, true, false, 20))
	
	for i := 0; i < num; i++ {
		users[i] = &store.Users{
			Username: gofakeit.Username(),
			Email:    gofakeit.Email(),
			Password: password,
			Role: store.Roles{
				Name: "user",
			},
		}

	}

	return users
}

func generatePosts(num int, users []*store.Users) []*store.Posts {

	gofakeit.Seed(8675309)
	posts := make([]*store.Posts, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		tags := make([]string, 3)

		for j := 0; j < len(tags); j++ {
			tags[j] = fmt.Sprintf("#%s", gofakeit.BuzzWord())
		}

		posts[i] = &store.Posts{
			UserID:  user.ID,
			Title:   gofakeit.BookTitle(),
			Content: gofakeit.Sentence(3),
			Tags:    tags,
		}
	}

	return posts
}

func generateComments(num int, posts []*store.Posts, users []*store.Users) []*store.Comment {
	gofakeit.Seed(8675309)

	comments := make([]*store.Comment, num)

	for i := 0; i < num; i++ {
		user:= users[rand.Intn(len(users))]
		post:= posts[rand.Intn(len(posts))]

		comments[i] = &store.Comment{
			PostID: int(post.ID),
			UserID: int(user.ID),
			Content: gofakeit.Sentence(3),
		}
	}
	
	return comments
}

func generateFollowers(num int, users []*store.Users)[]*store.Follower{

	gofakeit.Seed(8675309)
	
	follower := make([]*store.Follower, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		follower[i] = &store.Follower{
			FollowerID: int64(gofakeit.IntRange(1, int(user.ID))),
			UserID: user.ID,
		}
	}

	

	
	return follower
}