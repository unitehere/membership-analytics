require 'elasticsearch'
require 'manticore'
require 'json'
require 'date'
require "yaml"

require_relative 'models/ssn'

CONFIG = YAML.load(File.read("config.yml"))

ES_CLIENT = Elasticsearch::Client.new(CONFIG['es_config'])
INDEX = 'members-' + Date.today.strftime("%Y.%m.%d") # eg "members-2017.09.06"
# INDEX = 'members-test'

def decrypt_and_replace
  response = ES_CLIENT.search(index: INDEX, scroll:"3m", size: 10000,
  body: {
    query: {
      bool: {
        must: {
         match_all: { }
        },
        filter: {
          exists: { field: "demographics.ssn" }
        }
      }
    }
  })
  # process the first batch, generate a scroll_id to scroll through everything in that index
  initial_batch = response['hits']['hits'].map do |doc|
    SSN.new(doc['_id'], doc['_source']['demographics']['ssn'])
  end
  decrypt_ssns(initial_batch)
  replace_ssns(initial_batch)

  # Call the `scroll` API until empty results are returned
  while response = ES_CLIENT.scroll(scroll_id: response['_scroll_id'], scroll: '5m') and not response['hits']['hits'].empty? do
    next_batch = response['hits']['hits'].map { |doc| SSN.new(doc['_id'], doc['_source']['demographics']['ssn']) }
    decrypt_ssns(next_batch)
    replace_ssns(next_batch)
  end

  shift_members_alias
end

## takes in an array of SSN and decrypt each encrypted ssn
def decrypt_ssns(ssns)
  payload = ssns.map { |ssn| ssn.encrypted}
  decrypted_data = JSON.parse(Manticore.post("https://fpe.unitehere.org/v1/ark/ff1/decrypt", 
    body: JSON.generate({ values: payload }),
    headers: { Authorization: CONFIG['fpe_key'] }
    ).body)["values"]
  ssns.each { |ssn| ssn.decrypted = decrypted_data.shift }
  return ssns
end

def replace_ssns(ssns)
  query = ssns.map do |ssn|
    { update:
      { _index: INDEX,
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
  ## add members to today's index
  ES_CLIENT.indices.put_alias(index: INDEX, name: 'members')
  ## remove members from yesterday's index
  ES_CLIENT.indices.delete_alias(index: 'members-' + (Date.today-1).strftime("%Y.%m.%d"), name: 'members')
end

decrypt_and_replace
