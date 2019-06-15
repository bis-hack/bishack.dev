var params = {
  TableName: 'posts',
  KeySchema: [ // The type of of schema.  Must start with a HASH type, with an optional second RANGE.
    { // Required HASH type attribute
      AttributeName: 'id',
      KeyType: 'HASH',
    },
    { // Optional RANGE key type for HASH + RANGE tables
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
      AttributeName: 'username',
      AttributeType: 'S', // (S | N | B) for string, number, binary
    },
    {
      AttributeName: 'created',
      AttributeType: 'N', // (S | N | B) for string, number, binary
    },
  ],
  ProvisionedThroughput: { // required provisioned throughput for the table
    ReadCapacityUnits: 1,
    WriteCapacityUnits: 1,
  },
  GlobalSecondaryIndexes: [ // optional (list of GlobalSecondaryIndex)
    {
      IndexName: 'username',
      KeySchema: [
        { // Required HASH type attribute
          AttributeName: 'username',
          KeyType: 'HASH',
        },
        { // Optional RANGE key type for HASH + RANGE secondary indexes
          AttributeName: 'created',
          KeyType: 'RANGE',
        }
      ],
      Projection: { // attributes to project into the index
        ProjectionType: 'INCLUDE', // (ALL | KEYS_ONLY | INCLUDE)
        NonKeyAttributes: [ // required / allowed only for INCLUDE
          'attribute_name_1',
          // ... more attribute names ...
        ],
      },
      ProvisionedThroughput: { // throughput to provision to the index
        ReadCapacityUnits: 1,
        WriteCapacityUnits: 1,
      },
    },
    // ... more global secondary indexes ...
  ],
};

dynamodb.createTable(params, function(err, data) {
  if (err) ppJson(err); // an error occurred
  else ppJson(data); // successful response
});
