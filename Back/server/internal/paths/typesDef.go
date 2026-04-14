package paths

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

// --- STRUCTS BASE ---

type Static struct {
	Route string
	url   string
}

// devuelve la ruta en crudo (solo texto sin parámetros)
func (s Static) URL() string { return s.url }

// devuelve la ruta con parámetros de consulta añadidos
func (s Static) URLWithQuery(q map[string]string) string { return buildQuery(s.url, q) }

// --- STRUCT PARA RUTAS DINÁMICAS ---
type Dynamic[T any] struct {
	Route   string
	Builder func(T) string
}

// devuelve la ruta con el parámetro insertado
func (d Dynamic[T]) URL(param T) string { return d.Builder(param) }

// devuelve la ruta con el parámetro insertado y parámetros de consulta añadidos
func (d Dynamic[T]) URLWithQuery(param T, q map[string]string) string {
	return buildQuery(d.Builder(param), q)
}

// --- STRUCT PARA RUTAS ESTÁNDAR ---

type StandardResource struct {
	Prefix string
	Index  Static
	New    Static
	Edit   Dynamic[int64]
	Delete Dynamic[int64]
	Search Static
}

// --- FUNCIONES CONSTRUCTORAS ---

func newStatic(prefix, route string) Static {
	return Static{
		Route: route,
		url:   path.Join(prefix, route),
	}
}

func NewResource(prefix string) StandardResource {
	return StandardResource{
		Prefix: prefix,
		Index:  newStatic(prefix, ""),
		New:    newStatic(prefix, "/new"),
		Edit: Dynamic[int64]{
			Route:   "/edit/:id",
			Builder: func(id int64) string { return path.Join(prefix, "edit", fmt.Sprintf("%d", id)) },
		},
		Delete: Dynamic[int64]{
			Route:   "/:id",
			Builder: func(id int64) string { return path.Join(prefix, fmt.Sprintf("%d", id)) },
		},
		Search: newStatic(prefix, "/search"),
	}
}

func buildQuery(baseURL string, query map[string]string) string {
	if len(query) == 0 {
		return baseURL
	}
	vals := url.Values{}
	for k, v := range query {
		vals.Set(k, v)
	}
	queryString := vals.Encode()
	separator := "?"
	if strings.Contains(baseURL, "?") {
		separator = "&"
	}
	return fmt.Sprintf("%s%s%s", baseURL, separator, queryString)
}

//if you feel confused, here's the routes in basic string format:
/*
// standard CRUD routes
"/client"           // Index
"/client/new"       // New
"/client/edit/:id"  // Edit
"/client/:id"       // Delete

"/proyect"          // Index
"/proyect/new"      // New
"/proyect/edit/:id" // Edit
"/proyect/:id"      // Delete

"/file"             // Index
"/file/new"         // New
"/file/edit/:id"    // Edit
"/file/:id"         // Delete

// Specials routes
"/login"
"/register"
"/logout"
"/overview"         // Dashboard Index
"/overview/settings" // Dashboard Settings
*/
