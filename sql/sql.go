package sql

import (
	"fmt"
	"gen/common"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"gorm.io/rawsql"
)

const (
	SQL_FILE_PATH = "SQL_FILE_PATH"
)

var StartCmd = &cobra.Command{
	Use:          "sql",
	Short:        "Run the gen form sql",
	Example:      "cmd sql",
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, _ []string) error {
		var err error
		sqlPaths := os.Getenv(SQL_FILE_PATH)
		if sqlPaths == "" {
			return fmt.Errorf("请在.env配置中配置好sql文件路径SQL_FILE_PATH，使用英文分号\";\"间隔")
		}

		filePaths := make([]string, 0)
		sps := strings.Split(sqlPaths, ";")
		for _, sqlPath := range sps {
			if sqlPath != "" {
				filePaths = append(filePaths, sqlPath)
			}
		}

		common.DB, err = gorm.Open(rawsql.New(rawsql.Config{
			DriverName: "mysql",
			FilePath:   filePaths,
		}))
		if err != nil {
			return err
		}
		return nil
	},
}
