require 'elasticsearch'
require 'manticore'
require 'json'
require 'date'
require 'yaml'

require_relative 'models/ssn'

# rubocop:disable Metrics/MethodLength
# rubocop:disable Metrics/AbcSize

CONFIG = YAML.load(File.read('config.yml'))

ES_CLIENT = Elasticsearch::Client.new(CONFIG['es_config'])
INDEX = 'members-' + Date.today.strftime('%Y.%m.%d') # eg "members-2017.09.06"
# INDEX = 'members-test'

def decrypt_and_replace
  response = ES_CLIENT.search(index: INDEX, scroll: '3m', size: 10_000,
                              body: {
                                query: {
                                  bool: {
                                    must: {
                                      match_all: {}
                                    },
                                    filter: {
                                      exists: { field: 'demographics.ssn' }
                                    }
                                  }
                                }
                              })
  # process the first batch, generate a scroll_id to scroll through everything
  # in that index
  initial_batch = response['hits']['hits'].map do |doc|
    SSN.new(doc['_id'], doc['_source']['demographics']['ssn'])
  end
  decrypt_ssns(initial_batch)
  replace_ssns(initial_batch)

  # Call the `scroll` API until empty results are returned
  while (response = ES_CLIENT.scroll(scroll_id: response['_scroll_id'],
                                     scroll: '5m')) &&
        !response['hits']['hits'].empty?
    next_batch = response['hits']['hits'].map do |doc|
      SSN.new(doc['_id'], doc['_source']['demographics']['ssn'])
    end
    decrypt_ssns(next_batch)
    replace_ssns(next_batch)
  end

  shift_members_alias
end

## takes in an array of SSN and decrypt each encrypted ssn
def decrypt_ssns(ssns)
  payload = ssns.map(&:encrypted)
  decrypted_data = JSON.parse(
    Manticore.post('https://fpe.unitehere.org/v1/ark/ff1/decrypt',
                   body: JSON.generate(values: payload),
                   headers: { Authorization: CONFIG['fpe_key'] }).body
  )['values']
  ssns.each { |ssn| ssn.decrypted = decrypted_data.shift }
  ssns
end

def replace_ssns(ssns)
  query = ssns.map do |ssn|
    {
      update:
      {
        _index: INDEX,
        _type: 'member',
        _id: ssn.id,
        data: {
          doc: {
            demographics: {
              ssn: ssn.decrypted
            }
          }
        }
      }
    }
  end
  p ES_CLIENT.bulk(body: query)
end

def shift_members_alias
  ES_CLIENT.indices.update_aliases body: {
    actions: [
      { add: { index: INDEX, alias: 'members' } },
      {
        remove:
        {
          index: 'members-' + (Date.today - 1).strftime('%Y.%m.%d'),
          alias: 'members'
        }
      }
    ]
  }
end

decrypt_and_replace
