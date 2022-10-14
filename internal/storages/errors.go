package storages

const KeyError = DBKeyError("Key does not exist")

type DBKeyError string

func (e DBKeyError) Error() string {
	return string(e)
}
