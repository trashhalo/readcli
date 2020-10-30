# readcli

Tool that lets you read website content* on the command line. 

## Usage

```shell
readcli https://medium.com/compass-true-north/go-unit-testing-at-compass-3a7cb85ab54a
```

![](./sample.png)

## Install

### homebrew

```
brew install trashhalo/homebrew-brews/readcli
```

### prebuilt packages

Prebuilt packages can be found at the [releases page](https://github.com/trashhalo/imgcat/readcli)

## Website Content

The algorithm is as follows:
1. Use [go-readability](https://github.com/go-shiori/go-readability) to download a stripped down version of the website.
2. Use [html-to-markdown](https://github.com/JohannesKaufmann/html-to-markdown) to convert the clean html to markdown.
3. Use [glamour](https://github.com/charmbracelet/glamour) to render the markdown content.

This limits the tool to only sites that pass go-readability.

## Sites that work well

* Any medium post

## What about images in markdown content?

Stay tuned. https://github.com/trashhalo/imgcat/issues/11
