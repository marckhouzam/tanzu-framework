// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/clusterclient"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/log"
)

func WaitForPackagesInstallation(clusterClient clusterclient.Client,  packageInstallNames []string, namespace string, packageInstallTimeout time.Duration) error {
	// Start waiting for all packages in parallel using group.Wait
	// Note: As PackageInstall resources are created in the cluster itself
	// we are using currentClusterClient which will point to correct cluster
	group, _ := errgroup.WithContext(context.Background())

	for _, packageName := range packageInstallNames {
		pn := packageName
		log.V(3).Warningf("Waiting for package: %s", pn)
		group.Go(
			func() error {
				err := clusterClient.WaitForPackageInstall(pn, namespace, packageInstallTimeout)
				if err != nil {
					log.V(3).Warningf("Failure while waiting for package '%s'", pn)
				} else {
					log.V(3).Infof("Successfully reconciled package: %s", pn)
				}
				return err
			})
	}

	err := group.Wait()
	if err != nil {
		return errors.Wrap(err, "Failure while waiting for packages to be installed")
	}

	return nil
}