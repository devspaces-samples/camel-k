/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"context"
	"github.com/apache/camel-k/pkg/util/test"
	"github.com/spf13/cobra"
	"os"
	"testing"
)

func kamelTestPostAddCommandInit(rootCmd *cobra.Command) *cobra.Command {
	rootCmd, _ = kamelPostAddCommandInit(*rootCmd)
	return rootCmd
}

func kamelTestPreAddCommandInit() (RootCmdOptions, *cobra.Command) {
	options := RootCmdOptions{
		Context: context.Background(),
	}
	rootCmd := kamelPreAddCommandInit(options)
	rootCmd.Run = test.EmptyRun
	return options, rootCmd
}

func TestLoadFromCommandLine(t *testing.T) {
	options, rootCmd := kamelTestPreAddCommandInit()

	runCmdOptions := addTestRunCmd(options, rootCmd)

	rootCmd = kamelTestPostAddCommandInit(rootCmd)

	_, err := test.ExecuteCommand(rootCmd, "run", "route.java", "--env", "VAR1=value,othervalue", "--env", "VAR2=value2")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(runCmdOptions.EnvVars) != 2 {
		t.Errorf("Properties expected to contain: \n %v elements\nGot:\n %v elemtns\n", 2, len(runCmdOptions.EnvVars))
	}
	if runCmdOptions.EnvVars[0] != "VAR1=value,othervalue" || runCmdOptions.EnvVars[1] != "VAR2=value2" {
		t.Errorf("EnvVars expected to be: \n %v\nGot:\n %v\n", "[VAR1=value,othervalue VAR=value2]", runCmdOptions.EnvVars)
	}
}

func TestLoadFromEnvVar(t *testing.T) {
	//shows how to include a "," character inside an env value see VAR1 value
	os.Setenv("KAMEL_RUN_ENVS", "\"VAR1=value,\"\"othervalue\"\"\",VAR2=value2")

	options, rootCmd := kamelTestPreAddCommandInit()

	runCmdOptions := addTestRunCmd(options, rootCmd)

	rootCmd = kamelTestPostAddCommandInit(rootCmd)

	_, err := test.ExecuteCommand(rootCmd, "run", "route.java")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(runCmdOptions.EnvVars) != 2 {
		t.Fatalf("Properties expected to contain: \n %v elements\nGot:\n %v elemtns\n", 2, len(runCmdOptions.EnvVars))
	}
	if runCmdOptions.EnvVars[0] != "VAR1=value,\"othervalue\"" || runCmdOptions.EnvVars[1] != "VAR2=value2" {
		t.Fatalf("EnvVars expected to be: \n %v\nGot:\n %v\n", "[VAR1=value,\"othervalue\" VAR=value2]", runCmdOptions.EnvVars)
	}
}