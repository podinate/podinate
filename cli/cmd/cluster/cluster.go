package cluster

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cnrancher/autok3s/cmd"
	"github.com/cnrancher/autok3s/cmd/addon"
	"github.com/cnrancher/autok3s/cmd/airgap"
	"github.com/cnrancher/autok3s/cmd/sshkey"
	"github.com/cnrancher/autok3s/pkg/cli/kubectl"
	"github.com/cnrancher/autok3s/pkg/common"
	"github.com/cnrancher/autok3s/pkg/metrics"
	"github.com/docker/docker/pkg/reexec"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	gitVersion   = "dev"
	gitCommit    string
	gitTreeState string
	buildDate    string
)

func init() {
	reexec.Register("kubectl", kubectl.Main)

	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	common.CfgPath = filepath.Join(home, ".config/podinate")
}

func NewCommand() *cobra.Command {
	args := os.Args[0]
	os.Args[0] = filepath.Base(os.Args[0])
	if reexec.Init() {
		return nil
	}
	os.Args[0] = args

	rootCmd := cmd.Command()

	rootCmd.Use = "cluster"

	rootCmd.AddCommand(cmd.CompletionCommand(), cmd.VersionCommand(gitVersion, gitCommit, gitTreeState, buildDate),
		cmd.ListCommand(), cmd.CreateCommand(), cmd.JoinCommand(), cmd.KubectlCommand(), cmd.DeleteCommand(),
		cmd.SSHCommand(), cmd.DescribeCommand(), cmd.ServeCommand(), cmd.ExplorerCommand(), cmd.UpgradeCommand(),
		cmd.TelemetryCommand(), airgap.Command(), sshkey.Command(), cmd.DashboardCommand(), addon.Command())

	rootCmd.PersistentPreRun = func(c *cobra.Command, args []string) {
		common.InitLogger(logrus.StandardLogger())
		common.MetricsPrompt(c)
		common.SetupPrometheusMetrics(gitVersion)
		go metrics.Report()
		if c.Use == cmd.ServeCommand().Use {
			metrics.ReportEach(c.Context(), 1*time.Hour)
		}
	}
	rootCmd.PersistentPostRun = func(c *cobra.Command, args []string) {
		metrics.Report()
	}

	return rootCmd
}
