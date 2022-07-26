// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kapppkgv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	runtanzuv1alpha3 "github.com/vmware-tanzu/tanzu-framework/apis/run/v1alpha3"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/clusterclient"
	"github.com/vmware-tanzu/tanzu-framework/pkg/v1/tkg/log"
)

// GetClusterBootstrap returns ClusterBootstrap object for the given clustername in the management cluster
func  GetClusterBootstrap(managementClusterClient clusterclient.Client, clusterName, namespace string) (*runtanzuv1alpha3.ClusterBootstrap, error) {
	clusterBootstrap := &runtanzuv1alpha3.ClusterBootstrap{}
	err := managementClusterClient.GetResource(clusterBootstrap, clusterName, namespace,nil,&clusterclient.PollOptions{Interval: clusterclient.CheckResourceInterval, Timeout: clusterclient.PackageInstallTimeout})
	return clusterBootstrap, err
}

// GetCorePackagesFromCB returns addon's core packages details from the given ClsuterBootstrap object
func GetCorePackagesFromCB(clusterBootstrap *runtanzuv1alpha3.ClusterBootstrap, namespace string) []kapppkgv1alpha1.Package{
	var packages []kapppkgv1alpha1.Package
	if clusterBootstrap.Spec.Kapp != nil{
		kappPkgShortName, kappPkgName, kappPkgVersion := getPackageDetailsFromCBS(clusterBootstrap.Spec.Kapp.RefName)
		packages = append(packages, kapppkgv1alpha1.Package{ObjectMeta: metav1.ObjectMeta{Name: kappPkgShortName, Namespace: namespace},
			Spec: kapppkgv1alpha1.PackageSpec{RefName: kappPkgName, Version: kappPkgVersion}})
	}
	if clusterBootstrap.Spec.CNI != nil{
		kappPkgShortName, kappPkgName, kappPkgVersion := getPackageDetailsFromCBS(clusterBootstrap.Spec.CNI.RefName)
		packages = append(packages, kapppkgv1alpha1.Package{ObjectMeta: metav1.ObjectMeta{Name: kappPkgShortName, Namespace: namespace},
			Spec: kapppkgv1alpha1.PackageSpec{RefName: kappPkgName, Version: kappPkgVersion}})
	}
	if clusterBootstrap.Spec.CSI != nil{
		kappPkgShortName, kappPkgName, kappPkgVersion := getPackageDetailsFromCBS(clusterBootstrap.Spec.CSI.RefName)
		packages = append(packages, kapppkgv1alpha1.Package{ObjectMeta: metav1.ObjectMeta{Name: kappPkgShortName, Namespace: namespace},
			Spec: kapppkgv1alpha1.PackageSpec{RefName: kappPkgName, Version: kappPkgVersion}})
	}
	if clusterBootstrap.Spec.CPI != nil{
		kappPkgShortName, kappPkgName, kappPkgVersion := getPackageDetailsFromCBS(clusterBootstrap.Spec.CPI.RefName)
		packages = append(packages, kapppkgv1alpha1.Package{ObjectMeta: metav1.ObjectMeta{Name: kappPkgShortName, Namespace: namespace},
			Spec: kapppkgv1alpha1.PackageSpec{RefName: kappPkgName, Version: kappPkgVersion}})
	}
	
	return packages
}

// getPackageDetailsFromCBS parse given clusterbootstrap package reference and returns package short name, package name and version details
func getPackageDetailsFromCBS(CBSRefName string) (pkgShortName, pkgName, pkgVersion string) {
	pkgShortName = strings.Split(CBSRefName, ".")[0]
	pkgName = strings.Join(strings.Split(CBSRefName, ".")[0:4], ".")
	pkgVersion = strings.Join(strings.Split(CBSRefName, ".")[4:], ".")
	return
}

// MonitorAddonsCorePackageInstallation monitors addon's core packages (kapp, cni, csi and cpi) and returns error if any while monitoring packages or any packages are not installed successfully. First it monitors kapp package in management cluster then it monitors other core packages in workload cluster.
func MonitorAddonsCorePackageInstallation(regionalClusterClient clusterclient.Client, workloadClusterClient clusterclient.Client, packages []kapppkgv1alpha1.Package, packageInstallTimeout time.Duration) error {
	if len(packages) == 0{
		return nil
	}
	var corePackages []string
	var kappPackage []string
	for _, p := range packages{
		if strings.Contains(p.ObjectMeta.Name, "kapp-controller") {
			kappPackage =  append(kappPackage, p.ObjectMeta.Name)
		}else{
			corePackages = append(corePackages, p.ObjectMeta.Name)
		}
	}
	if len(kappPackage) > 0{
		err := WaitForPackagesInstallation(regionalClusterClient, kappPackage,packages[0].ObjectMeta.Namespace, packageInstallTimeout)
		if err != nil{
			return err
		}
	}
	return WaitForPackagesInstallation(workloadClusterClient, corePackages, packages[0].ObjectMeta.Namespace, packageInstallTimeout)
}

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