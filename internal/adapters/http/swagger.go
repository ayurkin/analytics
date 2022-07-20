package http

import _ "analytics/api/swagger/public"

// Dependencies:
// go get -u -t github.com/swaggo/swag/cmd/swag

// @Title Analytics microservice
// @Version 1.0.0
// @host    localhost:3000
// @BasePath /analytics/v1
// @Schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Cookie
