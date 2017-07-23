# sqlite3-cli

## Prepare database.

```console
$ sqlite3 db.sqlite3
sqlite> .header on
sqlite> .mode column
sqlite> create table tasks (
   ...>   id integer primary key,
   ...>   name text
   ...> );
sqlite> select * from tasks;
sqlite> insert into tasks(name) values('Create go-prompt-toolkit.');
sqlite> insert into tasks(name) values('Use sqlite3 from golang');
sqlite> select * from tasks;
id          name
----------  ----------------------
1           Create go-prompt-toolkit.
2           Use sqlite3 from golang
sqlite> .quit
```

