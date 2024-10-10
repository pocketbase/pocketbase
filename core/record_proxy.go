package core

// RecordProxy defines an interface for a Record proxy/project model,
// aka. custom model struct that acts on behalve the proxied Record to
// allow for example typed getter/setters for the Record fields.
//
// To implement the interface it is usually enough to embed the [BaseRecordProxy] struct.
type RecordProxy interface {
	// ProxyRecord returns the proxied Record model.
	ProxyRecord() *Record

	// SetProxyRecord loads the specified record model into the current proxy.
	SetProxyRecord(record *Record)
}

var _ RecordProxy = (*BaseRecordProxy)(nil)

// BaseRecordProxy implements the [RecordProxy] interface and it is intended
// to be used as embed to custom user provided Record proxy structs.
type BaseRecordProxy struct {
	*Record
}

// ProxyRecord returns the proxied Record model.
func (m *BaseRecordProxy) ProxyRecord() *Record {
	return m.Record
}

// SetProxyRecord loads the specified record model into the current proxy.
func (m *BaseRecordProxy) SetProxyRecord(record *Record) {
	m.Record = record
}
