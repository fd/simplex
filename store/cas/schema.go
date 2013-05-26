package cas

import (
	"database/sql"
	"simplex.sh/store/sqlutil"
)

func update_schema(db *sql.DB) error {
	var (
		present bool
		err     error
	)

	present, err = sqlutil.TableExists(db, "cas_objects")
	if err != nil {
		return err
	}

	if !present {
		_, err = db.Exec(
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

	present, err = sqlutil.IndexExists(db, "cas_objects", "cas_objects_addr_idx")
	if err != nil {
		return err
	}

	if !present {
		_, err = db.Exec(
			`
      CREATE
      INDEX cas_objects_addr_idx
      ON cas_objects (address);
      `,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
