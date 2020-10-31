// Package markdown are a set of common markdownparser utils used in
// pllaces like sevki.org/troff
package markdown

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/araddon/dateparse"
	"github.com/russross/blackfriday/v2"
	"sevki.org/x/oututil"
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
			magic := []byte(`title=`)
			// BUG(sevki): we get things that are not title blocks here
			isTitle := strings.Index(string(n.Literal), string(magic)) == 0
			if err := toml.Unmarshal(n.Literal, &r.title); isTitle && err != nil {
				return blackfriday.Terminate
			}
			r.parsingTitle = false
			return blackfriday.GoToNext
		}
	case blackfriday.Image:
		if entering {
			s := strings.Split(string(n.LinkData.Destination), "=")
			filename := strings.TrimSpace(s[0])
			sizebit := "0x0"
			if len(s) > 1 {
				sizebit = s[1]
			}
			width, height := parseSize(sizebit)
			if path.Ext(filename) == ".dot" {
				bytez, err := dot(filename, int(width), int(height))
				if err != nil {
					return blackfriday.Terminate
				}
				w.Write(bytez)
				return blackfriday.SkipChildren
			}
			n.LinkData.Destination = []byte(filename)
		}
	}
	return r.Renderer.RenderNode(w, n, entering)
}

func (r *renderer) Post() *Post {
	return r.title
}

// NewRenderer takes a blackfriday.Renderer and returns a Renderer
func NewRenderer(r blackfriday.Renderer) Renderer { return &renderer{r, false, nil} }
func parseSize(s string) (w int, h int) {
	fmt.Sscanf(strings.TrimSpace(s), "%dx%d", &w, &h)
	return
}

func dot(filename string, w int, h int) ([]byte, error) {
	psname := strings.Replace(filename, ".dot", ".svg", -1)
	_, psname = path.Split(psname)

	args := []string{"-T", "svg", filename}

	if w > 0 && h > 0 {
		//	args = append(args, "-size", fmt.Sprintf("%dx%d", w, h))
	}
	cmd := exec.CommandContext(context.Background(), "dot", args...)
	stdOut := bytes.NewBuffer(nil)
	stdErr := bytes.NewBuffer(nil)
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr
	fmt.Printf("running dot: %v\n", args)
	if err := cmd.Run(); err != nil {
		oututil.Indent(stdErr, 1)
		fmt.Fprintf(stdErr, `dot: %v
args= %v
%s
`,
			err,
			args,
			stdErr.String(),
		)
		return nil, errors.New(stdErr.String())
	}
	return stdOut.Bytes(), nil
}

func convertToPs(filename string, w int, h int) (string, error) {
	psname := strings.Replace(filename, ".ps", ".svg", -1)
	_, psname = path.Split(psname)

	args := []string{"convert"}

	if w > 0 && h > 0 {
		args = append(args, "-size", fmt.Sprintf("%dx%d", w, h))
	}
	args = append(args, filename, psname)
	cmd := exec.CommandContext(context.Background(), "magick", args...)
	stdOut := bytes.NewBuffer(nil)
	stdErr := bytes.NewBuffer(nil)
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr
	fmt.Printf("running magick: %v\n", args)
	if err := cmd.Run(); err != nil {
		oututil.Indent(stdErr, 1)
		fmt.Fprintf(stdErr, `magick: %v
args= %v
%s
`,
			err,
			args,
			stdErr.String(),
		)
		return "", errors.New(stdErr.String())
	}
	return psname, nil
}
