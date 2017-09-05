use UNITEUAT;

SELECT n.id, STUFF(
  (
    SELECT ', {' +
    '"imis_id": "' + n1.ID + '",' +
    CASE WHEN n1.ORG_CODE IS NULL THEN '' ELSE '"org_code": "' + COALESCE(n1.ORG_CODE,'') + '",' END +
    CASE WHEN n1.MEMBER_TYPE IS NULL THEN '' ELSE '"member_type": "' + COALESCE(n1.MEMBER_TYPE,'') + '",' END +
    CASE WHEN n1.CATEGORY IS NULL THEN '' ELSE '"category": "' + COALESCE(n1.CATEGORY,'') + '",' END +
    CASE WHEN n1.[STATUS] IS NULL THEN '' ELSE '"status": "' + COALESCE(n1.[STATUS],'') + '",' END +
    CASE WHEN n1.TITLE IS NULL THEN '' ELSE '"title": "' + COALESCE(n1.TITLE,'') + '",' END +
    CASE WHEN n1.COMPANY IS NULL THEN '' ELSE '"company": "' + STRING_ESCAPE(COALESCE(n1.COMPANY,''), 'json') + '",' END +
    CASE WHEN n1.PREFIX IS NULL THEN '' ELSE '"prefix": "' + COALESCE(n1.PREFIX,'') + '",' END +
    CASE WHEN n1.FIRST_NAME IS NULL THEN '' ELSE '"first_name": "' + COALESCE(n1.FIRST_NAME,'') + '",' END +
    CASE WHEN n1.MIDDLE_NAME IS NULL THEN '' ELSE '"middle_name": "' + COALESCE(MIDDLE_NAME,'') + '",' END +
    CASE WHEN n1.LAST_NAME IS NULL THEN '' ELSE '"last_name": "' + COALESCE(n1.LAST_NAME,'') + '",' END +
    CASE WHEN n1.SUFFIX IS NULL THEN '' ELSE '"suffix": "' + COALESCE(n1.SUFFIX,'') + '",' END +
    CASE WHEN n1.DESIGNATION IS NULL THEN '' ELSE '"designation": "' + COALESCE(n1.DESIGNATION,'') + '",' END +
    CASE WHEN n1.CHAPTER IS NULL THEN '' ELSE '"chapter": "' + COALESCE(n1.CHAPTER,'') + '",' END +
    CASE WHEN n1.FUNCTIONAL_TITLE IS NULL THEN '' ELSE '"functional_title": "' + COALESCE(n1.FUNCTIONAL_TITLE,'') + '",' END +
    CASE WHEN n1.MEMBER_STATUS IS NULL THEN '' ELSE '"member_status": "' + COALESCE(n1.MEMBER_STATUS,'') + '"' END +
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
      CASE WHEN a1.TEAM IS NULL THEN '' ELSE '"team": "' + COALESCE(a1.TEAM,'') + '",' END +
      CASE WHEN a1.preferred IS NULL THEN '' ELSE '"preferred": ' + COALESCE(convert(varchar(1),a1.PREFERRED),'') + ',' END +
      CASE WHEN a1.ADDRESS_TYPE IS NULL THEN '' ELSE '"address_type": "' + COALESCE(a1.ADDRESS_TYPE,'') + '",' END +
      CASE WHEN a1.ADDRESS_1 IS NULL THEN '' ELSE '"address_1": "' + COALESCE(a1.ADDRESS_1,'') + '",' END +
      CASE WHEN a1.ADDRESS_2 IS NULL THEN '' ELSE '"address_2": "' + COALESCE(a1.ADDRESS_2,'') + '",' END +
      CASE WHEN a1.ADDRESS_3 IS NULL THEN '' ELSE '"address_3": "' + COALESCE(a1.ADDRESS_3,'') + '",' END +
      CASE WHEN a1.CITY IS NULL THEN '' ELSE '"city": "' + COALESCE(a1.CITY,'') + '",' END +
      CASE WHEN a1.STATE_PROVINCE IS NULL THEN '' ELSE '"state_province": "' + COALESCE(a1.STATE_PROVINCE,'') + '",' END +
      CASE WHEN a1.ZIP IS NULL THEN '' ELSE '"zip": "' + COALESCE(a1.ZIP,'') + '",' END + 
      CASE WHEN a1.COUNTRY IS NULL THEN '' ELSE '"country": "' + COALESCE(a1.COUNTRY,'') + '",' END +
      CASE WHEN a1.NOTE IS NULL THEN '' ELSE '"note": "' + COALESCE(a1.NOTE,'') + '",' END +
      CASE WHEN a1.VALIDATED IS NULL THEN '' ELSE '"validated": ' + COALESCE(convert(varchar(1), a1.VALIDATED),'') END + 
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
        CASE WHEN e1.TEAM IS NULL THEN '' ELSE '"team": "' + COALESCE(e1.TEAM,'') + '",' END +
        CASE WHEN e1.PREFERRED IS NULL THEN '' ELSE '"preferred": ' + COALESCE(convert(varchar(50),e1.PREFERRED),'') + ',' END +
        CASE WHEN e1.EMAIL_TYPE IS NULL THEN '' ELSE '"email_type": "' + COALESCE(e1.EMAIL_TYPE,'') + '",' END +
        CASE WHEN e1.EMAIL IS NULL THEN '' ELSE '"email": "' + COALESCE(e1.EMAIL,'') + '",' END + 
        CASE WHEN e1.NOTE IS NULL THEN '' ELSE '"note": "' + COALESCE(e1.NOTE,'') + '",' END +
        CASE WHEN e1.BAD IS NULL THEN '' ELSE '"bad": "' + COALESCE(convert(varchar(50), e1.BAD, 121),'') + '",' END +
        CASE WHEN e1.PERMISSION_TO_EMAIL IS NULL THEN '' ELSE '"permission_to_email": "' + COALESCE(e1.PERMISSION_TO_EMAIL,'') + '"' END +
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
          CASE WHEN p1.TEAM IS NULL THEN '' ELSE '"team": "' + COALESCE(p1.TEAM,'') + '",' END +
          CASE WHEN p1.PREFERRED IS NULL THEN '' ELSE '"preferred": ' + COALESCE(convert(varchar(50),p1.PREFERRED),'') + ',' END +
          CASE WHEN p1.PHONE_TYPE IS NULL THEN '' ELSE '"phone_type": "' + COALESCE(p1.PHONE_TYPE,'') + '",' END +
          CASE WHEN p1.PHONE IS NULL THEN '' ELSE '"phone": "' + COALESCE(p1.PHONE,'') + '",' END + 
          CASE WHEN p1.OPTIN_TEXT IS NULL THEN '' ELSE '"optin_text": "' + COALESCE(p1.OPTIN_TEXT,'') + '",' END +
          CASE WHEN p1.CONTACT_TIME IS NULL THEN '' ELSE '"contact_time": "' + COALESCE(p1.CONTACT_TIME,'') + '",' END +
          CASE WHEN p1.NOTES IS NULL THEN '' ELSE '"notes": "' + COALESCE(p1.NOTES,'') + '",' END +
          CASE WHEN p1.PERMISS_TO_TEXT_DATE IS NULL THEN '' ELSE COALESCE('"permiss_to_text_date": "' +convert(varchar(50), p1.PERMISS_TO_TEXT_DATE, 121) + '",','') END +
          CASE WHEN p1.PERMISS_TO_TEXT_SOURCE IS NULL THEN '' ELSE '"permiss_to_text_source": "' + COALESCE(p1.PERMISS_TO_TEXT_SOURCE,'') + '",' END +
          CASE WHEN p1.EXTENSION IS NULL THEN '' ELSE '"extension": "' + COALESCE(p1.EXTENSION,'') + '"' END +
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
            CASE WHEN q1.EMPLOYER_NAME IS NULL THEN '' ELSE '"employer_name": "' + STRING_ESCAPE(COALESCE(q1.EMPLOYER_NAME,''), 'json') + '",' END +
            CASE WHEN q1.EMPLOYER_ID IS NULL THEN '' ELSE '"employer_id": "' + COALESCE(q1.EMPLOYER_ID,'') + '",' END +
            CASE WHEN q1.EMPLOYEE_ID IS NULL THEN '' ELSE '"employee_id": "' + COALESCE(q1.EMPLOYEE_ID,'') + '",' END +
            CASE WHEN q1.PRIMARY_EMPLOYER IS NULL THEN '' ELSE '"primary_employer": ' + COALESCE(CONVERT(varchar(1), q1.PRIMARY_EMPLOYER),'') + ',' END + 
            CASE WHEN q1.EFFECTIVE_DATE IS NULL THEN '' ELSE COALESCE('"effective_date": "' + convert(varchar(50), q1.EFFECTIVE_DATE, 121) + '",','') END +
            CASE WHEN q1.THRU_DATE IS NULL THEN '' ELSE COALESCE('"thru_date": "' + convert(varchar(50), q1.THRU_DATE, 121) + '",','') END +
            CASE WHEN q1.CLASSIFICATION IS NULL THEN '' ELSE '"classification": "' + COALESCE(q1.CLASSIFICATION,'') + '",' END +
            CASE WHEN q1.[LOCATION] IS NULL THEN '' ELSE '"location": "' + COALESCE(q1.[LOCATION],'') + '",' END +
            CASE WHEN q1.DEPT IS NULL THEN '' ELSE '"dept": "' + COALESCE(q1.DEPT,'') + '",' END +
            CASE WHEN q1.[SHIFT] IS NULL THEN '' ELSE '"shift": "' + COALESCE(q1.[SHIFT],'') + '",' END +
            CASE WHEN q1.WORK_HOURS IS NULL THEN '' ELSE '"work_hours": "' + COALESCE(q1.WORK_HOURS,'') + '"' END + 
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
            CASE WHEN d1.SSN IS NULL THEN '' ELSE '"ssn": "' + COALESCE(d1.SSN,'') + '",' END +
            CASE WHEN d1.SSN_SRC IS NULL THEN '' ELSE '"ssn_src": "' + COALESCE(d1.SSN_SRC,'') + '",' END +
            CASE WHEN d1.ETHNICITY IS NULL THEN '' ELSE '"ethnicity": "' + COALESCE(d1.ETHNICITY,'') + '",' END +
            CASE WHEN d1.ETHNICITY_SRC IS NULL THEN '' ELSE '"ethnicity_src": "' + COALESCE(d1.ETHNICITY_SRC,'') + '",' END +
            CASE WHEN d1.COUNTRY_ORIGIN IS NULL THEN '' ELSE '"country_origin": "' + COALESCE(d1.COUNTRY_ORIGIN,'') + '",' END + 
            CASE WHEN d1.COUNTRY_ORIGIN_SRC IS NULL THEN '' ELSE '"country_origin_src": "' + COALESCE(d1.COUNTRY_ORIGIN_SRC,'') + '",' END +
            CASE WHEN d1.PRIMARY_LANGUAGE IS NULL THEN '' ELSE '"primary_language": "' + COALESCE(d1.PRIMARY_LANGUAGE,'') + '",' END +
            CASE WHEN d1.PRIMARY_LANGUAGE_SRC IS NULL THEN '' ELSE '"primary_language_src": "' + COALESCE(d1.PRIMARY_LANGUAGE_SRC,'') + '",' END +
            CASE WHEN d1.OTHER_LANGUAGE IS NULL THEN '' ELSE '"other_language": "' + COALESCE(d1.OTHER_LANGUAGE,'') + '",' END +
            CASE WHEN d1.OTHER_LANGUAGE_SRC IS NULL THEN '' ELSE '"other_language_src": "' + COALESCE(d1.OTHER_LANGUAGE_SRC,'') + '",' END +
            CASE WHEN d1.CITY_ORIGIN IS NULL THEN '' ELSE '"city_origin": "' + COALESCE(d1.CITY_ORIGIN,'') + '",' END +
            CASE WHEN d1.CITY_ORIGIN_SRC IS NULL THEN '' ELSE '"city_origin_src": "' + COALESCE(d1.CITY_ORIGIN_SRC,'') + '",' END + 
            CASE WHEN d1.GENDER IS NULL THEN '' ELSE '"gender": "' + COALESCE(d1.GENDER,'') + '",' END +
            CASE WHEN d1.GENDER_SRC IS NULL THEN '' ELSE '"gender_src": "' + COALESCE(d1.GENDER_SRC,'') + '",' END +
            CASE WHEN d1.OTHER_GENDER IS NULL THEN '' ELSE '"other_gender": "' + COALESCE(d1.OTHER_GENDER,'') + '"' END + 
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


