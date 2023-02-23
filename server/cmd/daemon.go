package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"server/configure"

	"github.com/powerpuffpenguin/easy-grpc/core/path"

	"github.com/powerpuffpenguin/easy-grpc/core"
	"github.com/powerpuffpenguin/easy-grpc/core/logger"
	"github.com/spf13/cobra"
)

func init() {
	var (
		filename    string
		debug, test bool
		basePath    = path.BasePath()

		addr string
	)

	cmd := &cobra.Command{
		Use:   `daemon`,
		Short: `run as daemon`,
		Run: func(cmd *cobra.Command, args []string) {
			// 加載配置
			cnf := configure.Default()
			e := cnf.Load(filename)
			if e != nil {
				log.Fatalln(e)
			}
			if addr != `` {
				cnf.HTTP.Addr = addr
			}
			// 測試配置
			if test {
				b, e := json.MarshalIndent(cnf, "", "\t")
				if e != nil {
					log.Fatalln(e)
				}
				fmt.Println(core.BytesToString(b))
				return
			}

			// 初始化日誌
			logger.Init(basePath, &cnf.Logger)

			// // init db
			// manipulator.Init(&cnf.DB)
			// sessionid.Init(&cnf.Session)

			// daemon.Run(&cnf.HTTP, debug)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&filename, `config`,
		`c`,
		path.Abs(basePath, filepath.Join(`etc`, `server.jsonnet`)),
		`configure file`,
	)
	flags.StringVarP(&addr, `addr`,
		`a`,
		``,
		`listen address`,
	)

	flags.BoolVarP(&debug, `debug`,
		`d`,
		false,
		`run as debug`,
	)
	flags.BoolVarP(&test, `test`,
		`t`,
		false,
		`test configure`,
	)
	rootCmd.AddCommand(cmd)
}
