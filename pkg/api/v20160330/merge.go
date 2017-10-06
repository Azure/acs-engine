package v20160330

import (
	"github.com/imdario/mergo"
)

// Merge existing containerService attribute into cs
func (cs *ContainerService) Merge(ecs *ContainerService) error {
	if ecs.Properties.WindowsProfile != nil {
		if cs.Properties.WindowsProfile == nil {
			cs.Properties.WindowsProfile = &WindowsProfile{}
		}
		if err := mergo.Merge(cs.Properties.WindowsProfile,
			*ecs.Properties.WindowsProfile); err != nil {
			return err
		}
	}
	return nil
}
