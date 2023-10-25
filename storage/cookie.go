package storage

import (
	"net/http"
	"time"
)

type CookieStore struct {
	cookieName string
}

func NewCookieStore(cookieName string) *CookieStore {
	return &CookieStore{cookieName: cookieName}
}

// Get retrieves the session data from the cookie
func (cs *CookieStore) Get(r *http.Request) (string, bool) {
	cookie, err := r.Cookie(cs.cookieName)
	if err != nil {
		return "", false
	}
	return cookie.Value, true
}

// Set sets the session data as a cookie
func (cs *CookieStore) Set(w http.ResponseWriter, value string, ttl time.Duration) {
	expiration := time.Now().Add(ttl)
	cookie := &http.Cookie{
		Name:    cs.cookieName,
		Value:   value,
		Expires: expiration,
		Path:    "/", // this makes the cookie available on all paths
		//HttpOnly: true, // disallow JavaScript access to the cookie
	}
	http.SetCookie(w, cookie)
}

// Delete deletes the session cookie
func (cs *CookieStore) Delete(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     cs.cookieName,
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}
