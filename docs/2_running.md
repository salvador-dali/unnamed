### Start the server

Run `go run index.go`, which will start a server on 8080 port. If you go to [http://localhost:8080/](http://localhost:8080/)
you will see `404 page not found` which is expected, because this is a [REST application](https://en.wikipedia.org/wiki/Representational_state_transfer).

Before you can have any meaningful interaction with a server, you have to initialize a database.
Run `01_setting_up.sql` and `02_populate.sql` from an `SQL` folder or just run a python file 
`set_up_database.py` which will initialize everything for you. You can install [Postico](https://eggerapps.at/postico/)
or [Navicat](https://www.navicat.com/products/navicat-for-postgresql) to view your database. 

You can interact with a server using cURL or better install a browser extension 
[PostMan](https://www.getpostman.com/). Import data from `unnamed.postman_collection` file. It has 
all the routes predefined with all required parameters.

### API design

 - use nouns, not verbs (verbs are GET / POST / PUT / DELETE)
 - GET never changes the state
 - only plural nouns
 - if resource is related to another resource - use this `user/:id/purchases/` - returns all purchases for some user
 - [filtering, sorting, field selection, paging](http://blog.mwaysolutions.com/2014/06/05/10-best-practices-for-better-restful-api/)
 - return status codes properly
 - no trailing slashes, it looks like majority of the people do not use them

