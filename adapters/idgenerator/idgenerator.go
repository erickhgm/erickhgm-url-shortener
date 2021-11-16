package idgenerator

import (
	"ehgm.com.br/url-shortener/domain/ports"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

// Struct that implements 'IdGenerator' interface
type idGenerator struct {
	idLength int
}

// Get an instance of 'IdGenerator' using this method
func NewIdGenerator(idLength int) ports.IdGenerator {
	return &idGenerator{idLength: idLength}
}

func (g *idGenerator) New() (string, error) {
	return gonanoid.New(g.idLength)
}
