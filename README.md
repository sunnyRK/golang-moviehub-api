# golang-moviehub-api

Implemented APIs that allows consumers to access the movies data. Most of the queries will be against local database when data isn’t available in the local database then query the IMDB api and store data into local database for future refrence.

# Prerequisite
- Go
- Postgresql

# Install
- install all the Go dependencies by running the following command

     go get -d ./...

# Run
To run write below command into terminal
   
    go run .

# Test in Postman collection
you can test all the apis from postman collcetion **golang-moviehuv-api.postman_collection.json**

