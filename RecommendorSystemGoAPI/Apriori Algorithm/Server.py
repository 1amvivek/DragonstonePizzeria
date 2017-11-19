from pymongo import MongoClient
from apriori import runApriori
import json

client1 = MongoClient("192.168.99.100:27018")
db = client1.cmpe281
coll = db.recommend
cursor = db.recommend.find()

client2 = MongoClient("192.168.99.100:27017")
db2 = client2.cmpe281


minSupport = 3
transactions = []

for document in cursor:
    transactions.append(document['cart'].split(","))

# print(transactions)


items, rules = runApriori(transactions, 0.10, 0.68)

print(rules[0])

rule = {}
for rule in rules:
	for item in rule:
		if type(item) is tuple:
			cart = ",".join(item[0])
			recommend = ",".join(item[1])
			rule = { "cart" : cart , 
							"recommend" : recommend }

			# print(rule)
			if db.rules.find(rule).count() > 0:
				print("document already exists") 
			else:				
				result = db.rules.insert_one(rule)
				print(result.inserted_id)
				result = db2.rules.insert_one(rule)
			break