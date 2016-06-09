Tried a couple of routers:

 - [HttpRouter](https://github.com/julienschmidt/httprouter) - does not support 
 [normal routes](https://github.com/julienschmidt/httprouter/issues/12)
 - [Denco](https://github.com/naoina/denco) - does not support DELETE, PUT

The one that supports the route structure and all http verbs is 
[httptreemux](https://github.com/dimfeld/httptreemux)
 
## API design

 - use nouns, not verbs (verbs are GET/POST/PUT/DELETE)
 - GET never changes the state
 - only plural nouns
 - if resource is related to another resource - use this `user/:id/purchases/` - returns all purchases for some user
 - [filtering, sorting, field selection, paging](http://blog.mwaysolutions.com/2014/06/05/10-best-practices-for-better-restful-api/)
 - return status codes properly
 
Root url: `/api/`

 - GET    `/brands`                show all brands
 - GET    `/brands/:id`            show a particular brand
 - POST   `/brands`                create a new brand
 - PUT    `/brands/:id`            update a particular brand
 - GET    `/tags`                  show all tags
 - GET    `/tags/:id`              show a particular tag
 - POST   `/tags`                  create a new tag
 - PUT    `/tags/:id`              update a particular tag
 - POST   `/users`                 create a new user
 - POST   `/users/login`           log in current user
 - POST   `/users/logout`          log out current user
 - GET    `/users/:id`             show a particular user
 - GET    `/users/me/email/:hash`  verify email address after registration
 - PUT    `/users/me/info`         update user info
 - PUT    `/users/me/avatar`       update your avatar
 - POST   `/users/me/follow/:id`   start following user
 - DELETE `/users/me/follow/:id`   stop following user
 - GET    `/users/:id/followers`   who follows this user
 - GET    `/users/:id/following`   who this user follow
 - GET    `/users/:id/purchases`   get all the purchases of a particular user
 - GET    `/purchases`             get all the purchases
 - GET    `/purchases/tag/:id`     get all the purchases with a particular tag
 - GET    `/purchases/brand/:id`   get all the purchases with a particular brand
 - POST   `/purchases`             create a new purchase
 - GET    `/purchases/:id`         get a purchase with a particular id
 - POST   `/purchases/:id/like`    like a purchase with a particular id
 - DELETE `/purchases/:id/like`    unlike a purchase with a particular id
 - POST   `/purchases/:id/ask`     ask a question about a particular purchase
 - POST   `/questions/:id/vote`    upvote a question
 - DELETE `/questions/:id/vote`    downvote a question
 - POST   `/questions/:id/answer`  answer a question
 - POST   `/answer/:id/vote`       upvote an answer
 - DELETE `/answer/:id/vote`       downvote an answer
 
How to [setup Go with Pycharm](http://stackoverflow.com/a/37698196/1090562)

How to organize Go code:
 - https://talks.golang.org/2014/organizeio.slide#9
 - https://medium.com/@benbjohnson/structuring-applications-in-go-3b04be4ff091#.p2aokrleg
 - https://github.com/otoolep/go-httpd
 - http://darian.af/post/the-anatomy-of-a-golang-project/
 - http://dave.cheney.net/2014/12/01/five-suggestions-for-setting-up-a-go-project