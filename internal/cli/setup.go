package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/fixi2/InfraTrack/internal/setup"
	"github.com/spf13/cobra"
)

func newSetupCmd() *cobra.Command {
	var (
		binDir        string
		scopeText     string
		completionRaw string
		noPath        bool
	)

	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Plan local InfraTrack installation and PATH integration",
		RunE: func(cmd *cobra.Command, _ []string) error {
			scope, completion, err := parseSetupInputs(scopeText, completionRaw)
			if err != nil {
				return err
			}
			plan, err := setup.BuildPlan(setup.PlanInput{
				Scope:      scope,
				BinDir:     binDir,
				NoPath:     noPath,
				Completion: completion,
			})
			if err != nil {
				return err
			}
			printSetupPlan(cmd, plan)
			return nil
		},
	}

	cmd.Flags().StringVar(&binDir, "bin-dir", "", "Install target directory for infratrack binary")
	cmd.Flags().StringVar(&scopeText, "scope", string(setup.ScopeUser), "Setup scope: user")
	cmd.Flags().BoolVar(&noPath, "no-path", false, "Do not modify PATH in setup plan")
	cmd.Flags().StringVar(&completionRaw, "completion", string(setup.CompletionNone), "Completion setup mode: none")

	cmd.AddCommand(newSetupStatusCmd(&binDir, &scopeText))
	cmd.AddCommand(newSetupApplyCmd())
	cmd.AddCommand(newSetupUndoCmd())
	return cmd
}

func newSetupStatusCmd(parentBinDir *string, parentScope *string) *cobra.Command {
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show setup status",
		RunE: func(cmd *cobra.Command, _ []string) error {
			scope, err := setup.ResolveScope(strings.TrimSpace(*parentScope))
			if err != nil {
				return err
			}
			status, err := setup.BuildStatus(scope, strings.TrimSpace(*parentBinDir))
			if err != nil {
				return err
			}

			if jsonOut {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(status)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "OS: %s\n", status.OS)
			fmt.Fprintf(cmd.OutOrStdout(), "Scope: %s\n", status.Scope)
			fmt.Fprintf(cmd.OutOrStdout(), "Current executable: %s\n", status.CurrentExe)
			fmt.Fprintf(cmd.OutOrStdout(), "Target bin dir: %s\n", status.BinDir)
			fmt.Fprintf(cmd.OutOrStdout(), "Target binary: %s\n", status.TargetBinaryPath)
			fmt.Fprintf(cmd.OutOrStdout(), "Installed: %s\n", yesNo(status.Installed))
			fmt.Fprintf(cmd.OutOrStdout(), "PATH contains bin dir: %s\n", yesNo(status.PathOK))
			fmt.Fprintf(cmd.OutOrStdout(), "Setup state found: %s\n", yesNo(status.StateFound))
			fmt.Fprintf(cmd.OutOrStdout(), "Pending finalize: %s\n", yesNo(status.PendingFinalize))
			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOut, "json", false, "Print machine-readable JSON status")
	return cmd
}

func newSetupApplyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "apply",
		Short: "Apply setup changes",
		RunE: func(_ *cobra.Command, _ []string) error {
			return errors.New("setup apply is not implemented in this build yet")
		},
	}
}

func newSetupUndoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "undo",
		Short: "Undo setup changes",
		RunE: func(_ *cobra.Command, _ []string) error {
			return errors.New("setup undo is not implemented in this build yet")
		},
	}
}

func parseSetupInputs(scopeText, completionText string) (setup.Scope, setup.CompletionMode, error) {
	scope, err := setup.ResolveScope(scopeText)
	if err != nil {
		return "", "", err
	}
	if scope != setup.ScopeUser {
		return "", "", errors.New("only --scope user is available in v0.5.0 setup")
	}
	completion, err := setup.ResolveCompletion(completionText)
	if err != nil {
		return "", "", err
	}
	return scope, completion, nil
}

func printSetupPlan(cmd *cobra.Command, plan setup.Plan) {
	fmt.Fprintf(cmd.OutOrStdout(), "Detected OS: %s\n", plan.OS)
	fmt.Fprintf(cmd.OutOrStdout(), "Scope: %s\n", plan.Scope)
	fmt.Fprintf(cmd.OutOrStdout(), "Current binary: %s\n", plan.CurrentExe)
	fmt.Fprintf(cmd.OutOrStdout(), "Target bin dir: %s\n", plan.TargetBinDir)
	fmt.Fprintf(cmd.OutOrStdout(), "Target binary: %s\n", plan.TargetBinaryPath)
	fmt.Fprintln(cmd.OutOrStdout(), "Actions:")
	for i, action := range plan.Actions {
		fmt.Fprintf(cmd.OutOrStdout(), "  %d. %s\n", i+1, action)
	}
	if len(plan.Notes) > 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "Notes:")
		for _, note := range plan.Notes {
			fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", note)
		}
	}
}

func yesNo(v bool) string {
	if v {
		return "yes"
	}
	return "no"
}
