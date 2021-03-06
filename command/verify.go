package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/mitchellh/cli"
)

type VerifyCmd struct {
	Version string
}

func NewVerify(v string) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return VerifyCmd{Version: v}, nil
	}
}

func (v VerifyCmd) Help() string {
	return v.Synopsis()
}

func (v VerifyCmd) Run(args []string) int {
	if len(args) == 0 {
		fmt.Println("Unable to verify version, version constraint not provided")
		return 1
	}

	valid, err := v.check(args[0])
	if err != nil {
		fmt.Println(err)
		return 1
	}
	fmt.Printf("%t", valid)
	return 0
}

func (v VerifyCmd) Synopsis() string {
	return `
	Usage: gtm verify <version constraint>
	Verify gtm satisfies the version constraint
	`
}

func (v VerifyCmd) check(constraint string) (bool, error) {
	// Our version tags can have a 'v' prefix
	// Strip v prefix if it exists because it's not valid for a Semantic version
	cleanVersion := v.Version
	if strings.HasPrefix(strings.ToLower(cleanVersion), "v") {
		cleanVersion = cleanVersion[1:]
	}

	ver, err := version.NewVersion(cleanVersion)
	if err != nil {
		return false, err
	}

	c, err := version.NewConstraint(constraint)
	if err != nil {
		return false, err
	}

	return c.Check(ver), nil
}
