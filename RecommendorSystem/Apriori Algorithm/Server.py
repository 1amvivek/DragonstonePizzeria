from pymongo import MongoClient
from apriori import runApriori
import json

client = MongoClient("192.168.99.100:27017")
db = client.cmpe281
coll = db.recommend
cursor = db.recommend.find()


minSupport = 3
transactions = []

for document in cursor:
    transactions.append(document['cart'])

# print(transactions)


items, rules = runApriori(transactions, 0.10, 0.68)

print(rules)

rule = {}
for rule in rules:
	for item in rule:
		if type(item) is tuple:
			rule = { "cart" : [ current for current in item[0] ], 
							"recommend" : [ current for current in item[1] ] }

			if db.rules.find(rule).count() > 0:
				print("document already exists") 
			else:				
				result = db.rules.insert_one(rule)
				print(result.inserted_id)
			break