package acsengine

import (
	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/api/v20160330"
	"github.com/Azure/acs-engine/pkg/api/vlabs"
)

// DCOSNodeType represents the type of DCOS Node
type DCOSNodeType string

// VlabsContainerService is the type we read and write from file
// needed because the json that is sent to ARM and acs-engine
// is different from the json that the ACS RP Api gets from ARM
type VlabsContainerService struct {
	api.TypeMeta
	*vlabs.ContainerService
}

// V20160330ContainerService is the type we read and write from file
// needed because the json that is sent to ARM and acs-engine
// is different from the json that the ACS RP Api gets from ARM
type V20160330ContainerService struct {
	api.TypeMeta
	*v20160330.ContainerService
}
