# RandomLinuxMemesBot
Discord bot for fetching random linux memes from reddit


## How to setup new deployment
1. Copy `.env.example` file to `.env`:
```
cp .env.example .env
```

2. Create your reddit credentials on https://www.reddit.com/prefs/apps and add to `.env`:
```
REDDIT_CLIENT_ID=someid
REDDIT_CLIENT_SECRET=somesecret
```

3. Add your postgresql user, password and db variables to `.env`:
```
POSTGRES_USER=sentry
POSTGRES_PASSWORD=somepassword
POSTGRES_DB=sentry
```

4. Run postgresql and keydb:
```
docker-compose up -d postgres keydb
```

5. Generate secret key for sentry:
```
docker-compose run --rm sentry config generate-secret-key
```

6. Add sentry key to `.env`
```
SENTRY_SECRET_KEY=somekey
```

7. Run migration for sentry db. There will be prompt to create initial user. Create new user and save credentials somewhere
```
docker-compose run --rm sentry upgrade
```

8. Run sentry:
```
docker-compose up -d sentry
```

9. Go to http://your_ip:9000, log in to sentry and create new Go project. Copy sentry DSN to `.env`:
```
SENTRY_DSN=http://pub_key@sentry:9000/project_id
```

10. Run remaining services:
```
docker-compose up -d
```
