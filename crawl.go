package main


import (
    "fmt"
    "log"
    "net/http"
    "net/url"
    "io/ioutil"
    "golang.org/x/net/html"
    "bytes"
)

func get_html(string_url string) []byte {
	resp, err := http.Get(string_url)

	if err != nil {
		fmt.Println("error on http get")
		return nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("error on read body")
		return nil
	}


	return body
}

// extracts the href value from an anchor tag's
// href attribute. If no href is present we return an
// empty string.
//
// precondition: the token must be an anchor tag

func extract_href(token html.Token) string {
	for _, a := range token.Attr {
	    if a.Key == "href" {
	        return a.Val
	        break
	    }
	}
	return ""
}


// returns true for a valid absolute url or a valid relative url.
//
// e.g. https://aran.site -> true
//		/posts/nlp 		  -> true
// 		aran.site 		  -> false
//		posts/nlp 		  -> false

func is_valid_url(string_url string) bool {
	_, err := url.ParseRequestURI(string_url)
	// parse the request uri using url lib.

	if err != nil {
		// if an error occurs during parsing, we know the url is invalid.
	    return false
	}

	return true
	// no error has occured, the url is valid.
}

// filter out invalid urls, uses the above is_valid_url function
// if a url is relative we append its value to the original url
// and replace it in the urls array

func filter_urls(orig_url string, string_urls []string) []string {
	
	var valid_urls []string
	// declaring the array of valid urls

	for _, elem := range string_urls {
		// iterate through each url in the urls array

		if is_valid_url(elem) {
			// check if the url is valid

			if elem[0] == '/' {
				// if the url is relative, append it to the 
				// original url.
				elem = orig_url + elem
			}

			// append the element to the array of valid urls.
			valid_urls = append(valid_urls, elem)
		}
	}

	return valid_urls
}

// get all of the valid href links from a given webpage.

func get_links(string_url string) []string {
	var links []string
	
	htmlbytes := get_html(string_url)
	// get the bytes for the given url

	if htmlbytes != nil{
		r := bytes.NewReader(htmlbytes)
		// construct an ioReader from the bytes

		tokenizer := html.NewTokenizer(r)
		// construct the tokenizer from the ioReader

		for {
			// loop through every token in the document

			tokentype := tokenizer.Next()
			// get the next token

			switch tokentype{
			case html.ErrorToken:
				// the error token is hit when we reach the
				// end of the document.

				return filter_urls(string_url, links)
				// return the array of links
			
			case html.StartTagToken:

				token := tokenizer.Token()
				// get the token from the tokenizer

				if token.Data == "a" {
					// we've hit an anchor tag
					
					link := extract_href(token)
					// extract the href value (if one exists)

					if link != "" {
						// check if the tag has an href attribute

						links = append(links, link)
						// append the href to the array of links	
					}
				}
			}
		}
	}

	return []string{"error"}
}

func crawl(crawl_depth int, string_url string) []string {
	if crawl_depth == 0 {
		return []string{}
	}

	links := get_links(string_url)
	
	return links 
}


func main() {

	crawl(1, "https://aran.site")

	fmt.Println("Server running on 25565")

	http.HandleFunc("/crawl", func(writer http.ResponseWriter, req *http.Request) {
		fmt.Fprint(writer, crawl(1, "https://github.com/aranscope"))
	})

	http.HandleFunc("/test", func(writer http.ResponseWriter, req *http.Request) {
		fmt.Fprint(writer, "Hello World")
	})

    log.Fatal(http.ListenAndServe(":25565", nil))
}
