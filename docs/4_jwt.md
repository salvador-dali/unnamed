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

 - 1 `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzkyOTI4MzEsImlkIjoxfQ.E3KRJgFfpKHgexw13grm9-neaXrlb7sLjk5Q9XsBeRY`
 - 2 `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk1NTU2NzAsImlkIjoyfQ.-o8iN6TXLqeyUR8bkJ3WCfDr7527BZ9aHY12qCfOCvE`
 - 3 `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk1NTU3MTksImlkIjozfQ.Agi-2KpwE-J8B4wUwOz5n-5mcg8P9cUF9qqCwsL2USI`
 - 4 `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk1NTU3NDMsImlkIjo0fQ.ceGmymRfiO2sv-WV-_7z63FePcdZ36wrQmugHtyI94g`
 - 5 `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk1NTU3NzQsImlkIjo1fQ.FMx5hJQ-KdV1lCrOhP_UrKXhKvY1DfNeDzsnO2wlGwI`
 - 6 `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk1NTU3ODUsImlkIjo2fQ.sTQ9HMqrpaP1R6tl7mgrCPjbr52-qWpensYB2IsoaNo`
 - 7 `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk1NTU4MDAsImlkIjo3fQ.DhJpM75XmrvJet37OhEff0jN3ZBrpoBMbUoSOaCaqTM`
 - 8 `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk1NTU4MTEsImlkIjo4fQ.vF0Vo_Mpha7FcYhu7BraRfJqsn8hMBednlFGTMumAhk`
 - 9 `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Mzk1NTU4MjIsImlkIjo5fQ.huTzZZ2ToM1wflgT42oirBRwnyTZbtAJZw6hm6-aJck`


