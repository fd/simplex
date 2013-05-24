package router

import (
	"simplex.sh/store/sqlutil"
)

func (w *Writer) update_schema() error {
	var (
		present bool
		err     error
	)

	present, err = sqlutil.TableExists(w.tx, "shttp_routes")
	if err != nil {
		return err
	}

	if !present {
		_, err = w.tx.Exec(
			`
      CREATE TABLE shttp_routes (
        id           SERIAL,
        cas_key      BYTEA,
        path         VARCHAR(255),
        host         VARCHAR(128),
        content_type VARCHAR(50),
        language     VARCHAR(20),
        status       INTEGER,
        headers      BYTEA,
        address      BYTEA,

        PRIMARY KEY (id),
        CHECK (path IS NOT NULL),
        CHECK (host IS NOT NULL),
        CHECK (headers IS NOT NULL),
        UNIQUE (cas_key),
        UNIQUE (path, host, content_type, language)
      );
      `,
		)
		if err != nil {
			return err
		}
	}

	present, err = sqlutil.IndexExists(w.tx, "shttp_routes", "shttp_routes_path_idx")
	if err != nil {
		return err
	}

	if !present {
		_, err = w.tx.Exec(
			`
      CREATE
      INDEX shttp_routes_path_idx
      ON shttp_routes (path);
      `,
		)
		if err != nil {
			return err
		}
	}

	present, err = sqlutil.IndexExists(w.tx, "shttp_routes", "shttp_routes_host_idx")
	if err != nil {
		return err
	}

	if !present {
		_, err = w.tx.Exec(
			`
      CREATE
      INDEX shttp_routes_host_idx
      ON shttp_routes (host);
      `,
		)
		if err != nil {
			return err
		}
	}

	present, err = sqlutil.IndexExists(w.tx, "shttp_routes", "shttp_routes_lang_idx")
	if err != nil {
		return err
	}

	if !present {
		_, err = w.tx.Exec(
			`
      CREATE
      INDEX shttp_routes_lang_idx
      ON shttp_routes (language);
      `,
		)
		if err != nil {
			return err
		}
	}

	present, err = sqlutil.IndexExists(w.tx, "shttp_routes", "shttp_routes_type_idx")
	if err != nil {
		return err
	}

	if !present {
		_, err = w.tx.Exec(
			`
      CREATE
      INDEX shttp_routes_type_idx
      ON shttp_routes (content_type);
      `,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
