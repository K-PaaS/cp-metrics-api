package db

//var Conn *sql.DB
//
//func init() {
//	datasource := getDataSource1()
//	var err error
//
//	Conn, err = sql.Open("mysql", datasource)
//
//	if err != nil {
//		log.Fatal(err)
//	}
//	if err := Connection.Ping(); err != nil {
//		log.Fatal(err)
//	}
//}
//
//func getDataSource1() string {
//	databaseAddr = common.ConfInfo["database.mariadb.address"]
//	databasePort = common.ConfInfo["database.mariadb.port"]
//	databaseName = common.ConfInfo["database.mariadb.name"]
//	databaseId = common.ConfInfo["database.mariadb.id"]
//	databasePassword = common.ConfInfo["database.mariadb.password"]
//
//	dataSource := fmt.Sprint(databaseId, ":", databasePassword, "@(", databaseAddr, ":", databasePort, ")/", databaseName, "?parseTime=true")
//	fmt.Println("dataSource :: ", dataSource)
//	return dataSource
//}
