from database.db import db

class Implant(db.DynamicDocument):
    implant_id = db.StringField(required=True, unique=True)
    

class Task(db.DynamicDocument):
    task_id = db.StringField(required=True, unique=True)
    implant_id = db.StringField(required=True)
    
class Result(db.DynamicDocument):
    task_id = db.StringField(required=True, unique=True)
    implant_id = db.StringField(required=True)
    
