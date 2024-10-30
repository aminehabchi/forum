# forum

- Check we should not see static files to be handled ??

```
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
```

- Does the server use the right HTTP methods ??

```
/ *
* I removed it
* it's not a good way to verify methods
* http.HandleFunc("GET /", forum.Home) // bad way
*/
=> it should be handled inside function logic
```

- Are all the pages working? (Absence of 404 page?)

```
=> do not exist
```
