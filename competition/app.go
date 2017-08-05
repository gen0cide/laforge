package competition

type AppMap map[string]*App

type App struct {
	Name string
}

func GetApps() AppMap {
	return make(map[string]*App)
}
