package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

// cli
type CliUiData struct {
	Root string
	Dirs []Directory
}

// form page
type FormPageDatas struct {
	RedirectURL string
}

type WebUiData struct {
	Dirtitle     string
	PreviousPage string
	Queries      string
	Directories  []Directory
	ProgessPah   []ProgessPah
	SortOptions  []SortOption
}

type SortOption struct {
	Title    string
	Name     string
	Selected bool
}

type ProgessPah struct {
	Title    string
	Url      string
	SlashPre bool
}

var sortOptions = []SortOption{
	{
		Title: "a-z",
		Name:  "nameasc",
	},
	{
		Title: "z-a",
		Name:  "namedesc",
	},
	{
		Title: "smallar to larger",
		Name:  "sizeasc",
	},
	{
		Title: "larger to smaller",
		Name:  "sizedesc",
	},
	{
		Title: "directory first",
		Name:  "bydir",
	},
	{
		Title: "file first",
		Name:  "byfile",
	},
}

func ServeWebUi(w http.ResponseWriter, r *http.Request) {
	var err error
	var datas WebUiData
	datas.Directories, err = directories(w, r)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// previus page sananagans
	if split := strings.Split(r.URL.EscapedPath(), "/"); len(split) > 2 {
		datas.PreviousPage = strings.Join(split[:len(split)-1], "/")
	} else {
		datas.PreviousPage = "/"
	}

	datas.ProgessPah = possiblePahts(r)
	datas.Queries = allQueries(r)
	datas.Dirtitle = datas.ProgessPah[len(datas.ProgessPah)-1].Title
	datas.SortOptions = sortOptions

	// gigure out which sort is in use
	if s := r.FormValue("sort"); s != "" {
		for i, options := range datas.SortOptions {
			if options.Name == s {
				datas.SortOptions[i].Selected = true
				break
			}
		}
	}

	// while debuging re-read tihe files for eatch request
	if debug {
		indexTemplate, err = template.ParseFiles(
			"tailwind/src/index.html",
			"tailwind/src/form.html",
			"tailwind/src/index/components.html",
			"tailwind/src/index/list.html",
			"tailwind/src/index/index.js",
		)
		if err != nil {
			log.Println(err)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html")
	if err := indexTemplate.ExecuteTemplate(w, "main", datas); err != nil {
		log.Println(err)
		return
	}
}
