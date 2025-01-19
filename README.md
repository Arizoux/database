# DATABASE

Lightweight Database geschrieben in golang
- proprietaeres bin file .db format 



## Authors
- **Name:** Merdan
- **Name:** Felix
- **Name:** Dijar


## INSTRUCTIONS
```
go run main.go help
```

## CREATE
```
   go run main.go create <database_name>
```
   - legt file in cdatabase an (temporaer erstmal da ) 

## DEBUG
```
  go run main.go debug <database_name>
```
- konvertiert hexdump der bin file .db zu lesbaren text
- Magic Number (Database version und identifikation),  dblen (Laenge des database namen), dbname (Name der Database)

## Warum bin file format ?

- kein parsing noetig -> hoere effizienz
- identifikation der richtigen version und dateityp
- universell

## Aufbau bin
https://www.sqlite.org/fileformat.html
- Header
- Tables
- Definitionen