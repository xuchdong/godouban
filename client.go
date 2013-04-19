package douban

import (
    "time"
    "net/url"
    "net/http"
    "io/ioutil"
    "encoding/json"
)

const (
    ApiHost         string = "https://api.douban.com"
    AuthHost        string = "https://www.douban.com"
    TokenUrl        string = AuthHost + "/service/auth2/token"
    AuthorizeUrl    string = AuthHost + "/service/auth2/auth"
)

type Token struct {
    AccessToken     string
    RefreshToken    string
    Expiry          time.Time
}

func (t *Token) Expired() bool {
    if t.Expiry.IsZero() {
        return false
    }
    return t.Expiry.Before(time.Now())
}


type Client struct {
    ClientID        string
    ClientSecret    string
    CallBack        string
}

func (c *Client) GetAuthUrl() (string, error) {
    url_, err := url.Parse(AuthorizeUrl)
    if err != nil {
        return "", err
    }

    q := url.Values{
        "client_id":        {c.ClientID},
        "redirect_uri":     {c.CallBack},
        "response_type":    {"code"},
    }.Encode()
    if url_.RawQuery == "" {
        url_.RawQuery = q
    } else {
        url_.RawQuery += "&" + q
    }

    return url_.String(), nil
}

func (c *Client) GetAccessToken(code string) (*Token, error) {
    v := url.Values{
        "client_id":        {c.ClientID},
        "client_secret":    {c.ClientSecret},
        "redirect_uri":     {c.CallBack},
        "grant_type":       {"authorization_code"},
        "code":             {code},
    }

    resp, err := http.PostForm(TokenUrl, v)
    if err != nil {
        return nil, err
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var m struct {
        Access      string          `json:"access_token"`
        Refresh     string          `json:"refresh_token"`
        ExpiresIn   time.Duration   `json:"expires_in"`
    }

    if err := json.Unmarshal(body, &m); err != nil {
        return nil, err
    }

    m.ExpiresIn *= time.Second

    var tok Token
    tok.AccessToken = m.Access
    if len(m.Refresh) > 0 {
        tok.RefreshToken = m.Refresh
    }
    if m.ExpiresIn == 0 {
        tok.Expiry = time.Time{}
    } else {
        tok.Expiry = time.Now().Add(m.ExpiresIn)
    }

    return &tok, nil
}

