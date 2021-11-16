package ports

type UrlCounter interface {
	IncrementCounter(id string)
}
