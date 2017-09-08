## What this do

1. When ran, decrypt and replace all ssns under the day's member index from logstash. eg "members-2017.09.06".
2. When complete, it will run an alias script to point today's index as "members".

## Getting Started
1. Ensure the gem bundler exists for your ruby installation, use `gem install bundler` if not.
2. Run `bundle` to get the required gems as written in Gemfile.
3. Replace the keys inside config.example.yml with your credentials
4. Run ruby script.rb, or jruby script.rb

## Note
1. Run a reindex query to generate a `members-test` index to test. Example:
  POST /_reindex
  {
    "source": {
      "index": "members"
    },
    "dest": {
      "index": "members-test"
    }
  }
2. Some keys in the yaml are symbols (with the :) infront of it. The elastic search gem doesn't take string keys for some reason.
