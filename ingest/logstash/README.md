# Logstash Setup

Logstash is installed and run on the ASI cloud. This is generally a Windows server, so this directory contains a `install-plugins.bat` file that can be run to install the necessary plugins.

Configurations in this directory assume that Logstash and the SQL Server JDBC driver directory in this project are placed in
`C:\Elastic\`. If that needs to be modified, modify the Logstash configuration. In addition, you will need to fill in any credentials
that are not present in the configuration file.

To copy the files to the directory, you can run the `deploy-files.bat` script.
