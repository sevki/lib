// Package markdown are a set of common markdownparser utils used in
// pllaces like sevki.org/troff
package markdown

import (
	"io"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v2"
)

// Time is markdown post block time format
type Time struct {
	time.Time
}

// YYYY-MM-DD hh:mm:ss tz
const fuzzyFormat = "2006-01-02 15:04:05-07:00"

// UnmarshalText satisfies textunmarshaller
func (t *Time) UnmarshalText(text []byte) error {
	postTime, err := dateparse.ParseAny(string(text))
	if err != nil {
		return err
	}
	t.Time = postTime
	return nil
}

// Author information
type Author struct {
	Name        string
	Affiliation string
	Email       string
	Twitter     string
	Github      string
}

// Post holds info about the post
type Post struct {
	Title    string
	Date     Time
	Slug     string
	Authors  []Author
	Abstract string
	Tags     []string
}

// Renderer wraps the blackfriday.Renderer interface
type Renderer interface {
	blackfriday.Renderer
	Post() *Post
}

type renderer struct {
	blackfriday.Renderer
	parsingTitle bool
	title        *Post
}

func (r *renderer) RenderNode(w io.Writer, n *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	switch n.Type {
	case blackfriday.Heading:
		r.parsingTitle = n.IsTitleblock
		if n.IsTitleblock {
			if entering {
				return blackfriday.GoToNext
			}
			return blackfriday.SkipChildren
		}
	case blackfriday.Text:
		if r.parsingTitle {
			magic := []byte(`title:`)
			// BUG(sevki): we get things that are not title blocks here
			isTitle := strings.Index(string(n.Literal), string(magic)) == 0
			if err := yaml.Unmarshal(n.Literal, &r.title); isTitle && err != nil {
				panic(err.Error())
				return blackfriday.Terminate
			}
			r.parsingTitle = false
			return blackfriday.SkipChildren
		}
	}
	return r.Renderer.RenderNode(w, n, entering)
}

func (r *renderer) Post() *Post {
	return r.title
}

// NewRenderer takes a blackfriday.Renderer and returns a Renderer
func NewRenderer(r blackfriday.Renderer) Renderer { return &renderer{r, false, nil} }
