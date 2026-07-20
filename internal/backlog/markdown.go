package backlog

import (
	"html"
	"html/template"
	"regexp"
	"strconv"
	"strings"
)

// renderMarkdown renders a safe inline-markdown subset (bold, italic, inline
// code, links) to HTML. It is the ONLY place bklg emits template.HTML (raw), so
// safety is by construction:
//
//  1. The input is HTML-escaped FIRST (html.EscapeString), so no repo-authored
//     markup can survive — the only tags in the output are the ones this
//     function adds.
//  2. The tags it adds are a fixed whitelist (strong/em/code/a); link hrefs are
//     scheme-checked (safeURL) so javascript:/data: etc. never become live links.
//
// Unmatched markers stay literal (matched-pairs only), so a title truncated
// mid-`**` can't produce an unclosed tag. Block constructs (lists/headings) are
// out of scope for v1.
func renderMarkdown(s string) template.HTML {
	out := html.EscapeString(s)

	// 1) Code spans first, lifted out to placeholders so their contents aren't
	//    re-processed as bold/italic/links. Contents are already escaped.
	var codes []string
	out = codeRe.ReplaceAllStringFunc(out, func(m string) string {
		codes = append(codes, "<code>"+codeRe.FindStringSubmatch(m)[1]+"</code>")
		return "\x00" + strconv.Itoa(len(codes)-1) + "\x00"
	})

	// 2) Links — only whitelisted schemes become <a>; anything else stays literal.
	out = linkRe.ReplaceAllStringFunc(out, func(m string) string {
		sub := linkRe.FindStringSubmatch(m)
		text, url := sub[1], sub[2]
		if !safeURL(url) {
			return m
		}
		return `<a href="` + url + `" rel="noopener noreferrer">` + text + `</a>`
	})

	// 3) Bold before italic (so ** isn't eaten by *). Asterisk markers only —
	//    underscore emphasis is intentionally unsupported so snake_case and
	//    __dunder__ identifiers (common in these docs) aren't mangled.
	out = boldStarRe.ReplaceAllString(out, "<strong>$1</strong>")
	out = italStarRe.ReplaceAllString(out, "<em>$1</em>")

	// 4) Restore code spans.
	for i, c := range codes {
		out = strings.Replace(out, "\x00"+strconv.Itoa(i)+"\x00", c, 1)
	}
	return template.HTML(out)
}

var (
	codeRe     = regexp.MustCompile("`([^`]+)`")
	linkRe     = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	boldStarRe = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	italStarRe = regexp.MustCompile(`\*([^*]+)\*`)
)

// safeURL reports whether a link target is safe to place in an href. Only
// http(s)/mailto and root-relative (/…) or fragment (#…) targets, or a
// scheme-less relative path, are allowed; unknown schemes (javascript:, data:,
// vbscript:, …) are rejected so a repo can't inject a script link.
func safeURL(u string) bool {
	u = strings.TrimSpace(u)
	low := strings.ToLower(u)
	switch {
	case strings.HasPrefix(low, "http://"), strings.HasPrefix(low, "https://"), strings.HasPrefix(low, "mailto:"):
		return true
	case strings.HasPrefix(u, "/"), strings.HasPrefix(u, "#"):
		return true
	case !strings.Contains(u, ":"):
		return true // scheme-less relative path
	default:
		return false
	}
}
