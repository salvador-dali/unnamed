Service uses [JWT](https://jwt.io/) for authentication. Because a lot of APIs requires you to be 
authenticated, it makes sense to install a [jwt plugin](https://chrome.google.com/webstore/detail/jwt-debugger/ppmmlchacdbknfphdeafcbmklcghghmd/related)
 
Frontend developer does not really need to know a lot about JWT. It is sufficient to know that 
this is a base64 encoded string, that should not be shared with anyone (something like a cookie).

You get this string when a person log in to a site and have to destroy it when a user logs out. If
a user is logged in, you have to send this string in the header with every his request.

### Brief description of JWT

JWT consists of three base64 strings separated by a dot. Something like `HEADER.PAYLOAD.SIGNATURE`.
It is used to transmit some information securely over the unsecure channel. The JWT has all the
information needed to unpack and verify it's correctness.

 - `HEADER` is a json object consisting of only two keys: `{"typ": "JWT", "alg": "hs256"}`. `Typ` is always
 the same, but a sender can select any algorithm to encrypt his token.
 
 - `PAYLOAD` is a json object which can consist of any user-defined fields. An application can send 
 any information there. There are also a couple of reserved claims (which are standard).
 
 - `SIGNATURE` is the thing which makes all of this secure. It is an `HMACSHA256` encoding of the 
 header, payload and the secret (only servers knows the secret). So it looks like this: 
 `HMACSHA256(base64(HEADER) + "." + base64(PAYLOAD) + "." + secret)`.
 
A secret allows a sender to encrypt the message and makes sure that the client can't tamper the 
information.

### Advantages of JWT:

 - small (100 - 200 bytes)
 - self-contained (everything is there, no need to search anything in the database)
 - no need to worry about [CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/Access_control_CORS)
 
 
### Frontend usage
During the login process a client provides username and password and the server returns either a 
JWT token or some error message. It is up to a client to store this JWT token (in local storage or 
something similar) and to transmit this token at every next request. Client should transmit it
in a `token` header (`curl -X POST -H "token: youJwtToken" ...`).

If a user logs out, client should delete a token.

A list of valid tokens for every user with a very far-in-the-future expiration date (more about 
expiration date later)

 - 1 `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk2MzMyMTEsImlhdCI6MTQ2NjgzMzIxMSwiaWQiOjF9.LYey3jgBd70QYjygbZvoPqXGJHj90nZ8VUm2yeVlVVo`
 - 2 `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk2MzMyMTEsImlhdCI6MTQ2NjgzMzIxMSwiaWQiOjJ9.dbBN08ZNdGhKbPhFRSccRWvMgSxSTjlM3wC7K2oz3_M`
 - 3 `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk2MzMyMTEsImlhdCI6MTQ2NjgzMzIxMSwiaWQiOjN9.WF7GGKA2XB3Th5lztqseW1fixf9XApTYpwDhcvq_sDw`
 - 4 `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk2MzMyMTEsImlhdCI6MTQ2NjgzMzIxMSwiaWQiOjR9.cI2Ie6KDVQhWk1VRuP_UzE1HpKFfyT0jgTe9J2g7pJA`
 - 5 `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk2MzMyMTEsImlhdCI6MTQ2NjgzMzIxMSwiaWQiOjV9.bmmgOyeN700onUcVfJcFT4dn5XyNY7rdUfpYDhlfdOc`
 - 6 `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk2MzMyMTEsImlhdCI6MTQ2NjgzMzIxMSwiaWQiOjZ9.vqP4oem2PeQpzBBC2enSXYrKg2xDcPa8iXcJToSmWHs`
 - 7 `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk2MzMyMTEsImlhdCI6MTQ2NjgzMzIxMSwiaWQiOjd9.KlfEaHwqWLMGVA9MUIu_z8oSNaXbioJ6_mgftlbWpeI`
 - 8 `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk2MzMyMTEsImlhdCI6MTQ2NjgzMzIxMSwiaWQiOjh9.0u293Hl2-cJawLI1JlEcE1fYBB6yrkMvKUiGHy61-2A`
 - 9 `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk2MzMyMTEsImlhdCI6MTQ2NjgzMzIxMSwiaWQiOjl9.JS9Xc135ndkunTa2oKess5KCX4WVCcvAkI7bVsV4YVo`

The second part is the base64 encoded payload of the JWT. In this case it is `{"exp":1639633211,"iat":1466833211,"id":3}`.
If the current time is bigger than the expiration time, the token will be invalidated on the server.
`PROJ_JWT_EXP_DAYS` is a variable that is responsible for number of days the token is valid for.
Right now it is 2 days. A client can ask to extend a token (call `GET /users/login/extend` with
an old token in the header). This allows to create a new token for a user in the current one.

It is up to a client to decide when to ask for an extension. Reasonable heuristics is to call it 
when one half of the time is left (`currentTime > (exp + iat) / 2`).
