# Simple Note

A useless note taking app just to prove the following points:

- I can create a go web app😅
- It helps with doing something that could be useful🤭

## Features

* Create an account with email, password, profile picture
* create, update & delete simple notes(with title & body)
* Delete user account
* Login, Logout
* Middlewares to protect private routes
* View all users registered on the app

## Stack

* Go
* Html & CSS
* Bootstrap 5.1 (Uses CDN)
* Bolt DB (go key/value pair database)

## Setup

You need to have Golang installed on your machine with a minimum version of 1.13. Then clone this repo and run:

```
go get ./...
```
This will install all the required packages

Run the following command to boot up the server and start using the application on http://127.0.0.1:3000

```
go run main.go
```

## License

Licensed under [the MIT License](LICENSE)