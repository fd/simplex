package cas

import (
	"database/sql"
)

func UpdateSchema(txn *sql.Tx) error {
	var (
		err   error
		count int64
	)

	err = txn.QueryRow(
		`SELECT COUNT(table_name) FROM information_schema.tables WHERE table_schema = $1 AND table_name = $2;`,
		"public", "cas_objects",
	).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		_, err = txn.Exec(
			`
      CREATE TABLE cas_objects (
        address  BYTEA NOT NULL,
        content  BYTEA,
        external VARCHAR,

        PRIMARY KEY (address),
        CHECK (octet_length(address) = 20),
        CHECK (content IS NOT NULL OR external IS NOT NULL)
      );
      `,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
