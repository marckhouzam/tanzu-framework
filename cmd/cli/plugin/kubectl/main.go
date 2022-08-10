// Copyright 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/aunum/log"
	"github.com/spf13/cobra"

	cliv1alpha1 "github.com/vmware-tanzu/tanzu-framework/apis/cli/v1alpha1"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/cli/command/plugin"
)

var descriptor = cliv1alpha1.PluginDescriptor{
	Name:           "kubectl",
	Description:    "Full kubectl functionality in tanzu",
	Group:          cliv1alpha1.ExtraCmdGroup,
	Aliases:        []string{"k", "kctl", "kube"},
	CompletionType: cliv1alpha1.NativePluginCompletion,
}

func main() {
	p, err := plugin.NewPlugin(&descriptor)
	if err != nil {
		log.Fatal(err)
	}
	p.Cmd.DisableFlagParsing = true
	p.Cmd.Args = cobra.ArbitraryArgs
	p.Cmd.CompletionOptions.DisableDefaultCmd = true

	p.Cmd.RunE = func(cmd *cobra.Command, args []string) error {
		path, err := exec.LookPath("kubectl")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Unable to find 'kubectl' on your system.  Please install it and make sure it is on $PATH")
			return err
		}

		execCmd := exec.Command(path, args...)

		execCmd.Stdin = os.Stdin
		execCmd.Stderr = os.Stderr
		execCmd.Stdout = os.Stdout

		return execCmd.Run()
	}

	p.Cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		path, err := exec.LookPath("kubectl")
		if err != nil {
			cobra.CompErrorln("Unable to find 'kubectl' on your system.  Please install it and make sure it is on $PATH")
			return nil, cobra.ShellCompDirectiveError
		}

		finalArgs := []string{cobra.ShellCompRequestCmd}
		finalArgs = append(finalArgs, args...)
		finalArgs = append(finalArgs, toComplete)
		execCmd := exec.Command(path, finalArgs...)
		execCmd.Stdin = os.Stdin
		execCmd.Stderr = os.Stderr
		buf := new(bytes.Buffer)
		execCmd.Stdout = buf

		err = execCmd.Run()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var completions []string
		for _, comp := range strings.Split(buf.String(), "\n") {
			// Remove any empty lines
			if len(comp) > 0 {
				completions = append(completions, comp)
			}
		}

		// Check the last line of output for the completion directive
		// of the form :<integer>
		directive := cobra.ShellCompDirectiveDefault
		if len(completions) > 0 {
			lastLine := completions[len(completions)-1]
			if len(lastLine) > 1 && lastLine[0] == ':' {
				if strInt, err := strconv.Atoi(lastLine[1:]); err == nil {
					directive = cobra.ShellCompDirective(strInt)
					completions = completions[:len(completions)-1]
				}
			}
		}
		return completions, directive
	}

	if err := p.Execute(); err != nil {
		os.Exit(1)
	}
}
