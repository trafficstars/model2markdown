package model2markdown

import (
	"fmt"
)

type ErrConflict struct {
	a libraryStruct
	b libraryStruct
}

func NewErrConflict(a, b libraryStruct) ErrConflict {
	return ErrConflict{
		a,
		b,
	}
}

func (err ErrConflict) Error() string {
	return fmt.Sprintf(`conflicting structures: %v:%v and %v:%v`,
		err.a.File.Path, err.a.Struct.Name, err.b.File.Path, err.b.Struct.Name)
}