package api

var defaultPodSecurityPolicyConfig = map[string]string{}

func (cs *ContainerService) setPodSecurityPolicyConfig() {
	o := cs.Properties.OrchestratorProfile

	// If no user-configurable scheduler config values exists, use the defaults
	if o.KubernetesConfig.PodSecurityPolicyConfig == nil {
		o.KubernetesConfig.PodSecurityPolicyConfig = defaultPodSecurityPolicyConfig
	} else {
		for key, val := range defaultPodSecurityPolicyConfig {
			// If we don't have a user-configurable scheduler config for each option
			if _, ok := o.KubernetesConfig.PodSecurityPolicyConfig[key]; !ok {
				// then assign the default value
				o.KubernetesConfig.PodSecurityPolicyConfig[key] = val
			}
		}
	}

}
