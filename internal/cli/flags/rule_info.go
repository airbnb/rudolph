package flags

import "github.com/spf13/cobra"

type RuleInfoFlags struct {
	RuleType   *RuleType
	Identifier *string
	FilePath   *string
}

func (r *RuleInfoFlags) AddRuleInfoFlags(cmd *cobra.Command) {
	var (
		ruleTypeArg   RuleType
		identifierArg string
		filepathArg   string
	)

	// Flag specifying the binary
	cmd.Flags().StringVarP(&filepathArg, "filepath", "f", "", `The filepath of a binary/application. Provide exactly one of [--filepath|--sha]`)
	cmd.Flags().StringVarP(&identifierArg, "identifier", "i", "", `The Identifier/SHA256 for a file, application, teamID, or signingID`)

	// rule-type should be one of "binary" or "cert" ("bin" and "certificate" also work)
	cmd.Flags().VarP(&ruleTypeArg, "rule-type", "t", `type of rule being applied. valid options are: "binary", "bin", "certificate", "cert", "teamid", "signingid"`)
	_ = cmd.MarkFlagRequired("rule-type")

	// If we want to make the `rule-type` flag optional with a default (say "binary"),
	// we can remove the previous 2 lines and instead use the following
	// cmd.Flags().VarP(newRuleTypeValue("binary", &ruleTypeArg), "rule-type", "t", `type of rule being applied. valid options are: "binary", "bin", "certificate", or "cert"`)

	// rule-policy is to specify the policy for edit commands

	r.RuleType = &ruleTypeArg
	r.Identifier = &identifierArg
	r.FilePath = &filepathArg
}
