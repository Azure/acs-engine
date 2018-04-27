package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/Azure/acs-engine/pkg/acsengine"
	"github.com/Azure/acs-engine/pkg/acsengine/transform"
	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/i18n"
	"github.com/Jeffail/gabs"
	"github.com/leonelquinteros/gotext"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	generateName             = "generate"
	generateShortDescription = "Generate an Azure Resource Manager template"
	generateLongDescription  = "Generates an Azure Resource Manager template, parameters file and other assets for a cluster"
)

type generateCmd struct {
	apimodelPath      string
	outputDirectory   string // can be auto-determined from clusterDefinition
	caCertificatePath string
	caPrivateKeyPath  string
	classicMode       bool
	noPrettyPrint     bool
	parametersOnly    bool
	set               []string

	// derived
	containerService *api.ContainerService
	apiVersion       string
	locale           *gotext.Locale
}

type setFlagValue struct {
	stringValue   string
	intValue      int64
	arrayValue    bool
	arrayIndex    int
	arrayProperty string
	arrayName     string
}

func newGenerateCmd() *cobra.Command {
	gc := generateCmd{}

	generateCmd := &cobra.Command{
		Use:   generateName,
		Short: generateShortDescription,
		Long:  generateLongDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := gc.validate(cmd, args); err != nil {
				log.Fatalf(fmt.Sprintf("error validating generateCmd: %s", err.Error()))
			}

			if err := gc.mergeAPIModel(); err != nil {
				log.Fatalf(fmt.Sprintf("error merging API model in generateCmd: %s", err.Error()))
			}

			if err := gc.loadAPIModel(cmd, args); err != nil {
				log.Fatalf(fmt.Sprintf("error loading API model in generateCmd: %s", err.Error()))
			}

			return gc.run()
		},
	}

	f := generateCmd.Flags()
	f.StringVar(&gc.apimodelPath, "api-model", "", "")
	f.StringVar(&gc.outputDirectory, "output-directory", "", "output directory (derived from FQDN if absent)")
	f.StringVar(&gc.caCertificatePath, "ca-certificate-path", "", "path to the CA certificate to use for Kubernetes PKI assets")
	f.StringVar(&gc.caPrivateKeyPath, "ca-private-key-path", "", "path to the CA private key to use for Kubernetes PKI assets")
	f.StringArrayVar(&gc.set, "set", []string{}, "set values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	f.BoolVar(&gc.classicMode, "classic-mode", false, "enable classic parameters and outputs")
	f.BoolVar(&gc.noPrettyPrint, "no-pretty-print", false, "skip pretty printing the output")
	f.BoolVar(&gc.parametersOnly, "parameters-only", false, "only output parameters files")

	return generateCmd
}

func (gc *generateCmd) validate(cmd *cobra.Command, args []string) error {
	var err error

	gc.locale, err = i18n.LoadTranslations()
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("error loading translation files: %s", err.Error()))
	}

	if gc.apimodelPath == "" {
		if len(args) == 1 {
			gc.apimodelPath = args[0]
		} else if len(args) > 1 {
			cmd.Usage()
			return errors.New("too many arguments were provided to 'generate'")
		} else {
			cmd.Usage()
			return errors.New("--api-model was not supplied, nor was one specified as a positional argument")
		}
	}

	if _, err := os.Stat(gc.apimodelPath); os.IsNotExist(err) {
		return fmt.Errorf(fmt.Sprintf("specified api model does not exist (%s)", gc.apimodelPath))
	}

	return nil
}

func (gc *generateCmd) mergeAPIModel() error {
	var err error
	log.Infoln("test")

	// if --set flag has been used
	if gc.set != nil && len(gc.set) > 0 {
		m := make(map[string]setFlagValue)
		mapValues(m, gc.set)

		// overrides the api model and generates a new file
		gc.apimodelPath, err = mergeValuesWithAPIModel(gc.apimodelPath, m)
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("error merging --set values with the api model: %s", err.Error()))
		}

		log.Infoln(fmt.Sprintf("new api model file has been generated during merge: %s", gc.apimodelPath))
	}

	return nil
}

func (gc *generateCmd) loadAPIModel(cmd *cobra.Command, args []string) error {
	var caCertificateBytes []byte
	var caKeyBytes []byte
	var err error

	apiloader := &api.Apiloader{
		Translator: &i18n.Translator{
			Locale: gc.locale,
		},
	}
	gc.containerService, gc.apiVersion, err = apiloader.LoadContainerServiceFromFile(gc.apimodelPath, true, false, nil)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("error parsing the api model: %s", err.Error()))
	}

	if gc.outputDirectory == "" {
		if gc.containerService.Properties.MasterProfile != nil {
			gc.outputDirectory = path.Join("_output", gc.containerService.Properties.MasterProfile.DNSPrefix)
		} else {
			gc.outputDirectory = path.Join("_output", gc.containerService.Properties.HostedMasterProfile.DNSPrefix)
		}
	}

	// consume gc.caCertificatePath and gc.caPrivateKeyPath

	if (gc.caCertificatePath != "" && gc.caPrivateKeyPath == "") || (gc.caCertificatePath == "" && gc.caPrivateKeyPath != "") {
		return errors.New("--ca-certificate-path and --ca-private-key-path must be specified together")
	}
	if gc.caCertificatePath != "" {
		if caCertificateBytes, err = ioutil.ReadFile(gc.caCertificatePath); err != nil {
			return fmt.Errorf(fmt.Sprintf("failed to read CA certificate file: %s", err.Error()))
		}
		if caKeyBytes, err = ioutil.ReadFile(gc.caPrivateKeyPath); err != nil {
			return fmt.Errorf(fmt.Sprintf("failed to read CA private key file: %s", err.Error()))
		}

		prop := gc.containerService.Properties
		if prop.CertificateProfile == nil {
			prop.CertificateProfile = &api.CertificateProfile{}
		}
		prop.CertificateProfile.CaCertificate = string(caCertificateBytes)
		prop.CertificateProfile.CaPrivateKey = string(caKeyBytes)
	}

	return nil
}

func (gc *generateCmd) run() error {
	log.Infoln(fmt.Sprintf("Generating assets into %s...", gc.outputDirectory))

	ctx := acsengine.Context{
		Translator: &i18n.Translator{
			Locale: gc.locale,
		},
	}
	templateGenerator, err := acsengine.InitializeTemplateGenerator(ctx, gc.classicMode)
	if err != nil {
		log.Fatalln("failed to initialize template generator: %s", err.Error())
	}

	template, parameters, certsGenerated, err := templateGenerator.GenerateTemplate(gc.containerService, acsengine.DefaultGeneratorCode, false)
	if err != nil {
		log.Fatalf("error generating template %s: %s", gc.apimodelPath, err.Error())
		os.Exit(1)
	}

	if !gc.noPrettyPrint {
		if template, err = transform.PrettyPrintArmTemplate(template); err != nil {
			log.Fatalf("error pretty printing template: %s \n", err.Error())
		}
		if parameters, err = transform.BuildAzureParametersFile(parameters); err != nil {
			log.Fatalf("error pretty printing template parameters: %s \n", err.Error())
		}
	}

	writer := &acsengine.ArtifactWriter{
		Translator: &i18n.Translator{
			Locale: gc.locale,
		},
	}
	if err = writer.WriteTLSArtifacts(gc.containerService, gc.apiVersion, template, parameters, gc.outputDirectory, certsGenerated, gc.parametersOnly); err != nil {
		log.Fatalf("error writing artifacts: %s \n", err.Error())
	}

	return nil
}

func mapValues(m map[string]setFlagValue, values []string) {
	if values == nil || len(values) == 0 {
		return
	}

	for _, value := range values {
		splittedValues := strings.Split(value, ",")
		if len(splittedValues) > 1 {
			mapValues(m, splittedValues)
		} else {
			keyValueSplitted := strings.Split(value, "=")
			key := keyValueSplitted[0]
			stringValue := keyValueSplitted[1]

			flagValue := setFlagValue{}

			if asInteger, err := strconv.ParseInt(stringValue, 10, 64); err == nil {
				flagValue.intValue = asInteger
			} else {
				flagValue.stringValue = stringValue
			}

			// use regex to find array[index].property pattern in the key
			re := regexp.MustCompile(`(.*?)\[(.*?)\]\.(.*?)$`)
			match := re.FindStringSubmatch(key)

			// it's an array
			if len(match) != 0 {
				i, err := strconv.ParseInt(match[2], 10, 32)
				if err != nil {
					log.Warnln(fmt.Sprintf("array index is not specified for property %s", key))
				} else {
					arrayIndex := int(i)
					flagValue.arrayValue = true
					flagValue.arrayName = match[1]
					flagValue.arrayIndex = arrayIndex
					flagValue.arrayProperty = match[3]
					m[key] = flagValue
				}
			} else {
				m[key] = flagValue
			}
		}
	}
}

func mergeValuesWithAPIModel(apiModelPath string, m map[string]setFlagValue) (string, error) {
	// load the apiModel file from path
	fileContent, err := ioutil.ReadFile(apiModelPath)
	if err != nil {
		return "", err
	}

	// parse the json from file content
	jsonObj, err := gabs.ParseJSON(fileContent)
	if err != nil {
		return "", err
	}

	// update api model definition with each value in the map
	for key, flagValue := range m {
		// working on an array
		if flagValue.arrayValue {
			log.Infoln(fmt.Sprintf("--set flag array value detected. Name: %s, Index: %b, PropertyName: %s", flagValue.arrayName, flagValue.arrayIndex, flagValue.arrayProperty))
			arrayValue := jsonObj.Path(fmt.Sprint("properties.", flagValue.arrayName))
			if flagValue.stringValue != "" {
				arrayValue.Index(flagValue.arrayIndex).SetP(flagValue.stringValue, flagValue.arrayProperty)
			} else {
				arrayValue.Index(flagValue.arrayIndex).SetP(flagValue.intValue, flagValue.arrayProperty)
			}
		} else {
			if flagValue.stringValue != "" {
				jsonObj.SetP(flagValue.stringValue, fmt.Sprint("properties.", key))
			} else {
				jsonObj.SetP(flagValue.intValue, fmt.Sprint("properties.", key))
			}
		}
	}

	// generate a new file
	tmpFile, err := ioutil.TempFile("", "mergedApiModel")
	if err != nil {
		return "", err
	}

	tmpFileName := tmpFile.Name()
	err = ioutil.WriteFile(tmpFileName, []byte(jsonObj.String()), os.ModeAppend)
	if err != nil {
		return "", err
	}

	return tmpFileName, nil
}
