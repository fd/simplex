package example

import (
	"fmt"
	"html/template"
	"simplex.sh/container"
	"simplex.sh/paginate"
	"simplex.sh/shttp"
	"simplex.sh/static"
	"time"
)

var _ = container.App(func(app *container.Application) {
	app.Name = "example"

	app.Generator = static.GeneratorFunc(Generate)

	app.ExtraHosts = []string{"example.dev.", "*."}
})

type Article struct {
	Id          int
	Title       string
	Body        string
	PublishedAt time.Time
}

func Generate(tx *static.Tx) {
	var (
		articles *static.C
	)

	articles = tx.Coll("articles", &Article{})
	articles = articles.Select(published).Sort(by_publish_date)

	generate_index_pages(articles)
	generate_article_pages(articles)
}

func published(a *Article) bool {
	if a.Title == "" {
		return false
	}

	if a.PublishedAt.IsZero() {
		return false
	}

	if a.PublishedAt.After(time.Now()) {
		return false
	}

	return true
}

func by_publish_date(a, b interface{}) bool {
	var (
		article_a = a.(*Article)
		article_b = b.(*Article)
	)

	return article_a.PublishedAt.After(article_b.PublishedAt)
}

func generate_index_pages(articles *static.C) {
	index_pages := paginate.Paginate(articles, 50)

	shttp.Render(index_pages, func(page *paginate.Page, w shttp.Writer) error {
		if page.Number == 1 {
			w.Path("/articles")
		}

		w.Path(fmt.Sprintf("/articles/page-%d", page.Number))

		return index_tmpl.Execute(w, page)
	})
}

func generate_article_pages(articles *static.C) {
	shttp.Render(articles, func(article *Article, w shttp.Writer) error {
		w.Path(fmt.Sprintf("/articles/%d", article.Id))

		return article_tmpl.Execute(w, article)
	})
}

var index_tmpl = template.Must(template.New("hello").Parse(`
  <ul>
    {{range $.Elements}}
      <li><a href="/articles/{{.Id}}">{{.Title}}</a></li>
    {{end}}
  </ul>
`))

var article_tmpl = template.Must(template.New("hello").Parse(`
  <h1>{{.Title}}</h1>
  <p>{{.Body}}</p>
`))
