package cmd

func dbDriverPlugin(driverDb string) string {
	dbDriver, ok := dbDriver[driverDb]

	if !ok {
		return ""
	}

	return dbDriver
}

func getPlugins(driverDb string) []string {
	var plugins []string

	plugins = append(plugins, "github.com/joho/godotenv")
	dbPlug := dbDriverPlugin(driverDb)

	if dbPlug != "" {
		plugins = append(plugins, dbPlug)
	}

	return plugins
}
