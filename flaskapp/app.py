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
    return 'OK'


@app.route('/search/<term>')
def search(term):
    search = SearchClient(app)
    search_func = getattr(search, term)
    search_func(request.args.get('q', ''))
    app.logger.debug(search.debug())
    response = search.execute()
    return jsonify(response.hits.hits)
