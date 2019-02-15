package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"unicode/utf8"
	"words_of_boobs/generator"
)

func writeLog(value string) {
	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	_, err = f.WriteString(value + "\n")
	if err != nil {
		log.Println(err)
	}
	f.Close()
}

func Start(port int) error {
	http.Handle("/", http.FileServer(http.Dir("./html")))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./results"))))
	http.HandleFunc("/api/generate", func(w http.ResponseWriter, r *http.Request) {
		var (
			values         []string
			filename, text string
			fontName       string
			imgSetName     string
			ok             bool
			width          int = 0
			err            error
		)

		params := r.URL.Query()
		if values, ok = params["text"]; !ok {
			w.Write([]byte("error: no text parameter"))
			return
		}
		text = string(values[0])
		if utf8.RuneCountInString(text) > 15 {
			w.Write([]byte("error: maximum text length = 15 symbols"))
			return
		}

		if values, ok = params["width"]; ok {
			if width, err = strconv.Atoi(values[0]); err != nil {
				log.Println(err)
				w.Write([]byte("error: incorrect width value " + values[0]))
				return
			}
		}
		if width > 10000 {
			w.Write([]byte("error: maximum width = 10000"))
			return
		}

		if values, ok = params["font"]; !ok {
			w.Write([]byte("error: no font parameter"))
			return
		}
		fontName = string(values[0])
		if utf8.RuneCountInString(fontName) > 24 {
			w.Write([]byte("error: maximum font name length = 24 symbols"))
			return
		}

		if values, ok = params["setname"]; !ok {
			w.Write([]byte("error: no setname parameter"))
			return
		}
		imgSetName = string(values[0])
		if utf8.RuneCountInString(imgSetName) > 15 {
			w.Write([]byte("error: maximum set of img name length = 15 symbols"))
			return
		}

		filename, err = generator.GenerateImageForText(text, fontName, imgSetName, 1000, width)
		if err != nil {
			log.Println(err)
			w.Write([]byte("error: something wrong"))
			return
		}

		writeLog(text)
		log.Printf("generated %s for text '%s' with width=%d\n", filename, text, width)
		w.Write([]byte(filename))
	})

	http.HandleFunc("/api/sets", func(w http.ResponseWriter, r *http.Request) {
		type sets struct {
			Imgs  []string `json:"imgs"`
			Fonts []string `json:"fonts"`
		}

		s := sets{
			Imgs: []string{
				"boobs",
				"it",
				"stiic",
			},
			Fonts: []string{
				"NotoSans-Bold.ttf",
				"NotoSans-BoldItalic.ttf",
				"NotoSans-Italic.ttf",
				"NotoSans-Regular.ttf",
				"NotoSerif-Bold.ttf",
				"Oswald-Bold.ttf",
				"Oswald-ExtraLight.ttf",
				"Oswald-Light.ttf",
				"Oswald-Medium.ttf",
				"Oswald-Regular.ttf",
				"Oswald-SemiBold.ttf",
				"Oxygen-Bold.ttf",
				"Oxygen-Light.ttf",
				"Oxygen-Regular.ttf",
				"SourceSansPro-Bold.ttf",
				"Symbola.ttf",
			},
		}
		res, _ := json.Marshal(s)

		w.Write(res)
	})

	log.Printf("started on %d port\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
