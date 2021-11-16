package ports

type IdGenerator interface {
	New() (string, error)
}
