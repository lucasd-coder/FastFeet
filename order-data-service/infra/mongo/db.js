
rs.status();
db = db.getSiblingDB('fast-feet');
db.createUser({user: 'admin', pwd: 'admin123', roles: [ { role: 'root', db: 'admin' } ]});
db.getCollection("orders").createIndex(
	{ deliverymanId: 1}
)

db.getCollection("orders").createIndex(
	{ deliverymanId: 1, _id: 1},
	{unique: true}
)