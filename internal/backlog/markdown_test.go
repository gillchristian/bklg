package backlog

import (
	"strings"
	"testing"
)

// AC1: inline markdown renders to whitelisted tags.
func TestRenderMarkdownInline(t *testing.T) {
	cases := map[string]string{
		"**bold**":            "<strong>bold</strong>",
		"*italic*":            "<em>italic</em>",
		"`code`":              "<code>code</code>",
		"[t](https://ex.com)": `<a href="https://ex.com" rel="noopener noreferrer">t</a>`,
		"[t](/rel/path)":      `<a href="/rel/path" rel="noopener noreferrer">t</a>`,
	}
	for in, want := range cases {
		if got := string(renderMarkdown(in)); !strings.Contains(got, want) {
			t.Errorf("renderMarkdown(%q) = %q, want to contain %q", in, got, want)
		}
	}
}

// AC2: escape-first — no repo HTML survives as live markup.
func TestRenderMarkdownEscapesFirst(t *testing.T) {
	got := string(renderMarkdown(`<script>alert(1)</script> & "q" 'x' <img onerror=y>`))
	if strings.Contains(got, "<script>") || strings.Contains(got, "<img") {
		t.Errorf("live tag survived: %q", got)
	}
	for _, want := range []string{"&lt;script&gt;", "&amp;", "&#34;", "&#39;"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing escaped %q in %q", want, got)
		}
	}
}

// AC2: unsafe link schemes never become live links.
func TestRenderMarkdownUnsafeLink(t *testing.T) {
	for _, in := range []string{
		"[click](javascript:alert(1))",
		"[x](data:text/html,<script>alert(1)</script>)",
		"[x](vbscript:msgbox)",
		"[x](JavaScript:alert(1))",
	} {
		got := string(renderMarkdown(in))
		if strings.Contains(got, "href=") {
			t.Errorf("unsafe link produced an href: renderMarkdown(%q) = %q", in, got)
		}
	}
	if got := string(renderMarkdown("[ok](https://example.com)")); !strings.Contains(got, `href="https://example.com"`) {
		t.Errorf("safe link dropped: %q", got)
	}
}

// AC2/AC1: code-span contents are literal, not re-processed as markdown.
func TestRenderMarkdownCodeProtected(t *testing.T) {
	got := string(renderMarkdown("`**not bold**`"))
	if !strings.Contains(got, "<code>**not bold**</code>") {
		t.Errorf("code content should be literal: %q", got)
	}
	if strings.Contains(got, "<strong>") {
		t.Errorf("bold rendered inside code span: %q", got)
	}
}

// AC4: an unmatched marker (e.g. left by truncation) stays literal — no unclosed tag.
func TestRenderMarkdownUnclosed(t *testing.T) {
	got := string(renderMarkdown("a long **bold that never closes…"))
	if strings.Contains(got, "<strong>") {
		t.Errorf("unclosed ** should stay literal, got %q", got)
	}
}

// Underscore emphasis is intentionally unsupported (don't mangle identifiers).
func TestRenderMarkdownNoUnderscoreEmphasis(t *testing.T) {
	for _, in := range []string{"some_var_name", "__init__", "a _b_ c"} {
		if got := string(renderMarkdown(in)); strings.Contains(got, "<em>") || strings.Contains(got, "<strong>") {
			t.Errorf("renderMarkdown(%q) = %q, underscores must not emphasize", in, got)
		}
	}
}

func TestSafeURL(t *testing.T) {
	for _, u := range []string{"https://x.com", "http://x", "mailto:a@b.c", "/abs", "#frag", "rel/path", "./x"} {
		if !safeURL(u) {
			t.Errorf("safeURL(%q) = false, want true", u)
		}
	}
	for _, u := range []string{"javascript:alert(1)", "data:text/html,x", "vbscript:x", "JAVASCRIPT:x", " javascript:x"} {
		if safeURL(u) {
			t.Errorf("safeURL(%q) = true, want false", u)
		}
	}
}
