# Covid19-bot 

This is a project made to spread awareness about my city's COVID-19 situation.

## Required environment variables

```
CONSUMER_KEY=
CONSUMER_SECRET=
ACCESS_TOKEN=
ACCESS_SECRET=
```

### Not required but recommended.

```
CITY_LAT=
CITY_LONG=
```

## Executing

```
cp .env-dist .env
```

Add required values

```
source .env && go run main.go
```