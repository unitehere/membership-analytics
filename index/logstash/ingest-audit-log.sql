SELECT(
  SELECT audit_log_id AS audit_log_id,
         old_value AS old_value,
         new_value AS new_value,
	       log_type AS log_type,
	       record_id AS imis_id,
	       record_seqn AS imis_seqn,
	       entity AS entity,
	       property AS property,
	       audit_date AS audit_date,
		     'false' AS ssn_decrypted
  FOR JSON PATH, WITHOUT_ARRAY_WRAPPER
)
AS audit_log
FROM dbo.UH_AUDIT_LOG
WHERE AUDIT_DATE >= DATEADD(day, -11117, GETDATE()) AND ENTITY='UH_DEMO' AND PROPERTY='SSN'
