Follow new users on Twitter using Intercom
==========================================

If you're using Intercom, this tiny app will automatically follow your users two days after they register.

## Installation

Add this App to Heroku and supply the following environment variables:

* `TWITTER_ACCESS_TOKEN`
* `TWITTER_ACCESS_TOKEN_SECRET`
* `INTERCOM_APP_ID`
* `INTERCOM_API_KEY`
* `REDIS_URL`
* `TWITTER_CONSUMER_KEY`
* `TWITTER_CONSUMER_SECRET`

Log into Intercom and add create the webhook (https://docs.intercom.io/integrations/webhooks). It should point to https://my-app-on-heroky.herokuapp.com/webhook

That's it.
