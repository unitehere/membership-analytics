#!/bin/bash
if [[ $EUID -ne 0 ]]; then
  echo "This script must be run as root." 
  exit 1
fi

if [[ -z "${ES_HOME}" ]]; then
  read -p "Enter the elasticsearch-plugin installation path (default: /usr/share/elasticsearch/bin):" ES_PLUGIN_DIR
  ES_PLUGIN_DIR=${ES_PLUGIN_DIR:-/usr/share/elasticsearch/bin}
else
  ES_PLUGIN_DIR="${ES_HOME}/bin"
fi

cd $ES_PLUGIN_DIR
./elasticsearch-plugin install analysis-phonetic
./elasticsearch-plugin install -b com.floragunn:search-guard-5:5.6.7-19

cd -

service elasticsearch restart
