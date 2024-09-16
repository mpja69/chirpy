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

#### Authentication and similar
``` 
POST /api/login                 # Let a user login with email + password
POST /api/refresh               # Refresh a users access token (JWT token)
POST /api/revoke                # Revoke a users refresh token
PUT /api/users                  # Edit a users information
```

### Detailed instructions
#### POST /api/chirps                
Create a new chirp


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

#### POST /api/users                 
Create a user.

Request body:
```
{
    "password":
    "email": 
}
```
Response body:
```
{
    "id":
    "email": 
}
```

#### GET /api/users                  
Get a list of users

#### GET /api/users/{userId}         
Get a specific user

#### POST /api/polka/webhooks        
Used by a payment company, to upgrade a paying customer 

#### POST /api/login                 
Let a user login with email + password

Request body:
```
{
    "password":
    "email": 
}
```
Response body:
```
{
    "id":
    "email": 
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

Request body:
```
{
    "password":
    "email": 
}
```
Header:
```
Authorization: <JWT_TOKEN>

