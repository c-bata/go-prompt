from bottle import route, run, request

@route('/')
def hello():
    return "Hello World!"

@route('/ping')
def hello():
    return "pong!"

@route('/register', method='POST')
def register():
    name = request.json.get("name")
    return "Hello %s!" % name

if __name__ == "__main__":
    run(host='localhost', port=8000, debug=True)
