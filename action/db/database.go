package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aobt/sqlmapper"
	"github.com/fatih/structs"
	_ "github.com/go-sql-driver/mysql"
	"kube-metric-collector/common"
	"kube-metric-collector/model"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var databaseURL string
var databaseName string
var databaseId string
var databasePassword string

var Connection *sql.DB

func init() {
	datasource := getDataSource()
	var err error

	Connection, err = sql.Open("mysql", datasource)

	if err != nil {
		log.Fatal(err)
	}
	//if err := Connection.Ping(); err != nil {
	//	log.Fatal(err)
	//}
}

func getDataSource() string {

	profile := "local"
	if len(os.Getenv("PROFILE")) > 0 {
		profile = os.Getenv("PROFILE")
	}

	if profile == "prod" {
		databaseURL = os.Getenv("DATABASE_URL")
		databaseName = os.Getenv("DATABASE_NAME")
		databaseId = os.Getenv("DATABASE_TERRAMAN_ID")
		databasePassword = os.Getenv("DATABASE_TERRAMAN_PASSWORD")
	} else {
		databaseURL = common.ConfInfo["database.url"]
		databaseName = common.ConfInfo["database.name"]
		databaseId = common.ConfInfo["database.id"]
		databasePassword = common.ConfInfo["database.password"]
	}

	dataSource := fmt.Sprint(databaseId, ":", databasePassword, "@(", databaseURL, ")/", databaseName, "?parseTime=true")
	//fmt.Println("dataSource :: ", dataSource)
	return dataSource
}

func GetClusterList() []model.CpClusters {
	ctx := context.Background()
	tableName := "cp_clusters"

	rowArrAll, _ := QueryAll(ctx, nil, Connection, tableName)

	return rowArrAll
}

func QueryAll(ctx context.Context, tx *sql.Tx, db *sql.DB, tableName string) ([]model.CpClusters, error) {
	var row model.CpClusters
	fm, err := sqlmapper.NewFieldsMap(tableName, &row)
	if err != nil {
		return nil, err
	}

	objptrs, err := fm.SQLSelectAllRows(ctx, tx, db)
	if err != nil {
		return nil, err
	}

	var objs []model.CpClusters
	for i, olen := 0, len(objptrs); i < olen; i++ {
		objs = append(objs, *objptrs[i].(*model.CpClusters))
	}

	return objs, nil
}

func GetClusterNodeStatue() []model.NodeStateInfo {
	ctx := context.Background()
	tableName := "cp_metric_node_status"

	rowArrAll, _ := QueryAllNodeStatus(ctx, nil, Connection, tableName)

	return rowArrAll
}

func QueryAllNodeStatus(ctx context.Context, tx *sql.Tx, db *sql.DB, tableName string) ([]model.NodeStateInfo, error) {
	var row model.NodeStateInfo
	fm, err := sqlmapper.NewFieldsMap(tableName, &row)
	if err != nil {
		return nil, err
	}

	objptrs, err := fm.SQLSelectAllRows(ctx, tx, db)
	if err != nil {
		return nil, err
	}

	var objs []model.NodeStateInfo
	for i, olen := 0, len(objptrs); i < olen; i++ {
		objs = append(objs, *objptrs[i].(*model.NodeStateInfo))

	}

	return objs, nil
}

//func GetClusterInfo() ([]model.ClusterItems, model.ResourceCount) {
//	var count model.ResourceCount
//	ctx := context.Background()
//	tableName := "cp_metric_cluster_status"
//
//	rowArrAll, count, _ := getClusterItems(ctx, nil, Connection, tableName)
//
//	return rowArrAll, count
//}

func GetClusterInfo(req model.ReqData) ([]model.ClusterItems, model.ResourceCount) {
	var count model.ResourceCount
	ctx := context.Background()
	tableName := "cp_metric_cluster_status"

	//model.ReqData
	rowArrAll, count, _ := getClusterItems(ctx, nil, Connection, tableName, req)

	return rowArrAll, count
}

func getKeyList(objptrs []interface{}, req model.ReqData) []interface{} {
	var rst []interface{}

	for i, olen := 0, len(objptrs); i < olen; i++ {
		var row model.TableClusters
		row.ClusterId = objptrs[i].(*model.TableClusterStatus).ClusterId
		colTime := objptrs[i].(*model.TableClusterStatus).UpdateTime

		// 요청한 Cluster인 경우에만 진행
		for _, clusterId := range req.ClusterId {

			// 시간차 구하기
			currentTime := time.Now().Format(time.RFC3339)
			currentFormatTime, _ := time.Parse(time.RFC3339, currentTime)
			collTime, _ := time.Parse(time.RFC3339, colTime)
			minute := int32(currentFormatTime.Sub(collTime).Minutes())

			fmt.Println(row.ClusterId, ":: ", collTime, ":: ", minute, "(min)")

			if row.ClusterId == clusterId && minute < 30 {
				rst = append(rst, objptrs[i])
			}
		}
	}

	return rst
}

func getClusterItems(ctx context.Context, tx *sql.Tx, db *sql.DB, tableName string, req model.ReqData) ([]model.ClusterItems, model.ResourceCount, error) {
	var row model.TableClusterStatus
	var count model.ResourceCount
	fm, err := sqlmapper.NewFieldsMap(tableName, &row)
	if err != nil {
		return nil, count, err
	}

	objptrs, err := fm.SQLSelectAllRows(ctx, tx, db)
	if err != nil {
		return nil, count, err
	}

	clusterIdList := getKeyList(objptrs, req)

	//var objs []model.ClusterItems
	objs := make([]model.ClusterItems, len(clusterIdList))
	count.Cluster = len(clusterIdList)

	for i, olen := 0, len(clusterIdList); i < olen; i++ {

		var row model.TableClusters
		row.ClusterId = clusterIdList[i].(*model.TableClusterStatus).ClusterId
		//row.ClusterId = objptrs[i].(*model.TableClusterStatus).ClusterId

		fm, err := sqlmapper.NewFieldsMap("cp_clusters", &row)
		if err != nil {
			continue
		}
		objptr, err := fm.SQLSelectByPriKey(ctx, tx, db)
		if err != nil {
			return nil, count, err
		}

		objs[i].ClusterId = row.ClusterId
		objs[i].ClusterName = objptr.(*model.TableClusters).ClusterName
		objs[i].ClusterProviderType = objptr.(*model.TableClusters).ProviderType
		objs[i].Version = clusterIdList[i].(*model.TableClusterStatus).KubeletVersion

		pv := strings.Split(clusterIdList[i].(*model.TableClusterStatus).PvCnt, "/")
		objs[i].PvCount.Count, _ = strconv.Atoi(pv[0])
		objs[i].PvCount.All, _ = strconv.Atoi(pv[1])
		count.Pv += objs[i].PvCount.All

		pvc := strings.Split(clusterIdList[i].(*model.TableClusterStatus).PvcCnt, "/")
		objs[i].PvcCount.Count, _ = strconv.Atoi(pvc[0])
		objs[i].PvcCount.All, _ = strconv.Atoi(pvc[1])
		count.Pvc += objs[i].PvcCount.All

		pod := strings.Split(clusterIdList[i].(*model.TableClusterStatus).PodCnt, "/")
		objs[i].PodCount.Count, _ = strconv.Atoi(pod[0])
		objs[i].PodCount.All, _ = strconv.Atoi(pod[1])
		count.Pod += objs[i].PodCount.All

		namespace := strings.Split(clusterIdList[i].(*model.TableClusterStatus).NamespaceCnt, "/")
		objs[i].NamespaceCount.Count, _ = strconv.Atoi(namespace[0])
		objs[i].NamespaceCount.All, _ = strconv.Atoi(namespace[1])
		count.Namespace += objs[i].NamespaceCount.All

		node := strings.Split(clusterIdList[i].(*model.TableClusterStatus).NodeCnt, "/")
		objs[i].NodeCount.Count, _ = strconv.Atoi(node[0])
		objs[i].NodeCount.All, _ = strconv.Atoi(node[1])
		count.Node += objs[i].NodeCount.All

		//objs[i].Usage.Cpu = objptrs[i].(*model.TableClusterStatus).CpuRatio
		//objs[i].Usage.Memory = objptrs[i].(*model.TableClusterStatus).MemRatio

		// 소수점 자르기
		objs[i].Usage.Cpu, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", clusterIdList[i].(*model.TableClusterStatus).CpuRatio), 64)
		objs[i].Usage.Memory, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", clusterIdList[i].(*model.TableClusterStatus).MemRatio), 64)
		//nodeCount, Version 추가

		//objs = append(objs, *objptrs[i].(*model.ClusterItems))

	}

	return objs, count, nil
}

func GetTopCPU(targetId string) []model.TopNodeCPU {
	//db, err := sql.Open("mysql", "terraman:Paasta!2022@(15.164.195.107:31306)/cpdev?parseTime=true")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer db.Close()

	query := "SELECT b.cluster_id , a.name, b.node_name , b.cpu_ratio, b.cpu_raw, b.mem_ratio, b.mem_raw  " +
		"\nFROM cp_clusters a, cp_metric_node_status b " +
		"\nWHERE a.cluster_id = b.cluster_id" +
		"\n  AND a.cluster_id in (" + targetId + ")" +
		"\nORDER BY b.cpu_ratio desc" +
		"\nLIMIT 5"

	fmt.Println(query)

	query = strings.ReplaceAll(query, "'", "'")

	rows, err := Connection.Query(query)

	if err != nil {
		log.Fatal(err)
	}
	var result []model.TopNodeCPU
	var clusterId, name, nodeName string
	var cpuRatio, memRatio float64
	var cpuRaw, memRaw int64

	for rows.Next() {
		err := rows.Scan(&clusterId, &name, &nodeName, &cpuRatio, &cpuRaw, &memRatio, &memRaw)
		if err != nil {
			log.Fatal(err)
		}
		cpu := model.Cpu{
			Usage:   int(cpuRaw / 1000000),
			Percent: int(cpuRatio),
		}
		mem := model.Memory{
			Usage:   int(memRaw / 1000),
			Percent: int(memRatio),
		}
		info := model.TopNodeCPU{
			ClusterName: name,
			ClusterId:   clusterId,
			Name:        nodeName,
			Cpu:         cpu,
			Memory:      mem,
		}

		result = append(result, info)

	}

	return result
}

func GetTopMem(targetId string) []model.TopNodeMEM {
	//db, err := sql.Open("mysql", "terraman:Paasta!2022@(15.164.195.107:31306)/cpdev?parseTime=true")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer db.Close()

	rows, err := Connection.Query(
		"SELECT b.cluster_id , a.name, b.node_name , b.cpu_ratio, b.cpu_raw, b.mem_ratio, b.mem_raw  " +
			"\nFROM cp_clusters a, cp_metric_node_status b " +
			"\nWHERE a.cluster_id = b.cluster_id" +
			"\n  AND a.cluster_id in (" + targetId + ")" +
			"\nORDER BY b.mem_ratio desc" +
			"\nLIMIT 5")

	if err != nil {
		log.Fatal(err)
	}
	var result []model.TopNodeMEM
	var clusterId, name, nodeName string
	var cpuRatio, memRatio float64
	var cpuRaw, memRaw int64

	for rows.Next() {
		err := rows.Scan(&clusterId, &name, &nodeName, &cpuRatio, &cpuRaw, &memRatio, &memRaw)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(clusterId, " :: ", cpuRatio)
		cpu := model.Cpu{
			Usage:   int(cpuRaw / 1000000),
			Percent: int(cpuRatio),
		}
		mem := model.Memory{
			Usage:   int(memRaw / 1000),
			Percent: int(memRatio),
		}
		info := model.TopNodeMEM{
			ClusterName: name,
			ClusterId:   clusterId,
			Name:        nodeName,
			Cpu:         cpu,
			Memory:      mem,
		}

		result = append(result, info)

	}

	return result
}

func InsertData(data model.CpMetricNodeStatus) {

	ctx := context.Background()

	selData, _ := QueryByKey(ctx, nil, Connection, data.ClusterId, data.NodeName)
	if selData == nil {
		fmt.Println("No data in database(ClusterId = ", data.ClusterId, "NodeName =", data.NodeName, ")")
		return
	}
	log.Printf("\nClusterId = %s\n NodeName = %s\n ", data.ClusterId, data.NodeName)

	dataMap := structToMap(*selData)

	fmt.Println(dataMap)

	if dataMap != nil { // update
		err := Update(ctx, nil, &data)

		if err != nil {
			fmt.Println("fail to InsertNodeMetric()")
		}
	} else { // insert
		err := InsertNodeMetric(ctx, nil, Connection, data)

		if err != nil {
			fmt.Println("fail to InsertNodeMetric()")
		}
	}

}

// Update by primary key (field[0])
func Update(ctx context.Context, tx *sql.Tx, row *model.CpMetricNodeStatus) error {

	fm, err := sqlmapper.NewFieldsMap("cp_metric_node_status", row)
	if err != nil {
		return err
	}

	err = fm.SQLUpdateByPriKey(ctx, tx, Connection)
	if err != nil {
		return err
	}

	return nil
}

/*
 * struct 를 map 형태로 변환한다.
 */
func structToMap(data model.CpMetricNodeStatus) map[string]interface{} {
	res := structs.New(data)
	m := res.Map()

	return m
}

func QueryByKey(ctx context.Context, tx *sql.Tx, db *sql.DB, fieldKey string, fieldKey2 string) (*model.CpMetricNodeStatus, error) {
	var row model.CpMetricNodeStatus
	row.ClusterId = fieldKey
	row.NodeName = fieldKey2
	fm, err := sqlmapper.NewFieldsMap("cp_metric_node_status", &row)
	if err != nil {
		return nil, err
	}

	objptr, err := fm.SQLSelectByPriKey(ctx, tx, db)
	if err != nil {
		return nil, err
	}

	return objptr.(*model.CpMetricNodeStatus), nil
}

func InsertNodeMetric(ctx context.Context, tx *sql.Tx, db *sql.DB, rows model.CpMetricNodeStatus) error {
	fm, err := sqlmapper.NewFieldsMap("cp_metric_node_status", &rows)
	if err != nil {
		return err
	}
	err = fm.SQLInsert(ctx, tx, db)
	if err != nil {
		return err
	}
	return nil
}
