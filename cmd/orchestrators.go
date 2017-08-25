package cmd

import (
	"fmt"
	"strings"

	"github.com/Azure/acs-engine/pkg/acsengine"
	"github.com/Azure/acs-engine/pkg/api/common"
	"github.com/spf13/cobra"
)

const (
	cmdName             = "orchestrators"
	cmdShortDescription = "provide info about supported orchestrators"
	cmdLongDescription  = "provide info about versions of supported orchestrators"

	kubernetes = "Kubernetes"
	dcos       = "DCOS"
	swarm      = "Swarm"
	dockerCE   = "DockerCE"
)

type orchestratorsCmd struct {
	// user input
	orchestrator string
	release      string
}

type orchestratorsFunc func(string) error

func newOrchestratorsCmd() *cobra.Command {
	oc := orchestratorsCmd{}

	command := &cobra.Command{
		Use:   cmdName,
		Short: cmdShortDescription,
		Long:  cmdLongDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return oc.run(cmd, args)
		},
	}

	f := command.Flags()
	f.StringVar(&oc.orchestrator, "orchestrator", "", "orchestrator name (optional) ")
	f.StringVar(&oc.release, "release", "", "orchestrator release (optional)")

	return command
}

func (oc *orchestratorsCmd) validate() error {
	switch {
	case strings.EqualFold(oc.orchestrator, kubernetes):
		oc.orchestrator = kubernetes
	case strings.EqualFold(oc.orchestrator, dcos):
		oc.orchestrator = dcos
	case strings.EqualFold(oc.orchestrator, dockerCE):
		oc.orchestrator = dockerCE
	case strings.EqualFold(oc.orchestrator, swarm):
		oc.orchestrator = swarm
	case oc.orchestrator == "":
		if len(oc.release) > 0 {
			return fmt.Errorf("Must specify orchestrator for release '%s'", oc.release)
		}
	default:
		return fmt.Errorf("Unsupported orchestrator '%s'", oc.orchestrator)
	}
	return nil
}

func (oc *orchestratorsCmd) run(cmd *cobra.Command, args []string) error {
	if err := oc.validate(); err != nil {
		return err
	}
	funcmap := map[string]orchestratorsFunc{
		kubernetes: kubernetesInfo,
		dcos:       dcosInfo,
		swarm:      swarmInfo,
		dockerCE:   dockerceInfo,
	}
	if len(oc.orchestrator) == 0 {
		for _, f := range funcmap {
			if err := f(oc.release); err != nil {
				return err
			}
		}
		return nil
	}
	return funcmap[oc.orchestrator](oc.release)
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
	fmt.Printf("Swarm{Version: %s, Docker Compose Version: %s}\n", acsengine.SwarmVersion, acsengine.SwarmDockerComposeVersion)
	return nil
}

func dockerceInfo(release string) error {
	fmt.Printf("DockerCE{Version: %s, Docker Compose Version: %s}\n", acsengine.DockerCEVersion, acsengine.DockerCEDockerComposeVersion)
	return nil
}
