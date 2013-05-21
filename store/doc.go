package store

/*

- The store is backed by Postgres.


Tables:
- environments (id, name)
- objects      (id, environments_id, collection, value)
- cas_objects  (address, value, external)
- shttp_routes (environments_id, path, host, content_type, language, headers, address)

*/
