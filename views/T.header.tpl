{{define "header"}}
<!DOCTYPE HTML>
<html>
    <head>
        <meta charset="utf-8">
        <meta http-equiv="X-UA-Compatible" content="IE=edge">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <link rel="shortcut icon" href="/static/img/favicon.png">
        <link href="/static/css/bootstrap.min.css" rel="stylesheet">
        {{if .IsSignin}}
        <link href="/static/css/sign.css" rel="stylesheet">
        {{else}}
        <link href="/static/css/jumbotron.css" rel="stylesheet">
        {{end}}
        <title>{{.Title}}</title>
    </head>
    <body>
{{end}}