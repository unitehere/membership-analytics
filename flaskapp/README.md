# Membership Analytics

This is the search and analytics repository for the membership system.

## Getting Started
1. Install Pipenv `pip install Pipenv`
2. Install project dev dependencies `pipenv install -d`
3. Start up a project shell to do your work `pipenv shell`
3. Run the Flask app in debug/reloading mode `flask run --reload --debugger`

## Testing

In a pipenv shell

`nosetests`

## Configuration
You can create your own `config/local.py` configuration file to override the `config/common.py` settings. Development/test environments have their own overrides in config, so if you are developing make sure you override from those like so:

```python
# config/local.py
from .development import *  # import base settings

SECRET_KEY = 'foobarbaz'
```

All setting variables are available at [Flask Config](http://flask.pocoo.org/docs/0.12/config/#builtin-configuration-values)

## Deployment
After ensuring your configuration is correct, make sure you have your AWS CLI installed correctly before you deploy. Run `eb init` if it's the first time you are deploying and select to deploy to
the `us-west-2` region, the `membership-analytics` application, and the `membership-analytics`
environment.

After you initialize everything, run the `deploy.sh` script to deploy the application to AWS.
