# HTTP API written in Go using Echo framework
## Running

To run the project, run the following command:

```
docker-compose up --build
```

You should be able to access the API at http://localhost:1323

If you change `database.sql` file, you need to reinitate the database by running:

```
docker-compose down --volumes
```

## Testing

To run test, run the following command:

```
make test
```
