package apidto

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func findProjectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := wd
	for {
		backDir := filepath.Join(dir, "Back")
		frontDir := filepath.Join(dir, "Front")
		if info, err := os.Stat(backDir); err == nil && info.IsDir() {
			if info2, err := os.Stat(frontDir); err == nil && info2.IsDir() {
				return dir, nil
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("project root not found (no Back/ and Front/ directories found from %s)", wd)
		}
		dir = parent
	}
}

func getFeatureInfo(feature Feature) (*FeatureInfo, error) {
	root, err := findProjectRoot()
	if err != nil {
		return nil, err
	}

	info := &FeatureInfo{Name: feature}

	request, response, err := findDTOs(root, feature)
	if err == nil {
		if request != nil {
			info.Requests = append(info.Requests, *request)
		}
		if response != nil {
			info.Responses = append(info.Responses, *response)
		}
	}

	frontendRoot := filepath.Join(root, "Front", "client", "src")
	if fi, cerr := findFrontendDTOs(frontendRoot, feature); cerr == nil {
		for _, dto := range fi {
			dto.Source = "frontend"
			if strings.Contains(dto.Name, "Request") || strings.Contains(dto.Name, "Update") {
				info.Requests = append(info.Requests, dto)
			} else if strings.Contains(dto.Name, "Data") || strings.Contains(dto.Name, "Response") {
				info.Responses = append(info.Responses, dto)
			}
		}
	}

	endpoints, err := findEndpoints(root, feature)
	if err == nil {
		info.Endpoints = append(info.Endpoints, endpoints...)
	}

	return info, nil
}


