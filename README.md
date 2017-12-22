# redash-query-importer

Import queries from Definition file into [Redash](https://redash.io/)

## Installation

Download the latest binary from https://github.com/kirikiriyamama/redash-query-importer/releases

## Usage

### Synopsis

```
redash-query-importer -f FILE
```

### Options

```
Usage of redash-query-importer:
  -f string
        Path to definition file
```

### Definition file

```yml
api:
  base: https://redash.example.com
  key: xxxxx # NOT Query API Key but User API Key
queries:
  - name: name
    description: description # optional
    data_source_id: 1
    query: select 1
    schedule: 3600 # (periodic execution, in seconds) or 15:00 (daily execution), optional
```
