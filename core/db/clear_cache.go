package db

func ClearBeans(tablename string) {
	engine := Engine()
	cacher := engine.GetCacher(tablename)
	if cacher != nil {
		cacher.ClearBeans(tablename)
	}
}
func ClearIds(tablename string) {
	engine := Engine()
	cacher := engine.GetCacher(tablename)
	if cacher != nil {
		cacher.ClearIds(tablename)
	}
}
func ClearCache(tablename string) {
	engine := Engine()
	cacher := engine.GetCacher(tablename)
	if cacher != nil {
		cacher.ClearIds(tablename)
		cacher.ClearBeans(tablename)
	}
}
