package cache

type Cache interface {
	Order() Order
	Check() error
	Close() error
}
