# golang-moviehub-api

Implemented APIs that allows consumers to access the movies data. Most of the queries will be against local database when data isn’t available in the local database then query the IMDB api and store data into local database for future refrence.

# Prerequisite
Local database: Postgresql

# Install
Gorilla mux: 
    
    go get github.com/gorilla/mux

Postgres Library

    go get github.com/lib/pq

# Run
To run write below command into terminal
   
    go run .

# Test in Postman collection
you can test all the apis from postman collcetion **golang-moviehuv-api.postman_collection.json**

