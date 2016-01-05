# influxdbUrl

This is a helper project to let frontend code connect to backend influxdb databse.
By default, this service is bind to port 18080

follow the instructions here to install go on your environment:
https://golang.org/doc/install



build all go file in influxdbUrl project
```
go install influxdbUrl/
```

Copy and paste the result generated from encryptionGenerator project to credential.config, in the following format:
```
u=${YOUR_ENCRYPTED_USERNAME}
p=${YOUR_ENCRYPTED_PASSWORD}
l=${YOUR_ENCRYPTED_URL}
d=${YOUR_ENCRYPTED_DATABASE_NAME}
```

Use the following command to start the backend service
```
go run influxdbUrl/influxDBUrl.go
```

Note: because this program use current working path for config file, it can only work correctly using the above relative path. Do not run it from a different path

http syntax:
```
curl -H "Content-Type: application/json" -X POST -d "{\"timeStart\": \"2015-11-03\",\"timeEnd\": \"2018-11-04\",\"limit\": 1000,\"metric\": \"uptime_ms_cumulative\" }"  -i  http://localhost:18080/influxdbUrl
```

On the server side, it will print the logs of sql command constructed.
On the client side, it will print out the http response, it there are any.


After the above config, the last step is to change the file permmision of credential.config to read only
```
chmod 444 credential.config
```

To make the program run as backend jobs
```
nohup $GOPATH/bin/influxdbUrl &
```