package main

import (
	"kuto/config"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	m "kuto/models"

	"kuto/pkg"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	db := pkg.NewDatabase(config.DBHost, config.DBTable, config.DBUser, config.DBPwd)
	var sources []m.ArticleSource
	db.Select(&sources, m.ColumnArticleSourceHost+"=?", "foxnews.com")
	if len(sources) == 0 {
		panic("no foxnews source in database")
	}
	sid := sources[0].ID

	var categories []m.ArticleCategory
	db.Select(&categories, "")

	if len(os.Args) == 2 {
		u, err := url.Parse(os.Args[1])
		if err != nil {
			panic(err)
		}

		log.Println("url=", u)
		if strings.Contains(u.Host, "foxnews.com") {
			article := &m.Article{}
			article.ArticleSourceID = sid

			//获取category
			f := strings.Index(string(u.Path[1:]), "/")
			if f < 0 {
				panic("no flag found")
			}

			flag := string(u.Path[1 : f+1])
			for _, category := range categories {
				if strings.Contains(category.Keywords, flag) {
					article.ArticleCategoryID = category.ID
					break
				}
			}

			if article.ArticleCategoryID == 0 {
				log.Fatal("no category found, flag=" + flag)
			} else {
				request, err := http.NewRequest("GET", u.String(), nil)
				if err != nil {
					panic(err)
				}
				request.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 5.0; SM-G900P Build/LRX21T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Mobile Safari/537.36")
				response, err := http.DefaultClient.Do(request)
				if err != nil {
					panic(err)
				}

				doc, err := goquery.NewDocumentFromReader(response.Body)
				sel := doc.Find(".headline")
				if sel.Length() != 1 {
					panic("head title mismatch")
				}

				article.Title = sel.Text()
				article.Author = doc.Find(".article-source").Parent().Text()
				article.SourceURL = u.String()

				body := doc.Find(".article-body")
				imgs := body.Find("p")

				h, _ := imgs.Html()
				log.Println("imgs", h)
				log.Println(article)
			}
		} else {
			log.Fatal("not foxnews")
		}
	} else {
		log.Fatal("no url specific")
	}
	// res, err := http.Get("http://metalsucks.net")
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
