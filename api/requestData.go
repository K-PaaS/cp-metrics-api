package api

import (
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	"kube-metric-collector/action"
)
import "kube-metric-collector/model"

func RequestResty(clusterInfo model.ClusterInfo) (model.NodeMetricModel, error) {
	client := resty.New()

	token := fmt.Sprint("bearer ", clusterInfo.ClusterToken)
	url := fmt.Sprint(clusterInfo.ClusterApiUrl, "/apis/metrics.k8s.io/v1beta1/nodes")

	//fmt.Println(clusterInfo.ClusterId, ":: ", url, ":: ", token)
	fmt.Println(clusterInfo.ClusterId, ":: ", url)

	var rst model.NodeMetricModel

	_, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).R().
		SetHeader("Accept", "*/*").
		SetHeader("Authorization", token).
		SetResult(&rst).
		Get(url)

	if err != nil {
		fmt.Println(err)
	}

	return rst, err
}

func ResPing(ClusterId string) int {
	var clusterInfo model.ClusterInfo
	data, err := action.GetProviderSecretData(ClusterId)
	if err != nil {
		fmt.Println(err)
	}

	clusterInfo.ClusterId = ClusterId
	clusterInfo.ClusterApiUrl = fmt.Sprintf("%v", data["clusterApiUrl"])
	clusterInfo.ClusterToken = fmt.Sprintf("%v", data["clusterToken"])

	StatusCode, err := GetPing(clusterInfo)

	return StatusCode

}

func GetPing(clusterInfo model.ClusterInfo) (int, error) {
	client := resty.New()

	token := fmt.Sprint("bearer ", clusterInfo.ClusterToken)
	url := fmt.Sprint(clusterInfo.ClusterApiUrl)

	resp, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).R().
		SetHeader("Accept", "*/*").
		SetHeader("Authorization", token).
		EnableTrace().
		Get(url)
	fmt.Println(resp.StatusCode())
	ti := resp.Request.TraceInfo()
	fmt.Println(ti.ResponseTime)

	if err != nil {
		fmt.Println(err)
	}

	return resp.StatusCode(), err
}

func GetNodeData(clusterInfo model.ClusterInfo) model.NodeModel {
	client := resty.New()

	token := fmt.Sprint("bearer ", clusterInfo.ClusterToken)
	url := fmt.Sprint(clusterInfo.ClusterApiUrl, "/api/v1/nodes")

	var rst model.NodeModel
	_, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).R().
		SetHeader("Accept", "*/*").
		SetHeader("Authorization", token).
		SetResult(&rst).
		Get(url)

	if err != nil {
		fmt.Println(err)
	}

	return rst
}

func GetNodeInfo(clusterInfo model.ClusterInfo) (int, int, string) {
	client := resty.New()

	token := fmt.Sprint("bearer ", clusterInfo.ClusterToken)
	url := fmt.Sprint(clusterInfo.ClusterApiUrl, "/api/v1/nodes")

	var rst model.NodeModel
	_, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).R().
		SetHeader("Accept", "*/*").
		SetHeader("Authorization", token).
		SetResult(&rst).
		Get(url)

	if err != nil {
		fmt.Println(err)
	}

	tCnt := len(rst.Items)
	var rCnt int
	var kubeletVersion string
	for _, data := range rst.Items {
		//fmt.Println(i, " ", data.Status.Phase)
		kubeletVersion = data.Status.NodeInfo.KubeletVersion
		for _, nodeStatus := range data.Status.Conditions {
			if nodeStatus.Type == "Ready" {
				if nodeStatus.Status == "True" {
					rCnt++
				}
			}
		}
	}

	return tCnt, rCnt, kubeletVersion
}

func GetNameSpaceCnt(clusterInfo model.ClusterInfo) (int, int) {
	client := resty.New()

	token := fmt.Sprint("bearer ", clusterInfo.ClusterToken)
	url := fmt.Sprint(clusterInfo.ClusterApiUrl, "/api/v1/namespaces")

	var rst model.NameSpace
	_, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).R().
		SetHeader("Accept", "*/*").
		SetHeader("Authorization", token).
		SetResult(&rst).
		Get(url)

	if err != nil {
		fmt.Println(err)
	}

	tCnt := len(rst.Item)
	var rCnt int
	for _, data := range rst.Item {
		//fmt.Println(i, " ", data.Status.Phase)
		if data.Status.Phase == "Active" {
			rCnt++
		}
	}

	return tCnt, rCnt
}

func GetPodCnt(clusterInfo model.ClusterInfo) (int, int) {
	client := resty.New()

	token := fmt.Sprint("bearer ", clusterInfo.ClusterToken)
	url := fmt.Sprint(clusterInfo.ClusterApiUrl, "/api/v1/pods")

	var rst model.NameSpace
	_, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).R().
		SetHeader("Accept", "*/*").
		SetHeader("Authorization", token).
		SetResult(&rst).
		Get(url)

	if err != nil {
		fmt.Println(err)
	}

	tCnt := len(rst.Item)
	var rCnt int

	for _, data := range rst.Item {
		status := true
		for _, con := range data.Status.ConStatus {
			if !con.Ready {
				status = false
			}
		}
		if !status {
			rCnt++
		}
	}

	return tCnt, tCnt - rCnt
}

func GetPvCnt(clusterInfo model.ClusterInfo) (int, int) {
	client := resty.New()

	token := fmt.Sprint("bearer ", clusterInfo.ClusterToken)
	url := fmt.Sprint(clusterInfo.ClusterApiUrl, "/api/v1/persistentvolumes")

	var rst model.NameSpace
	_, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).R().
		SetHeader("Accept", "*/*").
		SetHeader("Authorization", token).
		SetResult(&rst).
		Get(url)

	if err != nil {
		fmt.Println(err)
	}

	tCnt := len(rst.Item)
	var rCnt int
	for _, data := range rst.Item {
		//fmt.Println(i, " ", data.Status.Phase)
		if data.Status.Phase == "Bound" {
			rCnt++
		}
	}

	return tCnt, rCnt
}

func GetPvcCnt(clusterInfo model.ClusterInfo) (int, int) {
	client := resty.New()

	token := fmt.Sprint("bearer ", clusterInfo.ClusterToken)
	url := fmt.Sprint(clusterInfo.ClusterApiUrl, "/api/v1/persistentvolumeclaims")

	var rst model.NameSpace
	_, err := client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).R().
		SetHeader("Accept", "*/*").
		SetHeader("Authorization", token).
		SetResult(&rst).
		Get(url)

	if err != nil {
		fmt.Println(err)
	}

	tCnt := len(rst.Item)
	var rCnt int
	for _, data := range rst.Item {
		//fmt.Println(i, " ", data.Status.Phase)
		if data.Status.Phase == "Bound" {
			rCnt++
		}
	}

	return tCnt, rCnt
}
