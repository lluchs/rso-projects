package main

import (
	"fmt"
	"html/template"
	"os"
	"sort"
	"time"

	"github.com/turnage/graw/reddit"
)

// Project holds information on an ongoing RSO project.
type Project struct {
	Title     template.HTML // already escaped from the Reddit API
	Organizer string
	URL       string
	StartDate string // ISO 8601
	EndDate   string // ISO 8601

	Registers             []string // sorted nicely
	InstrumentsByRegister map[string][]Instrument
}

// ProjectsByEndDate implements sort.Interface for soring by EndDate.
type ProjectsByEndDate []Project

func (a ProjectsByEndDate) Len() int           { return len(a) }
func (a ProjectsByEndDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ProjectsByEndDate) Less(i, j int) bool { return a[i].EndDate < a[j].EndDate }

func instrumentsByRegister(instruments []Instrument) map[string][]Instrument {
	m := make(map[string][]Instrument)
	for _, instr := range instruments {
		m[instr.Register] = append(m[instr.Register], instr)
	}
	return m
}

func createHTMLPage(posts []reddit.Post) {
	var projects []Project

	for _, post := range posts {
		if !isProject(&post) {
			continue
		}
		deadline := findDeadline(post.SelfText, int64(post.CreatedUTC))
		// Filter out finished projects.
		if time.Now().After(deadline) {
			continue
		}
		byreg := instrumentsByRegister(findInstruments(post.SelfText))
		var registers []string
		for _, reg := range Registers {
			if _, ok := byreg[reg]; ok {
				registers = append(registers, reg)
			}
		}
		p := Project{
			Title:                 template.HTML(post.Title),
			Organizer:             post.Author,
			URL:                   post.URL,
			StartDate:             time.Unix(int64(post.CreatedUTC), 0).Format("2006-01-02"),
			EndDate:               deadline.Format("2006-01-02"),
			Registers:             registers,
			InstrumentsByRegister: byreg,
		}
		projects = append(projects, p)
	}

	sort.Sort(ProjectsByEndDate(projects))

	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := os.Create("static/index.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = tmpl.Execute(f, map[string]interface{}{
		"Projects": projects,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}
