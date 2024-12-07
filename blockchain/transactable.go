package blockchain

type Transactable interface {
	run() (*string, error)
	print()
}
