package sync

import (
	"fmt"
	"gitee.com/iotdrive/tools/ctable"
	"gitee.com/iotdrive/tools/version"
	"os"

	"gitee.com/iotdrive/tools/logs"
)

func init() {
	app := ctable.GetAppName()
	for _, s := range os.Args {
		if s == "-sync" {
			ctable.Save2file(makeRows())
			version.SaveVersion("1.0.0", "DL698 Server", "upload", "物联网")
			os.Exit(0)
		}
		if s == "-init" {
			rows := makeRows()
			err := ctable.Sync(rows, app)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		}
	}
	//启动后若配置不存在，插入默认配置信息
	_, err := ctable.ConfCount(app)
	if err != nil {
		rows := makeRows()
		err = ctable.Sync(rows, app)
		if err != nil {
			logs.Error("Query configure failed:", err)
			os.Exit(1)
		}
	}
}
func makeRows() (rows []*ctable.ConfTable) {
	rows = make([]*ctable.ConfTable, 0)
	rs := ctable.AddDefaultWebRows()
	if rs != nil {
		rows = append(rows, rs...)
	}
	rs = ctable.AddDefaultRows()
	if rs != nil {
		rows = append(rows, rs...)
	}

	srs := "INT8|UINT8|INT16|UINT16|INT32|UINT32|FLOAT32|FLOAT64|INT64|UINT64|OCTET_STR|DateTime_S|Scaler_Uint"
	rs = ctable.AddDefaultPointRows(srs)
	rows = append(rows, rs...)

	rs = addDL698Conf()
	if rs != nil {
		rows = append(rows, rs...)
	}
	return
}

func addDL698Conf() []*ctable.ConfTable {
	rows := make([]*ctable.ConfTable, 0)
	row := ctable.AddRow("监听端口", "DL698_PORT", "28080",
		"", "addServer", "", "")
	rows = append(rows, row)

	row = ctable.AddRow("DL698_saAddress", "SA地址", "12345678",
		"", "addServer", "", "")
	rows = append(rows, row)

	return rows
}
