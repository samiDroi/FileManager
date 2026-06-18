package paths

import "path"

// 1. RUTAS ESTÁNDAR (Usando el motor que creamos en builder.go)
var (
	Client   = NewResource("/client")
	Project  = NewResource("/project")
	File     = NewResource("/files")
	TypeFile = NewResource("/type-files")
	User     = NewResource("/users")
)

var FileSpecials = struct {
	GetByIds Dynamic[[]int]
}{}

// 2. RUTAS ESPECIALES (Las que no siguen el patrón CRUD)
var Auth = struct {
	Prefix      string
	Login       Static
	Register    Static
	Logout      Static
	RefreshTkn  Static
	ForgotPwd   Static
	RecoveryPwd Static
	changePwd   Static
}{Prefix: ""}
var ClientsSidebar = struct {
	Prefix string
	Index  Static
}{
	Prefix: "/clients-projects",
	Index:  newStatic("/clients-projects", ""),
}
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
	Auth.ForgotPwd = newStatic(Auth.Prefix, "/forgot-password")
	Auth.RecoveryPwd = newStatic(Auth.Prefix, "/recovery-password")
	Auth.changePwd = newStatic(Auth.Prefix, "/change-password")
	// Dashboard
	Dashboard.Index = newStatic(Dashboard.Prefix, "")
	Dashboard.Settings = newStatic(Dashboard.Prefix, "/settings")
	// File Specials
	FileSpecials.GetByIds = Dynamic[[]int]{
		Route:   "/project/:id",
		Builder: func(ids []int) string { return path.Join(File.Prefix, "/project/:id") },
	}
}
