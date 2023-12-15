package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gen/common"
	"gen/db"
	"gen/sql"

	"github.com/gotomicro/ego/core/elog"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
	"gorm.io/gen"
	"gorm.io/gorm"
	_ "gorm.io/plugin/soft_delete"
)

var dataMap = map[string]func(gorm.ColumnType) (dataType string){
	"int":    func(columnType gorm.ColumnType) (dataType string) { return "int32" },
	"bigint": func(columnType gorm.ColumnType) (dataType string) { return "int64" },
	"json":   func(columnType gorm.ColumnType) string { return "json.RawMessage" },
	// bool mapping
	"tinyint": func(columnType gorm.ColumnType) (dataType string) {
		ct, _ := columnType.ColumnType()
		if strings.HasPrefix(ct, "tinyint(1)") {
			return "bool"
		}
		return "int"
	},
}

const (
	OUT_PATH         = "OUT_PATH"
	MODEL_PKG_PATH   = "MODEL_PKG_PATH"
	MODEL_TABLE_NAME = "MODEL_TABLE_NAME" // 模型表名
)

func main() {
	rootCMD := &cobra.Command{
		Use:   "[project-name]gorm-gen",
		Short: "A generator for Cobra based Applications",
		Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	}
	rootCMD.AddCommand(db.StartCmd)
	rootCMD.AddCommand(sql.StartCmd)
	// rootCMD.AddCommand(migrate.MigrateCmd)
	nowTime := time.Now()
	if err := rootCMD.Execute(); err != nil {
		elog.Errorf("启动失败，耗时%d秒", time.Since(nowTime)/time.Second)
		panic(err)
	}

	if common.DB == nil {
		panic("gorm db 不存在，请根据对应命令执行")
	}
	// 读取环境变量配置
	// 采用懒加载
	// if err := godotenv.Load(); err != nil {
	// 	panic(".env文件环境变量未配置")
	// }

	genConfig := gen.Config{
		WithUnitTest: true,
		// generate model global configuration
		//FieldNullable: false, // generate pointer when field is nullable
		FieldSignable: true, // detect integer field's unsigned type, adjust generated data type
		//FieldCoverable:    true, // generate pointer when field has default value
		FieldWithIndexTag: true, // generate with gorm index tag
		FieldWithTypeTag:  true, // generate with gorm column type tag
	}

	genConfig.OutPath = "./dal/query"
	if os.Getenv(OUT_PATH) != "" {
		genConfig.OutPath = os.Getenv(OUT_PATH)
	}

	genConfig.ModelPkgPath = "./dal/model"
	if os.Getenv(MODEL_PKG_PATH) != "" {
		genConfig.ModelPkgPath = os.Getenv(MODEL_PKG_PATH)
	}

	g := gen.NewGenerator(genConfig)

	g.UseDB(common.DB)

	// specify diy mapping relationship
	g.WithDataTypeMap(dataMap)

	// generate all field with json tag end with "_example"
	// g.WithJSONTagNameStrategy(func(c string) string { return c + "_example" })
	//g.WithJSONTagNameStrategy(func(c string) string { return "-" })

	/*
		holidaySettingsTable := g.GenerateModel("holiday_settings",
			gen.FieldType("deleted_at", "soft_delete.DeletedAt"),
		)

		holidaySettingsExtTable := g.GenerateModel("holiday_settings_dayoffwork") // gen.FieldType("deleted_at", "soft_delete.DeletedAt"),

		holidayAffairsSettingsTable := g.GenerateModel("holiday_affairs_settings",
			gen.FieldType("deleted_at", "soft_delete.DeletedAt"),
		)

		holidayAffairsTable := g.GenerateModel("holiday_affairs") //gen.FieldType("deleted_at", "soft_delete.DeletedAt"),

		g.ApplyBasic(
			holidaySettingsTable,
			holidaySettingsExtTable,
			holidayAffairsSettingsTable,
			holidayAffairsTable,
		)
		// g.ApplyBasic(g.GenerateAllTable()...) // generate all table in db server

	*/

	// holidayDistributedTable := g.GenerateModel("holiday_distributed_record")

	// g.ApplyBasic(
	// 	holidayDistributedTable,
	// )

	if os.Getenv(MODEL_TABLE_NAME) == "" {
		fmt.Println("请在.env配置中配上你需要生成model的表名")
	}
	tableNames := os.Getenv(MODEL_TABLE_NAME)
	tableNameList := strings.Split(tableNames, ";")
	for _, name := range tableNameList {
		genmodel := g.GenerateModel(name)
		g.ApplyBasic(
			genmodel,
		)
	}
	g.Execute()
}
