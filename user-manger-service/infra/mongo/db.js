
rs.status();
db = db.getSiblingDB('fast-feet');
db.createUser({user: 'admin', pwd: 'admin123', roles: [ { role: 'root', db: 'admin' } ]});
db.getCollection("users").createIndex(
	{ email: 1},
	{unique: true}
)

db.getCollection("users").createIndex(
	{ cpf: 1},
	{unique: true}
)

db.getCollection("users").createIndex(
	{ userId: 1},
	{unique: true}
)