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

func Start(sets *generator.Generator, port int) error {
	http.Handle("/", http.FileServer(http.Dir("./html")))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./results"))))
	http.HandleFunc("/api/generate", func(w http.ResponseWriter, r *http.Request) {
		var (
			values         []string
			filename, text string
			fontName       string
			imgCategory    string
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
			fontName = generator.DEFAULT_FONT
		} else {
			fontName = string(values[0])
		}
		if utf8.RuneCountInString(fontName) > 24 || !sets.IsFont(fontName) {
			w.Write([]byte("error: wrong font name"))
			return
		}

		if values, ok = params["category"]; !ok {
			imgCategory = generator.DEFAULT_IMG
		} else {
			imgCategory = string(values[0])
		}
		if utf8.RuneCountInString(imgCategory) > 15 || !sets.IsImageSet(imgCategory) {
			w.Write([]byte("error: wrong img category name"))
			return
		}

		filename, err = generator.GenerateImageForText(text, fontName, imgCategory, 1000, width)
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
		type trSets struct {
			Imgs  []string `json:"imgs"`
			Fonts []string `json:"fonts"`
		}
		st := trSets{}
		imgs := sets.GetImages()
		fonts := sets.GetFonts()
		st.Fonts = make([]string, 0, len(imgs))
		st.Imgs = make([]string, 0, len(fonts))
		for k := range imgs {
			st.Imgs = append(st.Imgs, k)
		}
		for k := range fonts {
			st.Fonts = append(st.Fonts, k)
		}

		res, _ := json.Marshal(st)
		w.Write(res)
	})

	log.Printf("started on %d port\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
