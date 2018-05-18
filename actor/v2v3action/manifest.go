package v2v3action

import (
	"code.cloudfoundry.org/cli/actor/versioncheck"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccversion"
	"code.cloudfoundry.org/cli/types"
	"code.cloudfoundry.org/cli/util/manifest"
)

func (actor *Actor) CreateApplicationManifestByNameAndSpace(appName string, appSpace string) (manifest.Application, Warnings, error) {
	var allWarnings Warnings

	manifestApp, v2warnings, err := actor.V2Actor.CreateApplicationManifestByNameAndSpace(appName, appSpace)
	allWarnings = append(allWarnings, v2warnings...)
	if err != nil {
		return manifest.Application{}, allWarnings, err
	}

	currentVersion := actor.V3Actor.CloudControllerAPIVersion()
	minimumVersion := ccversion.MinVersionManifestBuildpacksV3

	meetsV3Version, err := versioncheck.IsMinimumAPIVersionMet(currentVersion, minimumVersion)
	if err != nil {
		return manifest.Application{}, allWarnings, err
	}
	if meetsV3Version {
		v3App, v3warnings, v3Err := actor.V3Actor.GetApplicationByNameAndSpace(appName, appSpace)
		allWarnings = append(allWarnings, v3warnings...)
		if v3Err != nil {
			return manifest.Application{}, allWarnings, v3Err
		}

		manifestApp.Buildpack = types.FilteredString{}
		manifestApp.Buildpacks = v3App.LifecycleBuildpacks
	}

	return manifestApp, allWarnings, err
}

func (Actor) WriteApplicationManifest(manifestApp manifest.Application, manifestPath string) error {
	return manifest.WriteApplicationManifest(manifestApp, manifestPath)
}