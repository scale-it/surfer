package surfer

import (
	"errors"
	"net/http"
	"time"
)

func set_cookie(h *Handler, name string, value string, expires int64) {
	cookie := &http.Cookie{
		Name:  name,
		Value: value,
		Path:  "/",
	}
	if expires > 0 {
		d := time.Duration(expires) * time.Second
		cookie.Expires = time.Now().Add(d)
	}
	http.SetCookie(h.Response, cookie)
}

// Sets a cookie. Duration is the amount of time in seconds. 0 = browser defaults
func (this *Handler) SetCookie(name string, value interface{}, expires int64) error {
	b, err := serialize(value)
	if err != nil {
		return err
	}
	set_cookie(this, name, string(b), expires)
	return nil
}

func (this *Handler) GetCookie(name string) (v interface{}, err error) {
	cookie, err := this.Request.Cookie(name)
	if err != nil {
		return nil, err
	}
	err = deserialize([]byte(cookie.Value), v)
	return
}

// Delete cookie specified by `name`
func ClearCookie(name string) {
	set_cookie(name, "", -3600*24*365) // Set the cookie expires one year ago
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
	if b, err := this.SecCookie.Encode(name, value); err != nil {
		set_cookie(this, name, b, expires)
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
