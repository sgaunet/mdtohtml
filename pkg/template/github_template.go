package template

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed github-markdown.css
var githubCSS string

// GitHubTemplate implements HTMLTemplate with GitHub-style CSS.
type GitHubTemplate struct {
	css string
}

// NewGitHubTemplate creates a new GitHub-style HTML template.
func NewGitHubTemplate() *GitHubTemplate {
	return &GitHubTemplate{
		css: githubCSS,
	}
}

// NewGitHubTemplateWithCSS creates a new GitHub-style HTML template with custom CSS.
func NewGitHubTemplateWithCSS(css string) *GitHubTemplate {
	return &GitHubTemplate{
		css: css,
	}
}

// Wrap wraps HTML content with a complete HTML document structure.
func (t *GitHubTemplate) Wrap(content, title string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>%s</title>
</head>
<body>
%s
</body>
</html>`, title, content)
}

// InjectCSS injects CSS into an HTML document.
func (t *GitHubTemplate) InjectCSS(html, css string) string {
	if css == "" {
		css = t.css
	}
	cssInjection := fmt.Sprintf("<style>\n%s\n</style>\n</head>", css)
	return strings.Replace(html, "</head>", cssInjection, 1)
}