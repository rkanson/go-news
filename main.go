package main

import (
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/darkoatanasovski/htmltags"
	"github.com/mmcdole/gofeed"
	"google.golang.org/api/option"

	"golang.org/x/net/context"
)

// Article Model for News Article Storage
type Article struct {
	publisher, title, description, link, date string
}

// NewsPageData Model
type NewsPageData struct {
	news []map[string]interface{}
}

func main() {

	ctx, client := setupFirebase()

	publications := map[string]string{
		"BBC":    "http://feeds.bbci.co.uk/news/world/rss.xml",
		"CNN":    "http://rss.cnn.com/rss/cnn_us.rss",
		"NYT":    "http://www.nytimes.com/services/xml/rss/nyt/HomePage.xml",
		"Huffpo": "https://www.huffpost.com/section/front-page/feed?x=1",
		"Fox":    "http://feeds.foxnews.com/foxnews/latest",
		"USA":    "http://rssfeeds.usatoday.com/UsatodaycomNation-TopStories",
		// "Reuters":  "http://feeds.reuters.com/Reuters/domesticNews",
		"Politico": "http://www.politico.com/rss/politicopicks.xml",
		"Yahoo":    "https://www.yahoo.com/news/rss",
	}

	for k, v := range publications {
		fmt.Println("Fetching:", k)
		news := getNews(k, v)
		saveNews(ctx, news, client)
	}
}

func setupFirebase() (context.Context, *firestore.Client) {
	ctx := context.Background()
	sa := option.WithCredentialsFile("./serviceAccountKey.json")
	conf := &firebase.Config{ProjectID: "go-news-rk"}

	app, err := firebase.NewApp(ctx, conf, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	return ctx, client
}

func saveNews(ctx context.Context, news []Article, client *firestore.Client) {
	for _, v := range news {
		_, _, err := client.Collection(v.publisher).Add(ctx, map[string]interface{}{
			"title":       v.title,
			"description": v.description,
			"link":        v.link,
			"date":        v.date,
		})

		if err != nil {
			log.Fatalf("Failed adding: %v", err)
		}
	}
}

func getNews(publisher string, url string) []Article {
	var news []Article
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)

	if err != nil {
		log.Fatalln("Error:", err)
	}

	for i := 0; i < len(feed.Items); i++ {
		stripTitle, _ := htmltags.Strip(feed.Items[i].Title, []string{}, true)
		stripDesc, _ := htmltags.Strip(feed.Items[i].Description, []string{}, true)
		stripLink, _ := htmltags.Strip(feed.Items[i].Link, []string{}, true)
		stripDate, _ := htmltags.Strip(feed.Items[i].Published, []string{}, true)
		article := Article{
			publisher:   publisher,
			title:       stripTitle.ToString(),
			description: stripDesc.ToString(),
			link:        stripLink.ToString(),
			date:        stripDate.ToString(),
		}
		news = append(news, article)
	}

	return news
}
