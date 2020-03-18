package markdown

import (
	"reflect"
	"testing"

	"github.com/russross/blackfriday/v2"
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
			md: `% title: markdown troff renderer
% authors:
% - name: Sevki
%   email: s@sevki.org
%   affiliation: funemployment
% date: Sat Feb 22 15:18:37 GMT 2020
% 
% tags: [acme, plan9]
% abstract: Creating documents using plan9 troff 
# Heading 1`,
			fields: fields{
				Renderer: blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{}),
			},
			want: &Post{
				Title:    "markdown troff renderer",
				Abstract: "Creating documents using plan9 troff",
				Authors: []Author{
					Author{
						Name:        "Sevki",
						Email:       "s@sevki.org",
						Affiliation: "funemployment",
					},
				},
				Tags: []string{"acme", "plan9"},
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

			if got := r.Post(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("renderer.Post() = %v, want %v", got, tt.want)
			}
		})
	}
}
