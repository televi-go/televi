package pages

import (
	"github.com/televi-go/televi/core/metrics/grouping"
	"github.com/televi-go/televi/util"
	"html/template"
	"io"
	"time"
)

type UserWithAction struct {
	Action        string
	CommittedAt   time.Time
	FormattedTime string
	FirstName     string
	LastName      string
	UserName      string
	IsPremium     bool
}

func (userWithAction *UserWithAction) GetTime() time.Time {
	return userWithAction.CommittedAt
}

func (userWithAction *UserWithAction) InflateFormatted(formatted string) {
	userWithAction.FormattedTime = formatted
}

const actionBodySource = `
{{define "BODY"}}

{{if .HasMultipleActions}}

<div class="horizontal-selector">
	{{range .Actions}}
	<a class="variant" href="?action={{.Target}}" {{if eq .Target .SelectedAction}} data-selected {{end}}>
		<div class="v-content">
			{{.Label}}
			{{if gt .Count 0}}
				<div class="badge">
					{{.Count}}
				</div>
			{{end}}
		</div>
	</a>
	{{end}}
</div>

{{end}}

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
	{{range .Data}}
		<div class="grouping-entry">
			<div class="timepoint">
				{{ .FormattedTime}}
			</div>
			{{if .Action}}
				<div style="width:7rem">{{ .Action}}</div>
			{{end}}
			<div class="username">
				{{ .UserName}}
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

type ActionLink struct {
	Label          string
	Target         string
	SelectedAction string
	Count          int
}

type ActionsViewData struct {
	HasMultipleActions bool
	Actions            []ActionLink
	Title              string
	SelectedAction     string
	Groups             []grouping.DateGroup[*UserWithAction]
}

var actionsPageTemplate = template.Must(ScaffoldTemplate().Parse(actionBodySource))

var allActionsAct = []ActionLink{
	{Label: "All",
		Target: ""},
}

type ActionInfo struct {
	Action string
	Count  int
}

func mergeWithAll(source []ActionInfo, selectedAction string) []ActionLink {
	result := make([]ActionLink, len(source)+1)

	count := 0

	for i, info := range source {
		result[i+1] = ActionLink{
			Label:          info.Action,
			Target:         info.Action,
			SelectedAction: selectedAction,
			Count:          info.Count,
		}
		count += info.Count
	}

	result[0] = ActionLink{
		Label:          "All",
		Target:         "",
		SelectedAction: selectedAction,
		Count:          count,
	}

	return result
}

func ExecuteActionsPage(writer io.Writer, data []UserWithAction, selectedAction string, uniqueActions []ActionInfo) error {
	ptrArray := util.MakePointerArr(data)
	groups := grouping.GroupByDate(ptrArray)

	if selectedAction != "" {
		for i := 0; i < len(ptrArray); i++ {
			ptrArray[i].Action = ""
		}
	}

	viewData := ActionsViewData{
		SelectedAction:     selectedAction,
		HasMultipleActions: len(uniqueActions) > 0,
		Groups:             groups,
		Title:              "Actions",
		Actions:            mergeWithAll(uniqueActions, selectedAction),
	}
	return actionsPageTemplate.Execute(writer, viewData)
}
