package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly"
)

type Item struct {
  Link  string `json:"link"`
  Name  string `json:"Name"`
  Price  string `json:"Price"`
  Instock  string `json:"instock"`
}

func timer(name string) func(){
  start := time.Now()
  return func() {
    fmt.Println("%s took %v\n", name, time.Since(start))
  }
}

func main() {
  defer timer("main")()
  c := colly.NewCollector(colly.Async(true))

  items := []Item{}

  c.OnHTML("div.side_categories li ul li", func(h *colly.HTMLElement){
    link := h.ChildAttr("a","href")
    c.Visit(h.Request.AbsoluteURL(link))
  })

  c.OnHTML("", func(h *colly.HTMLElement) {
    c.Visit(h.Request.AbsoluteURL(h.Attr("href")))
  })

  c.OnHTML("article.product_pod", func(h *colly.HTMLElement){
    i:= Item {
      Link:   h.ChildAttr("a","href"),
      Name:   h.ChildAttr("h3 a","title"),
      Price:   h.ChildText("p.price_color"),
      Instock:   h.ChildText("p.instock"),
    }
    items = append(items, i)
  })

  c.OnHTML("li.next", func(h *colly.HTMLElement) {
		c.Visit(h.Request.AbsoluteURL(h.Attr("href")))
	})


  c.OnRequest(func (r *colly.Request){
    fmt.Println("visiting", r.URL)
  })

  //c.Visit("https://books.toscrape.com/catalogue/page-1.html")
  c.Visit("https://books.toscrape.com/catalogue/category/books/travel_2/index.html")
  c.Wait()
  
  
  // for i := 1; i <= 5; i++ {
	// 	pageURL := fmt.Sprintf("https://books.toscrape.com/catalogue/page-%d.html", i)
	// 	c.Visit(pageURL)
	// }
  
  data, err := json.MarshalIndent(items, " ","")
  if err != nil {
    log.Fatal()
  }
  fmt.Println(string(data))

}