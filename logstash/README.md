# Logstash Setup

Logstash is installed and run on the ASI cloud. This is generally a Windows server, so this directory contains a `plugins.txt` file
listing all the plugins that should be installed.

Configurations in this directory assume that Logstash and the SQL Server JDBC driver directory in this project are placed in
`C:\Elastic\`. If that needs to be modified, modify the Logstash configuration. In addition, you will need to fill in any credentials
that are not present in the configuration file.
