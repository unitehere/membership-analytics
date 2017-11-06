# SSN Audit Log Class
class SSNAuditLog
  attr_accessor :decrypted_old_value
  attr_accessor :decrypted_new_value
  attr_reader :id, :encrypted_old_value, :encrypted_new_value

  def initialize(id, encrypted_old_value, encrypted_new_value)
    @id = id
    @encrypted_old_value = encrypted_old_value
    @encrypted_new_value = encrypted_new_value
  end
end
