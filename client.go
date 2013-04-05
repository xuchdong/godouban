package douban

import (
    "net/http"
    "net/url"
    "errors"
    "io/ioutil"
    "fmt"
)

type Client struct {
    ClientID        string
    ClientSecret    string
    CallBack        string
    AccessToken     string
    User            string
    UserId          string
    RefreshToken    string
}


func (c *Client) GetAuthUrl() (string, error) {
    auth_url := "https://www.douban.com/service/auth2/auth"
    r, err := url.Parse(auth_url)
    if err != nil {
        return "", errors.New("parse url errror")
    }

    p := url.Values{}
    p.Set("client_id", c.ClientID)
    p.Set("redirect_uri", c.CallBack)
    p.Set("response_type", "code")
    // Don't set scope, state
    r.RawQuery = p.Encode()

    return r.String(), nil
}

func (c *Client) GetAccessToken(code string) error {
    r := "https://www.douban.com/service/auth2/token"
    v := url.Values{}
    v.Set("client_id", c.ClientID)
    v.Set("client_secret", c.ClientSecret)
    v.Set("redirect_uri", c.CallBack)
    v.Set("grant_type", "authorization_code")
    v.Set("code", code)
    resp, err := http.PostForm(r, v)
    if err != nil {
        return errors.New("post error")
    }
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return errors.New("read error")
    }
    fmt.Printf("%s\n", body)
    return nil
}

func (c *Client) GetUserInfo() error {
    rs := "https://api.douban.com/v2/user/:" + c.User
    r, err := url.Parse(rs)
    if err != nil {
        return errors.New("parse error")
    }

    v := url.Values{}
    v.Set("access_token", c.AccessToken)
    v.Set("apikey", c.ClientID)
    r.RawQuery = v.Encode()

    fmt.Printf(r.String())
    resp, err := http.Get(r.String())
    if err != nil {
        return errors.New("get error")
    }
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return errors.New("read error")
    }
    fmt.Printf("%s\n", body)
    return nil
}


