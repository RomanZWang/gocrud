# GO AWS DynamoDB CRUD Server
## Usage
### Running the server
```
go build && ./crud
```
### DAO Unit tests
```
go test
```
### Sending API requests
| Transaction Type | Example                                                                                                |
|------------------|--------------------------------------------------------------------------------------------------------|
| CREATE           | curl -d '{"id":"jim", "Number": 100, "AssociatedTypes": ["hello","there"]}' http://localhost:8420/data |
| READ             | [browser link](http://localhost:8420/data/jim)                                                                         |
| UPDATE           | curl -X PUT -d '{"id":"jim", "Number": 99, "AssociatedTypes": ["hello"]}'  http://localhost:8420/data            |
| DELETE           | curl -X DELETE http://localhost:8420/data/jim                                                          |