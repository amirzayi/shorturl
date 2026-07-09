# URL Shortener

A lightweight and fast URL shortener written in Go.

## Features

- Shorten long URLs into compact short codes
- Redirect short URLs to their original destination
- Pluggable storage:
  - Redis for persistence and distributed usage
  - In-memory for local/testing environments
  - SQL based database(mysql, postgres, sqlite)
- Pluggable Key Generator:
  - No Duplicated
  - Random Timebased collision avoided
  - Custom DI(uuid, third-party services, ...)
- Simple REST API
- Configurable via environment variables
- Short link generator api access control using user&password

## Installation

Clone the repository:

```bash
git clone https://github.com/amirzayi/shorturl.git
cd shorturl
```

Build the binary

```bash
go build -o shortener .
```

### USAGE

## Environment Variables

The application uses the following environment variables for configuration:

| Variable           | Description                     | Required                    | Default |
| ------------------ | ------------------------------- | --------------------------- | ------- |
| `DRIVER`           | Database driver type            | Yes                         | -       |
| `DRIVER_URL`       | Database connection string      | No if used memory or sqlite | -       |
| `PREFIX`           | Database Storage Prefix Name    | No                          | -       |
| `HTTP_PORT`        | HTTP server port                | No                          | `8070`  |
| `LIMIT_PER_MINUTE` | Seek api calls limit per minute | No                          | -       |

#### Supported Drivers

| Driver   | Example Connection String                                                     |
| -------- | ----------------------------------------------------------------------------- |
| memory   | save on memory, will wipe on restart application, no configuration need       |
| sqlite   | db.sqlite file storage                                                        |
| redis    | `redis://(username:password)@host:port/0`                                     |
| mysql    | `mysql://root:password@tcp(host:port)/db_name?parseTime=true&charset=utf8mb4` |
| postgres | `postgres://root:password@host:port/db_name?sslmode=disable`                  |

### URL Shortener api authentication

```bash
nano credentials.json
```

```json
{
  "user": "password",
  "amir": "mirzayi"
}
```

### Custom shortener key generator using Dependency Injection

> Need to get rid of collision? use [**Random Time Based Key Generator**](./keygen/timebased.go) or your own **Generator** by Dependency Injection in [main.go](../main.go) file as below:

```go
shortener := handler.NewShortener(storage, func(ctx context.Context) (string, error) {
	id, err := uuid.NewV6()
	return id.String(), err
})
```

---

Using environment variables directly to Launch binary on port 8070 with Redis Storage and 15 api call per minute limit

```bash
LIMIT_PER_MINUTE=15 DRIVER=redis DRIVER_URL=redis://amir:mirzayi@127.0.0.1:6379/shortener_db PREFIX=shorts HTTP_PORT=8070 ./shortener
```

Docker image is available on [docker hub](https://hub.docker.com/r/amirzayi/shorturl)

```bash
docker pull amirzayi/shorturl:latest
docker run -p 8070:8000 -it --name shorturl --rm -e DRIVER=memory shorturl:latest
```
