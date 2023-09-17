# Goascii
I wanted to play around with Go a bit, so I decided to write a picture-to-ASCII art converter.

Setup:
```
// `air` for live reloading
curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s

// `tailwindcss` for building our CSS
npm install -g tailwindcss
```

Run for local development:
```
// using `air` for live reloading
curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
air

// Without live reloading
go run main.go
```

Build the CSS:
```
tailwindcss -i index.css -o public/build.css --minify
```

Deploy to Fly.io:
```
flyctl deploy
```