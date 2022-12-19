package model

type CpClusters struct {
	ClusterId    string `sql:"cluster_id"`
	ClusterName  string `sql:"name"`
	ClusterType  string `sql:"cluster_type"`
	CreateTime   string `sql:"created"`
	LastModified string `sql:"last_modified"`
	ProviderType string `sql:"provider_type"`
	Description  string `sql:"description"`
	Status       string `sql:"status"`
}

type ClusterInfo struct {
	ClusterId     string
	ClusterApiUrl string
	ClusterToken  string
	StatusCode    int
}

type ReqData struct {
	ClusterId []string `json:"cluster_id"`
}

type ReqCluster struct {
	ClusterId string `json:"cluster_id"`
}

type ResponsePing struct {
	ClusterId  string `json:"cluster_id"`
	StatusCode int    `json:"status_code"`
}
