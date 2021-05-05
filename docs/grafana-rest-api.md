https://grafana.com/grafana/plugins/simpod-json-datasource/
https://github.com/grafana/simple-json-datasource

To work with this datasource the backend needs to implement 4 endpoints:

GET / with 200 status code response. Used for "Test connection" on the datasource config page.
POST /search returning available metrics when invoked.
POST /query returning metrics based on input.
POST /annotations returning annotations.

# search body example
{
  "target": "query field value",
  "type": "timeseries" or "table"
}

# query body example
{
  "app": "dashboard",
  "requestId": "Q100",
  "timezone": "browser",
  "panelId": 23763571993,
  "dashboardId": null,
  "range": {
    "from": "2021-04-27T08:11:21.672Z",
    "to": "2021-04-27T14:11:21.672Z",
    "raw": {
      "from": "now-6h",
      "to": "now"
    }
  },
  "timeInfo": "",
  "interval": "30s",
  "intervalMs": 30000,
  "targets": [
    {
      "refId": "A",
      "data": "", // could also be {key: "", operator: "=", value: x}
      "target": "",
      "type": "timeseries",
      "datasource": "JSON"
    }
  ],
  "maxDataPoints": 830,
  "scopedVars": {
    "__interval": {
      "text": "30s",
      "value": "30s"
    },
    "__interval_ms": {
      "text": "30000",
      "value": 30000
    }
  },
  "startTime": 1619532681818,
  "rangeRaw": {
    "from": "now-6h",
    "to": "now"
  },
  "adhocFilters": []
}

# timeseries response example
[
  {
    "target":"foo",
    "datapoints":[
      [10,1257893400000],
      [20,1257892800000]
    ]
  },
  {
    "target":"doo",
    "datapoints":[
      [100,1257893400000],
      [200,1257892800000]
    ]
  }
]
