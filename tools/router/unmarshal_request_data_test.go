package router_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/tools/router"
)

func pointer[T any](val T) *T {
	return &val
}

func TestUnmarshalRequestData(t *testing.T) {
	t.Parallel()

	mapData := map[string][]string{
		"number1": {"1"},
		"number2": {"2", "3"},
		"number3": {"2.1", "-3.4"},
		"number4": {"0", "-0", "0.0001"},
		"string0": {""},
		"string1": {"a"},
		"string2": {"b", "c"},
		"string3": {
			"0.0",
			"-0.0",
			"000.1",
			"000001",
			"-000001",
			"1.6E-35",
			"-1.6E-35",
			"10e100",
			"1_000_000",
			"1.000.000",
			" 123 ",
			"0b1",
			"0xFF",
			"1234A",
			"Infinity",
			"-Infinity",
			"undefined",
			"null",
		},
		"bool1":        {"true"},
		"bool2":        {"true", "false"},
		"mixed":        {"true", "123", "test"},
		"@jsonPayload": {`{"json_a":null,"json_b":123}`, `{"json_c":[1,2,3]}`},
	}

	structData := map[string][]string{
		"stringTag":         {"a", "b"},
		"StringPtr":         {"b"},
		"StringSlice":       {"a", "b", "c", ""},
		"stringSlicePtrTag": {"d", "e"},
		"StringSliceOfPtr":  {"f", "g"},

		"boolTag":         {"true"},
		"BoolPtr":         {"true"},
		"BoolSlice":       {"true", "false", ""},
		"boolSlicePtrTag": {"false", "false", "true"},
		"BoolSliceOfPtr":  {"false", "true", "false"},

		"int8Tag":         {"-1", "2"},
		"Int8Ptr":         {"3"},
		"Int8Slice":       {"4", "5", ""},
		"int8SlicePtrTag": {"5", "6"},
		"Int8SliceOfPtr":  {"7", "8"},

		"int16Tag":         {"-1", "2"},
		"Int16Ptr":         {"3"},
		"Int16Slice":       {"4", "5", ""},
		"int16SlicePtrTag": {"5", "6"},
		"Int16SliceOfPtr":  {"7", "8"},

		"int32Tag":         {"-1", "2"},
		"Int32Ptr":         {"3"},
		"Int32Slice":       {"4", "5", ""},
		"int32SlicePtrTag": {"5", "6"},
		"Int32SliceOfPtr":  {"7", "8"},

		"int64Tag":         {"-1", "2"},
		"Int64Ptr":         {"3"},
		"Int64Slice":       {"4", "5", ""},
		"int64SlicePtrTag": {"5", "6"},
		"Int64SliceOfPtr":  {"7", "8"},

		"intTag":         {"-1", "2"},
		"IntPtr":         {"3"},
		"IntSlice":       {"4", "5", ""},
		"intSlicePtrTag": {"5", "6"},
		"IntSliceOfPtr":  {"7", "8"},

		"uint8Tag":         {"1", "2"},
		"Uint8Ptr":         {"3"},
		"Uint8Slice":       {"4", "5", ""},
		"uint8SlicePtrTag": {"5", "6"},
		"Uint8SliceOfPtr":  {"7", "8"},

		"uint16Tag":         {"1", "2"},
		"Uint16Ptr":         {"3"},
		"Uint16Slice":       {"4", "5", ""},
		"uint16SlicePtrTag": {"5", "6"},
		"Uint16SliceOfPtr":  {"7", "8"},

		"uint32Tag":         {"1", "2"},
		"Uint32Ptr":         {"3"},
		"Uint32Slice":       {"4", "5", ""},
		"uint32SlicePtrTag": {"5", "6"},
		"Uint32SliceOfPtr":  {"7", "8"},

		"uint64Tag":         {"1", "2"},
		"Uint64Ptr":         {"3"},
		"Uint64Slice":       {"4", "5", ""},
		"uint64SlicePtrTag": {"5", "6"},
		"Uint64SliceOfPtr":  {"7", "8"},

		"uintTag":         {"1", "2"},
		"UintPtr":         {"3"},
		"UintSlice":       {"4", "5", ""},
		"uintSlicePtrTag": {"5", "6"},
		"UintSliceOfPtr":  {"7", "8"},

		"float32Tag":         {"-1.2"},
		"Float32Ptr":         {"1.5", "2.0"},
		"Float32Slice":       {"1", "2.3", "-0.3", ""},
		"float32SlicePtrTag": {"-1.3", "3"},
		"Float32SliceOfPtr":  {"0", "1.2"},

		"float64Tag":         {"-1.2"},
		"Float64Ptr":         {"1.5", "2.0"},
		"Float64Slice":       {"1", "2.3", "-0.3", ""},
		"float64SlicePtrTag": {"-1.3", "3"},
		"Float64SliceOfPtr":  {"0", "1.2"},

		"timeTag":         {"2009-11-10T15:00:00Z"},
		"TimePtr":         {"2009-11-10T14:00:00Z", "2009-11-10T15:00:00Z"},
		"TimeSlice":       {"2009-11-10T14:00:00Z", "2009-11-10T15:00:00Z"},
		"timeSlicePtrTag": {"2009-11-10T15:00:00Z", "2009-11-10T16:00:00Z"},
		"TimeSliceOfPtr":  {"2009-11-10T17:00:00Z", "2009-11-10T18:00:00Z"},

		// @jsonPayload fields
		"@jsonPayload": {
			`{"payloadA":"test", "shouldBeIgnored": "abc"}`,
			`{"payloadB":[1,2,3], "payloadC":true}`,
		},

		// unexported fields or `-` tags
		"unexperted":                           {"test"},
		"SkipExported":                         {"test"},
		"unexportedStructFieldWithoutTag.Name": {"test"},
		"unexportedStruct.Name":                {"test"},

		// structs
		"StructWithoutTag.Name": {"test1"},
		"exportedStruct.Name":   {"test2"},

		// embedded
		"embed_name":         {"test3"},
		"embed2.embed_name2": {"test4"},
	}

	type embed1 struct {
		Name string `form:"embed_name"  json:"embed_name"`
	}

	type embed2 struct {
		Name string `form:"embed_name2" json:"embed_name2"`
	}

	//nolint
	type TestStruct struct {
		String           string `form:"stringTag" query:"stringTag2"`
		StringPtr        *string
		StringSlice      []string
		StringSlicePtr   *[]string `form:"stringSlicePtrTag"`
		StringSliceOfPtr []*string

		Bool           bool `form:"boolTag" query:"boolTag2"`
		BoolPtr        *bool
		BoolSlice      []bool
		BoolSlicePtr   *[]bool `form:"boolSlicePtrTag"`
		BoolSliceOfPtr []*bool

		Int8           int8 `form:"int8Tag" query:"int8Tag2"`
		Int8Ptr        *int8
		Int8Slice      []int8
		Int8SlicePtr   *[]int8 `form:"int8SlicePtrTag"`
		Int8SliceOfPtr []*int8

		Int16           int16 `form:"int16Tag" query:"int16Tag2"`
		Int16Ptr        *int16
		Int16Slice      []int16
		Int16SlicePtr   *[]int16 `form:"int16SlicePtrTag"`
		Int16SliceOfPtr []*int16

		Int32           int32 `form:"int32Tag" query:"int32Tag2"`
		Int32Ptr        *int32
		Int32Slice      []int32
		Int32SlicePtr   *[]int32 `form:"int32SlicePtrTag"`
		Int32SliceOfPtr []*int32

		Int64           int64 `form:"int64Tag" query:"int64Tag2"`
		Int64Ptr        *int64
		Int64Slice      []int64
		Int64SlicePtr   *[]int64 `form:"int64SlicePtrTag"`
		Int64SliceOfPtr []*int64

		Int           int `form:"intTag" query:"intTag2"`
		IntPtr        *int
		IntSlice      []int
		IntSlicePtr   *[]int `form:"intSlicePtrTag"`
		IntSliceOfPtr []*int

		Uint8           uint8 `form:"uint8Tag" query:"uint8Tag2"`
		Uint8Ptr        *uint8
		Uint8Slice      []uint8
		Uint8SlicePtr   *[]uint8 `form:"uint8SlicePtrTag"`
		Uint8SliceOfPtr []*uint8

		Uint16           uint16 `form:"uint16Tag" query:"uint16Tag2"`
		Uint16Ptr        *uint16
		Uint16Slice      []uint16
		Uint16SlicePtr   *[]uint16 `form:"uint16SlicePtrTag"`
		Uint16SliceOfPtr []*uint16

		Uint32           uint32 `form:"uint32Tag" query:"uint32Tag2"`
		Uint32Ptr        *uint32
		Uint32Slice      []uint32
		Uint32SlicePtr   *[]uint32 `form:"uint32SlicePtrTag"`
		Uint32SliceOfPtr []*uint32

		Uint64           uint64 `form:"uint64Tag" query:"uint64Tag2"`
		Uint64Ptr        *uint64
		Uint64Slice      []uint64
		Uint64SlicePtr   *[]uint64 `form:"uint64SlicePtrTag"`
		Uint64SliceOfPtr []*uint64

		Uint           uint `form:"uintTag" query:"uintTag2"`
		UintPtr        *uint
		UintSlice      []uint
		UintSlicePtr   *[]uint `form:"uintSlicePtrTag"`
		UintSliceOfPtr []*uint

		Float32           float32 `form:"float32Tag" query:"float32Tag2"`
		Float32Ptr        *float32
		Float32Slice      []float32
		Float32SlicePtr   *[]float32 `form:"float32SlicePtrTag"`
		Float32SliceOfPtr []*float32

		Float64           float64 `form:"float64Tag" query:"float64Tag2"`
		Float64Ptr        *float64
		Float64Slice      []float64
		Float64SlicePtr   *[]float64 `form:"float64SlicePtrTag"`
		Float64SliceOfPtr []*float64

		// encoding.TextUnmarshaler
		Time           time.Time `form:"timeTag" query:"timeTag2"`
		TimePtr        *time.Time
		TimeSlice      []time.Time
		TimeSlicePtr   *[]time.Time `form:"timeSlicePtrTag"`
		TimeSliceOfPtr []*time.Time

		// @jsonPayload fields
		JSONPayloadA string `form:"shouldBeIgnored" json:"payloadA"`
		JSONPayloadB []int  `json:"payloadB"`
		JSONPayloadC bool   `json:"-"`

		// unexported fields or `-` tags
		unexported                      string
		SkipExported                    string `form:"-"`
		unexportedStructFieldWithoutTag struct {
			Name string `json:"unexportedStructFieldWithoutTag_name"`
		}
		unexportedStructFieldWithTag struct {
			Name string `json:"unexportedStructFieldWithTag_name"`
		} `form:"unexportedStruct"`

		// structs
		StructWithoutTag struct {
			Name string `json:"StructWithoutTag_name"`
		}
		StructWithTag struct {
			Name string `json:"StructWithTag_name"`
		} `form:"exportedStruct"`

		// embedded
		embed1
		embed2 `form:"embed2"`
	}

	scenarios := []struct {
		name   string
		data   map[string][]string
		dst    any
		tag    string
		prefix string
		error  bool
		result string
	}{
		{
			name:   "nil data",
			data:   nil,
			dst:    pointer(map[string]any{}),
			error:  false,
			result: `{}`,
		},
		{
			name:  "non-pointer map[string]any",
			data:  mapData,
			dst:   map[string]any{},
			error: true,
		},
		{
			name:  "unsupported *map[string]string",
			data:  mapData,
			dst:   pointer(map[string]string{}),
			error: true,
		},
		{
			name:  "unsupported *map[string][]string",
			data:  mapData,
			dst:   pointer(map[string][]string{}),
			error: true,
		},
		{
			name:   "*map[string]any",
			data:   mapData,
			dst:    pointer(map[string]any{}),
			result: `{"bool1":true,"bool2":[true,false],"json_a":null,"json_b":123,"json_c":[1,2,3],"mixed":[true,123,"test"],"number1":1,"number2":[2,3],"number3":[2.1,-3.4],"number4":[0,-0,0.0001],"string0":"","string1":"a","string2":["b","c"],"string3":["0.0","-0.0","000.1","000001","-000001","1.6E-35","-1.6E-35","10e100","1_000_000","1.000.000"," 123 ","0b1","0xFF","1234A","Infinity","-Infinity","undefined","null"]}`,
		},
		{
			name:   "valid pointer struct (all fields)",
			data:   structData,
			dst:    &TestStruct{},
			result: `{"String":"a","StringPtr":"b","StringSlice":["a","b","c",""],"StringSlicePtr":["d","e"],"StringSliceOfPtr":["f","g"],"Bool":true,"BoolPtr":true,"BoolSlice":[true,false,false],"BoolSlicePtr":[false,false,true],"BoolSliceOfPtr":[false,true,false],"Int8":-1,"Int8Ptr":3,"Int8Slice":[4,5,0],"Int8SlicePtr":[5,6],"Int8SliceOfPtr":[7,8],"Int16":-1,"Int16Ptr":3,"Int16Slice":[4,5,0],"Int16SlicePtr":[5,6],"Int16SliceOfPtr":[7,8],"Int32":-1,"Int32Ptr":3,"Int32Slice":[4,5,0],"Int32SlicePtr":[5,6],"Int32SliceOfPtr":[7,8],"Int64":-1,"Int64Ptr":3,"Int64Slice":[4,5,0],"Int64SlicePtr":[5,6],"Int64SliceOfPtr":[7,8],"Int":-1,"IntPtr":3,"IntSlice":[4,5,0],"IntSlicePtr":[5,6],"IntSliceOfPtr":[7,8],"Uint8":1,"Uint8Ptr":3,"Uint8Slice":"BAUA","Uint8SlicePtr":"BQY=","Uint8SliceOfPtr":[7,8],"Uint16":1,"Uint16Ptr":3,"Uint16Slice":[4,5,0],"Uint16SlicePtr":[5,6],"Uint16SliceOfPtr":[7,8],"Uint32":1,"Uint32Ptr":3,"Uint32Slice":[4,5,0],"Uint32SlicePtr":[5,6],"Uint32SliceOfPtr":[7,8],"Uint64":1,"Uint64Ptr":3,"Uint64Slice":[4,5,0],"Uint64SlicePtr":[5,6],"Uint64SliceOfPtr":[7,8],"Uint":1,"UintPtr":3,"UintSlice":[4,5,0],"UintSlicePtr":[5,6],"UintSliceOfPtr":[7,8],"Float32":-1.2,"Float32Ptr":1.5,"Float32Slice":[1,2.3,-0.3,0],"Float32SlicePtr":[-1.3,3],"Float32SliceOfPtr":[0,1.2],"Float64":-1.2,"Float64Ptr":1.5,"Float64Slice":[1,2.3,-0.3,0],"Float64SlicePtr":[-1.3,3],"Float64SliceOfPtr":[0,1.2],"Time":"2009-11-10T15:00:00Z","TimePtr":"2009-11-10T14:00:00Z","TimeSlice":["2009-11-10T14:00:00Z","2009-11-10T15:00:00Z"],"TimeSlicePtr":["2009-11-10T15:00:00Z","2009-11-10T16:00:00Z"],"TimeSliceOfPtr":["2009-11-10T17:00:00Z","2009-11-10T18:00:00Z"],"payloadA":"test","payloadB":[1,2,3],"SkipExported":"","StructWithoutTag":{"StructWithoutTag_name":"test1"},"StructWithTag":{"StructWithTag_name":"test2"},"embed_name":"test3","embed_name2":"test4"}`,
		},
		{
			name:  "non-pointer struct",
			data:  structData,
			dst:   TestStruct{},
			error: true,
		},
		{
			name:  "invalid struct uint value",
			data:  map[string][]string{"uintTag": {"-1"}},
			dst:   &TestStruct{},
			error: true,
		},
		{
			name:  "invalid struct int value",
			data:  map[string][]string{"intTag": {"abc"}},
			dst:   &TestStruct{},
			error: true,
		},
		{
			name:  "invalid struct bool value",
			data:  map[string][]string{"boolTag": {"abc"}},
			dst:   &TestStruct{},
			error: true,
		},
		{
			name:  "invalid struct float value",
			data:  map[string][]string{"float64Tag": {"abc"}},
			dst:   &TestStruct{},
			error: true,
		},
		{
			name:  "invalid struct TextUnmarshaler value",
			data:  map[string][]string{"timeTag": {"123"}},
			dst:   &TestStruct{},
			error: true,
		},
		{
			name: "custom tagKey",
			data: map[string][]string{
				"tag1": {"a"},
				"tag2": {"b"},
				"tag3": {"c"},
				"Item": {"d"},
			},
			dst: &struct {
				Item string `form:"tag1" query:"tag2" json:"tag2"`
			}{},
			tag:    "query",
			result: `{"tag2":"b"}`,
		},
		{
			name: "custom prefix",
			data: map[string][]string{
				"test.A":     {"1"},
				"A":          {"2"},
				"test.alias": {"3"},
			},
			dst: &struct {
				A string
				B string `form:"alias"`
			}{},
			prefix: "test",
			result: `{"A":"1","B":"3"}`,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			err := router.UnmarshalRequestData(s.data, s.dst, s.tag, s.prefix)

			hasErr := err != nil
			if hasErr != s.error {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.error, hasErr, err)
			}

			if hasErr {
				return
			}

			raw, err := json.Marshal(s.dst)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(raw, []byte(s.result)) {
				t.Fatalf("Expected dst \n%s\ngot\n%s", s.result, raw)
			}
		})
	}
}

// note: extra unexported checks in addition to the above test as there
// is no easy way to print nested structs with all their fields.
func TestUnmarshalRequestDataUnexportedFields(t *testing.T) {
	t.Parallel()

	//nolint:all
	type TestStruct struct {
		Exported string

		unexported string
		// to ensure that the reflection doesn't take tags with higher priority than the exported state
		unexportedWithTag string `form:"unexportedWithTag" json:"unexportedWithTag"`
	}

	dst := &TestStruct{}

	err := router.UnmarshalRequestData(map[string][]string{
		"Exported": {"test"}, // just for reference

		"Unexported":        {"test"},
		"unexported":        {"test"},
		"UnexportedWithTag": {"test"},
		"unexportedWithTag": {"test"},
	}, dst, "", "")

	if err != nil {
		t.Fatal(err)
	}

	if dst.Exported != "test" {
		t.Fatalf("Expected the Exported field to be %q, got %q", "test", dst.Exported)
	}

	if dst.unexported != "" {
		t.Fatalf("Expected the unexported field to remain empty, got %q", dst.unexported)
	}

	if dst.unexportedWithTag != "" {
		t.Fatalf("Expected the unexportedWithTag field to remain empty, got %q", dst.unexportedWithTag)
	}
}
