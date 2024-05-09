# otp-service

## About this service

This service is written in `Golang` and leverages`RabbmitMQ` as a message broker to consume events produced by auth-service and produce its own events. Also, it makes of `Twilio` to send randomized generated 6 digit OTP sms to user's phone number.


## Configure service dependencies

We first need to configure below required components:

### RabbitMQ
This service also uses RabbitMQ to produce and consume service events to communicate with auth-service. As an example, below command can be used to run RabbitMQ container. *Make sure you have Docker installed locally*  

```sh
docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.13-management
```  

`NOTE: If you have already setup RabbitMQ container, then make sure both auth-service and otp-service are connected to same container instance.`  

### Twilio
We also need a Twilio account to send sms. Twilio provides a free trial account with 15$ credits to use. Please visit the URL: https://www.twilio.com/docs/messaging/guides/how-to-use-your-free-trial-account to setup your own trial account. Then from the Twilio Dashboard console, get the details of your :
- `Account SID`
- `Auth Token`
- `Public phone number` provided by Twilio (*This may require further setup from the dashboard.*)


### .env  
Open the .env file in the root of the `otp-service` folder and enter below required RabbitMQ and Twilio config details.  

`AMQP_ADDRESS` - This is RabbitMQ local running URL containing its host, user/password and port.  

`TWILIO_ACCOUNT_SID` - Get it from Twilio dashboard  

`TWILIO_AUTH_TOKEN` - Get it from Twilio dashboard

`TWILIO_PHONE_NUMBER` - Get it from Twilio dashboard

```sh
AMQP_ADDRESS=amqp://guest:guest@localhost:5672/
TWILIO_ACCOUNT_SID=
TWILIO_AUTH_TOKEN=
TWILIO_PHONE_NUMBER=
```

## Run the service 
To run the service we need to install Go dependencies i.e., third-party packages used. CD into the root of the project directory. And run below commands sequentially:

**Install dependencies**
```sh
go mod download
```
**Then, run :**
```sh
go run main.go
```  

If everything setup correctly, you will see something like this:
![alt text](image.png)

