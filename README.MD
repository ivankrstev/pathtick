### Arguments

| Argument | Type and accepted values | Default | Description |
|---|---|---:|---|
| [target-directory] | Positional path | . | Directory to scan |
| --title | String | Target directory’s name | Title shown in the document |
| --output-name, -n | String | PathTick | Output filename stem; format supplies the extension |
| --output-dir, -o | Directory path | . | Where exported files are saved |
| --format, -f | md, html, both | md | Export format |
| --stdout | Boolean | false | Write the export to standard output |
| --overwrite | Boolean | false | Allow replacing existing export files |
| --depth, -d | Integer; -1 means unlimited | -1 | Maximum directory traversal depth |
| --sort, -s | none, asc, desc | none | Order entries in the export(case-insensitive) |
| --skip-symlinks | Boolean | false | Exclude symbolic links |
| --skip-hidden | Boolean | false | Exclude entries whose names start with a dot |
