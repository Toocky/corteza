package automation

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type (
	htmlHandler struct {
		reg htmlHandlerRegistry
	}
)

func HtmlHandler(reg htmlHandlerRegistry) *htmlHandler {
	h := &htmlHandler{
		reg: reg,
	}

	h.register()
	return h
}

// blockType represents the structural type of a section item.
type blockType string

const (
	blockHeading   blockType = "heading"
	blockParagraph blockType = "paragraph"
	blockListItem  blockType = "list_item"
)

// sectionItem is a single structural block within a section.
type sectionItem struct {
	Type        blockType `json:"type"`
	Text        string    `json:"text"`
	IndentLevel int       `json:"indent_level"`
}

type customSection struct {
	SectionTitle string        `json:"section_title"`
	SectionItems []sectionItem `json:"section_items"`
}

type customSectionsResult struct {
	CustomSections []customSection `json:"custom_sections"`
}

func (h htmlHandler) toJson(_ context.Context, args *htmlToJsonArgs) (res *htmlToJsonResults, err error) {
	result, err := htmlToCustomSections(args.Html)
	if err != nil {
		return nil, err
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return &htmlToJsonResults{ResultJson: string(jsonBytes)}, nil
}

func (h htmlHandler) toMarkdown(_ context.Context, args *htmlToMarkdownArgs) (res *htmlToMarkdownResults, err error) {
	md, err := htmlToMarkdownString(args.Html)
	if err != nil {
		return nil, err
	}
	return &htmlToMarkdownResults{Markdown: md}, nil
}

// htmlToMarkdownString converts an HTML string to Markdown.
func htmlToMarkdownString(rawHTML string) (string, error) {
	doc, err := html.Parse(strings.NewReader(rawHTML))
	if err != nil {
		return "", err
	}

	body := findHTMLBody(doc)
	if body == nil {
		body = doc
	}

	var sb strings.Builder
	htmlNodeToMarkdown(&sb, body, 0, false, false)

	// Clean up excessive blank lines
	result := strings.TrimSpace(sb.String())
	for strings.Contains(result, "\n\n\n") {
		result = strings.ReplaceAll(result, "\n\n\n", "\n\n")
	}
	return result, nil
}

// htmlNodeToMarkdown recursively converts an HTML node tree to Markdown.
func htmlNodeToMarkdown(sb *strings.Builder, n *html.Node, listDepth int, ordered bool, inPre bool) {
	if n.Type == html.TextNode {
		text := n.Data
		if !inPre {
			text = strings.Join(strings.Fields(text), " ")
		}
		sb.WriteString(text)
		return
	}

	if n.Type != html.ElementNode {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			htmlNodeToMarkdown(sb, c, listDepth, ordered, inPre)
		}
		return
	}

	tag := n.Data

	switch tag {
	case "h1", "h2", "h3", "h4", "h5", "h6":
		level := int(tag[1] - '0')
		sb.WriteString("\n")
		sb.WriteString(strings.Repeat("#", level))
		sb.WriteString(" ")
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			htmlNodeToMarkdown(sb, c, listDepth, ordered, false)
		}
		sb.WriteString("\n\n")

	case "p":
		sb.WriteString("\n")
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			htmlNodeToMarkdown(sb, c, listDepth, ordered, false)
		}
		sb.WriteString("\n")

	case "br":
		sb.WriteString("  \n")

	case "strong", "b":
		sb.WriteString("**")
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			htmlNodeToMarkdown(sb, c, listDepth, ordered, inPre)
		}
		sb.WriteString("**")

	case "em", "i":
		sb.WriteString("*")
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			htmlNodeToMarkdown(sb, c, listDepth, ordered, inPre)
		}
		sb.WriteString("*")

	case "code":
		if inPre {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				htmlNodeToMarkdown(sb, c, listDepth, ordered, true)
			}
		} else {
			sb.WriteString("`")
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				htmlNodeToMarkdown(sb, c, listDepth, ordered, false)
			}
			sb.WriteString("`")
		}

	case "pre":
		sb.WriteString("\n```\n")
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			htmlNodeToMarkdown(sb, c, listDepth, ordered, true)
		}
		sb.WriteString("\n```\n")

	case "blockquote":
		var inner strings.Builder
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			htmlNodeToMarkdown(&inner, c, listDepth, ordered, false)
		}
		for _, line := range strings.Split(strings.TrimSpace(inner.String()), "\n") {
			sb.WriteString("> ")
			sb.WriteString(line)
			sb.WriteString("\n")
		}
		sb.WriteString("\n")

	case "a":
		href := htmlAttr(n, "href")
		sb.WriteString("[")
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			htmlNodeToMarkdown(sb, c, listDepth, ordered, false)
		}
		sb.WriteString("](")
		sb.WriteString(href)
		sb.WriteString(")")

	case "img":
		alt := htmlAttr(n, "alt")
		src := htmlAttr(n, "src")
		sb.WriteString("![")
		sb.WriteString(alt)
		sb.WriteString("](")
		sb.WriteString(src)
		sb.WriteString(")")

	case "ul":
		if listDepth == 0 {
			sb.WriteString("\n")
		}
		counter := 0
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "li" {
				counter++
				htmlRenderListItem(sb, c, listDepth, false, counter)
			}
		}
		if listDepth == 0 {
			sb.WriteString("\n")
		}

	case "ol":
		if listDepth == 0 {
			sb.WriteString("\n")
		}
		counter := 0
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "li" {
				counter++
				htmlRenderListItem(sb, c, listDepth, true, counter)
			}
		}
		if listDepth == 0 {
			sb.WriteString("\n")
		}

	case "hr":
		sb.WriteString("\n---\n\n")

	case "table":
		htmlTableToMarkdown(sb, n)

	default:
		// div, section, span, etc — just recurse
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			htmlNodeToMarkdown(sb, c, listDepth, ordered, inPre)
		}
	}
}

// htmlRenderListItem renders a single <li> as a Markdown list item.
func htmlRenderListItem(sb *strings.Builder, li *html.Node, depth int, ordered bool, counter int) {
	indent := strings.Repeat("  ", depth)

	if ordered {
		sb.WriteString(fmt.Sprintf("%s%d. ", indent, counter))
	} else {
		sb.WriteString(indent)
		sb.WriteString("- ")
	}

	// Render direct inline content (skip nested ul/ol)
	for c := li.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && (c.Data == "ul" || c.Data == "ol") {
			continue
		}
		htmlNodeToMarkdown(sb, c, depth+1, ordered, false)
	}
	sb.WriteString("\n")

	// Render nested lists
	for c := li.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && (c.Data == "ul" || c.Data == "ol") {
			htmlNodeToMarkdown(sb, c, depth+1, c.Data == "ol", false)
		}
	}
}

// htmlTableToMarkdown converts a <table> to a Markdown table.
func htmlTableToMarkdown(sb *strings.Builder, table *html.Node) {
	var rows [][]string

	var walkTable func(*html.Node)
	walkTable = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			var cells []string
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && (c.Data == "td" || c.Data == "th") {
					cells = append(cells, htmlCleanText(c))
				}
			}
			rows = append(rows, cells)
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walkTable(c)
		}
	}
	walkTable(table)

	if len(rows) == 0 {
		return
	}

	sb.WriteString("\n")
	// Header row
	sb.WriteString("| ")
	sb.WriteString(strings.Join(rows[0], " | "))
	sb.WriteString(" |\n")

	// Separator
	sep := make([]string, len(rows[0]))
	for i := range sep {
		sep[i] = "---"
	}
	sb.WriteString("| ")
	sb.WriteString(strings.Join(sep, " | "))
	sb.WriteString(" |\n")

	// Data rows
	for _, row := range rows[1:] {
		// Pad if row has fewer cells than header
		for len(row) < len(rows[0]) {
			row = append(row, "")
		}
		sb.WriteString("| ")
		sb.WriteString(strings.Join(row, " | "))
		sb.WriteString(" |\n")
	}
	sb.WriteString("\n")
}

// htmlAttr returns the value of the named attribute on an element node.
func htmlAttr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

// htmlToCustomSections parses an HTML string and returns a customSectionsResult.
//
// The highest heading level found (e.g. h2) is used as the section boundary
// and becomes section_title. Lower-level headings become heading blocks inside
// section_items. <p> tags become paragraph blocks. <li> elements become
// list_item blocks with indent_level derived from nesting depth.
func htmlToCustomSections(rawHTML string) (customSectionsResult, error) {
	doc, err := html.Parse(strings.NewReader(rawHTML))
	if err != nil {
		return customSectionsResult{}, err
	}

	body := findHTMLBody(doc)
	if body == nil {
		body = doc
	}

	nodes := flattenHTMLBody(body)
	topLevel := findTopHeadingLevel(nodes)

	var sections []customSection

	if topLevel == 0 {
		var items []sectionItem
		for _, n := range nodes {
			items = append(items, htmlNodeToItems(n, 0)...)
		}
		if len(items) > 0 {
			sections = append(sections, customSection{
				SectionTitle: "Untitled Section",
				SectionItems: items,
			})
		}
	} else {
		var current *customSection

		for _, n := range nodes {
			tag := htmlTagName(n)
			level := htmlHeadingLevel(tag)
			text := htmlCleanText(n)

			if level == topLevel {
				if current != nil {
					sections = append(sections, *current)
				}
				current = &customSection{
					SectionTitle: text,
					SectionItems: []sectionItem{},
				}
			} else if current != nil {
				current.SectionItems = append(current.SectionItems, htmlNodeToItems(n, topLevel)...)
			} else {
				current = &customSection{
					SectionTitle: "Untitled Section",
					SectionItems: []sectionItem{},
				}
				current.SectionItems = append(current.SectionItems, htmlNodeToItems(n, topLevel)...)
			}
		}

		if current != nil {
			sections = append(sections, *current)
		}
	}

	var filtered []customSection
	for _, s := range sections {
		if len(s.SectionItems) > 0 || s.SectionTitle != "Untitled Section" {
			filtered = append(filtered, s)
		}
	}

	if filtered == nil {
		filtered = []customSection{}
	}

	return customSectionsResult{CustomSections: filtered}, nil
}

func htmlNodeToItems(n *html.Node, topLevel int) []sectionItem {
	tag := htmlTagName(n)

	// Handle lists before extracting text — extracting text from a list node
	// would redundantly walk all nested children.
	if tag == "ul" || tag == "ol" {
		return parseHTMLList(n)
	}

	text := htmlCleanText(n)
	if text == "" {
		return nil
	}

	level := htmlHeadingLevel(tag)
	switch {
	case level > 0 && level > topLevel:
		return []sectionItem{{Type: blockHeading, Text: text, IndentLevel: 0}}
	case tag == "p":
		return []sectionItem{{Type: blockParagraph, Text: text, IndentLevel: 0}}
	}
	return nil
}

func parseHTMLList(n *html.Node) []sectionItem {
	var items []sectionItem
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "li" {
			items = append(items, parseHTMLListItem(c)...)
		}
	}
	return items
}

func parseHTMLListItem(li *html.Node) []sectionItem {
	indent := htmlListIndent(li)

	var textParts []string
	for c := li.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			textParts = append(textParts, c.Data)
		} else if c.Type == html.ElementNode {
			ct := c.Data
			if ct != "ul" && ct != "ol" {
				textParts = append(textParts, htmlExtractText(c))
			}
		}
	}

	directText := strings.Join(textParts, "")
	directText = strings.Join(strings.Fields(directText), " ")

	var items []sectionItem
	if directText != "" {
		items = append(items, sectionItem{
			Type:        blockListItem,
			Text:        directText,
			IndentLevel: indent,
		})
	}

	for c := li.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && (c.Data == "ul" || c.Data == "ol") {
			items = append(items, parseHTMLList(c)...)
		}
	}

	return items
}

func flattenHTMLBody(body *html.Node) []*html.Node {
	var nodes []*html.Node
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type != html.ElementNode {
			return
		}
		tag := n.Data
		if tag == "div" || tag == "section" || tag == "article" || tag == "main" {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				walk(c)
			}
		} else {
			nodes = append(nodes, n)
		}
	}
	for c := body.FirstChild; c != nil; c = c.NextSibling {
		walk(c)
	}
	return nodes
}

func findTopHeadingLevel(nodes []*html.Node) int {
	for level := 1; level <= 6; level++ {
		for _, n := range nodes {
			if htmlHeadingLevel(htmlTagName(n)) == level {
				return level
			}
		}
	}
	return 0
}

func htmlListIndent(n *html.Node) int {
	level := 0
	node := n.Parent
	for node != nil {
		t := htmlTagName(node)
		if t == "ul" || t == "ol" {
			level++
		}
		node = node.Parent
	}
	if level > 6 {
		return 6
	}
	return level
}

func htmlCleanText(n *html.Node) string {
	return strings.Join(strings.Fields(htmlExtractText(n)), " ")
}

func htmlExtractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var sb strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		sb.WriteString(htmlExtractText(c))
	}
	return sb.String()
}

func htmlHeadingLevel(tag string) int {
	if len(tag) == 2 && tag[0] == 'h' && tag[1] >= '1' && tag[1] <= '6' {
		return int(tag[1] - '0')
	}
	return 0
}

func htmlTagName(n *html.Node) string {
	if n == nil || n.Type != html.ElementNode {
		return ""
	}
	return n.Data
}

func findHTMLBody(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "body" {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found := findHTMLBody(c); found != nil {
			return found
		}
	}
	return nil
}
