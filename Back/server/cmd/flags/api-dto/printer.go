package apidto

import (
	"fmt"
	"io"
	"os"
	"strings"
)

var printWriter io.Writer = os.Stdout

func printFeature(info *FeatureInfo) {
	title := fmt.Sprintf("★ Feature: %s", info.Name)
	fmt.Fprintln(printWriter, title)
	fmt.Fprintln(printWriter, strings.Repeat("─", len(title)))
	fmt.Fprintln(printWriter)

	printDTOs(info)
	printEndpoints(info)
}

func printDTOs(info *FeatureInfo) {
	if len(info.Requests) > 0 {
		fmt.Fprintln(printWriter, "▸ REQUEST (Frontend → Backend)")
		for i := range info.Requests {
			dtoSummary(&info.Requests[i])
		}
		fmt.Fprintln(printWriter)
	}

	if len(info.Responses) > 0 {
		fmt.Fprintln(printWriter, "▸ RESPONSE (Backend → Frontend)")
		for i := range info.Responses {
			dtoSummary(&info.Responses[i])
		}
		fmt.Fprintln(printWriter)
	}
}

func dtoSummary(dto *DTO) {
	source := "frontend"
	if dto.Source == "backend" {
		source = "backend"
	}
	fmt.Fprintf(printWriter, "  %s  (%s)\n", dto.Name, source)
	for _, f := range dto.Fields {
		opt := ""
		if f.Optional {
			opt = "?"
		}
		fmt.Fprintf(printWriter, "    %s%s  %s\n", f.Name, opt, f.Type)
	}
}

func printEndpoints(info *FeatureInfo) {
	if len(info.Endpoints) == 0 {
		return
	}

	fmt.Fprintln(printWriter, "▸ ENDPOINTS")
	for _, ep := range info.Endpoints {
		if ep.Path == "" {
			continue
		}
		resp := ep.ResponseDTO
		if resp == "" {
			resp = "–"
		}
		fmt.Fprintf(printWriter, "  %-6s %-25s → %-18s → %s\n",
			ep.Method, ep.Path, ep.Handler, resp)
	}
}


