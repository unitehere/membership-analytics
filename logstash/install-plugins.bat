PUSHD C:\Elastic

CALL .\Logstash\bin\logstash-plugin.bat install logstash-input-jdbc
CALL .\Logstash\bin\logstash-plugin.bat install logstash-output-elasticsearch
CALL .\Logstash\bin\logstash-plugin.bat install logstash-filter-json
CALL .\Logstash\bin\logstash-plugin.bat install x-pack

POPD
