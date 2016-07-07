### Start the server

Run `go run index.go`, which will start a server on 8080 port. If you go to [http://localhost:8080/](http://localhost:8080/)
you will see `404 page not found` which is expected, because this is a [REST application](https://en.wikipedia.org/wiki/Representational_state_transfer).

Before you can have any meaningful interaction with a server, you have to initialize a database.
Run [01_setting_up.sql](../SQL/01_setting_up.sql) and [02_populate.sql](../SQL/01_setting_up.sql) 
or just `pip install psycopg2` and run a python file  [set_up_database.py](../SQL/set_up_database.py)
which will initialize everything for you. You can install [Postico](https://eggerapps.at/postico/)
or [Navicat](https://www.navicat.com/products/navicat-for-postgresql) to view your database. 

You can interact with a server using cURL or better install a browser extension 
[PostMan](https://www.getpostman.com/). Import data from [unnamed.postman_collection](unnamed.postman_collection) file. It has 
all the routes predefined with all required parameters.

### API design

 - use nouns, not verbs (verbs are GET / POST / PUT / DELETE)
 - GET never changes the state
 - your do not need to be authorized to see majority of GET results. For everything else you need
 to be [logged in](4_jwt.md)
 - only plural nouns
 - if resource is related to another resource - use this `user/:id/purchases/` - returns all purchases for some user
 - [filtering, sorting, field selection, paging](http://blog.mwaysolutions.com/2014/06/05/10-best-practices-for-better-restful-api/)
 - return status codes properly
 - no trailing slashes, it looks like majority of the people do not use them

Some of the routes requires you to upload an image ( *update information about yourself*, 
*create a purchase*, etc). The decision of how to do this is the following.

A person in the beginning uploads an image using a separate endpoint ( *image/purchase*, 
*image/avatar*, etc). These endpoints make all the necessary checks and image manipulations and return
a name of the resulting image. A client should remember this name and pass it together with all other
data in create purchase/avatar endpoint.

This achieves a faster speed for creating of the element because while the image is uploading a person
can work on writing other information.