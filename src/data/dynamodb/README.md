To get information about local DynamoDB installation:

<https://docs.aws.amazon.com/pt_br/amazondynamodb/latest/developerguide/DynamoDBLocal.DownloadingAndRunning.html>

Start local DynamoDB:\

    >java -Djava.library.path=./DynamoDBLocal_lib -jar DynamoDBLocal.jar -sharedDb &

Configure local aws environment(you must install aws client before): 

    >aws configure

    AWS Access Key ID [****************DYK7]: fakeMyKeyId
    AWS Secret Access Key [****************eKHo]: fakeSecretAccessKey
    Default region name [eu-west-1]: local
    Default output format [None]: 

Create the tables to store transactions:

    >aws dynamodb create-table \
        --table-name EncryptedTransaction \
        --attribute-definitions \
            AttributeName=Hash,AttributeType=B \
        --key-schema AttributeName=Hash,KeyType=HASH\
        --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1
    
    >aws dynamodb create-table \
        --table-name EncryptedRawTransaction \
        --attribute-definitions \
            AttributeName=Hash,AttributeType=B \
        --key-schema AttributeName=Hash,KeyType=HASH\
        --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1
    
    >aws dynamodb create-table \
        --table-name Peer \
        --attribute-definitions \
            AttributeName=URL,AttributeType=S \
        --key-schema AttributeName=URL,KeyType=HASH\
        --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1