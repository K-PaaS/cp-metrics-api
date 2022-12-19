package model

// NameSpace Resource
type NameSpace struct {
	Kind string `json:"kind"`
	Item []Item `json:"items"`
}

type Item struct {
	Metadata NSMeta   `json:"metadata"`
	Status   NSStatus `json:"status"`
}

type NSMeta struct {
	Name string `json:"name"`
	UID  string `json:"UID"`
}

type NSStatus struct {
	Phase         string            `json:"phase"`
	ConStatus     []ContainerStatus `json:"containerStatuses"`
	NodeCondition []NodeCondition   `json:"conditions"`
	NodeInfo      NodeInfo          `json:"nodeInfo"`
}

type NodeCondition struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}

type ContainerStatus struct {
	Name  string `json:"name"`
	Ready bool   `json:"ready"`
}
