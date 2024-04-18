// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

func (p *PHP) includes_database_php() error {
	if !p.globals.IN_GAME {
		p.die("Access denied")
	}

	p.require_once("includes/misc.php")

	return nil
}

type DB struct {
	adapter string
}

// Supported database types:
// 'sqlite' - SQLite 3 or greater (requires PHP 5.3 or greater)
// 'mysql' - MySQL 4.1 or greater
// 'pgsql' - PostgreSQL 8.1 or greater

func (p *PHP) db_open(dbType, sock, host, port, user, pass, name string) (*DB, error) {
	if len(sock) == 0 && len(host) == 0 {
		p.die("Invalid database configuration - must specify either hostname or UNIX socket")
	} else if len(sock) != 0 && len(host) != 0 {
		p.die("Invalid database configuration - cannot specify both hostname and UNIX socket")
	}

	var db *DB
	switch dbType {
	case "sqlite":
		db = &DB{adapter: "QM_PDO_SQLITE::open($sock, $host, $port, $user, $pass, $name)"}
	case "mysql":
		db = &DB{adapter: "QM_PDO_MYSQL::open($sock, $host, $port, $user, $pass, $name)"}
	case "pgsql":
		db = &DB{adapter: "QM_PDO_PGSQL::open($sock, $host, $port, $user, $pass, $name)"}
	default:
		p.die("An unsupported database driver has been specified!")
	}

	return db, nil
}

// Attempt to lock all of the specified entities
func (p *PHP) db_lockentities(ents, owner, specials any) error {
	panic("not implemented")
}
