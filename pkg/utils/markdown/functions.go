package markdown

import (
	"bytes"
	"encoding/xml"
	"io"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"

	"spacetime/pkg/utils/types"
)

var htmlPolicy *bluemonday.Policy

func init() {
	// Use bluemonday policy for life of program
	htmlPolicy = bluemonday.NewPolicy()

	// Disallow images

	// Allow elements
	htmlPolicy.AllowElements(
		"h1", "h2", "h3", "h4", "h5", "h6",
		"p", "span", "em", "strong", "del", "sup", "sub",
		"ul", "ol", "li",
		"table", "thead", "tbody", "th", "tr", "td",
		"blockquote",
		"pre", "code",
		"hr", "br")

	// Allow links
	htmlPolicy.AllowAttrs("href", "title").OnElements("a")

	// Validate URLs using net/url.Parse and require 'mailto:', 'http://' or 'https://'
	htmlPolicy.AllowURLSchemes("mailto", "http", "https")
	htmlPolicy.AllowRelativeURLs(true) // links within crowdspec
	htmlPolicy.RequireParseableURLs(true)
	htmlPolicy.RequireNoFollowOnLinks(true)
	htmlPolicy.AddTargetBlankToFullyQualifiedLinks(true)

	// Allow code blocks with language class
	htmlPolicy.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")

	// Allow simple styles to dress up tables
	htmlPolicy.AllowAttrs("style").OnElements("table", "thead", "tbody", "th", "td")

	// style="border-spacing:50px;" and style="border-spacing:3em;" on table
	var matchValue = regexp.MustCompile(`^(\d{1,3})(px|rem|em)(?: (\d{1,3})(px|rem|em))?$`)
	var checkRange = func(a string, min, max int) bool {
		v, err := types.AtoInt(a)
		if err != nil {
			return false
		}
		return v >= min && v <= max
	}
	const maxPx = 50
	const maxEm = 3
	htmlPolicy.AllowStyles("border-spacing").MatchingHandler(func(s string) bool {
		if matchValue.MatchString(s) {
			var parts = matchValue.FindStringSubmatch(s)
			// first unit will be defined as px|rem|em
			if parts[2] == "px" {
				if !checkRange(parts[1], 0, maxPx) {
					return false
				}
			} else if parts[2] == "rem" || parts[2] == "em" {
				if !checkRange(parts[1], 0, maxEm) {
					return false
				}
			}
			// second unit may be defined
			if parts[4] == "px" {
				if !checkRange(parts[3], 0, maxPx) {
					return false
				}
			} else if parts[4] == "rem" || parts[4] == "em" {
				if !checkRange(parts[3], 0, maxEm) {
					return false
				}
			}
			return true
		}
		return false
	}).OnElements("table")

	// style="width:1000px;" and style="width:5%;" on table, th, td
	htmlPolicy.AllowStyles("width").MatchingHandler(func(w string) bool {
		w = strings.TrimSpace(w)
		if strings.HasSuffix(w, "%") {
			i, err := types.AtoInt(w[:len(w)-1])
			return err == nil && i >= 5 && i <= 100
		} else if strings.HasSuffix(w, "px") {
			i, err := types.AtoInt(w[:len(w)-2])
			return err == nil && i >= 25 && i <= 1000
		}
		return false
	}).OnElements("table", "th", "td")

	// style="white-space:nowrap;" on applicable elements
	htmlPolicy.AllowStyles("white-space").MatchingEnum(
		"nowrap", "pre", "pre-wrap", "pre-line",
	).OnElements("table", "thead", "tbody", "tr", "th", "td")

	// style="text-align:left;" on applicable elements
	htmlPolicy.AllowStyles("text-align").MatchingEnum(
		"left", "center", "right",
	).OnElements("table", "thead", "tbody", "tr", "th", "td")

	// style="vertival-align:top;" on applicable elements
	htmlPolicy.AllowStyles("vertical-align").MatchingEnum(
		"top", "middle", "bottom",
	).OnElements("table", "thead", "tbody", "tr", "th", "td")
}

func RenderMarkdown(markdown string) (string, error) {

	var output = blackfriday.Run([]byte(markdown),
		blackfriday.WithExtensions(
			blackfriday.SpaceHeadings|
				blackfriday.Autolink|
				blackfriday.Strikethrough|
				blackfriday.FencedCode|
				blackfriday.Tables|
				blackfriday.BackslashLineBreak))

	output = htmlPolicy.SanitizeBytes(output)

	// Validate HTML
	var decoder = xml.NewDecoder(bytes.NewBuffer(output))
	for {
		// xml.Decoder.Token guarantees start/end tags are properly nested and closed
		_, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
	}

	// log.Println("Output:" + string(output))

	return string(output), nil
}
