package acsengine

import (
	"testing"
)

func TestWriteTLSArtifacts(t *testing.T) {

	writer := &acsengine.ArtifactWriter{
		Translator: &i18n.Translator{
			Locale: nil,
		},
	}

	a.WriteTLSArtifacts(containerService *api.ContainerService, apiVersion, template, parameters, artifactsDir string, certsGenerated bool, parametersOnly bool)
}
