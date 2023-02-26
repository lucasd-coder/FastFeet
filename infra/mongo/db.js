
rs.status();
db = db.getSiblingDB('user-manger-service');
db.createUser({user: 'admin', pwd: 'admin123', roles: [ { role: 'root', db: 'admin' } ]});
db.getCollection("user").createIndex(
	{ email: 1},
	{unique: true}
)

db.getCollection("user").createIndex(
	{ cpf: 1},
	{unique: true}
)

db = db.getSiblingDB('user-manger-service-test');
db.createUser(
  {user: 'test', pwd: 'test123', 
  roles:[ { 
      role: 'dbOwner', 
      db:  'user-manger-service-test'
      }
    ]
  }
);
