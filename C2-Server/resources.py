import uuid
import json

from flask import request, Response
from flask_restful import Resource
from database.db import initialize_db
from database.models import Task, Result
from flask_jwt_extended import jwt_required, get_jwt_identity


class Tasks(Resource):
    # ListTasks
    @jwt_required()
    def get(self):
        implant_id = get_jwt_identity()

        if implant_id:
            Tasks = Task.objects(implant_id=str(implant_id)).to_json()
            return Response(Tasks, mimetype="application/json", status=200)
        else:
            return Response("Unauthorized", mimetype="application/json", status=401)

    # AddTasks
    @jwt_required()
    def post(self):
        print("POST")
        implant_id = get_jwt_identity()
        
        if not implant_id:
            return Response("Unauthorized", mimetype="application/json", status=401)
        body = request.get_json()
        obj_nbr = len(body)
        for i in range(len(body)):
            body[i]['task_id'] = str(uuid.uuid4())
            body[i]['implant_id'] = str(implant_id)
            Task(**body[i]).save()
            
            
        return Response(Task.objects.skip(Task.objects.count() - obj_nbr).to_json(), mimetype="application/json", status=201)
        
        
class Results(Resource):
    @jwt_required()
    def get(self):
        implant_id = get_jwt_identity()
        if implant_id:
            res = Result.objects(implant_id=str(implant_id)).to_json()
            return Response(res, mimetype="application/json", status=200)
        else:
            return Response("Unauthorized", mimetype="application/json", status=401)
    
    @jwt_required()
    def post(self):
        implant_id = get_jwt_identity()
        body = request.get_json()
        print(body['task_id'])
        
        if not implant_id:
            return Response("Unauthorized", mimetype="application/json", status=401)
        if 'task_id' not in body:
            return Response("Task ID not provided", mimetype="application/json", status=400)
        body['implant_id'] = str(implant_id)
        task = Task.objects(task_id=body['task_id']).first()
        if not task:
            return Response("Task ID not found", mimetype="application/json", status=404)
        body['task_obj'] = task.to_json()
        task.delete()
        
        res = Result(**body).save()
        
        return Response(res.to_json(), mimetype="application/json", status=200)
        
