from json import load
from unittest import TestCase

from application import app
from query import SearchClient

term = 'foobar'


class TestSearch(TestCase):
    def setUp(self):
        self.search = SearchClient(app)

    def _load(self, name):
        return load(open('fixtures/%s.json' % name))

    def test_classification(self):
        self.search.classification(term)
        self.assertDictEqual(self._load('classification'),
                             self.search._search.to_dict())

    def test_employer_id(self):
        self.search.employer_id(term)
        self.assertDictEqual(self._load('employer_id'),
                             self.search._search.to_dict())

    def test_employer_name(self):
        self.search.employer_name(term)
        self.assertDictEqual(self._load('employer_name'),
                             self.search._search.to_dict())

    def test_first_name(self):
        self.search.first_name(term)
        self.assertDictEqual(self._load('first_name'),
                             self.search._search.to_dict())

    def test_last_name(self):
        self.search.last_name(term)
        self.assertDictEqual(self._load('last_name'),
                             self.search._search.to_dict())

    def test_ssn(self):
        self.search.ssn(term)
        self.assertDictEqual(self._load('ssn'),
                             self.search._search.to_dict())

    def test_state_province(self):
        self.search.state_province(term)
        self.assertDictEqual(self._load('state_province'),
                             self.search._search.to_dict())
