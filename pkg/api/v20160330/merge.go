package v20160330

import (
	"github.com/imdario/mergo"
)

// Merge existing containerService attribute into cs
func (cs *ContainerService) Merge(ecs *ContainerService) error {
	if err := mergo.Merge(cs.Properties.WindowsProfile,
		*ecs.Properties.WindowsProfile); err != nil {
		return err
	}
	return nil
}
