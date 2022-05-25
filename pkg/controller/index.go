package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Home
// @Produce html
// @Tags home
// @Router /index [get]
func Index(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(
		`<html>
	<head>
		<title>Weave Server</title>
	</head>
	<body>
		<h1>Hello Weave</h1>
		<ul>
			<li><a href="/swagger/index.html">swagger</a></li>
			<li><a href="/metrics">metrics</a></li>
			<li><a href="/healthz">healthz</a></li>
	  	</ul>
		<hr>
		<center>Weave/1.0</center>
	</body>
<html>`))
}
