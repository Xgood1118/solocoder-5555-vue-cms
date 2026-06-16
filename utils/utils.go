package utils

import (
	"crypto/rand"
	"encoding/base64"
	"html"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mozillazg/go-pinyin"
)

func GenerateCSRFToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func SanitizeHTML(s string) string {
	return html.EscapeString(s)
}

func Slugify(title string) string {
	if title == "" {
		return ""
	}

	args := pinyin.NewArgs()
	args.Style = pinyin.Normal
	pinyinResult := pinyin.Pinyin(title, args)

	var parts []string
	for _, p := range pinyinResult {
		if len(p) > 0 {
			parts = append(parts, p[0])
		}
	}

	slug := strings.Join(parts, "-")
	slug = strings.ToLower(slug)
	slug = strings.ReplaceAll(slug, " ", "-")

	var builder strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			builder.WriteRune(r)
		}
	}
	slug = builder.String()

	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	slug = strings.Trim(slug, "-")

	if slug == "" {
		slug = "article"
	}

	return slug
}

func GenerateUniqueSlug(title string, exists func(string) bool) string {
	baseSlug := Slugify(title)
	if !exists(baseSlug) {
		return baseSlug
	}

	counter := 2
	for {
		candidate := baseSlug + "-" + itoa(counter)
		if !exists(candidate) {
			return candidate
		}
		counter++
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var result []byte
	for n > 0 {
		result = append([]byte{byte('0' + n%10)}, result...)
		n /= 10
	}
	return string(result)
}

func TruncateText(text string, maxLen int) string {
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", " ")

	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}

	if utf8.RuneCountInString(text) <= maxLen {
		return text
	}

	runes := []rune(text)
	if len(runes) > maxLen {
		return string(runes[:maxLen]) + "..."
	}
	return text
}

func ExtractFirstParagraph(content string) string {
	content = strings.TrimSpace(content)

	if strings.HasPrefix(content, "<") {
		return extractFirstHTMLParagraph(content)
	}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "```") {
			return line
		}
	}
	return content
}

func extractFirstHTMLParagraph(html string) string {
	start := strings.Index(html, "<p>")
	if start == -1 {
		return stripHTMLTags(html)
	}
	start += 3
	end := strings.Index(html[start:], "</p>")
	if end == -1 {
		return stripHTMLTags(html)
	}
	return stripHTMLTags(html[start : start+end])
}

func stripHTMLTags(s string) string {
	var builder strings.Builder
	inTag := false
	for _, r := range s {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			builder.WriteRune(r)
		}
	}
	return strings.TrimSpace(builder.String())
}

func IsValidEmail(email string) bool {
	if email == "" {
		return false
	}
	atIndex := strings.Index(email, "@")
	if atIndex == -1 || atIndex == 0 {
		return false
	}
	dotIndex := strings.LastIndex(email, ".")
	if dotIndex == -1 || dotIndex < atIndex || dotIndex == len(email)-1 {
		return false
	}
	return true
}

func CountLinks(text string) int {
	count := 0
	lower := strings.ToLower(text)
	count += strings.Count(lower, "http://")
	count += strings.Count(lower, "https://")
	count += strings.Count(lower, "www.")
	return count
}

func ContainsSensitiveWords(text string, words []string) bool {
	lower := strings.ToLower(text)
	for _, word := range words {
		if strings.Contains(lower, strings.ToLower(word)) {
			return true
		}
	}
	return false
}

func HighlightText(text string, keyword string) string {
	if keyword == "" {
		return text
	}
	lowerText := strings.ToLower(text)
	lowerKeyword := strings.ToLower(keyword)

	var result strings.Builder
	offset := 0
	for {
		idx := strings.Index(lowerText[offset:], lowerKeyword)
		if idx == -1 {
			result.WriteString(text[offset:])
			break
		}
		result.WriteString(text[offset : offset+idx])
		result.WriteString("<mark>")
		result.WriteString(text[offset+idx : offset+idx+len(keyword)])
		result.WriteString("</mark>")
		offset += idx + len(keyword)
	}
	return result.String()
}

func ValidatePassword(password string) bool {
	if len(password) < 6 {
		return false
	}
	hasLetter := false
	hasNumber := false
	for _, r := range password {
		if unicode.IsLetter(r) {
			hasLetter = true
		}
		if unicode.IsDigit(r) {
			hasNumber = true
		}
	}
	return hasLetter && hasNumber
}
