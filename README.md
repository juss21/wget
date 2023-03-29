# WGET

- Basic functionalities of wget linux command
- Downloading a single file and saving it under a different name
- Downloading and saving the file in a specific directory
- Set the download speed, limiting the rate of a download
- Downloading a file in background
- Downloading multiple files at the same time, by reading a file containing multiple donwload links asynchronously
- Main feature to download an entire website, [mirroring a website](https://en.wikipedia.org/wiki/Mirror_site)

## Usage
Build the project
```
go build .
./wget <url>
```

Run by hand
```
go run . <url>
```

Project help
```go run . -h 
./wget -h
```

**Authors:**
* Joel Meeks (kasepuu)
* Juss Märtson (juss)
* Rain-Ander Laagus (laagusra)
* Aleksandr Lavronenko (alavrone)
