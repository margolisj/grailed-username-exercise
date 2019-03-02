# Grailed Username Exercise

## Notes
I'd design a system providing some uuid type identifier to do profile linking (to alleviate the disallowed and collision issue) such that looking up a user involves going to grailed.com/user/<user-id>. Similarly, the colliding username issue could be solve on account creation.

I made disallowed username and collision user name functions like the example. I also created a collision username function with UUID as sort of discussed in my email. Also, the disallowed username function could introduce collisions. I chose not to check that and keep the actual functionality of the two functions separate. The only other thing to note is the dryrun functionality in ResolveUsernameCollisions doesn't run recursively, and thus there might collisions in the dry run information. I allowed this since this shouldn't be the case with the random string method.

Golang doesn't have true optional parameters. One of the strategies, which was used here, is using variadic parameters to handle 0 to any parameters.

## Tech usage
I've written a few personal apps in Golang but have done a good amount of work with it recently. I haven't had the chance to use it in a professional setting. I've used SQLite a bunch. I've done some SQL with go but I normally use something like gorm but that wasn't worth using due to 1) the Db already being built and 2) wanting to do as much of the lifting in actual SQL as possible.

## Running Instructions
I included the database for sake of simplicity.
```
go build
go test
```

## TODO:
1) Write a function that finds all users with disallowed usernames. Disallowed usernames can be found in the `disallowed_usernames` table.
2) Write a function that resolves all username collisions. E.g., two users with the username `foo` should become `foo` and `foo1`. The function accepts an optional "dry run" argument that will print the affected rows to the console, not commit the changes to the db.
3) Write a function that resolves all disallowed usernames. E.g., `grailed` becomes `grailed1`. The function accepts an optional "dry run" argument that will print the affected rows to the console, not commit the changes to the db.

## Libraries
* https://github.com/mattn/go-sqlite3
* https://github.com/DATA-DOG/go-txdb
* https://godoc.org/github.com/stretchr/testify/assert
