package username

import (
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/stretchr/testify/assert"
)

var testDbLocation = "./db/grailed-exercise.sqlite3"

func init() {
	txdb.Register("txdb", "sqlite3", testDbLocation)
}

func TestSelectUsers(t *testing.T) {
	db, err := initTestDB(testDbLocation, "TestSelectUsers")
	defer db.Close()
	if err != nil {
		t.Error(err)
	}
	row := db.QueryRow("SELECT * FROM users WHERE id=1")
	var user User
	err = row.Scan(&user.ID, &user.Username)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "wade.corkery4", user.Username)
}

func TestUpdateUsernames(t *testing.T) {
	db, err := initTestDB(testDbLocation, "TestUpdateUsernames")
	defer db.Close()
	if err != nil {
		t.Error(err)
	}
	testUsername := "testUser"
	updates := map[int]string{
		1: testUsername,
	}
	err = db.updateUsernames(&updates)
	if err != nil {
		t.Error(err)
	}

	row := db.QueryRow("SELECT * FROM users WHERE id=1")
	var user User
	err = row.Scan(&user.ID, &user.Username)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, testUsername, user.Username)
}
func TestGetUsernameCollisions(t *testing.T) {
	db, err := initTestDB(testDbLocation, "TestGetUsernameCollisions")
	defer db.Close()
	if err != nil {
		t.Error(err)
	}

	users, err := db.getUsernameCollisions()
	if err != nil {
		t.Error(err)
	}
	assert.Len(t, users, 336)
	assert.Equal(t, usernameCount{"abdiel1", 2}, *users[0])
}

func TestResolveUsernameCollisions(t *testing.T) {
	db, err := initTestDB(testDbLocation, "TestResolveUsernameCollisions")
	defer db.Close()
	if err != nil {
		t.Error(err)
	}
	err = db.ResolveUsernameCollisions()
	if err != nil {
		t.Error(err)
	}

	users, err := db.getUsernameCollisions()
	if err != nil {
		t.Error(err)
	}
	assert.Len(t, users, 0)
}

func TestResolveUsernameCollisionsDryRun(t *testing.T) {
	db, err := initTestDB(testDbLocation, "TestResolveUsernameCollisionsDryRun")
	defer db.Close()
	if err != nil {
		t.Error(err)
	}
	err = db.ResolveUsernameCollisions(true)
	if err != nil {
		t.Error(err)
	}

	// Username collisions should still exist
	users, err := db.getUsernameCollisions()
	if err != nil {
		t.Error(err)
	}
	assert.Len(t, users, 336)
	assert.Equal(t, usernameCount{"abdiel1", 2}, *users[0])
}

func TestResolveUsernameCollisionsUUID(t *testing.T) {
	db, err := initTestDB(testDbLocation, "TestResolveUsernameCollisionsUUID")
	defer db.Close()
	if err != nil {
		t.Error(err)
	}
	err = db.ResolveUsernameCollisionsUUID()
	if err != nil {
		t.Error(err)
	}

	users, err := db.getUsernameCollisions()
	if err != nil {
		t.Error(err)
	}
	assert.Len(t, users, 0)
}

func TestResolveDisallowedUsernames(t *testing.T) {
	db, err := initTestDB(testDbLocation, "TestResolveDisallowedUsernames")
	defer db.Close()
	if err != nil {
		t.Error(err)
	}
	err = db.ResolveDisallowedUsernames(false)
	if err != nil {
		t.Error(err)
	}

	users, err := db.GetUsersWithDisallowedUsernames()
	if err != nil {
		t.Error(err)
	}
	assert.Len(t, users, 0)
}
func TestResolveDisallowedUsernamesDryRun(t *testing.T) {
	db, err := initTestDB(testDbLocation, "TestResolveDisallowedUsernamesDryRun")
	defer db.Close()
	if err != nil {
		t.Error(err)
	}
	err = db.ResolveDisallowedUsernames(true)
	if err != nil {
		t.Error(err)
	}

	users, err := db.GetUsersWithDisallowedUsernames()
	if err != nil {
		t.Error(err)
	}
	assert.Len(t, users, 25)
}

func TestGetUsersWithDisallowedUsernames(t *testing.T) {
	db, err := initTestDB(testDbLocation, "TestGetUsersWithDisallowedUsernames")
	defer db.Close()
	if err != nil {
		t.Error(err)
	}

	users, err := db.GetUsersWithDisallowedUsernames()
	if err != nil {
		t.Error(err)
	}
	assert.Len(t, users, 25)
}
