package store

import (

  "timeseries/lib/telemetry"
)

type mock struct {}

func (m *mock) Insert(rcd telemetry.Record) error {
  return nil
}

func (m *mock) Keys() ([]string, error) {
  return []string{}, nil
}

func (m *mock) Query(q telemetry.Query) ([]telemetry.Record, error) {
  return []telemetry.Record{}, nil
}

func (m *mock) Close() {}

func Mock() *mock {
  return &mock{}
}
