package apidto

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func parseTSInterfaces(content string) ([]DTO, error) {
	interfacePattern := regexp.MustCompile(`export\s+(interface|type)\s+(\w+)\s*(?:=\s*)?\{([^}]*)\}`)
	fieldPattern := regexp.MustCompile(`\s*(\w+)\??\s*[?:]\s*([^;]+?)(?:;|\n)`)

	matches := interfacePattern.FindAllStringSubmatch(content, -1)
	var dtos []DTO

	for _, m := range matches {
		if len(m) < 4 {
			continue
		}
		name := m[2]
		body := m[3]

		dto := DTO{
			Name: name,
		}

		fieldMatches := fieldPattern.FindAllStringSubmatch(body, -1)
		for _, fm := range fieldMatches {
			if len(fm) < 3 {
				continue
			}
			fieldName := fm[1]

			fieldType := strings.TrimSpace(fm[2])
			fieldType = strings.TrimRight(fieldType, ";,")
			fieldType = strings.TrimSpace(fieldType)

			optional := strings.Contains(fm[0], "?")

			dto.Fields = append(dto.Fields, Field{
				Name:     fieldName,
				Type:     fieldType,
				Optional: optional,
			})
		}

		if len(dto.Fields) > 0 {
			dtos = append(dtos, dto)
		}
	}

	return dtos, nil
}

func findFrontendDTOs(frontendRoot string, feature Feature) ([]DTO, error) {
	dtoFileMap := map[Feature]string{
		FeatureClient:   "client",
		FeatureProject:  "project",
		FeatureFile:     "file",
		FeatureTypeFile: "fileType",
	}

	fileName, ok := dtoFileMap[feature]
	if !ok {
		return nil, nil
	}

	dtoDir := filepath.Join(frontendRoot, "app", "core", "domain", "dtos")
	dtoFile := filepath.Join(dtoDir, fileName+".ts")

	data, err := os.ReadFile(dtoFile)
	if err != nil {
		return nil, nil
	}

	return parseTSInterfaces(string(data))
}
