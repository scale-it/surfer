package surfer

import (
	"errors"
	"net/http"
	"time"
)

// Sets a cookie. Duration is the amount of time in seconds. 0 = browser defaults
func (this *Handler) SetCookie(name string, value string, expires int64) {
	cookie := &http.Cookie{
		Name:  name,
		Value: value,
		Path:  "/",
	}
	if expires > 0 {
		d := time.Duration(expires) * time.Second
		cookie.Expires = time.Now().Add(d)
	}
	http.SetCookie(this.Response, cookie)
}

func (this *Handler) GetCookie(name string) (string, error) {
	cookie, err := this.Request.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// Sets a secure cookie. Befor using this function you need to initialise Handler.SecCookie value
// using securecookie.New function.
// Duration is the amount of time in seconds. 0 = browser defaults
// Returns nil if succeed, false otherwise
func (this *Handler) SetSecureCookie(name string, value interface{}, expires int64) (err error) {
	if this.SecCookie == nil {
		err = errors.New("Attempt to use secure cookie without initialising Handler.SecCookie variable. This should be initialised in by in user's Prepare method")
		this.App.Log.Fatal(err.Error())
		return
	}
	if encoded, err := this.SecCookie.Encode(name, value); err != nil {
		this.SetCookie(name, encoded, expires)
	}
	return
}

func (this *Handler) GetSecureCookie(name string) (interface{}, error) {
	if this.SecCookie == nil {
		err := errors.New("Attempt to use secure cookie without initialising Handler.SecCookie variable. This should be initialised in by in user's Prepare method")
		this.App.Log.Fatal(err.Error())
		return nil, err
	}
	cookie, err := this.Request.Cookie(name)
	if err == nil {
		value := make(map[string]string)
		err = this.SecCookie.Decode(name, cookie.Value, &value)
		return value, err
	}
	return nil, err
}
