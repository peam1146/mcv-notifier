package service

import (
	"net/http"
	"net/url"
	"sync"
)

type JarService interface {
	SetCookies(u *url.URL, cookies []*http.Cookie)
	Cookies(u *url.URL) []*http.Cookie
}

type jarService struct {
	lock    sync.Mutex
	cookies map[string][]*http.Cookie
}

func NewJarService() JarService {
	jar := new(jarService)
	jar.cookies = make(map[string][]*http.Cookie)
	return jar
}

func (jar *jarService) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.lock.Lock()
	jar.cookies[u.Host] = cookies
	jar.lock.Unlock()
}

func (jar *jarService) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies[u.Host]
}
