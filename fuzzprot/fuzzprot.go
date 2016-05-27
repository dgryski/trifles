package fuzzprot

import (
	"errors"
	"strconv"
)

type user struct {
	Type string
	Name string
	Age  int
}

func Unpack(data []byte) ([]user, error) {

	count := data[0]
	users := make([]user, count)

	data = data[1:]

	var uidx int
	for len(data) > 0 {
		switch data[0] {
		case 0:
			uidx++
			data = data[1:]
		case 1:
			users[uidx].Type, data = grabString(data[1:])
		case 2:
			users[uidx].Name, data = grabString(data[1:])
		case 3:
			var err error
			data = data[1:]
			users[uidx].Age, err = strconv.Atoi(string(data[:2]))
			if err != nil {
				return nil, err
			}
			data = data[2:]
		default:
			return nil, errors.New("unknown field type")
		}
	}

	return users, nil
}

func grabString(data []byte) (string, []byte) {
	l, data := data[0], data[1:]
	return string(data[:l]), data[l:]
}
