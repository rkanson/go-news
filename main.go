package main

import (
	"log"

	firebase "firebase.google.com/go"
	"github.com/darkoatanasovski/htmltags"
	"github.com/mmcdole/gofeed"
	"google.golang.org/api/option"

	"golang.org/x/net/context"
)

// Article Model for News Article Storage
type Article struct {
	title, description, link, date string
}

// Publication Model for Storing arrays of Articles
type Publication struct {
	articles Article
}

func main() {

	projectID := "go-news-rk"
	ctx := context.Background()
	sa := option.WithCredentialsFile("./serviceAccountKey.json")
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	news := getNews("http://feeds.bbci.co.uk/news/world/rss.xml")

	for _, v := range news {
		_, data, err := client.Collection("users").Add(ctx, v)

		if data != nil {
			log.Fatalf("Failed adding: %v", data)
		}

		if err != nil {
			log.Fatalf("Failed adding: %v", err)
		}
	}

	// getNews("http://rss.cnn.com/rss/cnn_us.rss")

}

func getNews(url string) []Article {
	var news []Article
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(url)

	for i := 0; i < len(feed.Items); i++ {
		stripTitle, _ := htmltags.Strip(feed.Items[i].Title, []string{}, true)
		stripDesc, _ := htmltags.Strip(feed.Items[i].Description, []string{}, true)
		stripLink, _ := htmltags.Strip(feed.Items[i].Link, []string{}, true)
		stripDate, _ := htmltags.Strip(feed.Items[i].Published, []string{}, true)
		article := Article{
			title:       stripTitle.ToString(),
			description: stripDesc.ToString(),
			link:        stripLink.ToString(),
			date:        stripDate.ToString(),
		}
		news = append(news, article)
	}

	return news
}
