import tweepy
import pickle
import os
import requests

auth = tweepy.OAuthHandler("", "")
auth.set_access_token("", "")

api = tweepy.API(auth)

tweets = {}

if os.path.exists('tweets.pkl'):
    with open('tweets.pkl', 'rb') as f:
        tweets = pickle.load(f)

public_tweets = api.search("october.app", count=100, result_type='recent')

for tweet in public_tweets:
    if not tweet.id in tweets.keys():
        tweets[tweet.id] = tweet
        post_data = {
            'username': 'Twitter',
            'channel': '#twitter',
            'text': '(<https://twitter.com/statuses/' + tweet.id_str + '|Open Tweet>) *' + tweet.user.name + ':* ' + tweet.text,
            'icon_url': 'https://cdn2.iconfinder.com/data/icons/metro-uinvert-dock/256/Twitter_NEW.png'
        }
        post_response = requests.post(url='https://hooks.slack.com/services/...', json=post_data)

with open('tweets.pkl', 'wb') as f:
    pickle.dump(tweets, f, pickle.HIGHEST_PROTOCOL)
