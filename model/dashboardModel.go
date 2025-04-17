package model

type DashboardResult struct {
	ResultCode      string         `json:"resultCode"`
	ResultMessage   string         `json:"resultMessage"`
	HttpStatusCode  int            `json:"httpStatusCode"`
	DetailMessage   string         `json:"detailMessage"`
	ItemMetaData    string         `json:"itemMetaData"`
	ClusterStatus   int            `json:"clusterStatus"`
	NamespaceStatus int            `json:"namespaceStatus"`
	PvcStatus       int            `json:"pvcStatus"`
	PvStatus        int            `json:"pvStatus"`
	PodStatus       int            `json:"podStatus"`
	ClusterItems    []ClusterItems `json:"items"`
	TopNodeCPU      []TopNodeCPU   `json:"topNodeCPU"`
	TopNodeMEM      []TopNodeMEM   `json:"topNodeMem"`
}

type ClusterItems struct {
	ClusterId           string         `json:"clusterId"`
	ClusterName         string         `json:"clusterName"`
	ClusterType         string         `json:"clusterType"`
	ClusterProviderType string         `json:"clusterProviderType"`
	Version             string         `json:"version"`
	NodeCount           NodeCount      `json:"nodeCount"`
	NamespaceCount      NamespaceCount `json:"namespaceCount"`
	PodCount            PodCount       `json:"podCount"`
	PvCount             PvCount        `json:"pvCount"`
	PvcCount            PvcCount       `json:"pvcCount"`
	Usage               Usage          `json:"usage"`
}

type TopNodeCPU struct {
	ClusterName string `json:"clusterName"`
	ClusterId   string `json:"clusterId"`
	Name        string `json:"name"`
	Cpu         Cpu    `json:"cpu"`
	Memory      Memory `json:"memory"`
}

type TopNodeMEM struct {
	ClusterName string `json:"clusterName"`
	ClusterId   string `json:"clusterId"`
	Name        string `json:"name"`
	Cpu         Cpu    `json:"cpu"`
	Memory      Memory `json:"memory"`
}

type NodeCount struct {
	Count int `json:"count"`
	All   int `json:"all"`
}

type NamespaceCount struct {
	Count int `json:"count"`
	All   int `json:"all"`
}

type PodCount struct {
	Count int `json:"count"`
	All   int `json:"all"`
}

type PvCount struct {
	Count int `json:"count"`
	All   int `json:"all"`
}

type PvcCount struct {
	Count int `json:"count"`
	All   int `json:"all"`
}

type Usage struct {
	Cpu    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
}

type Cpu struct {
	Usage   int `json:"usage"`
	Percent int `json:"percent"`
}

type Memory struct {
	Usage   int `json:"usage"`
	Percent int `json:"percent"`
}

// table cp_metric_cluster_status
type TableClusterStatus struct {
	ClusterId      string  `sql:"cluster_id"`
	KubeletVersion string  `sql:"kubelet_version"`
	NodeCnt        string  `sql:"node_cnt"`
	NamespaceCnt   string  `sql:"namespace_cnt"`
	PvCnt          string  `sql:"pv_cnt"`
	PvcCnt         string  `sql:"pvc_cnt"`
	PodCnt         string  `sql:"pod_cnt"`
	CpuRatio       float64 `sql:"cpu_ratio"`
	MemRatio       float64 `sql:"mem_ratio"`
	UpdateTime     string  `sql:"update_time"`
}

type TableNodeStatus struct {
	ClusterId  string  `sql:"cluster_id"`
	NodeName   string  `sql:"node_name"`
	Cpu        float64 `sql:"cpu_ratio"`
	Memory     float64 `sql:"mem_ratio"`
	UpdateTime string  `sql:"update_time"`
}

type TableClusters struct {
	ClusterId    string `sql:"cluster_id"`
	ClusterName  string `sql:"name"`
	ClusterType  string `sql:"cluster_type"`
	ProviderType string `sql:"provider_type"`
}

type ResourceCount struct {
	Cluster   int
	Node      int
	Namespace int
	Pod       int
	Pv        int
	Pvc       int
}
