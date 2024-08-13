package main

import (
	"tent/constant"
	"tent/installer"
	"tent/logger"
)

func main() {
	logger.Tent("Tent", "v"+constant.VERSION)
	build := installer.Selector()
	if build == "" {
		logger.Error("No build selected")
		return
	}

	installer.Install(build)
}
