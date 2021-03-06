# golang-url-shortener
A simple GO url shortener standalone web app.

With a simple REST-interface to created, edit and delete slugs with its corresponding URL. The slugs are saved in a SQLite database.


## How to install
Download this repository, create `.env` file, and add `APP_KEY` and `DEFAULT_URL` values. run `go get` and the script is ready to be used.

## How to run
Run the main.no script

`go run main.go`

## How to use
When adding a slug with a url, use the root url followed by the slug, like https://example.com/slug, and it will redirect to the url provided with the slug. The redirect is a 307 temporary redirect.



## REST interface

### GET (all slugs)

`/api/?api_key=your-api-key`

### GET (single slug)

`/api/slug?api_key=your-api-key`

### POST (create new slug)

`/api/`

#### json body:
`{"api_key": "your-api-key", "slug": "slug", "url": "https://..."}`

### PUT (update existing slug)

`/api/`

#### json body: 
`{"api_key": "your-api-key", "slug": "slug", "url": "https://..."}`


### DELETE (delete single slug)

`/api/`

#### json body: 
`{"api_key": "your-api-key", "slug": "slug"}`


### GET (slug redirect)

`/{slug}`
