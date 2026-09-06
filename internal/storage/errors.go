package storage

import "errors"

var ErrNotFound = errors.New("not found")

var ErrAliasExists = errors.New("alias already exists")
