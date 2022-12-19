package main

import (
	"kube-metric-collector/api"
	"kube-metric-collector/app"
	_ "kube-metric-collector/docs"
	"time"
)

// @title MetricCollector API
// @version 0.1.0
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @contact.email support@swagger.io
func main() {
	println("Let's start to metric collector")

	go api.ProcessREST()

	for {
		app.Process()
		time.Sleep(time.Second * 60)

	}

}
