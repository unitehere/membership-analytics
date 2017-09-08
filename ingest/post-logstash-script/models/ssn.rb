class SSN
  attr_accessor :decrypted
  attr_reader :id, :encrypted

  def initialize(id, encrypted)
    @id = id
    @encrypted = encrypted
  end
end