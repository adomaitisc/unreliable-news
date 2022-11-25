package main

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

type Article struct {
	ID	  int64 `json:"article_id"`
	Title string `json:"article_title"`
	Body  string `json:"article_body"`
	Image string `json:"article_image_url"`
	CreatedAt string `json:"article_created_at"`
}

func GetArticles(c *gin.Context) {
	// Collect new articles
	CollectNewArticle()

	query := "SELECT * FROM articles"
	res, err := db.Query(query)
	if err != nil {
		log.Fatal("(GetArticles) db.Query: ", err)
	}
	defer res.Close()

	articles := []Article{}
	for res.Next() {
		var article Article
		err := res.Scan(&article.ID, &article.Title, &article.Body, &article.Image, &article.CreatedAt)
		if err != nil {
			log.Fatal("(GetArticles) res.Scan: ", err)
		}
		articles = append(articles, article)
	}

	c.JSON(http.StatusOK, articles)
}

func GetSingleArticle(c *gin.Context) {
	articleId := c.Param("articleId")
	articleId = strings.ReplaceAll(articleId, "/", "")
	articleIdInt, err := strconv.Atoi(articleId)
	if err != nil {
		log.Fatal("(GetSingleArticle) strconv.Atoi: ", err)
	}

	var article Article
	query := `SELECT * FROM articles WHERE article_id = ?`
	err = db.QueryRow(query, articleIdInt).Scan(&article.ID, &article.Title, &article.Body, &article.Image, &article.CreatedAt)
	if err != nil {
		log.Fatal("(GetSingleArticle) db.Exec: ", err)
	}

	c.JSON(http.StatusOK, article)
}

func PostArticle(newArticle Article) {
	query := `INSERT INTO articles (article_title, article_body, article_image_url) VALUES (?, ?, ?)`

	res, err := db.Exec(query, newArticle.Title, newArticle.Body, newArticle.Image)
	if err != nil {
		log.Fatal("(CreateArticle) db.Exec: ", err)
	}
	newArticle.ID, err = res.LastInsertId()
	if err != nil {
		log.Fatal("(CreateProduct) res.LastInsertId: ", err)
	}
}

func CollectNewArticle() {
	// Query URL
	url := "https://www.bloomberg.com/search?query=elon%20musk%20twitter"
	
	// Setting up http client
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("(CollectNewArticle) http.NewRequest: ", err)
	}

	// Setting up headers
	req.Header = http.Header{
		"cookie": {`seen_uk=1; exp_pref=AMER; agent_id=b324de52-da76-4bbd-8ac1-5a842e67cb8a; session_id=11e4dee7-9d01-4856-9e57-7fcdd0b6558c; session_key=4e01b7c3294ef8e7111d55cc39290420d16cf032; gatehouse_id=9b62a5ad-7d60-4e7a-8ae5-75f7b55e0313; geo_info={"countryCode":"US","country":"US","cityId":"4930956","provinceId":"6254926","field_p":"176C9","field_d":"starry-inc.net","field_mi":4,"field_n":"hf","trackingRegion":"US","cacheExpiredTime":1669958263048,"region":"US","fieldMI":4,"fieldN":"hf","fieldD":"starry-inc.net","fieldP":"176C9"}|1669958263048; geo_info={"country":"US","region":"US","cityId":"4930956","provinceId":"6254926","fieldP":"176C9","fieldD":"starry-inc.net","fieldMI":4,"fieldN":"hf"}|1669958263111; _reg-csrf=s:ZQ7IS7JmHQRKR227J_8EgBsc.9UJ7flKgxcn4sqxHIvLc4RCg8JGs7wkY1QaHAkP0/6Q; ccpaUUID=1f79ac6e-15ac-43c0-844a-3b7e62470b00; dnsDisplayed=true; ccpaApplies=true; signedLspa=false; _sp_krux=false; _sp_v1_ss=1:H4sIAAAAAAAAAItWqo5RKimOUbLKK83J0YlRSkVil4AlqmtrlXRGldFSWSwAQNXmRIcBAAA=; sampledUser=false; bbgconsentstring=req1fun1pad1; __gads=ID=f477a785aa964cdf:T=1669353463:S=ALNI_MbPekbZz1SIpBISDVXJM6qmFFf0kA; __gpi=UID=000008b2bc172c87:T=1669353463:RT=1669353463:S=ALNI_MaBFoBV9uRhyz9pWUHx0uXjSVL_QA; _sp_v1_uid=1:437:17015dbf-b1e3-421c-a693-09d4136131af; _sp_v1_data=2:517482:1668943320:0:11:0:11:0:0:_:-1; _reg-csrf-token=pC8ulZRA-PfI9rLTYmFWUEkWASlpKBNMOLAw; _user-data={"status":"anonymous","newsletterIds":[]}`},
		"user-agent": {`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36`},
		"sec-ch-ua": {`" Not A;Brand";v="99", "Chromium";v="91", "Google Chrome";v="91"`},
		"sec-ch-ua-mobile": {`?0`},
		"sec-fetch-dest": {`document`},
		"sec-fetch-mode": {`navigate`},
		"sec-fetch-site": {`none`},
		"sec-fetch-user": {`?1`},
		"upgrade-insecure-requests": {`1`},
		"accept": {`text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9`},
		"accept-language": {`en-US,en;q=0.9`},
		"cache-control": {`max-age=0`},
		"accept-encoding": {`gzip, deflate, br`},
		"dnt": {`1`},
		"if-none-match": {`W/"b52f4-WWkWHmU+yXnB3hNwVWCmYx35fuk""`},
	}

	// Making request
	res, err := client.Do(req)	
	if err != nil {
		log.Fatal("(CollectNewArticle) client.Do: ", err)
	}
	defer res.Body.Close()

	// Reading response
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal("(CollectNewArticle) goquery.NewDocumentFromReader: ", err)
	}

	// Collecting first article url
	articleLink, _ := doc.Find(".headline__3a97424275").First().Attr("href")

	// Fetching article
	req, err = http.NewRequest("GET", articleLink, nil)
	if err != nil {
		log.Fatal("(CollectNewArticle) http.NewRequest: ", err)
	}

	// Setting up headers
	req.Header = http.Header{
		"cookie": {`seen_uk=1; exp_pref=AMER; agent_id=b324de52-da76-4bbd-8ac1-5a842e67cb8a; session_id=11e4dee7-9d01-4856-9e57-7fcdd0b6558c; session_key=4e01b7c3294ef8e7111d55cc39290420d16cf032; gatehouse_id=9b62a5ad-7d60-4e7a-8ae5-75f7b55e0313; geo_info={"countryCode":"US","country":"US","cityId":"4930956","provinceId":"6254926","field_p":"176C9","field_d":"starry-inc.net","field_mi":4,"field_n":"hf","trackingRegion":"US","cacheExpiredTime":1669958263048,"region":"US","fieldMI":4,"fieldN":"hf","fieldD":"starry-inc.net","fieldP":"176C9"}|1669958263048; geo_info={"country":"US","region":"US","cityId":"4930956","provinceId":"6254926","fieldP":"176C9","fieldD":"starry-inc.net","fieldMI":4,"fieldN":"hf"}|1669958263111; _reg-csrf=s:ZQ7IS7JmHQRKR227J_8EgBsc.9UJ7flKgxcn4sqxHIvLc4RCg8JGs7wkY1QaHAkP0/6Q; ccpaUUID=1f79ac6e-15ac-43c0-844a-3b7e62470b00; dnsDisplayed=true; ccpaApplies=true; signedLspa=false; _sp_krux=false; _sp_v1_ss=1:H4sIAAAAAAAAAItWqo5RKimOUbLKK83J0YlRSkVil4AlqmtrlXRGldFSWSwAQNXmRIcBAAA=; sampledUser=false; bbgconsentstring=req1fun1pad1; __gads=ID=f477a785aa964cdf:T=1669353463:S=ALNI_MbPekbZz1SIpBISDVXJM6qmFFf0kA; __gpi=UID=000008b2bc172c87:T=1669353463:RT=1669353463:S=ALNI_MaBFoBV9uRhyz9pWUHx0uXjSVL_QA; _sp_v1_uid=1:437:17015dbf-b1e3-421c-a693-09d4136131af; _sp_v1_data=2:517482:1668943320:0:11:0:11:0:0:_:-1; _reg-csrf-token=pC8ulZRA-PfI9rLTYmFWUEkWASlpKBNMOLAw; _user-data={"status":"anonymous","newsletterIds":[]}`},
		"user-agent": {`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36`},
		"sec-ch-ua": {`" Not A;Brand";v="99", "Chromium";v="91", "Google Chrome";v="91"`},
		"sec-ch-ua-mobile": {`?0`},
		"sec-fetch-dest": {`document`},
		"sec-fetch-mode": {`navigate`},
		"sec-fetch-site": {`none`},
		"sec-fetch-user": {`?1`},
		"upgrade-insecure-requests": {`1`},
		"accept": {`text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9`},
		"accept-language": {`en-US,en;q=0.9`},
		"cache-control": {`max-age=0`},
		"accept-encoding": {`gzip, deflate, br`},
		"dnt": {`1`},
		"if-none-match": {`W/"b52f4-WWkWHmU+yXnB3hNwVWCmYx35fuk""`},
	}

	// Making request
	res, err = client.Do(req)
	if err != nil {
		log.Fatal("(CollectNewArticle) client.Do: ", err)
	}

	// Reading response
	doc, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal("(CollectNewArticle) goquery.NewDocumentFromReader: ", err)
	}

	// Collecting article title
	articleTitle := doc.Find("h1").Text()
	articleImage, _ := doc.Find(".lazy-image__09ca8e3c.lazy-img__image").First().Attr("src")

	// Check if article exists in database
	if ArticleExists(articleTitle) {
		log.Print("(CollectNewArticle) Article already exists in database")
		return
	}

	// Collecting article content
	doc.Find("[data-component-props='ArticleBody']").Html()

	// Collect article content
	substring := doc.Text()[strings.Index(doc.Text(), `<div`):]
	substring = substring[:strings.Index(substring, `",`)]

	// create a doc from the string
	contentDoc, err := goquery.NewDocumentFromReader(strings.NewReader(substring))
	if err != nil {
		log.Fatal("(CollectNewArticle) goquery.NewDocumentFromReader: ", err)
	}

	pre := contentDoc.Find("p").Text()

	// prettify the article body
	pre = strings.ReplaceAll(pre, `<p>`, "")
	pre = strings.ReplaceAll(pre, `<\/p>`, "")
	pre = strings.ReplaceAll(pre, `<a>`, "")
	pre = strings.ReplaceAll(pre, `<\/a>`, "")
	pre = strings.ReplaceAll(pre, `<span>`, "")
	pre = strings.ReplaceAll(pre, `<\/span>`, "")
	pre = strings.ReplaceAll(pre, `<strong>`, "")
	pre = strings.ReplaceAll(pre, `<\/strong>`, "")
	pre = strings.ReplaceAll(pre, `<em>`, "")
	pre = strings.ReplaceAll(pre, `<\/em>`, "")

	// post request to paraphrase the article
	articleBody := ParaphraseArticle(pre)

	// Creating new article
	newArticle := Article{
		Title: articleTitle,
		Body: articleBody,
		Image: articleImage,
	}

	// Inserting new article into database
	PostArticle(newArticle)
}

func ArticleExists(title string) bool {
	var article Article
	query := `SELECT * FROM articles WHERE article_title = ?`
	err := db.QueryRow(query, title).Scan(&article.ID, &article.Title, &article.Body, &article.Image, &article.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Fatal("(GetSingleArticle) db.Exec: ", err)
	}

	return true
}

func ParaphraseArticle(text string) string {
	// get api key
	apiKey := os.Getenv("PAK")
	url := "https://api.apilayer.com/paraphraser"

	// find if text is larger than 2000 characters
	for len(text) > 2000  {
		// remove first 5 sentences
		for i := 0; i < 5; i++ {
			text = text[strings.Index(text, ".")+1:]
		}
		// remove last 5 sentences
		for i := 0; i < 5; i++ {
			text = text[:strings.LastIndex(text, ".")]
		}
	}

	payload := strings.NewReader(text)
  
	client := &http.Client {}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		log.Fatal("(ParaphraseArticle) http.NewRequest: ", err)
	}
	req.Header.Set("apikey", apiKey)

	res, err := client.Do(req)
	if err != nil {
	  log.Fatal("(ParaphraseArticle) client.Do: ", err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
	  log.Fatal("(ParaphraseArticle) ioutil.ReadAll: ", err)
	}

	// return the paraphrased value only
	// remove until the first comma
	body = body[strings.Index(string(body), ",")+1:]
	// remove until the semi colon
	body = body[strings.Index(string(body), ":")+1:]
	// remove all the quotes and brackets
	bodyString := strings.ReplaceAll(string(body), (`"`), (""))
	bodyString = strings.ReplaceAll(bodyString, (`{`), (""))
	bodyString = strings.ReplaceAll(bodyString, (`}`), (""))

	return bodyString
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load env: ", err)
	}

	// Open a connection to the database
	db, err = sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		log.Fatal("failed to open db connection: ", err)
	}

	// Build router & define routes
	router := gin.Default()
	router.GET("/articles", GetArticles)
	router.GET("/articles/:articleId", GetSingleArticle)

	// Run the router
	router.Run()
}