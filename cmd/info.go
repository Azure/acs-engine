package cmd

import (
	"fmt"
	"strings"

	"github.com/Azure/acs-engine/pkg/api/common"
	"github.com/spf13/cobra"
)

const (
	infoName             = "info"
	infoShortDescription = "provide info about supported orchestrators"
	infoLongDescription  = "provide info about versions of supported orchestrators"

	kubernetes = "Kubernetes"
	dcos       = "DCOS"
	swarm      = "Swarm"
	dockerCE   = "DockerCE"

	// To be in sync with parts/configure-swarm-cluster.sh
	swarmVersion              = "1.1.0"
	swarmDockerComposeVersion = "1.6.2"
	// To be in sync with parts/configure-swarmmode-cluster.sh
	dockerceVersion              = "17.03"
	dockerceDockerComposeVersion = "1.14.0"
)

type infoCmd struct {
	// user input
	orchestrator string
	release      string
}

type infoFunc func(string) error

func newInfoCmd() *cobra.Command {
	ic := infoCmd{}

	infoCmd := &cobra.Command{
		Use:   infoName,
		Short: infoShortDescription,
		Long:  infoLongDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return ic.run(cmd, args)
		},
	}

	f := infoCmd.Flags()
	f.StringVar(&ic.orchestrator, "orchestrator", "", "orchestrator name (optional) ")
	f.StringVar(&ic.release, "release", "", "orchestrator release (optional)")

	return infoCmd
}

func (ic *infoCmd) validate() error {
	switch {
	case strings.EqualFold(ic.orchestrator, kubernetes):
		ic.orchestrator = kubernetes
	case strings.EqualFold(ic.orchestrator, dcos):
		ic.orchestrator = dcos
	case strings.EqualFold(ic.orchestrator, dockerCE):
		ic.orchestrator = dockerCE
	case strings.EqualFold(ic.orchestrator, swarm):
		ic.orchestrator = swarm
	case ic.orchestrator == "":
		if len(ic.release) > 0 {
			return fmt.Errorf("Must specify orchestrator for release '%s'", ic.release)
		}
	default:
		return fmt.Errorf("Unsupported orchestrator '%s'", ic.orchestrator)
	}
	return nil
}

func (ic *infoCmd) run(cmd *cobra.Command, args []string) error {
	if err := ic.validate(); err != nil {
		return err
	}
	funcmap := map[string]infoFunc{
		kubernetes: kubernetesInfo,
		dcos:       dcosInfo,
		swarm:      swarmInfo,
		dockerCE:   dockerceInfo,
	}
	if len(ic.orchestrator) == 0 {
		for _, f := range funcmap {
			if err := f(ic.release); err != nil {
				return err
			}
		}
		return nil
	}
	return funcmap[ic.orchestrator](ic.release)
}

func printInfo(orch, rel, ver, def string) {
	fmt.Printf("%s{Release: %s, Version: %s, Default: %t}\n", orch, rel, ver, rel == def)
}

func kubernetesInfo(release string) error {
	if len(release) == 0 {
		// print info for all supported versions
		for r, v := range common.KubeReleaseToVersion {
			printInfo(kubernetes, r, v, common.KubernetesDefaultRelease)
		}
	} else {
		// print info for the specified release
		ver, ok := common.KubeReleaseToVersion[release]
		if !ok {
			return fmt.Errorf("Kubernetes release %s is not supported", release)
		}
		printInfo(kubernetes, release, ver, common.KubernetesDefaultRelease)
	}
	return nil
}

func dcosInfo(release string) error {
	if len(release) == 0 {
		// print info for all supported versions
		for rel, ver := range common.DCOSReleaseToVersion {
			printInfo(dcos, rel, ver, common.DCOSDefaultRelease)
		}
	} else {
		// print info for the specified release
		ver, ok := common.DCOSReleaseToVersion[release]
		if !ok {
			return fmt.Errorf("DCOS release %s is not supported", release)
		}
		printInfo(dcos, release, ver, common.DCOSDefaultRelease)
	}
	return nil
}

func swarmInfo(release string) error {
	fmt.Printf("Swarm{Version: %s, Docker Compose Version: %s}\n", swarmVersion, swarmDockerComposeVersion)
	return nil
}

func dockerceInfo(release string) error {
	fmt.Printf("DockerCE{Version: %s, Docker Compose Version: %s}\n", dockerceVersion, dockerceDockerComposeVersion)
	return nil
}
