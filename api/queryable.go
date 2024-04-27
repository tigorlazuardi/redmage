package api

type Queryable interface {
	Get(string) string
}
