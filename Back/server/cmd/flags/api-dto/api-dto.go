package apidto

import "fmt"

type Feature string

const (
	FeatureAuth     Feature = "auth"
	FeatureClient   Feature = "client"
	FeatureProject  Feature = "project"
	FeatureFile     Feature = "file"
	FeatureTypeFile Feature = "type-files"
)

var AllFeatures = []Feature{
	FeatureAuth, FeatureClient, FeatureProject,
	FeatureFile, FeatureTypeFile,
}

type Field struct {
	Name     string
	Type     string
	Optional bool
}

type DTO struct {
	Name   string
	Fields []Field
	Role   string
	Source string
}

type Endpoint struct {
	Method      string
	Path        string
	Handler     string
	RequestDTO  string
	ResponseDTO string
}

type FeatureInfo struct {
	Name      Feature
	Requests  []DTO
	Responses []DTO
	Endpoints []Endpoint
}

func Run(featureName string, showAll bool) error {
	if showAll {
		for i, f := range AllFeatures {
			if i > 0 {
				fmt.Println()
			}
			info, err := getFeatureInfo(f)
			if err != nil {
				fmt.Fprintf(printWriter, "⚠ %s: %v\n", f, err)
				continue
			}
			printFeature(info)
		}
		return nil
	}

	feature := Feature(featureName)
	info, err := getFeatureInfo(feature)
	if err != nil {
		return fmt.Errorf("feature %q: %w", featureName, err)
	}
	printFeature(info)
	return nil
}
