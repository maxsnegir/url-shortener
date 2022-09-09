package storages

type URLDataBase map[string]string

func (db URLDataBase) Set(key string, value string) error {
	db[key] = value
	return nil
}

func (db URLDataBase) Get(key string) (string, error) {
	value, ok := db[key]
	if !ok {
		return "", KeyError
	}
	return value, nil
}

func NewURLDataBase() URLDataBase {
	return make(URLDataBase)
}
