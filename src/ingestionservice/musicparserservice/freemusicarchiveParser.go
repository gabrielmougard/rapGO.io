package main

import (
	"bytes"
	"errors"
	"io"
	"fmt"
	"math/rand"
	"time"
	"strconv"
	"io/ioutil"
	"net/http"
	"golang.org/x/net/html"
	"github.com/google/uuid"
)

/*
This parser is meant to be used with the website https://freemusicarchive.org

-- Purpose -- :
Concurrently scrap all the songs on a given URL and save
them on disk under .mp3 format or store their URLs in a database
for late download.

-- Usage -- :
Put the array of the base URLs which you want to
download (use Seed() method).
Then use the Start() method to begin scrapping.

*/

type MusicParser struct {
	baseURLS       []string
	rootURL        string
	batchSize      int
	directDownload bool
	limit          int
}

func (mp *MusicParser) Seed(baseURLs []string, directDownload bool, limit int) {
	mp.baseURLS       = baseURLs 
	mp.rootURL        = "https://freemusicarchive.org/genre/"
	mp.directDownload = directDownload
	mp.limit          = limit
}

func (mp *MusicParser) SetLimit(newlimit int) {
	mp.limit = newlimit
}

func (mp *MusicParser) GetLimit() int {
	return mp.limit
}

func (mp *MusicParser) Start() {
	//if set to true, download .mp3 directly. Else, store only URLs in database
	//'limit' parameter is the limit of songs to download.
	
	for pageNumber := 1; mp.GetLimit() > 0 ; pageNumber++{
		fmt.Println("parsing page number #"+strconv.FormatInt(int64(pageNumber),10))
		for j, URL := range mp.baseURLS {
			if URL != "" {
				baseSongsURL := mp.rootURL+URL+"?sort=track_date_published&d=1&page="+strconv.FormatInt(int64(pageNumber),10)
				mp3URLs := renderMP3Node(baseSongsURL)
				
				if len(mp3URLs) > 0 {
					if len(mp3URLs) > mp.limit {
						diff := len(mp3URLs)-mp.limit
						mp3URLs = randomTruncate(mp3URLs,diff) //delete randomly some URLs
					}
			
					//download here
					if mp.directDownload {
						batch,err := downloadMP3Files(mp3URLs)
						if err != nil {
							return
						}
						writeBatch(&batch)
						batch = nil //clear slice
					}
			
					mp.SetLimit(mp.GetLimit()-len(mp3URLs))
					if mp.GetLimit() < 0 {
						return
					}
				} else {
					mp.baseURLS[j] = ""
				}
			}	
		}

		nulCount := 0
		for _, URL := range mp.baseURLS {
			if URL != "" {
				nulCount += 1
			}
		}
		if nulCount == 0 {
			break //no more songs to download in all our baseURLS
		} 
	}
}

func renderMP3Node(pageURL string) (nodes []string) {
	resp, _ := http.Get(pageURL)
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	reader := bytes.NewReader(data)
	
	z := html.NewTokenizer(reader)

	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			//End of the document, we're done
			return
		case tt == html.StartTagToken:
			t := z.Token()
			isLink := t.Data == "a"
			if isLink {
				for _,attr := range t.Attr {
					if attr.Key == "class" && attr.Val =="icn-arrow" {
						for _, attr2 := range t.Attr {
							if attr2.Key == "href" {
								nodes = append(nodes,attr2.Val)
							} 
						}
					}
				}
			}
		}
	}
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

func downloadMP3Files(urls []string) ([][]byte, error) {
	fmt.Println("Downloading files...")
	done := make(chan []byte, len(urls))
	errch := make(chan error, len(urls))
	for _, URL := range urls {
		go func(URL string) {
			b, err := downloadFile(URL)
			if err != nil {
				errch <- err
				done <- nil
				return
			}
			done <- b
			fmt.Println("Download finished for "+URL)
			errch <- nil
		}(URL)
	}
	bytesArray := make([][]byte, 0)
	var errStr string
	for i := 0; i < len(urls); i++ {
		bytesArray = append(bytesArray, <-done)
		if err := <-errch; err != nil {
			errStr = errStr + " " + err.Error()
		}
	}
	var err error
	if errStr!=""{
		err = errors.New(errStr)
	}
	return bytesArray, err
}

func downloadFile(URL string) ([]byte, error) {
	fmt.Println("Start downloading : "+URL)
	response, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}
	var data bytes.Buffer
	_, err = io.Copy(&data, response.Body)
	if err != nil {
		return nil, err
	}
	return data.Bytes(), nil
}

func writeBatch(batch *[][]byte) {
	for _, song := range *batch {
		err := ioutil.WriteFile("data/song_"+uuid.New().String()+".mp3",song,0644)
		if err != nil {
			panic(err)
		}
	}
}

func randomTruncate(urls []string, diff int) []string {
	fmt.Println("Random truncating enabled.")
	r := rand.New(rand.NewSource(time.Now().Unix()))
	n := len(urls)-diff
	ret := make([]string,n)
	for i := 0; i < n; i++ {
		randIndex := r.Intn(len(urls))
		ret[i] = urls[randIndex]
		urls = append(urls[:randIndex],urls[randIndex+1:]...)
	}
	return ret
}