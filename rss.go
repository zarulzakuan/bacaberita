package main

import (
    "fmt"
    "log"
    "sync"
    "bytes"
    "net/http"
    "github.com/PuerkitoBio/goquery"
    "github.com/mmcdole/gofeed"
    "compress/gzip"
    b64 "encoding/base64"
    "io/ioutil"

)

const URL = "https://www.sinarharian.com.my/rssFeed/211"

func main(){
    fp := gofeed.NewParser()
    feed, _ := fp.ParseURL(URL)
    wg := new(sync.WaitGroup)
    
    for _, item := range feed.Items {
        wg.Add(1)
        
        fmt.Println(item.Title)
        fmt.Println(item.Link)
        fmt.Println(item.Published)
        go getSourceCode(item.Link, wg)
    }
    wg.Wait()
}

func getSourceCode(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("HTML code of %s ...\n", url)
	resp, err := http.Get(url)
	// handle the error if there is one
	if err != nil {
		panic(err)
	}
	// do this now so it won't be forgotten
	defer resp.Body.Close()
    
    // Create a goquery document from the HTTP response
    document, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        log.Fatal("Error loading HTTP response body. ", err)
    }

    // Find all links and process them with the function
    // defined earlier
    document.Find("#content-holder").Each(getContent)

    
}

func getContent(index int, s *goquery.Selection) {
    
    // See if the href attribute exists on the element
    content := s.Find("p").Text()
    sEnc := b64.StdEncoding.EncodeToString([]byte(content))
    var b bytes.Buffer
    gz := gzip.NewWriter(&b)


    defer gz.Close()
    if _, err := gz.Write([]byte(sEnc)); err != nil {
        panic(err)
    }
    if err := gz.Flush(); err != nil {
        panic(err)
    }
    if err := gz.Close(); err != nil {
        panic(err)
    }
    fmt.Println(b)
    

    gr, _ := gzip.NewReader(bytes.NewBuffer(b.Bytes()))
	defer gr.Close()
	c, _ := ioutil.ReadAll(gr)
    fmt.Println(c)
    sDec, _ := b64.StdEncoding.DecodeString(string(c))
    fmt.Println(string(sDec))

}
