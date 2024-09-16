# chirpy
This is a simple webserver. 

It's an RESTful server that let's the clients simulate sending "chirps", (think "tweets"). 


This project was built as part of an assignment on the backend course on Boot.dev.

## Usage
### Server
```
go build
./chirtpy [--debug]
```
If the `--debug` flag is given the database will be deleted (and created) upon start.

### API overiew

#### Chrips
```
POST /api/chirps                # Create a new chirp
GET /api/chirps                 # Get a list of created chirps
GET /api/chirps/{chirpId}       # Get a specific chirp
DELETE /api/chirps/{chirpId}    # Delete a specific chirp
```
#### Users
```
POST /api/users                 # Create a user
GET /api/users                  # Get a list of users
GET /api/users/{userId}         # Get a specific user
POST /api/polka/webhooks        # Used by a payment company, to upgrade a paying customer 
PUT /api/users                  # Edit a users information
```
#### Authentication and authorization
```
POST /api/login                 # Let a user login with email + password
POST /api/refresh               # Refresh a users access token (JWT token)
POST /api/revoke                # Revoke a users refresh token
```
### Detailed instructions
#### POST /api/chirps                
Create a new chirp

##### Request header
```
Autorization: Bearer <JWT_TOKEN>
```
##### Request body
```
{
    "body": "Lorem ipsum..."
}
```

#### GET /api/chirps                 
Get a list of created chirps
```
GET /api/chirps?author_id=<USER_ID>     # Get a list of chirps created by USER_ID
GET /api/chirps?sort=<asc|desc>         # Get a list of chirps sorted by creation date

```
#### GET /api/chirps/{chirp_id}       
Get a chirp with ID `{chirp_id}`


#### DELETE /api/chirps/{chirp_id}    
Delete a chirp with ID `{chirp_id}`

##### Request header
```
Autorization: Bearer <JWT_TOKEN>
```

##### Resonse status
If the chirp is found (has a valid ID) and is successfully deleted
```
Code: 204
Text: No content
```

If the user isn't authorized
```
Code: 403
Text: Unauthorized
```

#### POST /api/users                 
Create a user.

##### Request body:
```
{
    "password":
    "email": 
}
```
##### Response body:
```
{
    "id":
    "email": 
}
```

#### GET /api/users                  
Get a list of users

##### Response
```
[
    {
        "id": 1,
        "email": "name@example.com",
        "is_chirpy_red": false
    },
    {
        "id": 2,
        "email": "name2@example.com",
        "is_chirpy_red": true
    }
]
```

#### GET /api/users/{userId}         
Get a user 1: `GET /api/users/1`

##### Response
```

{
    "id": 1,
    "email": "name@example.com",
    "is_chirpy_red": false
}
```
#### POST /api/polka/webhooks        
Used by a payment company, to upgrade a paying customer 
##### Header
```
Authorization: ApiKey <POLKA_API_KEY>
```
##### Request
```
{
    "event": "user.upgraded",
    "data": {
        "user_id": 1
    }
}
```
##### Response Status
If the user is found
```
Code: 204
Text: No Content
```

If the user isn't found:
```
Code: 404
Text: Not Found
```

If the API key doesn't match:
```
Code: 401
Text: Unauthorized
```

#### POST /api/login                 
Let a user login with email + password

##### Request
```
{
    "password": "abc123"
    "email": "name@example.com"
}
```
##### Response
```
{
    "id": 1
    "email": "name@example.com" 
    "token": <JWT_TOKEN>
    "refresh_token": <REFRESH_TOKEN>
}
```

#### POST /api/refresh               
Refresh a users access token (JWT token)

#### POST /api/revoke                
Revoke a users refresh token

#### PUT /api/users                  
Edit a users information

##### Request
```
{
    "password":
    "email": 
}
```
##### Header
```
Authorization: <JWT_TOKEN>

