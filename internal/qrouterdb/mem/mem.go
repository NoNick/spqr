package mem

import (
	"sync"

	"github.com/pg-sharding/spqr/internal/qrouterdb"
	"github.com/pg-sharding/spqr/yacc/spqrparser"
	"golang.org/x/xerrors"
)

type QrouterDBMem struct {
	mu   sync.Mutex
	txmu sync.Mutex

	freq map[string]int
	krs  map[string]*spqrparser.KeyRange
}

func (q *QrouterDBMem) Begin() error {
	q.txmu.Lock()

	return nil
}

func (q *QrouterDBMem) Commit() error {
	q.txmu.Unlock()
	return nil
}

func (q *QrouterDBMem) Add(keyRange *spqrparser.KeyRange) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if _, ok := q.krs[keyRange.KeyRangeID]; ok {
		return xerrors.Errorf("key range %v already present in qdb", keyRange.KeyRangeID)
	}

	q.freq[keyRange.KeyRangeID] = 1
	q.krs[keyRange.KeyRangeID] = keyRange

	return nil
}

func (q *QrouterDBMem) Update(keyRange *spqrparser.KeyRange) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.krs[keyRange.KeyRangeID] = keyRange

	return nil
}

func (q *QrouterDBMem) Check(key int) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, kr := range q.krs {
		if kr.From <= key && key <= kr.To {
			return true
		}
	}

	return false
}

func NewQrouterDBMem() (*QrouterDBMem, error) {
	return &QrouterDBMem{
		freq: map[string]int{},
		krs:  map[string]*spqrparser.KeyRange{},
	}, nil
}

func (q *QrouterDBMem) Lock(keyRange *spqrparser.KeyRange) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if cnt, ok := q.freq[keyRange.KeyRangeID]; ok {
		q.freq[keyRange.KeyRangeID] = cnt + 1
	} else {
		q.freq[keyRange.KeyRangeID] = 1
	}

	q.krs[keyRange.KeyRangeID] = keyRange

	return nil
}

func (q *QrouterDBMem) UnLock(keyRange *spqrparser.KeyRange) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if cnt, ok := q.freq[keyRange.KeyRangeID]; !ok {
		return xerrors.Errorf("key range %v not locked", keyRange)
	} else if cnt > 1 {
		q.freq[keyRange.KeyRangeID] = cnt - 1
	} else {
		delete(q.freq, keyRange.KeyRangeID)
		delete(q.krs, keyRange.KeyRangeID)
	}

	return nil
}

var _ qrouterdb.QrouterDB = &QrouterDBMem{}
