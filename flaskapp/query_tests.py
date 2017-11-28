from json import load
from unittest import TestCase

from application import app
from query import MemberSearchClient, HistorySearchClient

term = 'foobar'

class TestMemberSearch(TestCase):
    def setUp(self):
        self.search = MemberSearchClient(app)

    def _load(self, name):
        return load(open('fixtures/%s.json' % name))

    def test_from(self):
        self.search.set_from('')
        self.assertFalse('from' in self.search._search.to_dict())
        self.search.set_from(133)
        self.assertEquals(self.search._search.to_dict()['from'], 133)

    def test_size(self):
        self.search.set_size('')
        self.assertFalse('size' in self.search._search.to_dict())
        self.search.set_size(200)
        self.assertEquals(self.search._search.to_dict()['size'], 200)

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

class TestHistorySearch(TestCase):
    def setUp(self):
        self.search = HistorySearchClient(app)

    def _load(self, name):
        return load(open('fixtures/%s.json' % name))

    def test_from(self):
        self.search.set_from('')
        self.assertFalse('from' in self.search._search.to_dict())
        self.search.set_from(133)
        self.assertEquals(self.search._search.to_dict()['from'], 133)

    def test_size(self):
        self.search.set_size('')
        self.assertFalse('size' in self.search._search.to_dict())
        self.search.set_size(200)
        self.assertEquals(self.search._search.to_dict()['size'], 200)

    def test_ssn(self):
        self.search.ssn(term)
        self.assertDictEqual(self._load('ssn_history'),
                             self.search._search.to_dict())
