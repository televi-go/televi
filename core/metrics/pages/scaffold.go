package pages

import "text/template"

var ScaffoldTemplate = template.Must(template.New("scaffold").Parse(pageScaffold))

const pageScaffold = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{ .Title}}</title>
    <link rel="stylesheet" href="https://www.unpkg.com/televi_assets_x@latest/css/main.css">
</head>
<body>
<div class="heading">
    <div class="content-wrap"
         style="position:relative; height: 80px; display:flex; align-items: center; isolation: isolate">
        <div id="menu" class="hidden"></div>
        <img src="https://www.unpkg.com/televi_assets_x@latest/images/logo.png"
             style="display: block; position: absolute; height:100%; top:0; left:-90px"
             alt="">
        <h1 style="margin:auto 0; flex:1">
            {{.Title}}
        </h1>
        <div class="material-symbols-rounded clickable" id="menu-button">
            menu
        </div>

    </div>
</div>
<div class="content-wrap">
	<div>
    {{block "BODY" .}}
	{{end}}
	</div>
</div>
<script src="https://www.unpkg.com/televi_assets_x@latest/js/index.js"></script>
</body>
<script>
    window.PAGES_DATA = [
        {
            icon: "group",
            label: "Clients",
            link: "/users"
        }
    ]
    window.PAGES_DATA = [
        {
            icon: "account_circle",
            label: "Profile",
            link: "/profile"
        }
    ]
</script>
</html>
`
