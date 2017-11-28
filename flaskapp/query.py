import json

from elasticsearch_dsl import Search
from elasticsearch_dsl.query import Fuzzy, Match, Wildcard, Bool


class SearchClient:
    extra = {
        'min_score': 0,
        'size': 200,
    }

    def __init__(self, app):
        self.app = app
        self._search = Search(using=app.es, index=app.config['ELASTICSEARCH_INDEX'])
        self._search = self._search.extra(**self.extra)

    def execute(self):
        return self._search.execute()

    def debug(self):
        return json.dumps(self._search.to_dict(), indent=4, sort_keys=False)

    def last_name(self, term):
        self._search._extra['min_score'] += 2
        self._search = self._search.query(Bool(
            should=[Fuzzy(last_name=term)],
            must=[Match(last_name__phonetic={'query': term, 'boost': 10, 'fuzziness': 1})]
        ))

    def ssn(self, term):
        self._search._extra['min_score'] += 2
        self._search = self._search.query(Bool(should=[
            Wildcard(demographics__ssn='*%s*' % term),
            Fuzzy(demographics__ssn=term)
        ]))

    def should_match(self, **kwargs):
        return self._search.query(Bool(should=[Match(**kwargs)]))

    def first_name(self, term):
        self._search._extra['min_score'] += 6
        self._search = self.should_match(first_name__phonetic=term)

    def employer_id(self, term):
        self._search._extra['min_score'] += 1
        self._search = self.should_match(employers__employer_id=term)

    def employer_name(self, term):
        self._search._extra['min_score'] += 1
        self._search = self.should_match(employers__employer_name=term)

    def classification(self, term):
        self._search._extra['min_score'] += 1
        self._search = self.should_match(employers__classification=term)

    def state_province(self, term):
        self._search._extra['min_score'] += 2
        self._search = self.should_match(addresses__state_province=term)
