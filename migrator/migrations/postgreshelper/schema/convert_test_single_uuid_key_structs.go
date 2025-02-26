// Code generated by pg-bindings generator. DO NOT EDIT.
package schema

import (
	"github.com/lib/pq"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
)

// ConvertTestSingleUUIDKeyStructFromProto converts a `*storage.TestSingleUUIDKeyStruct` to Gorm model
func ConvertTestSingleUUIDKeyStructFromProto(obj *storage.TestSingleUUIDKeyStruct) (*TestSingleUUIDKeyStructs, error) {
	serialized, err := obj.Marshal()
	if err != nil {
		return nil, err
	}
	model := &TestSingleUUIDKeyStructs{
		Key:         obj.GetKey(),
		Name:        obj.GetName(),
		StringSlice: pq.Array(obj.GetStringSlice()).(*pq.StringArray),
		Bool:        obj.GetBool(),
		Uint64:      obj.GetUint64(),
		Int64:       obj.GetInt64(),
		Float:       obj.GetFloat(),
		Labels:      obj.GetLabels(),
		Timestamp:   pgutils.NilOrTime(obj.GetTimestamp()),
		Enum:        obj.GetEnum(),
		Enums:       pq.Array(pgutils.ConvertEnumSliceToIntArray(obj.GetEnums())).(*pq.Int32Array),
		Serialized:  serialized,
	}
	return model, nil
}

// ConvertTestSingleUUIDKeyStructToProto converts Gorm model `TestSingleUUIDKeyStructs` to its protobuf type object
func ConvertTestSingleUUIDKeyStructToProto(m *TestSingleUUIDKeyStructs) (*storage.TestSingleUUIDKeyStruct, error) {
	var msg storage.TestSingleUUIDKeyStruct
	if err := msg.Unmarshal(m.Serialized); err != nil {
		return nil, err
	}
	return &msg, nil
}
