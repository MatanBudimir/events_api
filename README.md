Event management API

How to build and run the API?
1. First copy and adjust contents of `.env.example` to `.env` file
2. Run `docker-compose build`
3. Run `docker-compose up -d`
4. SSH into the database container called `db` and follow the next steps:
    - mysql -u root -p
    - enter the password
    - run `use api`
    - copy and paste the content of api.sql file
5. To generate a JWT token run `go run scripts/main.go jwt`
6. After that's done send requests using Postman or any other client of choice
7. URLs and example request bodies for Postman are available in [postman_example.md](https://github.com/MatanBudimir/events_api/blob/main/postman_example.md)

How to run the integration test?
1. Run the following command `cd pkg/api/handlers/meetings`
2. Run `go test`
