package main

import (
	"github.com/iostrovok/go-jsoncompare/jsoncompare"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var url1 string = "http://ergast.com/api/f1/2011.json"
var url2 string = "http://ergast.com/api/f1/2012.json"

func main() {
	b1, e1 := loadUrl(url1)
	if e1 != nil {
		log.Fatalln(e1)
	}

	b2, e2 := loadUrl(url2)
	if e2 != nil {
		log.Fatalln(e2)
	}

	list, err := jsoncompare.Compare(b1, b2)
	if err != nil {
		log.Fatalln(err)
	}

	leftOnly, rightOnly, noEqual, goodList := jsoncompare.SplitBySide(list)

	printList("GOOD: ", goodList)
	printList("Left Only: ", leftOnly)
	printList("Right Only: ", rightOnly)
	printList("No Equal: ", noEqual)

}

func printList(suff string, list []*jsoncompare.PathDiff) {
	for i, v := range list {

		viewPath := v.PathRight
		if v.PathLeft != "" {
			viewPath = v.PathLeft
		}

		fmt.Printf("%d. %s. path: %s\n", i, suff, viewPath)

		if !v.IsEqual {
			if v.ValueLeft != nil && v.ValueRight != nil {
				fmt.Printf("        %s <!=> %s\n", v.ValueLeft, v.ValueRight)
			} else if v.ValueLeft != nil {
				fmt.Printf("        %s\n", v.ValueLeft)
			} else if v.ValueRight != nil {
				fmt.Printf("        %s\n", v.ValueRight)
			}
		}
	}
}

func loadUrl(url string) ([]byte, error) {
	resp, err_get := http.Get(url)
	if err_get != nil {
		return nil, err_get
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
