#!/bin/bash
cd "$(dirname "$0")"
pipenv lock --requirements > requirements.txt
eb deploy
cd -
