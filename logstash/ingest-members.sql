SELECT n.id, '[' + STUFF(
  (
    SELECT '  {' +
    '"ID": "' + n1.ID + '",' +
    '"ORG_CODE": "' + COALESCE(n1.ORG_CODE,'') + '",' +
    '"MEMBER_TYPE": "' + COALESCE(n1.MEMBER_TYPE,'') + '",' +
    '"CATEGORY": "' + COALESCE(n1.CATEGORY,'') + '",' +
    '"STATUS": "' + COALESCE(n1.[STATUS],'') + '",' +
    '"TITLE": "' + COALESCE(n1.TITLE,'') + '",' +
    '"COMPANY": "' + COALESCE(n1.COMPANY,'') + '",' +
    '"PREFIX": "' + COALESCE(n1.PREFIX,'') + '",' +
    '"FIRST_NAME": "' + COALESCE(n1.FIRST_NAME,'') + '",' +
    '"MIDDLE_NAME": "' + COALESCE(MIDDLE_NAME,'') + '",' +
    '"LAST_NAME": "' + COALESCE(n1.LAST_NAME,'') + '",' +
    '"SUFFIX": "' + COALESCE(n1.SUFFIX,'') + '",' +
    '"DESIGNATION": "' + COALESCE(n1.DESIGNATION,'') + '",' +
    '"CHAPTER": "' + COALESCE(n1.CHAPTER,'') + '",' +
    '"FUNCTIONAL_TITLE": "' + COALESCE(n1.FUNCTIONAL_TITLE,'') + '",' +
    '"JOIN_DATE": "' + COALESCE(convert(varchar(50), n1.JOIN_DATE),'') + '",' + 
    '"MEMBER_STATUS": "' + COALESCE(n1.MEMBER_STATUS,'') + '",' +
    '"MEMBER_STATUS_DATE": "' + COALESCE(convert(varchar(50), n1.MEMBER_STATUS_DATE),'') + '"' +
    '},'
  FROM dbo.Name n1
  WHERE n1.id = n.id
  ORDER BY n1.id
  FOR XML PATH(''), TYPE
    ).value('.','varchar(max)')
    , 1, 2, '') + ']' AS name_array,
  '[' + STUFF(
    (
      SELECT '  {' +
      '"ID": "' + a1.ID + '",' +
      '"SEQN": "' + COALESCE(convert(varchar(50),a1.SEQN),'') + '",' +
      '"TEAM": "' + COALESCE(a1.TEAM,'') + '",' +
      '"PREFERRED": "' + COALESCE(convert(varchar(1),a1.PREFERRED),'') + '",' +
      '"ADDRESS_TYPE": "' + COALESCE(a1.ADDRESS_TYPE,'') + '",' +
      '"ADDRESS_1": "' + COALESCE(a1.ADDRESS_1,'') + '",' +
      '"ADDRESS_2": "' + COALESCE(a1.ADDRESS_2,'') + '",' +
      '"ADDRESS_3": "' + COALESCE(a1.ADDRESS_3,'') + '",' +
      '"CITY": "' + COALESCE(a1.CITY,'') + '",' +
      '"STATE_PROVINCE": "' + COALESCE(a1.STATE_PROVINCE,'') + '",' +
      '"ZIP": "' + COALESCE(a1.ZIP,'') + '",' + 
      '"COUNTRY": "' + COALESCE(a1.COUNTRY,'') + '",' +
      '"NOTE": "' + COALESCE(a1.NOTE,'') + '",' +
      '"VALIDATED": "' + COALESCE(convert(varchar(1), a1.VALIDATED),'') + '"' + 
      '},'
  FROM dbo.UH_ADDRESS a1
  WHERE a1.id = n.id
  ORDER BY a1.id
  FOR XML PATH(''), TYPE
      ).value('.','varchar(max)')
      , 1, 2, '') + ']' AS address_array,
  '[' + STUFF(
      (
        SELECT '  {' +
        '"ID": "' + e1.ID + '",' +
        '"SEQN": "' + COALESCE(convert(varchar(50),e1.SEQN),'') + '",' +
        '"TEAM": "' + COALESCE(e1.TEAM,'') + '",' +
        '"PREFERRED": "' + COALESCE(convert(varchar(50),e1.PREFERRED),'') + '",' +
        '"EMAIL_TYPE": "' + COALESCE(e1.EMAIL_TYPE,'') + '",' +
        '"EMAIL": "' + COALESCE(e1.EMAIL,'') + '",' + 
        '"NOTE": "' + COALESCE(e1.NOTE,'') + '",' +
        '"BAD": "' + COALESCE(convert(varchar(50), e1.BAD, 121),'') + '",' +
        '"PERMISSION_TO_EMAIL": "' + COALESCE(e1.PERMISSION_TO_EMAIL,'') + '",' +
        '},'
  FROM dbo.UH_EMAIL e1
  WHERE e1.id = n.id
  ORDER BY e1.id
  FOR XML PATH(''), TYPE
        ).value('.','varchar(max)')
        , 1, 2, '') + ']' AS email_array,
  '[' + STUFF(
        (
          SELECT '  {' +
          '"ID": "' + p1.ID + '",' +
          '"SEQN": "' + COALESCE(convert(varchar(50),p1.SEQN),'') + '",' +
          '"TEAM": "' + COALESCE(p1.TEAM,'') + '",' +
          '"PREFERRED": "' + COALESCE(convert(varchar(50),p1.PREFERRED),'') + '",' +
          '"PHONE_TYPE": "' + COALESCE(p1.PHONE_TYPE,'') + '",' +
          '"PHONE": "' + COALESCE(p1.PHONE,'') + '",' + 
          '"OPTIN_TEXT": "' + COALESCE(p1.OPTIN_TEXT,'') + '",' +
          '"CONTACT_TIME": "' + COALESCE(p1.CONTACT_TIME,'') + '",' +
          '"NOTES": "' + COALESCE(p1.NOTES,'') + '",' +
          '"PERMISS_TO_TEXT_DATE": "' + COALESCE(convert(varchar(50), p1.PERMISS_TO_TEXT_DATE, 121),'') + '",' +
          '"PERMISS_TO_TEXT_SOURCE": "' + COALESCE(p1.PERMISS_TO_TEXT_SOURCE,'') + '",' +
          '"EXTENSION": "' + COALESCE(p1.EXTENSION,'') + '"' +
          '},'
  FROM dbo.UH_PHONE p1
  WHERE p1.id = n.id
  ORDER BY p1.id
  FOR XML PATH(''), TYPE
          ).value('.','varchar(max)')
          , 1, 2, '') + ']' AS phone_array,
  '[' + STUFF(
          (
            SELECT '  {' +
            '"ID": "' + q1.ID + '",' +
            '"SEQN": "' + COALESCE(convert(varchar(50),q1.SEQN),'') + '",' +
            '"EMPLOYER_NAME": "' + COALESCE(q1.EMPLOYER_NAME,'') + '",' +
            '"EMPLOYER_ID": "' + COALESCE(q1.EMPLOYER_ID,'') + '",' +
            '"EMPLOYEE_ID": "' + COALESCE(q1.EMPLOYEE_ID,'') + '",' +
            '"PRIMARY_EMPLOYER": "' + COALESCE(CONVERT(varchar(1), q1.PRIMARY_EMPLOYER),'') + '",' + 
            '"EFFECTIVE_DATE": "' + COALESCE(convert(varchar(50), q1.EFFECTIVE_DATE, 121),'') + '",' +
            '"THRU_DATE": "' + COALESCE(convert(varchar(50), q1.THRU_DATE, 121),'') + '",' +
            '"CLASSIFICATION": "' + COALESCE(q1.CLASSIFICATION,'') + '",' +
            '"LOCATION": "' + COALESCE(q1.[LOCATION],'') + '",' +
            '"DEPT": "' + COALESCE(q1.DEPT,'') + '",' +
            '"SHIFT": "' + COALESCE(q1.[SHIFT],'') + '",' +
            '"WORK_HOURS": "' + COALESCE(q1.WORK_HOURS,'') + '"' + 
            '},'
  FROM dbo.UH_EMPLOYER q1
  WHERE q1.id = n.id
  ORDER BY q1.id
  FOR XML PATH(''), TYPE
            ).value('.','varchar(max)')
            , 1, 2, '') + ']' AS employer_array,
  '[' + STUFF(
            (
              SELECT '  {' +
              '"ID": "' + d1.ID + '",' +
              '"SSN": "' + COALESCE(d1.SSN,'') + '",' +
              '"SSN_SRC": "' + COALESCE(d1.SSN_SRC,'') + '",' +
              '"ETHNICITY": "' + COALESCE(d1.ETHNICITY,'') + '",' +
              '"ETHNICITY_SRC": "' + COALESCE(d1.ETHNICITY_SRC,'') + '",' +
              '"COUNTRY_ORIGIN": "' + COALESCE(d1.COUNTRY_ORIGIN,'') + '",' + 
              '"COUNTRY_ORIGIN_SRC": "' + COALESCE(d1.COUNTRY_ORIGIN_SRC,'') + '",' +
              '"PRIMARY_LANGUAGE": "' + COALESCE(d1.PRIMARY_LANGUAGE,'') + '",' +
              '"PRIMARY_LANGUAGE_SRC": "' + COALESCE(d1.PRIMARY_LANGUAGE_SRC,'') + '",' +
              '"OTHER_LANGUAGE": "' + COALESCE(d1.OTHER_LANGUAGE,'') + '",' +
              '"OTHER_LANGUAGE_SRC": "' + COALESCE(d1.OTHER_LANGUAGE,'') + '",' +
              '"CITY_ORIGIN": "' + COALESCE(d1.CITY_ORIGIN,'') + '",' +
              '"CITY_ORIGIN_SRC": "' + COALESCE(d1.CITY_ORIGIN_SRC,'') + '",' + 
              '"GENDER": "' + COALESCE(d1.GENDER,'') + '",' +
              '"GENDER_SRC": "' + COALESCE(d1.GENDER_SRC,'') + '",' +
              '"OTHER_GENDER": "' + COALESCE(d1.OTHER_GENDER,'') + '"' + 
              '},'
  FROM dbo.UH_DEMO d1
  WHERE d1.id = n.id
  ORDER BY d1.id
  FOR XML PATH(''), TYPE
              ).value('.','varchar(max)')
              , 1, 2, '') + ']' AS demo_array

FROM dbo.Name n
WHERE n.Member_Type = 'M'
ORDER BY n.id
