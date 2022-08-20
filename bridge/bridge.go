package bridge

type Provider struct {
}

type Bridge struct {
	Provoders map[string]Provider
}

func New() *Bridge {
	b := &Bridge{}

	return b
}
