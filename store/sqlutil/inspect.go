package sqlutil

import (
	"database/sql"
)

func IndexExists(txn *sql.Tx, table, index string) (bool, error) {
	var (
		err   error
		count int64
	)

	err = txn.QueryRow(
		`
      SELECT
        COUNT(*)
      FROM
        pg_catalog.pg_indexes
      WHERE schemaname = $1
        AND tablename = $2
        AND indexname = $3;
    `,
		"public", table, index,
	).Scan(&count)
	if err != nil {
		return false, err
	}

	return (count == 1), nil
}

func TableExists(txn *sql.Tx, table string) (bool, error) {
	var (
		err   error
		count int64
	)

	err = txn.QueryRow(
		`
      SELECT
        COUNT(*)
      FROM
        pg_catalog.pg_tables
      WHERE schemaname = $1
        AND tablename = $2;
    `,
		"public", table,
	).Scan(&count)
	if err != nil {
		return false, err
	}

	return (count == 1), nil
}
