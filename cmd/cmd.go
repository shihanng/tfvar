package cmd

import (
	"os"

	"github.com/cockroachdb/errors"
	"github.com/shihanng/tfvar/pkg/tfvar"
	"github.com/spf13/cobra"
	"github.com/zclconf/go-cty/cty"
	"go.uber.org/zap"
)

const (
	flagDebug     = "debug"
	flagNoDefault = "ignore-default"
	flagEnvVar    = "env-var"
)

func New() (*cobra.Command, func()) {
	r := &runner{}

	rootCmd := &cobra.Command{
		Use:   "tfvar [DIR]",
		Short: "A CLI tool that helps generate variable definitions for Terraform module",
		Long: `Generate variable definitions for Terraform module as one would write it
in .tfvars files. This tool works for both root modules and child modules.
`,
		PreRunE: r.preRootRunE,
		RunE:    r.rootRunE,
		Args:    cobra.ExactArgs(1),
	}

	rootCmd.PersistentFlags().BoolP(flagDebug, "d", false, "Print debug log on stderr")
	rootCmd.PersistentFlags().Bool(flagNoDefault, false, "Do not use defined default values")
	rootCmd.PersistentFlags().BoolP(flagEnvVar, "e", false, "Print output in export TF_VAR_image_id=ami-abc123 format")

	return rootCmd, func() {
		if r.log != nil {
			_ = r.log.Sync()
		}
	}
}

type runner struct {
	log *zap.SugaredLogger
}

func (r *runner) preRootRunE(cmd *cobra.Command, args []string) error {
	// Setup logger
	logConfig := zap.NewDevelopmentConfig()

	isDebug, err := cmd.PersistentFlags().GetBool(flagDebug)
	if err != nil {
		return errors.Wrap(err, "cmd: get flag --debug")
	}

	if !isDebug {
		logConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	logger, err := logConfig.Build()
	if err != nil {
		return errors.Wrap(err, "cmd: create new logger")
	}

	r.log = logger.Sugar()
	r.log.Debug("Logger initialized")

	return nil
}

func (r *runner) rootRunE(cmd *cobra.Command, args []string) error {
	vars, err := tfvar.Load(args[0])
	if err != nil {
		return err
	}

	ignoreDefault, err := cmd.PersistentFlags().GetBool(flagNoDefault)
	if err != nil {
		return errors.Wrap(err, "cmd: get flag --ignore-default")
	}

	if ignoreDefault {
		r.log.Debug("Replacing values with null")
		for i, v := range vars {
			vars[i].Value = cty.NullVal(v.Value.Type())
		}
	}

	isEnvVar, err := cmd.PersistentFlags().GetBool(flagEnvVar)
	if err != nil {
		return errors.Wrap(err, "cmd: get flag --env-var")
	}

	writer := tfvar.WriteAsTFVars

	if isEnvVar {
		r.log.Debug("Print outputs in environment variables format")
		writer = tfvar.WriteAsEnvVars
	}

	return writer(os.Stdout, vars)
}
