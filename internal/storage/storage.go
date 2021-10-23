package storage

type Storage struct {
	db    DataBase
	queue Queue
	cache Cache
}

func New(db DataBase, queue Queue, cache Cache) *Storage {
	return &Storage{
		db:    db,
		queue: queue,
		cache: cache,
	}
}

func (s *Storage) DB() DataBase {
	return s.db
}

func (s *Storage) Q() Queue {
	return s.queue
}

func (s *Storage) C() Cache {
	return s.cache
}
