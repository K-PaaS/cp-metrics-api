package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	echoSwagger "github.com/swaggo/http-swagger"
	"github.com/unrolled/render"
	"kube-metric-collector/action/db"
	"kube-metric-collector/common"
	"kube-metric-collector/model"
	"net/http"
)

var rd *render.Render

func ProcessREST() {
	mux := NewHandler()
	err := http.ListenAndServe(":8900", mux)
	if err != nil {
		fmt.Println("Listen error :: ", err)
		return
	}
}

func NewHandler() http.Handler {
	//mux := pat.New()

	r := mux.NewRouter()

	r.HandleFunc("/v1/metrics/cluster/ping/{clusterId}", getClusterPing).Methods("GET")
	r.HandleFunc("/v1/metrics/cluster/node", use(getClusterStatusHandler, basicAuth)).Methods("GET")
	r.HandleFunc("/getNodeStatus", use(getClusterStatusHandler, basicAuth)).Methods("GET")
	r.HandleFunc("/v1/metrics/cluster/dashboard", use(getDashboardData, basicAuth)).Methods("POST")
	r.HandleFunc("/getDashboardData", use(getDashboardData, basicAuth)).Methods("POST")

	r.PathPrefix("/swagger").Handler(echoSwagger.WrapHandler).Methods("GET")

	//amw := authenticationMiddleware{tokenUsers: make(map[string]string)}
	//amw.Populate()

	//mux.Use(amw.Middleware)

	//queue에 넣는 API, 백단에서는 Queue를 계속 관찰한다.
	//mux.Get("/getNodeStatus", getClusterStatusHandler)
	//mux.Post("/getDashboardData", getDashboardData)
	//mux.HandleFunc("/v1/metrics/cluster/node", use(getClusterStatusHandler, basicAuth))
	//mux.HandleFunc("/getNodeStatus", use(getClusterStatusHandler, basicAuth))
	//mux.HandleFunc("/v1/metrics/cluster/dashboard", use(getDashboardData, basicAuth))
	//mux.HandleFunc("/getDashboardData", use(getDashboardData, basicAuth))
	//
	//mux.Get("/v1/metrics/cluster/ping/{clusterId}", getClusterPing)
	//
	//mux.PathPrefix("/swagger").Handler(echoSwagger.WrapHandler)

	//mux.PathPrefix("/swagger2").Handler(echoSwagger.WrapHandler)

	return r
}

// @Summary Create model.ResponsePing
// @Description It is function that returns the cluster state.
// @Accept json
// @Produce json
// @Param ResponsePing body model.ResponsePing true "Infomation of ClusterID"
// @Success 200 {object} model.ResponsePing
// @Failure 400 {object} model.ResponsePing
// @Router /v1/metrics/cluster/ping/{clusterId} [get]
func getClusterPing(wr http.ResponseWriter, r *http.Request) {

	var info model.ResponsePing

	vars := mux.Vars(r)
	clusterId := vars["clusterId"]

	//var req model.ReqCluster
	//err := json.NewDecoder(r.Body).Decode(&req)
	//if err != nil {
	//	fmt.Println(err)
	//	wr.WriteHeader(http.StatusBadRequest)
	//	_, err := fmt.Fprint(wr, err)
	//	if err != nil {
	//		return
	//	}
	//}

	info.ClusterId = clusterId
	info.StatusCode = ResPing(clusterId)

	wr.Header().Add("Content-type", "application/json")
	wr.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(info)
	fmt.Fprint(wr, string(data))
}

// @Summary Create model.NodeStateInfo
// @Description It is created according to the IaaS information requested when creating the cluster.
// @Accept json
// @Produce json
// @Param iaasBody body model.NodeStateInfo true "Node State Info"
// @Success 200 {object} model.NodeStateInfo
// @Failure 400 {object} model.NodeStateInfo
// @Router /v1/metrics/cluster/node [get]
func getClusterStatusHandler(wr http.ResponseWriter, r *http.Request) {

	var status []model.NodeStateInfo

	status = db.GetClusterNodeStatue()

	//rd.JSON(wr, http.StatusOK, status)

	wr.Header().Add("Content-type", "application/json")
	wr.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(status)
	fmt.Fprint(wr, string(data))
}

// @Summary Create model.DashboardResult
// @Description It is created according to the IaaS information requested when creating the cluster.
// @Accept json
// @Produce json
// @Success 200 {object} model.DashboardResult
// @Failure 400 {object} model.DashboardResult
// @Router /v1/metrics/cluster/dashboard [get]
func getDashboardData(wr http.ResponseWriter, r *http.Request) {

	var req model.ReqData
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println(err)
		wr.WriteHeader(http.StatusBadRequest)
		_, err := fmt.Fprint(wr, err)
		if err != nil {
			return
		}
	}

	//json.NewDecoder(r.Body).

	for i, clusterId := range req.ClusterId {
		fmt.Println("request Cluster(", i+1, ") :: ", clusterId)
	}

	var items []model.ClusterItems
	var count model.ResourceCount

	result := model.DashboardResult{
		ResultCode:      "SUCCESS",
		ResultMessage:   "Processed successfully.",
		HttpStatusCode:  200,
		DetailMessage:   "Processed successfully.",
		ItemMetaData:    "",
		ClusterStatus:   len(req.ClusterId),
		NamespaceStatus: 0,
		PvcStatus:       0,
		PvStatus:        0,
		PodStatus:       0,
		ClusterItems:    nil,
		TopNodeCPU:      nil,
		TopNodeMEM:      nil,
	}

	targetId := common.MakeInQueryValue(req)
	items, count = db.GetClusterInfo(req)
	cpu := db.GetTopCPU(targetId)
	mem := db.GetTopMem(targetId)

	result.ClusterItems = items
	result.TopNodeCPU = cpu
	result.TopNodeMEM = mem
	//result.ClusterStatus = count.Cluster
	result.NamespaceStatus = count.Namespace
	result.PodStatus = count.Pod
	result.PvStatus = count.Pv
	result.PvcStatus = count.Pvc

	//err = rd.JSON(wr, http.StatusOK, result)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	wr.Header().Add("Content-type", "application/json")
	wr.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(result)
	fmt.Fprint(wr, string(data))
}

// @Summary Create model.DashboardResult
// @Description It is created according to the IaaS information requested when creating the cluster.
// @Accept json
// @Produce json
// @Success 200 {object} model.DashboardResult
// @Failure 400 {object} model.DashboardResult
// @Router /getDashboardData [get]
//func getDashboardData(wr http.ResponseWriter, r *http.Request) {
//
//	var items []model.ClusterItems
//	var count model.ResourceCount
//
//	result := model.DashboardResult{
//		ResultCode:      "SUCCESS",
//		ResultMessage:   "Processed successfully.",
//		HttpStatusCode:  200,
//		DetailMessage:   "Processed successfully.",
//		ItemMetaData:    "",
//		ClusterStatus:   0,
//		NamespaceStatus: 0,
//		PvcStatus:       0,
//		PvStatus:        0,
//		PodStatus:       0,
//		ClusterItems:    nil,
//		TopNodeCPU:      nil,
//		TopNodeMEM:      nil,
//	}
//
//	items, count = db.GetClusterInfo()
//	cpu := db.GetTopCPU()
//	mem := db.GetTopMem()
//
//	result.ClusterItems = items
//	result.TopNodeCPU = cpu
//	result.TopNodeMEM = mem
//	result.ClusterStatus = count.Cluster
//	result.NamespaceStatus = count.Namespace
//	result.PodStatus = count.Pod
//	result.PvStatus = count.Pv
//	result.PvcStatus = count.Pvc
//
//	//rd.JSON(wr, http.StatusOK, status)
//
//	wr.Header().Add("Content-type", "application/json")
//	wr.WriteHeader(http.StatusOK)
//	data, _ := json.Marshal(result)
//	fmt.Fprint(wr, string(data))
//}
