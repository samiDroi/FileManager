package apidto

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type parsedRoute struct {
	Method  string
	Path    string
	Handler string
}

func parseGoDTOs(filePath string) ([]DTO, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, data, parser.AllErrors)
	if err != nil {
		return nil, fmt.Errorf("parse error in %s: %w", filePath, err)
	}

	var dtos []DTO
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			ttype, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			st, ok := ttype.Type.(*ast.StructType)
			if !ok {
				continue
			}

			dto := DTO{
				Name: ttype.Name.Name,
			}

			for _, field := range st.Fields.List {
				if len(field.Names) == 0 {
					continue
				}
				for _, name := range field.Names {
					fi := Field{
						Name: name.Name,
						Type: goExprString(field.Type),
					}

					if field.Tag != nil {
						raw := field.Tag.Value
						raw = strings.Trim(raw, "`")
						tag := reflect.StructTag(raw)
						if jsonTag := tag.Get("json"); jsonTag != "" {
							parts := strings.Split(jsonTag, ",")
							if parts[0] != "" {
								fi.Name = parts[0]
							}
							for _, p := range parts[1:] {
								if p == "omitempty" {
									fi.Optional = true
								}
							}
						}
						if formTag := tag.Get("form"); formTag != "" {
							parts := strings.Split(formTag, ",")
							if parts[0] != "" {
								fi.Name = parts[0]
							}
							for _, p := range parts[1:] {
								if p == "omitempty" {
									fi.Optional = true
								}
							}
						}
					}

					dto.Fields = append(dto.Fields, fi)
				}
			}

			if len(dto.Fields) > 0 {
				dtos = append(dtos, dto)
			}
		}
	}
	return dtos, nil
}

func goExprString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + goExprString(t.X)
	case *ast.SelectorExpr:
		return goExprString(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		if t.Len == nil {
			return "[]" + goExprString(t.Elt)
		}
		return "[" + goExprString(t.Len) + "]" + goExprString(t.Elt)
	case *ast.MapType:
		return "map[" + goExprString(t.Key) + "]" + goExprString(t.Value)
	case *ast.BasicLit:
		return t.Value
	case *ast.InterfaceType:
		return "any"
	default:
		return fmt.Sprintf("%T", expr)
	}
}

func parseGoRoutes(filePath string) ([]parsedRoute, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, data, parser.AllErrors)
	if err != nil {
		return nil, fmt.Errorf("parse error in %s: %w", filePath, err)
	}

	var routes []parsedRoute
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "RegisterRoutes" {
			continue
		}

		ast.Inspect(fn.Body, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			method := sel.Sel.Name
			switch method {
			case "GET", "POST", "PUT", "DELETE", "PATCH":
			default:
				return true
			}

			if len(call.Args) < 2 {
				return true
			}

			path := exprRoutePath(call.Args[0])
			handler := exprHandlerName(call.Args[1])

			if path != "" && handler != "" {
				routes = append(routes, parsedRoute{
					Method:  method,
					Path:    path,
					Handler: handler,
				})
			}
			return true
		})
	}
	return routes, nil
}

func exprRoutePath(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.BasicLit:
		if t.Kind == token.STRING {
			return strings.Trim(t.Value, `"`)
		}
	case *ast.SelectorExpr:
		return exprRoutePath(t.X) + "." + t.Sel.Name
	case *ast.Ident:
		return t.Name
	case *ast.CallExpr:
		return "fn()"
	case *ast.BinaryExpr:
		left := exprRoutePath(t.X)
		right := exprRoutePath(t.Y)
		if left != "" && right != "" {
			return left + " + " + right
		}
	}
	return ""
}

func exprHandlerName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.SelectorExpr:
		return t.Sel.Name
	case *ast.Ident:
		return t.Name
	}
	return ""
}

func findDTOs(projectRoot string, feature Feature) (requestDTO, responseDTO *DTO, err error) {
	backendDir := filepath.Join(projectRoot, "Back", "server", "internal", "domain", "dtos")

	var backendFiles []string

	switch feature {
	case FeatureAuth:
		authDTOFile := filepath.Join(projectRoot, "Back", "server", "internal", "http", "handlers", "auth", "dto.go")
		backendFiles = append(backendFiles, authDTOFile)
	default:
		dtoFileMap := map[Feature]string{
			FeatureClient:   "client",
			FeatureProject:  "project",
			FeatureFile:     "file",
			FeatureTypeFile: "fileType",
		}
		if fn, ok := dtoFileMap[feature]; ok {
			backendFiles = append(backendFiles, filepath.Join(backendDir, fn+".go"))
		}
	}

	for _, f := range backendFiles {
		if _, ferr := os.Stat(f); ferr != nil {
			continue
		}
		dtos, derr := parseGoDTOs(f)
		if derr != nil {
			continue
		}
		for i, d := range dtos {
			if strings.Contains(d.Name, "Request") || strings.Contains(d.Name, "Update") {
				dtos[i].Source = "backend"
				dtos[i].Role = "request"
				requestDTO = &dtos[i]
			} else if strings.Contains(d.Name, "Response") {
				dtos[i].Source = "backend"
				dtos[i].Role = "response"
				responseDTO = &dtos[i]
			}
		}
	}

	return requestDTO, responseDTO, nil
}

func findEndpoints(projectRoot string, feature Feature) ([]Endpoint, error) {
	backendDir := filepath.Join(projectRoot, "Back", "server", "internal", "http", "handlers")

	handlerMap := map[Feature]string{
		FeatureAuth:     filepath.Join(backendDir, "auth"),
		FeatureClient:   filepath.Join(backendDir, "clients"),
		FeatureProject:  filepath.Join(backendDir, "projects"),
		FeatureFile:     filepath.Join(backendDir, "files"),
		FeatureTypeFile: backendDir,
	}

	handlerFileMap := map[Feature]string{
		FeatureAuth:     "authHandler.go",
		FeatureClient:   "client.go",
		FeatureProject:  "project.go",
		FeatureFile:     "file.go",
		FeatureTypeFile: "typeFile.go",
	}

	handlerDir, ok := handlerMap[feature]
	if !ok {
		return nil, nil
	}

	handlerFile, ok := handlerFileMap[feature]
	if !ok {
		return nil, nil
	}

	routeFile := filepath.Join(handlerDir, handlerFile)
	routes, err := parseGoRoutes(routeFile)
	if err != nil {
		return nil, err
	}

	prefixMap := map[Feature]string{
		FeatureAuth:     "",
		FeatureClient:   "/client",
		FeatureProject:  "/project",
		FeatureFile:     "/files",
		FeatureTypeFile: "/type-files",
	}

	routePathMap := map[Feature]map[string]string{
		FeatureClient: {
			"Index.Route":  "",
			"New.Route":    "/new",
			"Edit.Route":   "/edit/:id",
			"Delete.Route": "/:id",
			"Search.Route": "/search",
		},
		FeatureProject: {
			"Index.Route":  "",
			"New.Route":    "/new",
			"Update.Route": "/edit/:id",
			"Search.Route": "/search",
		},
		FeatureFile: {
			"Index.Route":          "",
			"GetByIds.Route":       "/project/:id",
			"Upload.Route":         "/upload",
			"Delete.Route":         "/:id",
			"Search.Route":         "/search",
			"Download.Route":       "/:id/download",
			"Edit.Route":           "/edit/:id",
		},
		FeatureTypeFile: {
			"Index.Route": "",
			"New.Route":   "/new",
		},
		FeatureAuth: {
			"Login.Route":      "/login",
			"Register.Route":   "/register",
			"Logout.Route":     "/logout",
			"RefreshTkn.Route": "/refresh",
		},
	}

	prefix := prefixMap[feature]
	pathMap := routePathMap[feature]

	var endpoints []Endpoint
	for _, r := range routes {
		p := resolvePath(r.Path, prefix, pathMap)
		endpoints = append(endpoints, Endpoint{
			Method:  r.Method,
			Path:    p,
			Handler: r.Handler,
		})
	}

	return endpoints, nil
}

func resolvePath(raw string, prefix string, known map[string]string) string {
	if known != nil {
		if p, ok := known[raw]; ok {
			return prefix + p
		}
	}

	if strings.HasPrefix(raw, "/") {
		return raw
	}

	if strings.Contains(raw, ".") {
		for k, v := range known {
			if strings.Contains(raw, k) {
				return prefix + v
			}
		}
	}

	if raw != "" {
		return prefix + "/" + raw
	}
	return prefix
}
