package Model

import "database/sql"

type XgqNode struct {
	IpAddr int64
}

func (n *XgqNode) CreatDataBase() {
	//打开数据库，如果不存在，则创建
	db, err := sql.Open("sqlite3", "./XgqNode.db")

	defer db.Close()

	if err != nil {
		return
	}

	//创建表
	sql_table := `
    CREATE TABLE IF NOT EXISTS userinfo(
        uid INTEGER PRIMARY KEY AUTOINCREMENT,
        username VARCHAR(64) NULL,
        departname VARCHAR(64) NULL,
        created DATE NULL
    );
    `

	db.Exec(sql_table)
}
