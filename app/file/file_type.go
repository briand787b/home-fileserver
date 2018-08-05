package file

import "github.com/pkg/errors"

const (
	typeMovie = iota
	typeSong
	typeBook
	typeGame
)

var validTypes []Type

func init() {

}

type Type struct {
	Code int
	Name string `json:"Type"`
	Dir  string `json:"Dir"`
}

func (t *Type) IsValid() {

}

func (t *Type) GetDir() string {

}

func (t *Type) SetDir(d string) error {
	return errors.New("not implemented")
}
