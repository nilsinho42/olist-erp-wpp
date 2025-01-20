from flask import Flask, request

app = Flask(__name__)

@app.route('/oauth/callback')
def oauth_callback():
    auth_code = request.args.get('code')
    if auth_code:
        return f"Authorization code: {auth_code}", 200
    return "Authorization code not found", 400

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
