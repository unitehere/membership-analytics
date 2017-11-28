from os import environ

from flask import Flask, request, jsonify
from elasticsearch import Elasticsearch
from flask_cors import CORS
from werkzeug.utils import ImportStringError

from query import SearchClient

app = Flask(__name__)
CORS(app)
app.config.from_object('config.%s' % environ.get('FLASK_ENV', 'development'))
try:
    app.config.from_object('config.local')
except ImportStringError:
    pass
app.es = Elasticsearch(app.config['ELASTICSEARCH_HOST'])


@app.route('/health')
def health_check():
    return 'Service is up and running.'


@app.route('/search/<term>')
def search(term):
    search = SearchClient(app)
    search.set_from(request.args.get('from', ''))
    search.set_size(request.args.get('size', ''))
    search_func = getattr(search, term)
    search_func(request.args.get('q', ''))
    response = search.execute()
    return jsonify(response.hits.hits)

if __name__ == '__main__':
    app.run(host='0.0.0.0')
