var params = {
  TableName: 'likes',
  KeySchema: [ // The type of of schema.  Must start with a HASH type, with an optional second RANGE.
    { // Required HASH type attribute
      AttributeName: 'id',
      KeyType: 'HASH',
    },
    { // Required HASH type attribute
      AttributeName: 'created',
      KeyType: 'RANGE',
    }
  ],
  AttributeDefinitions: [ // The names and types of all primary and index key attributes only
    {
      AttributeName: 'id',
      AttributeType: 'S', // (S | N | B) for string, number, binary
    },
    {
      AttributeName: 'created',
      AttributeType: 'N', // (S | N | B) for string, number, binary
    }
  ],
  ProvisionedThroughput: { // required provisioned throughput for the table
    ReadCapacityUnits: 1,
    WriteCapacityUnits: 1,
  }
};

dynamodb.createTable(params, function(err, data) {
  if (err) ppJson(err); // an error occurred
  else ppJson(data); // successful response
});

