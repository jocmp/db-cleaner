# DB Cleaner

A proof-of-concept database cleaner based on <https://github.com/DatabaseCleaner/database_cleaner>.

The `CleanDatabase` function reads the information catalog schema from SQL Server and
filters down to tables only, excluding views.

It iterates through each result and runs the `DELETE FROM` command on each table.

## Takeaways

This assumes full access to the database. It does not disable integrity constraints. This in
turn would disallow unsorted tables from getting deleted without first disabling constraints using
`NO CHECK`.

A faster way to execute this could be to either use the "transaction" strategy, which is to run the
entire test suite within a single database transaction or to `TRUNCATE` instead of deleting.

Finally, the cleaner reveals an interesting factor about the speed of Go development in general.
A Ruby on Rails developer expects a lot of magic in exchange for speed.
In Go, a tool like this could easily turn milliseconds-long tests into 4 or 5 second test runs.
Sometimes that trade-off is worth it, but other times tests can be refactored instead.

## Getting started

Download Go, connect to a database, and `make run`.


