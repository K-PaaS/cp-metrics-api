package model

type NodeMetricModel struct {
	Kind       string     `json:"kind"`
	ApiVersion string     `json:"apiVersion"`
	Items      []NodeItem `json:"items"`
}

type NodeItem struct {
	NodeInfo  MetricMetadata `json:"metadata"`
	Timestamp string         `json:"timestamp"`
	Window    string         `json:"window"`
	Usage     NodeUsage      `json:"usage"`
}

type MetricMetadata struct {
	Name              string `json:"name"`
	CreationTimestamp string `json:"creationTimestamp"`
}

type NodeUsage struct {
	Cpu    string `json:"cpu"`
	Memory string `json:"memory"`
}

//database
type CpMetricNodeStatus struct {
	ClusterId  string  `gorm:"column:cluster_id"`
	NodeName   string  `gorm:"column:node_name"`
	Cpu        float64 `gorm:"column:cpu_ratio"`
	CpuRaw     int64   `gorm:"column:cpu_raw"`
	Memory     float64 `gorm:"column:mem_ratio"`
	MemRaw     int64   `gorm:"column:mem_raw"`
	UpdateTime string  `gorm:"column:update_time"`
}

type CpMetricClusterStatus struct {
	ClusterId      string  `gorm:"column:cluster_id"`
	KubeletVersion string  `gorm:"column:kubelet_version"`
	NodeCnt        string  `gorm:"column:node_cnt"`
	NameSpaceCnt   string  `gorm:"column:namespace_cnt"`
	PvCnt          string  `gorm:"column:pv_cnt"`
	PvcCnt         string  `gorm:"column:pvc_cnt"`
	PodCnt         string  `gorm:"column:pod_cnt"`
	CpuRatio       float64 `gorm:"column:cpu_ratio"`
	MemRatio       float64 `gorm:"column:mem_ratio"`
	UpdateTime     string  `gorm:"column:update_time"`
}
