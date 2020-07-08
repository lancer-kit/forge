linters-settings:
  govet:
    check-shadowing: true
  golint:
    min-confidence: 0
  gocyclo:
    min-complexity: 25
  maligned:
    suggest-new: true
  dupl:
    threshold: 200
  goconst:
    min-len: 2
    min-occurrences: 2

run:
  skip-dirs:
    - templates

linters:
  enable:
    - gochecknoglobals
    - goconst
    - gofmt
    - lll
    - misspell
    - scopelint
    - gochecknoinits
    - golint
    - gocritic
    - stylecheck
    - goimports
    - gosec
    - unconvert
    - unparam
  disable:
    - maligned
    - dupl
    - nakedret
