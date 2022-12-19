package common

import (
	"kube-metric-collector/model"
	"strings"
)

func MakeInQueryValue(req model.ReqData) string {
	targetId := ""
	for _, id := range req.ClusterId {
		targetId += "'" + id + "',"
	}
	targetId = strings.TrimRight(targetId, ",")

	return targetId
}
