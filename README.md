# Go to typescript
Convert golang structs to typescript types
Based on [typescriptify-golang-structs](https://github.com/tkrajina/typescriptify-golang-structs).  

## Input variables

* `INPUT_FOLDER` - Path to the folder with .go files containing structs
* `OUTPUT_FILE` - Path for typescript file to be generated. For example `/src/types` will be processed to `/src/types.ts`

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
```

## Contributing

Pull requests are welcome at [VoVaVc/go-to-typescript](https://github.com/VoVaVc/go-to-typescript)

## License

The scripts and documentation in this project are released under the [MIT License](LICENSE)
