# Go to typescript
Convert golang structs to typescript types
Based on [typescriptify-golang-structs](https://github.com/tkrajina/typescriptify-golang-structs).  

## Input variables
| Name  | Required  | Default  | Description  |
| ------------- | ------------- | ------------- | ------------- |
| INPUT_FOLDER  | true  | no  | Path to the folder with .go files containing structs |
| OUTPUT_FILE  | true  | no  | Path for typescript file to be generated. For example `/src/types` will be processed to `/src/types.ts` |
| IDENT  | false  | tab  | Identation in generated .ts types file, tab by default |
| PREFIX  | false  | no  | Prefix for typescript interface, ex.: `Prefix_MyApiInterface` |
| SUFFIX  | false  | no  | Suffix for typescript interface, ex.: `MyApiInterface_Suffix`  |
| CREATE_FROM_METHOD  | false  | no  | Create classes with createFrom method to init an object |
| CREATE_CONSTRUCTOR  | false  | no  | Create classes with plain construct |
| DONT_EXPORT  | false  | no  | Don't export created interfaces |
| CREATE_INTERFACE  | false  | no  | Create typescript interfaces instead of classes |
| BACKUP_DIR  | false  | no  | Path to a folder where to store backups |

## Usage with docker-compose.yml

```yaml
version: "3.9"
services:
    generate-types:
        image: vovavc/go-to-typescript
        volumes:
            - ./sources:/sources
            - ./src:/src
        environment:
            - INPUT_FOLDER=/sources/path-to-your-golang-structs-folder/
            - OUTPUT_FILE=/src/types
            - CREATE_INTERFACE=true
```

## Contributing

Pull requests are welcome at [VoVaVc/go-to-typescript](https://github.com/VoVaVc/go-to-typescript)

## License

The scripts and documentation in this project are released under the [MIT License](LICENSE)
