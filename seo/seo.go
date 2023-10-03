package seo

import (
	"io/ioutil"
	"regexp"
)

func GenerateSEOSupportHTML(source string, destination string) {
	// Read the contents of the index.html file
	content, err := ioutil.ReadFile(source)
	if err != nil {
		panic(err)
	}

	// Define regular expressions for matching the dynamic metadata tags
	urlRegex := regexp.MustCompile(`<meta\s+name="url"\s+content="(.*?)"`)
	typeRegex := regexp.MustCompile(`<meta\s+name="type"\s+content="(.*?)"`)
	titleRegex := regexp.MustCompile(`<title>(.*?)<\/title>`)
	descRegex := regexp.MustCompile(`<meta\s+name="description"\s+content="(.*?)"`)
	imageRegex := regexp.MustCompile(`<meta\s+name="image"\s+content="(.*?)"`)
	ogUrlRegex := regexp.MustCompile(`<meta\s+property="og:url"\s+content="(.*?)"`)
	ogTitleRegex := regexp.MustCompile(`<meta\s+property="og:title"\s+content="(.*?)"`)
	ogDescRegex := regexp.MustCompile(`<meta\s+property="og:description"\s+content="(.*?)"`)
	ogImageRegex := regexp.MustCompile(`<meta\s+property="og:image"\s+content="(.*?)"`)
	titleRegexMeta := regexp.MustCompile(`<meta\s+name="title"\s+content="(.*?)"`)

	// Replace the dynamic metadata using the regular expressions
	newContent := urlRegex.ReplaceAllString(string(content), `<meta name="url" content="{{ .Url }}"`)
	newContent = typeRegex.ReplaceAllString(newContent, `<meta name="type" content="{{ .Type }}"`)
	newContent = titleRegex.ReplaceAllString(newContent, `<title>{{ .Title }}</title>`)
	newContent = titleRegexMeta.ReplaceAllString(newContent, `<meta name="title" content="{{ .Title }}"`)
	newContent = descRegex.ReplaceAllString(newContent, `<meta name="description" content="{{ .Description }}"`)
	newContent = imageRegex.ReplaceAllString(newContent, `<meta name="image" content="{{ .Image }}"`)
	newContent = ogUrlRegex.ReplaceAllString(newContent, `<meta property="og:url" content="{{ .Url }}"`)
	newContent = ogTitleRegex.ReplaceAllString(newContent, `<meta property="og:title" content="{{ .Title }}"`)
	newContent = ogDescRegex.ReplaceAllString(newContent, `<meta property="og:description" content="{{ .Description }}"`)
	newContent = ogImageRegex.ReplaceAllString(newContent, `<meta property="og:image" content="{{ .Image }}"`)

	// Write the modified content back to the index.html file
	err = ioutil.WriteFile(destination, []byte(newContent), 0644)
	if err != nil {
		panic(err)
	}
}
