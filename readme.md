# Running the Project

To run this project, follow these steps:

1. Open a terminal and navigate to the project directory.
2. Run the following command to start the server:

    go run cmd/main.go

3. Open a web browser and navigate to [http://localhost:8080](http://localhost:8080).

Note: Make sure you have Go installed on your machine and the GOPATH is set correctly.

# Setup MongoDB

If you don't have MongoDB installed, please follow these steps:

1. Go to the MongoDB Community Server download page and download the correct version for your operating system.
2. Follow the installation instructions for your operating system.
3. Once installed, start the MongoDB server by running `mongod` in a terminal.
4. MongoURI should be `mongodb://localhost:27017`
5. Database should be called `backend-challenge` and `backend-challenge-test` (for now)

# TODO

* Write unit tests for edge cases of the product service
* Write unit tests for the brain service
* 