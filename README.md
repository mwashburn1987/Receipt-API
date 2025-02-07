# Receipt-API
Simple Backend Receipt Checker


## Running With Docker
This assumes that you already have docker installed
Note: You may need to pull the latest golang image from docker. 
You can do this with ```docker pull golang:latest```

You can build the docker image with ```docker build -t receipt-api .```

You can then run this binary in a docker container with ```docker run -p 8080:8080 receipt-api``` 

## Running without Docker
You can run locally without docker by navigating to the root directory and
running ```go build```

After the build is complete you can run the application with ```./receipt-api```

## Testing the receipt API

Post Request: 
```
curl --request POST \
  --url http://localhost:8080/receipts/process \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/10.3.0' \
  --data '{
	"retailer": "M&M Corner Market",
	"purchaseDate": "2022-03-20",
	"purchaseTime": "14:33",
	"items": [
		{
			"shortDescription": "Gatorade",
			"price": "2.25"
		},
		{
			"shortDescription": "Gatorade",
			"price": "2.25"
		},
		{
			"shortDescription": "Gatorade",
			"price": "2.25"
		},
		{
			"shortDescription": "Gatorade",
			"price": "2.25"
		}
	],
	"total": "9.00"
}'
```

You can copy and run that curl directly in a terminal or using REST software such as postman/insomnia.

You should see a successful return with a 200 status code and a return body in the format of ```{"id": "1"}```

GET Request: 
You can then run a get request at 
```
curl --request GET \
  --url http://localhost:8080/receipts/1/point \
  --header 'User-Agent: insomnia/10.3.0'
  ```
  Just be sure to insert your receipt number into the url request.