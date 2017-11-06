require 'elasticsearch'
require 'manticore'
require 'json'
require 'date'
require 'yaml'

require_relative 'models/ssn_audit_log'

# rubocop:disable Metrics/MethodLength
# rubocop:disable Metrics/AbcSize

CONFIG = YAML.load(File.read('config.yml'))

ES_CLIENT = Elasticsearch::Client.new(CONFIG['es_config'])
INDEX = 'member-audit-logs'.freeze
# INDEX = 'members-test'

def decrypt_and_replace
  response = ES_CLIENT.search(index: INDEX, scroll: '3m', size: 10_000,
                              body: {
                                query: {
                                  bool: {
                                    should: [
                                      { exists: { field: 'old_value' } },
                                      { exists: { field: 'new_value' } }
                                    ],
                                    must: [
                                      { term: { entity: 'UH_DEMO' } },
                                      { term: { property: 'SSN' } },
                                      { term: { ssn_decrypted: false } }
                                    ]
                                  }
                                }
                              })
  # process the first batch, generate a scroll_id to scroll through everything
  # in that index
  initial_batch = response['hits']['hits'].map do |doc|
    SSNAuditLog.new(doc['_id'],
                    doc['_source']['old_value'],
                    doc['_source']['new_value'])
  end
  initial_batch = decrypt_audit_logs(initial_batch)
  replace_audit_logs(initial_batch)

  # Call the `scroll` API until empty results are returned
  while (response = ES_CLIENT.scroll(scroll_id: response['_scroll_id'],
                                     scroll: '5m')) &&
        !response['hits']['hits'].empty?
    next_batch = response['hits']['hits'].map do |doc|
      SSNAuditLog.new(doc['_id'],
                      doc['_source']['old_value'],
                      doc['_source']['new_value'])
    end
    next_batch = decrypt_audit_logs(next_batch)
    replace_audit_logs(next_batch)
  end
end

## takes in an array of SSN and decrypt each encrypted ssn
def decrypt_audit_logs(audit_logs)
  old_values_payload = audit_logs.map(&:encrypted_old_value).reject(&:empty?)
  decrypted_data = JSON.parse(
    Manticore.post('https://fpe.unitehere.org/v1/ark/ff1/decrypt',
                   body: JSON.generate(values: old_values_payload),
                   headers: { Authorization: CONFIG['fpe_key'] }).body
  )['values']
  audit_logs.each do |audit_log|
    unless audit_log.encrypted_old_value.empty?
      audit_log.decrypted_old_value = decrypted_data.shift
    end
  end

  new_values_payload = audit_logs.map(&:encrypted_new_value).reject(&:empty?)
  decrypted_data = JSON.parse(
    Manticore.post('https://fpe.unitehere.org/v1/ark/ff1/decrypt',
                   body: JSON.generate(values: new_values_payload),
                   headers: { Authorization: CONFIG['fpe_key'] }).body
  )['values']
  audit_logs.each do |audit_log|
    unless audit_log.encrypted_new_value.empty?
      audit_log.decrypted_new_value = decrypted_data.shift
    end
  end
  audit_logs
end

def replace_audit_logs(audit_logs)
  query = audit_logs.map do |audit_log|
    {
      update:
      {
        _index: INDEX,
        _type: 'member-audit-log',
        _id: audit_log.id,
        data: {
          doc: {
            old_value: audit_log.decrypted_old_value,
            new_value: audit_log.decrypted_new_value,
            ssn_decrypted: true
          }
        }
      }
    }
  end
  p ES_CLIENT.bulk(body: query)
end

decrypt_and_replace
