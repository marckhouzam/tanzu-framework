// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package helper

import (
	"fmt"
	"os"
)

// GetFakeKubeConfigFilePath returns fake kubeconfig file path
func GetFakeKubeConfigFilePath(testingDir, filePath string) string {
	f, err := os.CreateTemp(testingDir, "kube")
	if err != nil {
		fmt.Println("Error creating TempFile: ", err.Error())
	}
	copyFile(filePath, f.Name())
	return f.Name()
}

func copyFile(sourceFile, destFile string) {
	input, err := os.ReadFile(sourceFile)
	if err != nil {
		fmt.Println("Error ReadFile TempFile: ", err.Error())
	}
	_ = os.WriteFile(destFile, input, 0o600)
	if err != nil {
		fmt.Println("Error WriteFile TempFile: ", err.Error())
	}
}

// CreateTempTestingDirectory create temporary directory for testing
func CreateTempTestingDirectory() string {
	testingDir, _ := os.MkdirTemp("", "testing")
	return testingDir
}

// DeleteTempTestingDirectory deletes temporary directory
func DeleteTempTestingDirectory(testingDir string) {
	os.Remove(testingDir)
}
