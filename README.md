# forum

- https://medium.com/tag/programming
```
yep in the filter when You change it to uppercase still retrieve the data.
```

- ParseFiles():
```
-> making it as it was it better to show which template err's
-> just a little update on code to make it more simpler and readable
```

- CreateDB():
```
Example: I changed ON DELETE to N DELETE in posts query:
before => 
go run .
near "N": syntax error
after => 
go run .
failed to create posts table: near "N": syntax error

==> at least we know where error occurs exactly

-> PRIMARY KEY added for post_interactions and comment_interactions.
```

- HomeHandler:
```
- GetPosts:
-> in home: we get all posts
-> in commenting: we get specific post using post_id
-> in posting: we get all posts but based on my offset
-> in filter:
-> I think in all my queries I have a baseQuery that exists in all of them
baseQuery := `
    SELECT posts.id, posts.user_id, posts.title, posts.created_at, posts.content, users.uname 
    FROM posts
    JOIN users ON posts.user_id = users.id`

-> to be reviewed that it doesn't miss anything.
```


