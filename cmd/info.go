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

	Kubernetes = "Kubernetes"
	DCOS       = "DCOS"
	Swarm      = "Swarm"
	DockerCE   = "DockerCE"
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
	case strings.EqualFold(ic.orchestrator, Kubernetes):
		ic.orchestrator = Kubernetes
	case strings.EqualFold(ic.orchestrator, DCOS):
		ic.orchestrator = DCOS
	case strings.EqualFold(ic.orchestrator, DockerCE):
		ic.orchestrator = DockerCE
	case strings.EqualFold(ic.orchestrator, Swarm):
		ic.orchestrator = Swarm
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
		Kubernetes: kubernetesInfo,
		DCOS:       dcosInfo,
		Swarm:      swarmInfo,
		DockerCE:   dockerceInfo,
	}
	if len(ic.orchestrator) == 0 {
		for o, f := range funcmap {
			fmt.Printf("%s:\n", o)
			if err := f(ic.release); err != nil {
				return err
			}
			fmt.Println()
		}
		return nil
	}
	return funcmap[ic.orchestrator](ic.release)
}

func printInfo(rel, ver, def string) {
	fmt.Printf("Release: %s Version: %s", rel, ver)
	if rel == def {
		fmt.Println(" (default)")
	} else {
		fmt.Println()
	}
}

func kubernetesInfo(release string) error {
	if len(release) == 0 {
		// print info for all supported versions
		for r, v := range common.KubeReleaseToVersion {
			printInfo(r, v, common.KubernetesDefaultRelease)
		}
	} else {
		// print info for the specified release
		ver, ok := common.KubeReleaseToVersion[release]
		if !ok {
			return fmt.Errorf("Kubernetes release %s is not supported", release)
		}
		printInfo(release, ver, common.KubernetesDefaultRelease)
	}
	return nil
}

func dcosInfo(release string) error {
	if len(release) == 0 {
		// print info for all supported versions
		for rel, ver := range common.DCOSReleaseToVersion {
			printInfo(rel, ver, common.DCOSDefaultRelease)
		}
	} else {
		// print info for the specified release
		ver, ok := common.DCOSReleaseToVersion[release]
		if !ok {
			return fmt.Errorf("DCOS release %s is not supported", release)
		}
		printInfo(release, ver, common.DCOSDefaultRelease)
	}
	return nil
}

func swarmInfo(release string) error {
	fmt.Println("Version: 1.1.0 Docker Compose Version: 1.6.2")
	return nil
}

func dockerceInfo(release string) error {
	fmt.Println("Version: 17.03 Docker Compose Version: 1.14.0")
	return nil
}
