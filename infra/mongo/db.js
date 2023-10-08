
rs.status();
db.getSiblingDB('fast-feet').createUser({user: 'admin', pwd: 'admin123', roles: [ { role: 'root', db: 'admin' } ]});

db.getSiblingDB('fast-feet').getCollection("users").createIndex(
	{ email: 1},
	{unique: true}
)

db.getSiblingDB('fast-feet').getCollection("users").createIndex(
	{ cpf: 1},
	{unique: true}
)

db.getSiblingDB('fast-feet').getCollection("users").createIndex(
	{ userId: 1},
	{unique: true}
)

db.getSiblingDB('fast-feet').getCollection("orders").createIndex(
	{ deliverymanId: 1}
)

db.getSiblingDB('fast-feet').getCollection("orders").createIndex(
	{ deliverymanId: 1, _id: 1},
	{unique: true}
)