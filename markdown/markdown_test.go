package markdown

import (
	"testing"

	"github.com/russross/blackfriday/v2"
	"gopkg.in/d4l3k/messagediff.v1"
)

func Test_renderer_Post(t *testing.T) {
	tym := &Time{}
	if err := tym.UnmarshalText([]byte("Sat Feb 22 15:18:37 GMT 2020")); err != nil {
		panic(err)
	}
	type fields struct {
		Renderer blackfriday.Renderer
	}
	tests := []struct {
		name   string
		fields fields
		want   *Post
		md     string
	}{
		{
			name: "hello",
			md: `% title = "markdown troff renderer"
% date = "Sat Feb 22 15:18:37 GMT 2020"
% tags = [
%   "acme",
%   "plan9"
% ]
% abstract = "somehting something"
% [[authors]]
% name = "Sevki"
% email = "s@sevki.org"
% twitter = "sevki"
% github = "sevki"
% affiliation = "funemployment"

# Heading 1`,
			fields: fields{
				Renderer: blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{}),
			},
			want: &Post{
				Title: "markdown troff renderer",
				Authors: []Author{
					{
						Name:        "Sevki",
						Email:       "s@sevki.org",
						Twitter:     "sevki",
						Github:      "sevki",
						Affiliation: "funemployment",
					},
				},
				Tags:     []string{"acme", "plan9"},
				Abstract: "somehting something",
				Date:     *tym,
				Slug:     "",
			},
		},
		{
			name: "bldy-and-harvey",
			md: `% title = "bldy and Harvey"
% date = "Sat Feb 22 15:18:37 GMT 2020"
% tags = [
%   "bldy",
%   "Harvey"
% ]
% abstract = "bldy and Harvey"
% [[authors]]
% name = "Sevki"
% email = "s@sevki.org"
% twitter = "sevki"
% github = "sevki"
			
bldy has reached a milestone in Harvey. It can now compile a fully working version of Harvey for the amd64 arch. It has been [1 Year, 4 Months, 23 Days since](https://groups.google.com/d/msg/harvey/IwK8-gebgyw/SVfuwv2LAAAJ) I started working on bldy. There is a lot of room to grow, a lot to fix but for now we have a working system. 
Thanks to the entire harvey team for being patient with me and thanks to [Ron Minnich](https://github.com/rminnich) for all his help and guidance.
![](https://ffbyt.es/bldy-and-harvey/bldy-and-harvey.png)  
`,
			fields: fields{
				Renderer: blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{}),
			},
			want: &Post{
				Title:    "bldy and Harvey",
				Abstract: "bldy and Harvey",
				Authors: []Author{
					{
						Name:    "Sevki",
						Twitter: "sevki",
						Github:  "sevki",
						Email:   "s@sevki.org",
					},
				},
				Tags: []string{"bldy", "Harvey"},
				Date: *tym,
				Slug: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRenderer(tt.fields.Renderer)
			blackfriday.Run([]byte(tt.md),
				blackfriday.WithRenderer(r),
				blackfriday.WithExtensions(blackfriday.Titleblock|blackfriday.CommonExtensions),
			)
			got := r.Post()
			if diff, equal := messagediff.PrettyDiff(got, tt.want); !equal {
				t.Errorf("renderer.Post() = %v, want %v", got, tt.want)
				t.Error(diff)
			}
		})
	}
}
