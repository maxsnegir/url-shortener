package storage

// MapURLDataBase In-Memory хранилище
type MapURLDataBase map[string]string

func (db MapURLDataBase) Set(key string, value string) error {
	db[key] = value
	return nil
}

func (db MapURLDataBase) Get(key string) (string, error) {
	value, ok := db[key]
	if !ok {
		return "", KeyError
	}
	return value, nil
}

func (db MapURLDataBase) Shutdown() error {
	return nil
}

func NewMapURLDataBase() Storage {
	return make(MapURLDataBase)
}
