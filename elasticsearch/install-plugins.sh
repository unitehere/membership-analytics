#!/bin/bash
if [[ $EUID -ne 0 ]]; then
  echo "This script must be run as root." 
  exit 1
fi

if [[ -z "${ES_HOME}" ]]; then
  echo "Where is elasticsearch-plugin installed (e.g. /usr/share/elasticsearch/bin)?"
  read ES_PLUGIN_DIR
else
  ES_PLUGIN_DIR="${ES_HOME}/bin"
fi

cd $ES_PLUGIN_DIR
./elasticsearch-plugin install analysis-phonetic

cd -

service elasticsearch restart
