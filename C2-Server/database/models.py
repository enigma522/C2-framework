from database.db import db

class Implant(db.DynamicDocument):
    implant_id = db.StringField(required=True, unique=True)
    online = db.BooleanField(required=True,default=True)
    

class Task(db.DynamicDocument):
    task_id = db.StringField(required=True, unique=True)
    implant_id = db.StringField(required=True)
    
class Result(db.DynamicDocument):
    task_id = db.StringField(required=True, unique=True)
    implant_id = db.StringField(required=True)
    
class Profile(db.DynamicDocument):
    username = db.StringField(required=True, unique=True)
    password = db.StringField(required=True)
    
