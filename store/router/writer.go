package router

import (
	"database/sql"
	"encoding/json"
	"sync"
)

type Writer struct {
	tx   *sql.Tx
	stmt *sql.Stmt
	err  error
	keys map[string]bool
	mtx  sync.RWMutex
}

func OpenWriter(tx *sql.Tx) (*Writer, error) {
	w := &Writer{
		tx:   tx,
		keys: make(map[string]bool, 1000),
	}

	err := w.update_schema()
	if err != nil {
		return nil, err
	}

	{
		stmt, err := tx.Prepare(`
      INSERT
      INTO shttp_routes
        (cas_key, path, host, content_type, language, status, headers, address)
      VALUES
        ($1, $2, $3, $4, $5, $6, $7, $8);
    `)
		if err != nil {
			return nil, err
		}

		w.stmt = stmt
	}

	w.load_keys()

	return w, nil
}

func (w *Writer) Insert(rule *Rule) error {
	w.mtx.RLock()
	defer w.mtx.RUnlock()

	if w.err != nil {
		return w.err
	}

	rule.calculate_key()

	{ // don't write existing rules
		key_str := string(rule.Key)
		if _, p := w.keys[key_str]; p {
			w.keys[key_str] = true
			return nil
		}
		w.keys[key_str] = true
	}

	headers, err := json.Marshal(rule.Header)
	if err != nil {
		return err
	}

	_, err = w.stmt.Exec(
		[]byte(rule.Key),
		rule.Path,
		rule.Host,
		rule.ContentType,
		rule.Language,
		rule.Status,
		headers,
		[]byte(rule.Address),
	)
	return err
}

func (w *Writer) load_keys() {
	w.mtx.Lock()
	defer w.mtx.Unlock()

	rows, err := w.tx.Query(`SELECT cas_key FROM shttp_routes;`)
	if err != nil {
		w.err = err
		return
	}

	defer rows.Close()

	for rows.Next() {
		var (
			addr []byte
		)

		err := rows.Scan(&addr)
		if err != nil {
			w.err = err
			return
		}

		w.keys[string(addr)] = false
	}

	if err := rows.Err(); err != nil {
		w.err = err
		return
	}
}
