package db

import "server/configure"

func Init() {
	defaultFilesystem.onStart(configure.DefaultConfigure().Logger.Filename)
}
