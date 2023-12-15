package db

import (
	"fmt"
	"gen/common"
	"os"
	"runtime/debug"

	"github.com/gotomicro/ego/core/elog"
	"github.com/spf13/cobra" // for cobra.Command
)

const (
	MYSQL_DSN     = "MYSQL_DSN"
	SQL_PATH      = "SQL_PATH"
	SQL_FILE_PATH = "SQL_FILE_PATH"
)

var StartCmd = &cobra.Command{
	Use:          "db",
	Short:        "Run the gen form db database",
	Example:      "cmd db",
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, _ []string) error {
		defer func() {
			if err := recover(); err != nil {
				elog.Errorf("Recover error : %v;[stack]:%v", err, debug.Stack())
			}
		}()
		if os.Getenv(MYSQL_DSN) == "" || os.Getenv(SQL_PATH) == "" {
			return fmt.Errorf("请在.env中配置MYSQL_DSN数据库连接信息和SQL_PATH执行sql的文件路径")
		}

		common.DB = ConnectDB(os.Getenv(MYSQL_DSN))
		if err := ImportSQL(common.DB, os.Getenv(SQL_PATH)); err != nil {
			return fmt.Errorf("ImportSQL failed:%w", err)
		}

		return nil
	},
}
