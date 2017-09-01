use UNITEUAT;

SELECT n.id, STUFF(
  (
    SELECT ', {' +
    '"imis_id": "' + n1.ID + '",' +
    '"org_code": "' + COALESCE(n1.ORG_CODE,'') + '",' +
    '"member_type": "' + COALESCE(n1.MEMBER_TYPE,'') + '",' +
    '"category": "' + COALESCE(n1.CATEGORY,'') + '",' +
    '"status": "' + COALESCE(n1.[STATUS],'') + '",' +
    '"title": "' + COALESCE(n1.TITLE,'') + '",' +
    '"company": "' + STRING_ESCAPE(COALESCE(n1.COMPANY,''), 'json') + '",' +
    '"prefix": "' + COALESCE(n1.PREFIX,'') + '",' +
    '"first_name": "' + COALESCE(n1.FIRST_NAME,'') + '",' +
    '"middle_name": "' + COALESCE(MIDDLE_NAME,'') + '",' +
    '"last_name": "' + COALESCE(n1.LAST_NAME,'') + '",' +
    '"suffix": "' + COALESCE(n1.SUFFIX,'') + '",' +
    '"designation": "' + COALESCE(n1.DESIGNATION,'') + '",' +
    '"chapter": "' + COALESCE(n1.CHAPTER,'') + '",' +
    '"functional_title": "' + COALESCE(n1.FUNCTIONAL_TITLE,'') + '",' +
    '"member_status": "' + COALESCE(n1.MEMBER_STATUS,'') + '"' +
    '}'
  FROM dbo.Name n1
  WHERE n1.id = n.id
  ORDER BY n1.id
  FOR XML PATH(''), TYPE
    ).value('.','varchar(max)')
    , 1, 2, '') AS name_object,
  '[' + STUFF(
    (
      SELECT ', {' +
      '"imis_seqn": "' + COALESCE(convert(varchar(50),a1.SEQN),'') + '",' +
      '"team": "' + COALESCE(a1.TEAM,'') + '",' +
      '"preferred": ' + COALESCE(convert(varchar(1),a1.PREFERRED),'') + ',' +
      '"address_type": "' + COALESCE(a1.ADDRESS_TYPE,'') + '",' +
      '"address_1": "' + COALESCE(a1.ADDRESS_1,'') + '",' +
      '"address_2": "' + COALESCE(a1.ADDRESS_2,'') + '",' +
      '"address_3": "' + COALESCE(a1.ADDRESS_3,'') + '",' +
      '"city": "' + COALESCE(a1.CITY,'') + '",' +
      '"state_province": "' + COALESCE(a1.STATE_PROVINCE,'') + '",' +
      '"zip": "' + COALESCE(a1.ZIP,'') + '",' + 
      '"country": "' + COALESCE(a1.COUNTRY,'') + '",' +
      '"note": "' + COALESCE(a1.NOTE,'') + '",' +
      '"validated": ' + COALESCE(convert(varchar(1), a1.VALIDATED),'') + 
      '}'
  FROM dbo.UH_ADDRESS a1
  WHERE a1.id = n.id
  ORDER BY a1.id
  FOR XML PATH(''), TYPE
      ).value('.','varchar(max)')
      , 1, 2, '') + ']' AS address_array,
  '[' + STUFF(
      (
        SELECT ',  {' +
        '"imis_seqn": "' + COALESCE(convert(varchar(50),e1.SEQN),'') + '",' +
        '"team": "' + COALESCE(e1.TEAM,'') + '",' +
        '"preferred": ' + COALESCE(convert(varchar(50),e1.PREFERRED),'') + ',' +
        '"email_type": "' + COALESCE(e1.EMAIL_TYPE,'') + '",' +
        '"email": "' + COALESCE(e1.EMAIL,'') + '",' + 
        '"note": "' + COALESCE(e1.NOTE,'') + '",' +
        '"bad": "' + COALESCE(convert(varchar(50), e1.BAD, 121),'') + '",' +
        '"permission_to_email": "' + COALESCE(e1.PERMISSION_TO_EMAIL,'') + '"' +
        '}'
  FROM dbo.UH_EMAIL e1
  WHERE e1.id = n.id
  ORDER BY e1.id
  FOR XML PATH(''), TYPE
        ).value('.','varchar(max)')
        , 1, 2, '') + ']' AS email_array,
  '[' + STUFF(
        (
          SELECT ', {' +
          '"imis_seqn": "' + COALESCE(convert(varchar(50),p1.SEQN),'') + '",' +
          '"team": "' + COALESCE(p1.TEAM,'') + '",' +
          '"preferred": ' + COALESCE(convert(varchar(50),p1.PREFERRED),'') + ',' +
          '"phone_type": "' + COALESCE(p1.PHONE_TYPE,'') + '",' +
          '"phone": "' + COALESCE(p1.PHONE,'') + '",' + 
          '"optin_text": "' + COALESCE(p1.OPTIN_TEXT,'') + '",' +
          '"contact_time": "' + COALESCE(p1.CONTACT_TIME,'') + '",' +
          '"notes": "' + COALESCE(p1.NOTES,'') + '",' +
           COALESCE('"permiss_to_text_date": "' +convert(varchar(50), p1.PERMISS_TO_TEXT_DATE, 121) + '",','') +
          '"permiss_to_text_source": "' + COALESCE(p1.PERMISS_TO_TEXT_SOURCE,'') + '",' +
          '"extension": "' + COALESCE(p1.EXTENSION,'') + '"' +
          '}'
  FROM dbo.UH_PHONE p1
  WHERE p1.id = n.id
  ORDER BY p1.id
  FOR XML PATH(''), TYPE
          ).value('.','varchar(max)')
          , 1, 2, '') + ']' AS phone_array,
  '[' + STUFF(
          (
            SELECT ', {' +
            '"imis_seqn": "' + COALESCE(convert(varchar(50),q1.SEQN),'') + '",' +
            '"employer_name": "' + STRING_ESCAPE(COALESCE(q1.EMPLOYER_NAME,''), 'json') + '",' +
            '"employer_id": "' + COALESCE(q1.EMPLOYER_ID,'') + '",' +
            '"employee_id": "' + COALESCE(q1.EMPLOYEE_ID,'') + '",' +
            '"primary_employer": ' + COALESCE(CONVERT(varchar(1), q1.PRIMARY_EMPLOYER),'') + ',' + 
            COALESCE('"effective_date": "' + convert(varchar(50), q1.EFFECTIVE_DATE, 121) + '",','') +
            COALESCE('"thru_date": "' + convert(varchar(50), q1.THRU_DATE, 121) + '",','') +
            '"classification": "' + COALESCE(q1.CLASSIFICATION,'') + '",' +
            '"location": "' + COALESCE(q1.[LOCATION],'') + '",' +
            '"dept": "' + COALESCE(q1.DEPT,'') + '",' +
            '"shift": "' + COALESCE(q1.[SHIFT],'') + '",' +
            '"work_hours": "' + COALESCE(q1.WORK_HOURS,'') + '"' + 
            '}'
  FROM dbo.UH_EMPLOYER q1
  WHERE q1.id = n.id
  ORDER BY q1.id
  FOR XML PATH(''), TYPE
            ).value('.','varchar(max)')
            , 1, 2, '') + ']' AS employer_array,
  STUFF(
             (
            SELECT ', {' +
            '"ssn": "' + COALESCE(d1.SSN,'') + '",' +
            '"ssn_src": "' + COALESCE(d1.SSN_SRC,'') + '",' +
            '"ethnicity": "' + COALESCE(d1.ETHNICITY,'') + '",' +
            '"ethnicity_src": "' + COALESCE(d1.ETHNICITY_SRC,'') + '",' +
            '"country_origin": "' + COALESCE(d1.COUNTRY_ORIGIN,'') + '",' + 
            '"country_origin_src": "' + COALESCE(d1.COUNTRY_ORIGIN_SRC,'') + '",' +
            '"primary_language": "' + COALESCE(d1.PRIMARY_LANGUAGE,'') + '",' +
            '"primary_language_src": "' + COALESCE(d1.PRIMARY_LANGUAGE_SRC,'') + '",' +
            '"other_language": "' + COALESCE(d1.OTHER_LANGUAGE,'') + '",' +
            '"other_language_src": "' + COALESCE(d1.OTHER_LANGUAGE,'') + '",' +
            '"city_origin": "' + COALESCE(d1.CITY_ORIGIN,'') + '",' +
            '"city_origin_src": "' + COALESCE(d1.CITY_ORIGIN_SRC,'') + '",' + 
            '"gender": "' + COALESCE(d1.GENDER,'') + '",' +
            '"gender_src": "' + COALESCE(d1.GENDER_SRC,'') + '",' +
            '"other_gender": "' + COALESCE(d1.OTHER_GENDER,'') + '"' + 
            '}'
  FROM dbo.UH_DEMO d1
  WHERE d1.id = n.id
  ORDER BY d1.id
  FOR XML PATH(''), TYPE
              ).value('.','varchar(max)')
              , 1, 2, '') AS demographic_object

FROM dbo.Name n
WHERE n.Member_Type = 'M'
ORDER BY n.id
