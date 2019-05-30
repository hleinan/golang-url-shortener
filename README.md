# golang-url-shortener
A simple GO url shortener standalone web app.

With a simple REST-interface to created, edit and delete slugs with its URL.


## How to install
Download this repository, change the appKey and the defaultUrl

## How to run
Run the main.no script

`go run main.go`



## REST interface

### GET (all slugs)

`/api/?api_key=your-api-key`

### GET (single slug)

`https://dlfl.no/api/slug?api_key=your-api-key`

### POST (create new slug)

`https://dlfl.no/api/`

#### json body:
`{"api_key": "your-api-key", "slug": slug, "url": "https://..."}`

### PUT (update existing slug)

`https://dlfl.no/api/`

#### json body: 
`{"api_key": "your-api-key", "slug": slug, "url": "https://..."}`


### DELETE (delete single slug)

`https://dlfl.no/api/`

#### json body: 
`{"api_key": "your-api-key", "slug": slug, "url": "https://..."}`


### GET (slug redirect)

`https://dlfl.no/{slug}`
