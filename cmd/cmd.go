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
	flagAutoAssign = "auto-assign"
	flagDebug      = "debug"
	flagEnvVar     = "env-var"
	flagNoDefault  = "ignore-default"
)

func New() (*cobra.Command, func()) {
	r := &runner{}

	rootCmd := &cobra.Command{
		Use:   "tfvar [DIR]",
		Short: "A CLI tool that helps generate template for Terraform's variable definitions",
		Long: `Generate variable definitions template for Terraform module as
one would write it in .tfvars files.
`,
		PreRunE: r.preRootRunE,
		RunE:    r.rootRunE,
		Args:    cobra.ExactArgs(1),
	}

	rootCmd.PersistentFlags().BoolP(flagAutoAssign, "a", false, `Use values from environment variables TF_VAR_* and
variable definitions files e.g. terraform.tfvars[.json] *.auto.tfvars[.json]`)
	rootCmd.PersistentFlags().BoolP(flagDebug, "d", false, "Print debug log on stderr")
	rootCmd.PersistentFlags().BoolP(flagEnvVar, "e", false, "Print output in export TF_VAR_image_id=ami-abc123 format")
	rootCmd.PersistentFlags().Bool(flagNoDefault, false, "Do not use defined default values")

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
	dir := args[0]

	vars, err := tfvar.Load(dir)
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

	isAutoAssign, err := cmd.PersistentFlags().GetBool(flagAutoAssign)
	if err != nil {
		return errors.Wrap(err, "cmd: get flag --auto-assign")
	}

	unparseds := make(map[string]tfvar.UnparsedVariableValue)

	if isAutoAssign {
		r.log.Debug("Collecting values from environment variables")
		tfvar.CollectFromEnvVars(unparseds)

		autoFiles := tfvar.LookupTFVarsFiles(dir)

		for _, f := range autoFiles {
			if err := tfvar.CollectFromFile(f, unparseds); err != nil {
				return err
			}
		}
	}

	vars, err = tfvar.ParseValues(unparseds, vars)
	if err != nil {
		return err
	}

	writer := tfvar.WriteAsTFVars

	if isEnvVar {
		r.log.Debug("Print outputs in environment variables format")
		writer = tfvar.WriteAsEnvVars
	}

	return writer(os.Stdout, vars)
}
