godouban
========

Douban go's SDK

使用说明:
    
    package main
    
    import (
        "fmt"
        "douban"
    )
    
    func main(){
        client := douban.Client{
        }

        url_, _ := client.GetAuthUrl()
        fmt.Printf("Go to the following link in your browser:\n %s\n", url_)

        code := new([]byte)
        fmt.Scanf("Enter the verification code:%s", code)
        tok, _= client.GetAccessToken(string(*code))
    }

