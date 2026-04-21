package paths

// 1. RUTAS ESTÁNDAR (Usando el motor que creamos en builder.go)
var (
	Client  = NewResource("/client")
	Project = NewResource("/project")
	File    = NewResource("/file")
)

// 2. RUTAS ESPECIALES (Las que no siguen el patrón CRUD)
var Auth = struct {
	Prefix     string
	Login      Static
	Register   Static
	Logout     Static
	RefreshTkn Static
}{Prefix: ""}

var Dashboard = struct {
	Prefix   string
	Index    Static
	Settings Static
}{Prefix: "/overview"}

// 3. INICIALIZADOR DE RUTAS ESPECIALES
func init() {
	// Auth
	Auth.Login = newStatic(Auth.Prefix, "/login")
	Auth.Register = newStatic(Auth.Prefix, "/register")
	Auth.Logout = newStatic(Auth.Prefix, "/logout")
	Auth.RefreshTkn = newStatic(Auth.Prefix, "/refresh")
	// Dashboard
	Dashboard.Index = newStatic(Dashboard.Prefix, "")
	Dashboard.Settings = newStatic(Dashboard.Prefix, "/settings")

}
