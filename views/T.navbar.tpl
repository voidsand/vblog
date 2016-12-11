{{define "navbar"}}
<div class="navbar navbar-default navbar-fixed-top navbar-inverse" role="navigation">
    <div class="container">
        <div class="navbar-header">
            <button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar" aria-expanded="false" aria-controls="navbar">
                <span class="sr-only">Toggle navigation</span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
            </button>
            <a class="navbar-brand" href="/">我的博客</a>
        </div>
        <div id="navbar" class="navbar-collapse collapse">
            <ul class="nav navbar-nav">
                <li {{if .IsHome}}class="active"{{end}}><a href="/">首页</a></li>
                <li {{if .IsCategory}}class="active"{{end}}><a href="/category">分类</a></li>
                <li {{if .IsTopic}}class="active"{{end}}><a href="/topic">文章</a></li>
                <li {{if .IsAbout}}class="active"{{end}}><a href="/about">关于</a></li>
	        </ul>
	        <div class="pull-right">
                <ul class="nav navbar-nav">
                    {{if .LoginReady}}
                    <li><a href="/sign/out">退出</a></li>
                    {{else}}
                    <li><a href="/sign">登录</a></li>
                    {{end}}
                </ul>
            </div>
	    </div>
	</div>
</div>
{{end}}