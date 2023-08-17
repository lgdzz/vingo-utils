package mysql

import (
	"flag"
	"os"
	"strings"
)

func InitCli() {
	model := flag.String("model", "", "生成数据库模型，支持多个表生成，格式：table1,table2")
	flag.StringVar(model, "m", "", "生成数据库模型，支持多个表生成，格式：table1,table2")

	dbbook := flag.Bool("dbbook", false, "生成数据库字典")
	flag.BoolVar(dbbook, "d", false, "生成数据库字典")

	help := flag.Bool("h", false, "Show help")

	// 解析命令行参数
	flag.Parse()

	if *help {
		// 如果使用 -h 或 --help 标志，则显示帮助信息
		flag.Usage()
		os.Exit(0)
	}

	// 创建数据库字典
	if *dbbook {
		_ = BuildBook()
		os.Exit(0)
	}

	// 创建数据表模型文件
	if *model != "" {
		_, _ = CreateDbModel(strings.Split(*model, ",")...)
		os.Exit(0)
	}
}
