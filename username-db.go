package username

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQL driver
)

// DB username database object
type DB struct {
	*sql.DB
}

// InitDB initalizs the db connection object
func InitDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// initTestDB initalizes a testDb where the transactions are all localized.
func initTestDB(dataSourceName string, identifier string) (*DB, error) {
	db, err := sql.Open("txdb", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// User models the users table
type User struct {
	ID       int
	Username string
}

// updateUsernames updates a usernames with a provided making of id's to new usernames
func (db *DB) updateUsernames(usernameChanges *map[int]string) error {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("UPDATE users SET username=? WHERE id=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for id, username := range *usernameChanges {
		_, err = stmt.Exec(username, id)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()

	return nil
}

// ResolveUsernameCollisions resolves username collisions by adding a number to the end
// as per the example
func (db *DB) ResolveUsernameCollisions(dryrunOptional ...bool) error {
	dryrun := false
	if len(dryrunOptional) > 0 {
		dryrun = dryrunOptional[0]
	}

	users, err := db.getUsersWithCollisions()
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return nil
	}

	// Users will have been sorted by the previous sql call
	currentDisallowedName := users[0].Username
	counter := 1
	resolvedUsers := make(map[int]string)
	for _, user := range users {
		if user.Username != currentDisallowedName {
			currentDisallowedName = user.Username
			counter = 1
		}
		resolvedUsers[user.ID] = user.Username + string(counter)
		counter++
	}

	if dryrun {
		for id, name := range resolvedUsers {
			fmt.Printf("Id: %d ResolvedName: %s\n", id, name)
		}
		return nil
	}

	err = db.updateUsernames(&resolvedUsers)
	if err != nil {
		return err
	}
	return db.ResolveUsernameCollisions(false) // dryrun must be false at this point
}

// ResolveUsernameCollisionsUUID resolves username collisions by adding a number to the end
func (db *DB) ResolveUsernameCollisionsUUID(dryrunOptional ...bool) error {
	dryrun := false
	if len(dryrunOptional) > 0 {
		dryrun = dryrunOptional[0]
	}

	users, err := db.getUsersWithCollisions()
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return nil
	}

	// Users will have been sorted by the previous sql call
	// but doesn't matter for UUID
	currentDisallowedName := users[0].Username
	resolvedUsers := make(map[int]string)
	for _, user := range users {
		if user.Username != currentDisallowedName {
			currentDisallowedName = user.Username
		}
		resolvedUsers[user.ID] = fmt.Sprintf("%s-%s", user.Username, randStringBytesMaskImprSrc(20))
	}

	if dryrun {
		for id, name := range resolvedUsers {
			fmt.Printf("Id: %d ResolvedName: %s\n", id, name)
		}
		return nil
	}

	err = db.updateUsernames(&resolvedUsers)
	if err != nil {
		return err
	}
	return nil
}

// getUsersWithCollisions finds all users with collisions and returns their id's and usernames
func (db *DB) getUsersWithCollisions() ([]*User, error) {
	rows, err := db.Query(`
	SELECT a.*
	FROM users a
	JOIN (
		SELECT username, COUNT(*)
		FROM users 
		GROUP BY username
		HAVING count(*) > 1) b
	ON a.username = b.username
	ORDER BY a.username`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		var u User
		if err = rows.Scan(&u.ID, &u.Username); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}

// ResolveDisallowedUsernames finds all users with disallowed usernames and renames them
func (db *DB) ResolveDisallowedUsernames(dryrunOptional ...bool) error {
	dryrun := false
	if len(dryrunOptional) > 0 {
		dryrun = dryrunOptional[0]
	}

	users, err := db.GetUsersWithDisallowedUsernames()
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return nil
	}

	// Users will have been sorted by the previous sql call
	currentDisallowedName := users[0].Username
	counter := 1
	resolvedUsers := make(map[int]string)
	for _, user := range users {
		if user.Username != currentDisallowedName {
			currentDisallowedName = user.Username
			counter = 1
		}
		resolvedUsers[user.ID] = user.Username + string(counter)
		counter++
	}

	if dryrun {
		for id, name := range resolvedUsers {
			fmt.Printf("Id: %d ResolvedName: %s\n", id, name)
		}
	} else {
		err := db.updateUsernames(&resolvedUsers)
		if err != nil {
			return err
		}
	}

	return nil
}

//GetUsersWithDisallowedUsernames finds all users with usernames found in the disallowed_usernames table
func (db *DB) GetUsersWithDisallowedUsernames() ([]*User, error) {
	rows, err := db.Query(`
		SELECT * FROM users
		WHERE username IN
			(SELECT invalid_username FROM disallowed_usernames)
		ORDER BY username ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		var u User
		if err = rows.Scan(&u.ID, &u.Username); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}

type usernameCount struct {
	Username string
	Count    int
}

// GetUsernameCollisions retrieves the duplicated usernames and their counts. Used for testing
func (db *DB) getUsernameCollisions() ([]*usernameCount, error) {
	rows, err := db.Query("SELECT username, COUNT(*) count FROM users GROUP BY username HAVING count > 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := []*usernameCount{}
	for rows.Next() {
		var c usernameCount
		if err = rows.Scan(&c.Username, &c.Count); err != nil {
			return nil, err
		}
		counts = append(counts, &c)
	}

	return counts, nil
}
