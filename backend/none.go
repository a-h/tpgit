package backend

type None struct {
}

func (n None) GetLease() (id string, err error) {
	return "", nil
}

func (n None) ExtendLease(id string) (ok bool, err error) {
	return true, nil
}

func (n None) CancelLease() (err error) {
	return nil
}

func (n None) IsProcessed(hash string) (bool, error) {
	return false, nil
}

func (n None) MarkProcessed(hash string) error {
	return nil
}
