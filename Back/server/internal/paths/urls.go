package paths

// 1. RUTAS ESTÁNDAR (Usando el motor que creamos en builder.go)
var (
	Client  = NewResource("/client")
	Proyect = NewResource("/proyect")
	File    = NewResource("/file")
)

// 2. RUTAS ESPECIALES (Las que no siguen el patrón CRUD)
var Auth = struct {
	Prefix   string
	Login    Static
	Register Static
	Logout   Static
}{Prefix: "/auth"}

var Dashboard = struct {
	Prefix   string
	Index    Static
	Settings Static
}{Prefix: "/dashboard"}

// 3. INICIALIZADOR DE RUTAS ESPECIALES
func init() {
	// Auth
	Auth.Login = newStatic(Auth.Prefix, "/login")
	Auth.Register = newStatic(Auth.Prefix, "/register")
	Auth.Logout = newStatic(Auth.Prefix, "/logout")

	// Dashboard
	Dashboard.Index = newStatic(Dashboard.Prefix, "")
	Dashboard.Settings = newStatic(Dashboard.Prefix, "/settings")
}
