package _tests

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils/_tests"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

const (
	JfrogTestsHome      = ".jfrogTest"
	JfrogHomeEnv        = "JFROG_CLI_HOME"
	CliIntegrationTests = "github.com/jfrog/jfrog-client-go/tests"
)

func TestMain(m *testing.M) {
	InitArtifactoryServiceManager()
	result := m.Run()
	os.Exit(result)
}

func InitArtifactoryServiceManager() {
	flag.Parse()
	log.SetLogger(log.NewLogger(log.DEBUG, nil))
	createArtifactoryUploadManager()
	createArtifactorySearchManager()
	createArtifactoryDeleteManager()
	createArtifactoryDownloadManager()
	createArtifactorySecurityManager()
	createArtifactoryCreateLocalRepositoryManager()
	createArtifactoryUpdateLocalRepositoryManager()
	createArtifactoryCreateRemoteRepositoryManager()
	createArtifactoryUpdateRemoteRepositoryManager()
	createArtifactoryCreateVirtualRepositoryManager()
	createArtifactoryUpdateVirtualRepositoryManager()
	createArtifactoryDeleteRepositoryManager()
	createArtifactoryGetRepositoryManager()
	createArtifactoryReplicationCreateManager()
	createArtifactoryReplicationUpdateManager()
	createArtifactoryReplicationGetManager()
	createArtifactoryReplicationDeleteManager()
	if *DistUrl != "" {
		createDistributionManager()
	}
	createReposIfNeeded()
}

func TestUnitTests(t *testing.T) {
	homePath, err := filepath.Abs(JfrogTestsHome)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	setJfrogHome(homePath)
	packages := _tests.GetTestPackages("./../...")
	packages = _tests.ExcludeTestsPackage(packages, CliIntegrationTests)
	_tests.RunTests(packages, false)
	cleanUnitTestsJfrogHome(homePath)
}

func setJfrogHome(homePath string) {
	if err := os.Setenv(JfrogHomeEnv, homePath); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func cleanUnitTestsJfrogHome(homePath string) {
	os.RemoveAll(homePath)
	if err := os.Unsetenv(JfrogHomeEnv); err != nil {
		os.Exit(1)
	}
}
