package db

// 导入依赖包
import (
	"fmt"
	"github.com/lgdzz/vingo-utils/vingo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

// Db 连接池句柄
var Pool *gorm.DB

type MysqlConfig struct {
	Host     string `yaml:"host" json:"host"`
	Port     string `yaml:"port" json:"port"`
	Dbname   int    `yaml:"dbname" json:"dbname"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}

func InitMysqlService(config *MysqlConfig) {
	charset := "utf8mb4"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Dbname,
		charset)
	//这里 gorm.Open()函数与之前版本的不一样，大家注意查看官方最新gorm版本的用法
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
			logger.Config{
				SlowThreshold:             time.Second, // 慢 SQL 阈值
				LogLevel:                  logger.Warn, // 日志级别
				IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  true,        // 禁用彩色打印
			},
		),
		NowFunc: func() time.Time {
			tmp := time.Now().Local().Format(vingo.DatetimeFormat)
			now, _ := time.ParseInLocation(vingo.DatetimeFormat, tmp, time.Local)
			return now
		},
	})
	if err != nil {
		panic("Error to Db connection, err: " + err.Error())
	}

	// 连接池配置
	sqlDB, _ := db.DB()
	// 最大空闲数
	sqlDB.SetMaxIdleConns(10)
	// 最大连接数
	sqlDB.SetMaxOpenConns(100)
	// 连接最大存活时长
	sqlDB.SetConnMaxLifetime(60 * time.Minute)

	// 注册统一异常插件
	RegisterAfterQuery(db)
	RegisterAfterCreate(db)
	RegisterAfterUpdate(db)
	RegisterAfterDelete(db)

	Pool = db
}

func RegisterAfterQuery(db *gorm.DB) {
	err := db.Callback().Query().After("gorm:query").Register("gormerror:after_query", func(db *gorm.DB) {
		if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
			panic(&vingo.DbException{Message: db.Error.Error()})
		}
	})
	if err != nil {
		panic(fmt.Sprintf("插件注册失败: %v", err.Error()))
	}
}

func RegisterAfterCreate(db *gorm.DB) {
	err := db.Callback().Create().After("gorm:create").Register("gormerror:after_create", func(db *gorm.DB) {
		if db.Error != nil {
			panic(&vingo.DbException{Message: db.Error.Error()})
		}
	})
	if err != nil {
		panic(fmt.Sprintf("插件注册失败: %v", err.Error()))
	}
}

func RegisterAfterUpdate(db *gorm.DB) {
	err := db.Callback().Update().After("gorm:update").Register("gormerror:after_update", func(db *gorm.DB) {
		if db.Error != nil {
			panic(&vingo.DbException{Message: db.Error.Error()})
		}
	})
	if err != nil {
		panic(fmt.Sprintf("插件注册失败: %v", err.Error()))
	}
}

func RegisterAfterDelete(db *gorm.DB) {
	err := db.Callback().Delete().After("gorm:delete").Register("gormerror:after_delete", func(db *gorm.DB) {
		if db.Error != nil {
			panic(&vingo.DbException{Message: db.Error.Error()})
		}
	})
	if err != nil {
		panic(fmt.Sprintf("插件注册失败: %v", err.Error()))
	}
}
