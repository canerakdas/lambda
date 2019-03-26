package crawlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func Image(url string, name string) {

	if strings.Contains(url, "http") || strings.Contains(url, "https") {
		response, e := http.Get(url)
		if e != nil {
			fmt.Println(e)
		}
		defer response.Body.Close()

		//open a file for writing
		file, err := os.Create(name)
		if err != nil {
			fmt.Println(e)
		}
		defer file.Close()

		// Use io.Copy to just dump the response body to the file. This supports huge files
		_, err = io.Copy(file, response.Body)
		if err != nil {
			fmt.Println(e)
		}

	} else {
		fmt.Println("Protocol error:", url)
	}
}

/*
	Usage:
	path := []string{"/home/fishy/image/",strconv.Itoa(v.SteamAppid),".jpg" }
	Image(v.Background,strings.Join(path,""))
*/
