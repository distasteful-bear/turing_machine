from pprint import pprint

from google.cloud import firestore

db = firestore.Client(project="james-metz", database="turing-machine")

db.collection("users").document("user1").set({"name": "John Doe"})

user = db.collection("users").document("user1").get()
pprint(user.to_dict())
