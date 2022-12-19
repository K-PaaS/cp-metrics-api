package app

import (
	"fmt"
	"kube-metric-collector/action"
	"kube-metric-collector/action/db"
	"kube-metric-collector/action/gorm"
	"kube-metric-collector/api"
	"kube-metric-collector/model"
	"strconv"
	"time"
)

func getClusterList() []model.CpClusters {

	var clusterList []model.CpClusters

	clusterList = db.GetClusterList()
	for _, info := range clusterList {
		fmt.Println(info.ClusterId)
	}

	return clusterList
}

func Process() {
	var clusterList []model.CpClusters
	var clusterInfo model.ClusterInfo

	clusterList = gorm.GetClusterList()

	//clusterList = getClusterList()

	for _, info := range clusterList {
		data, err := action.GetProviderSecretData(info.ClusterId)
		if err != nil {
			continue
		}

		var mo model.CpMetricClusterStatus
		now := time.Now()
		custom := now.Format(time.RFC3339)

		clusterInfo.ClusterId = info.ClusterId
		clusterInfo.ClusterApiUrl = fmt.Sprintf("%v", data["clusterApiUrl"])
		clusterInfo.ClusterToken = fmt.Sprintf("%v", data["clusterToken"])

		cpu, mem, err := getClusterMetric(clusterInfo)
		if err != nil || cpu == 0.0 || mem == 0.0 {
			continue
		}
		mo.CpuRatio = cpu
		mo.MemRatio = mem
		mo.ClusterId = info.ClusterId
		mo.UpdateTime = custom
		fmt.Println("CPU :: ", cpu, "MEM :: ", mem)
		getClusterResource(clusterInfo, &mo)
		gorm.UpdateCluster(mo)

		time.Sleep(time.Millisecond * 300)
	}

}

func getClusterResource(clusterInfo model.ClusterInfo, result *model.CpMetricClusterStatus) {
	total, active, kubeletVersion := api.GetNodeInfo(clusterInfo)
	value := fmt.Sprint(active, "/", total)
	result.NodeCnt = value
	result.KubeletVersion = kubeletVersion

	fmt.Println(clusterInfo.ClusterId, ":: nodeCnt = ", active, "/", total)

	total, active = api.GetNameSpaceCnt(clusterInfo)
	//var bb model.NameSpace
	//bb := api.GeneralGet(clusterInfo)
	//api.GeneralGet(clusterInfo)
	value = fmt.Sprint(active, "/", total)
	result.NameSpaceCnt = value

	fmt.Println(clusterInfo.ClusterId, ":: nsCnt = ", active, "/", total)

	total, active = api.GetPodCnt(clusterInfo)
	value = fmt.Sprint(active, "/", total)
	result.PodCnt = value

	fmt.Println(clusterInfo.ClusterId, ":: podsCnt = ", active, "/", total)

	total, active = api.GetPvCnt(clusterInfo)
	value = fmt.Sprint(active, "/", total)
	result.PvCnt = value

	fmt.Println(clusterInfo.ClusterId, ":: PVCnt = ", active, "/", total)

	total, active = api.GetPvcCnt(clusterInfo)
	value = fmt.Sprint(active, "/", total)
	result.PvcCnt = value

	fmt.Println(clusterInfo.ClusterId, ":: PVCCnt = ", active, "/", total)

}

func getClusterMetric(clusterInfo model.ClusterInfo) (float64, float64, error) {

	var CpuN float64 = 1000000000
	rst, err := api.RequestResty(clusterInfo)
	if err != nil || len(rst.Items) == 0 {
		return 0.0, 0.0, err
	}
	nodeData := api.GetNodeData(clusterInfo)

	var ClCpuRatio, ClMemRatio, ClCpuRaw, ClMemRaw, ClCpuCore, ClMemSize float64

	fmt.Println("kind: ", rst.Kind)
	fmt.Println("----------------------------------")
	for i, _ := range rst.Items {

		for j, _ := range nodeData.Items {

			if rst.Items[i].NodeInfo.Name == nodeData.Items[j].MetaData.NodeName {
				//fmt.Println("nodeName: ", rst.Items[i].NodeInfo.Name)
				//fmt.Println("cpu: ", rst.Items[i].Usage.Cpu)
				//fmt.Println("memory: ", rst.Items[i].Usage.Memory)
				//fmt.Println("----------------------------------")

				cpu := rst.Items[i].Usage.Cpu
				cpu = cpu[:len(cpu)-1]
				cpuRaw, _ := strconv.ParseFloat(cpu, 64)
				ClCpuRaw += cpuRaw

				memory := rst.Items[i].Usage.Memory
				memory = memory[:len(memory)-2]
				memRaw, _ := strconv.ParseFloat(memory, 64)
				ClMemRaw += memRaw

				cpuCore, _ := strconv.ParseFloat(nodeData.Items[j].Status.Capacity.Cpu, 64)
				memory = nodeData.Items[j].Status.Capacity.Memory
				memory = memory[:len(memory)-2]
				memSize, _ := strconv.ParseFloat(memory, 64)
				ClCpuCore += cpuCore
				ClMemSize += memSize

				cpuRatio := (cpuRaw / (cpuCore * CpuN)) * 100
				//fmt.Println("cpuRatio :: ", cpuRatio)

				memRatio := (memRaw / memSize) * 100
				//fmt.Println("memRatio :: ", memRatio)

				now := time.Now()
				custom := now.Format("2006-01-02 15:04:05")

				data := model.CpMetricNodeStatus{
					ClusterId:  clusterInfo.ClusterId,
					NodeName:   rst.Items[i].NodeInfo.Name,
					Cpu:        cpuRatio,
					CpuRaw:     int64(cpuRaw),
					Memory:     memRatio,
					MemRaw:     int64(memRaw),
					UpdateTime: custom,
				}
				//db.InsertData(data)
				gorm.UpdateNode(data)
			}
		}

	}
	ClCpuRatio = (ClCpuRaw / (ClCpuCore * CpuN)) * 100
	ClMemRatio = (ClMemRaw / ClMemSize) * 100

	return ClCpuRatio, ClMemRatio, err
}
