package service

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/peam1146/mcv-notifier/src/model"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ScraperService interface {
	Scraping() ([]*model.Notification, error)
}

type scraperService struct {
	notificationService NotificationService
	usr                 string
	pwd                 string
}

func NewScraperService(ns NotificationService, usr, pwd string) ScraperService {
	return &scraperService{ns, usr, pwd}
}

const (
	BaseUrl = "https://www.mycourseville.com"
)

func getNotiIDFromUrl(url string) (int, bool) {
	u := strings.Split(url, "/")
	i, err := strconv.Atoi(u[len(u)-1])
	if err != nil {
		log.Println(err.Error())
		return 0, false
	}
	return i, true
}

func request(client *http.Client, url string) io.ReadCloser {
	res, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	return res.Body
}

func login(client *http.Client, loginUrl, token, usr, pwd string) io.ReadCloser {
	res, err := client.PostForm(loginUrl, url.Values{
		"_token":   {token},
		"username": {usr},
		"password": {pwd},
	})
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	return res.Body
}

func scraping(usr, pwd string) (notifications []*model.Notification) {
	jar := NewJarService()
	client := http.Client{
		Jar: jar,
	}

	//	get session cookie
	request(&client, BaseUrl)

	// get XSRF token
	resBody := request(&client, BaseUrl+"/api/oauth/authorize?response_type=code&client_id=mycourseville.com&redirect_uri=https://www.mycourseville.com&login_page=itchula")
	defer resBody.Close()

	doc, err := goquery.NewDocumentFromReader(resBody)
	if err != nil {
		log.Fatal(err)
	}

	_token, _ := doc.Find("form#cv-login-cvecologin-form > input[name=_token]").Attr("value")

	// login
	resBody = login(&client, BaseUrl+"/api/chulalogin", _token, usr, pwd)
	defer resBody.Close()

	// visit notification page
	resBody = request(&client, BaseUrl+"/?q=courseville/course/notification")
	defer resBody.Close()

	doc, err = goquery.NewDocumentFromReader(resBody)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("#courseville-panel-mid-inner > main > div.courseville-feed-item-container > a").Each(func(i int, selection *goquery.Selection) {
		content := selection.Find("div.courseville-feed-item-content")
		href, _ := selection.Attr("href")
		CourseImg, _ := selection.Find("div.courseville-feed-item-icon > img").Attr("src")
		if ID, ok := getNotiIDFromUrl(href); ok {
			notification := model.Notification{
				Model: &gorm.Model{
					ID: uint(ID),
				},
				Link:       BaseUrl + href,
				CourseImg:  fmt.Sprintf("%s/%s", BaseUrl, CourseImg),
				CourseName: content.Find("div.courseville-feed-item-course").Text(),
				Subject:    content.Find("div.courseville-feed-item-title").Text(),
				Date:       content.Find("div.courseville-feed-item-created").Text(),
			}
			notifications = append(notifications, &notification)
		}
	})
	return
}

func (sc *scraperService) Scraping() ([]*model.Notification, error) {
	notifications := scraping(sc.usr, sc.pwd)
	latestNotifications, err := sc.notificationService.SaveNotifications(notifications)
	if err != nil {
		return nil, err
	}
	return latestNotifications, nil
}
