package pages

import "html/template"

const joinedBodySource = `
{{define "BODY"}}


{{range .Groups}}

<div class="grouping-header">
	<div class="content">
		{{.Title}}
	</div>
	<div class="divider"></div>
	<div class="count">
		{{.Count}}
	</div>
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

var JoinedPageTemplate = template.Must(ScaffoldTemplate().Parse(joinedBodySource))

type JoinedAt struct {
	FirstName   string
	LastName    string
	FormattedAt string
	UiName      string
}

type Group struct {
	Title string
	Count int
	Users []JoinedAt
}

func MakeGroup(title string, users []JoinedAt) Group {
	return Group{
		Title: title,
		Count: len(users),
		Users: users,
	}
}

type JoinedPageViewData struct {
	Title  string
	Name   string
	Groups []Group
}
