package transform

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Jeffail/gabs"
	log "github.com/sirupsen/logrus"
)

// APIModelValue represents a value in the APIModel JSON file
type APIModelValue struct {
	stringValue   string
	intValue      int64
	arrayValue    bool
	arrayIndex    int
	arrayProperty string
	arrayName     string
}

// MapValues converts an arraw of rwa ApiModel values (like ["masterProfile.count=4","linuxProfile.adminUsername=admin"]) to a map
func MapValues(m map[string]APIModelValue, values []string) {
	if values == nil || len(values) == 0 {
		return
	}

	for _, value := range values {
		splittedValues := strings.Split(value, ",")
		if len(splittedValues) > 1 {
			MapValues(m, splittedValues)
		} else {
			keyValueSplitted := strings.Split(value, "=")
			key := keyValueSplitted[0]
			stringValue := keyValueSplitted[1]

			flagValue := APIModelValue{}

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

// MergeValuesWithAPIModel takes the path to an ApiModel JSON file, loads it and merges it with the values in the map to another temp file
func MergeValuesWithAPIModel(apiModelPath string, m map[string]APIModelValue) (string, error) {
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
