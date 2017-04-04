package backend

type InMemory struct {
	hashes map[string]interface{}
}

func NewInMemory() *InMemory {
	return &InMemory{
		hashes: make(map[string]interface{}),
	}
}

func (n *InMemory) GetLease() (id string, err error) {
	return "", nil
}

func (n *InMemory) ExtendLease(id string) (ok bool, err error) {
	return true, nil
}

func (n *InMemory) CancelLease() (err error) {
	return nil
}

func (n *InMemory) IsProcessed(hash string) (bool, error) {
	_, ok := n.hashes[hash]
	return ok, nil
}

func (n *InMemory) MarkProcessed(hash string) error {
	n.hashes[hash] = true
	return nil
}
