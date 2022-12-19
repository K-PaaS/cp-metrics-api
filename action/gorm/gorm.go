package gorm

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"kube-metric-collector/common"
	"kube-metric-collector/model"
	"os"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

type CpClusters struct {
	gorm.Model
	ClusterId    string `sql:"cluster_id"`
	ClusterName  string `sql:"name"`
	ClusterType  string `sql:"cluster_type"`
	ProviderType string `sql:"provider_type"`
	Description  string `sql:"description"`
	Status       string `sql:"status"`
}

var databaseURL string
var databaseName string
var databaseId string
var databasePassword string

var db *gorm.DB

func init() {
	db = connect()
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

func connect() *gorm.DB {
	dsn := getDataSource()
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		panic("Db 연결에 실패하였습니다.")
	}

	return db
}

func GetClusterList() []model.CpClusters {
	var clusterList []model.CpClusters
	//var aa CpClusters

	//db.Unscoped().Take(&clusterList)
	//fmt.Println(clusterList)
	result := db.Unscoped().Find(&clusterList)

	fmt.Println(result)

	//cnt := result.RowsAffected

	return clusterList

}

func UpdateNode(data model.CpMetricNodeStatus) {
	var nodeInfo model.CpMetricNodeStatus

	result := db.First(&nodeInfo, "cluster_id = ? AND node_name = ?", data.ClusterId, data.NodeName)
	//db.Where("cluster_id = ? AND node_name = ?", data.ClusterId, data.NodeName).Find(&nodeInfo)
	//db.Where(map[string]interface{}{"cluster_id": data.ClusterId, "node_name": data.NodeName}).Find(&nodeInfo)
	if result.RowsAffected == 1 {
		db.Model(&nodeInfo).
			Where("cluster_id = ? AND node_name = ? ", data.ClusterId, data.NodeName).
			Updates(map[string]interface{}{
				"cpu_ratio":   data.Cpu,
				"cpu_raw":     data.CpuRaw,
				"mem_ratio":   data.Memory,
				"mem_raw":     data.MemRaw,
				"update_time": data.UpdateTime})
	} else {
		//db.Select("cluster_id", "node_name", "cpu_ratio", "mem_ratio", "update_time").Create(&nodeInfo)
		insertInfo := model.CpMetricNodeStatus{
			ClusterId:  data.ClusterId,
			NodeName:   data.NodeName,
			Cpu:        data.Cpu,
			CpuRaw:     data.CpuRaw,
			Memory:     data.Memory,
			MemRaw:     data.MemRaw,
			UpdateTime: data.UpdateTime,
		}
		aa := db.Create(&insertInfo)
		fmt.Println(aa.RowsAffected)
	}

}

func UpdateCluster(data model.CpMetricClusterStatus) {
	var clusterInfo model.CpMetricClusterStatus

	result := db.First(&clusterInfo, "cluster_id = ?", data.ClusterId)

	if result.RowsAffected == 1 {
		db.Model(&clusterInfo).
			Where("cluster_id = ?", data.ClusterId).
			Updates(map[string]interface{}{
				"namespace_cnt":   data.NameSpaceCnt,
				"kubelet_version": data.KubeletVersion,
				"node_cnt":        data.NodeCnt,
				"pv_cnt":          data.PvCnt,
				"pvc_cnt":         data.PvcCnt,
				"pod_cnt":         data.PodCnt,
				"cpu_ratio":       data.CpuRatio,
				"mem_ratio":       data.MemRatio,
				"update_time":     data.UpdateTime,
			})
	} else {
		insertInfo := model.CpMetricClusterStatus{
			ClusterId:      data.ClusterId,
			NodeCnt:        data.NodeCnt,
			KubeletVersion: data.KubeletVersion,
			NameSpaceCnt:   data.NameSpaceCnt,
			PvCnt:          data.PvCnt,
			PvcCnt:         data.PvcCnt,
			PodCnt:         data.PodCnt,
			CpuRatio:       data.CpuRatio,
			MemRatio:       data.MemRatio,
			UpdateTime:     data.UpdateTime,
		}
		aa := db.Create(&insertInfo)
		fmt.Println(aa.RowsAffected)
	}
}

func Process() {

	//db = connect()

	var product Product
	//conn.First(&product, 1)

	// 테이블 자동 생성
	db.AutoMigrate(&Product{})

	// 생성
	db.Create(&Product{Code: "D42", Price: 100})

	// 읽기

	db.First(&product, 1)                 // primary key기준으로 product 찾기
	db.First(&product, "code = ?", "D42") // code가 D42인 product 찾기

	// 수정 - product의 price를 200으로
	db.Model(&product).Update("Price", 200)
	// 수정 - 여러개의 필드를 수정하기
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"})
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// 삭제 - product 삭제하기
	db.Delete(&product, 1)
}
