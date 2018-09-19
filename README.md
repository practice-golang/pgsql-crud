# PostgreSQL CRUD
Practice for PostgreSQL

## File, Directories
* src - Sources
* sql - Syntax practice
* pgdata.zip - pgsql data folder. you can set this folder and can use with root. (nothing password)

## Packages using
https://github.com/lib/pq

When At least one variable or more variables in dbinfo have empty, I met below error.  
In my case, I entered empty password.  
So, each values(=%s) sould be wrapped single quotes(='').  

Int which wrapped single quotes would show error.  
I don't know that the reason is caused from pgsql or lib/pq
```go
dbinfo := fmt.Sprintf(
    "host=%s port='%s' user='%s' password='%s' dbname='%s' sslmode='disable'",
    dbHost, dbPort, dbUser, dbPassword, dbName)
db, err := sql.Open("postgres", dbinfo)
chkErr(err)
defer db.Close()
```
```sh
panic: pq: database "user" does not exist
```

## Reference
* http://freeprog.tistory.com/248
* https://github.com/filewalkwithme/go-pg-crud
* https://astaxie.gitbooks.io/build-web-application-with-golang/en/05.4.html
