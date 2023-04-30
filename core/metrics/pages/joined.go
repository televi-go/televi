package pages

import "text/template"

const joinedBodySource = `
{{define "BODY"}}

<h2>
Bot <a href="https://t.me/{{.Name}}">@{{.Name}}</a>
</h2>

{{range .Groups}}

<div class="grouping-header">
	{{.Title}}
</div>

<div class="grouping">
	{{range .Users}}
		<div class="grouping-entry">
			<div class="timepoint">
				{{ .FormattedAt}}
			</div>
			<div class="username">
				{{ .UiName}}
			</div>
			<div class="name">
				{{ .FirstName}} {{ .LastName}}
			</div>
		</div>
	{{end}}
</div>

{{end}}

{{end}}
`

var JoinedPageTemplate = template.Must(ScaffoldTemplate.Parse(joinedBodySource))

type JoinedAt struct {
	FirstName   string
	LastName    string
	FormattedAt string
	UiName      string
}

type Group struct {
	Title string
	Users []JoinedAt
}

type JoinedPageViewData struct {
	Title  string
	Name   string
	Groups []Group
}
