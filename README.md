This service is to synchronize data with each other as container level

set data using API 
```
curl --location --request GET 'http://localhost:9001/set' \
--header 'Content-Type: application/json' \
--data '{
    "data": "sakib"
}'
```

get data using API
```
curl --location 'http://localhost:9002/get/adata'
```

delete data using API
```
curl --location 'http://localhost:9003/delete/key'
```


to start container 
```
make docker-start
```

to stop container 
```
make docker-stop
```
