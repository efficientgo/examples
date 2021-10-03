// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: rpc.proto

package export2

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type StoreType int32

const (
	StoreType_UNKNOWN StoreType = 0
	StoreType_QUERY   StoreType = 1
	StoreType_RULE    StoreType = 2
	StoreType_SIDECAR StoreType = 3
	StoreType_STORE   StoreType = 4
	StoreType_RECEIVE StoreType = 5
	// DEBUG represents some debug StoreAPI components e.g. thanos tools store-api-serve.
	StoreType_DEBUG StoreType = 6
)

// Enum value maps for StoreType.
var (
	StoreType_name = map[int32]string{
		0: "UNKNOWN",
		1: "QUERY",
		2: "RULE",
		3: "SIDECAR",
		4: "STORE",
		5: "RECEIVE",
		6: "DEBUG",
	}
	StoreType_value = map[string]int32{
		"UNKNOWN": 0,
		"QUERY":   1,
		"RULE":    2,
		"SIDECAR": 3,
		"STORE":   4,
		"RECEIVE": 5,
		"DEBUG":   6,
	}
)

func (x StoreType) Enum() *StoreType {
	p := new(StoreType)
	*p = x
	return p
}

func (x StoreType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (StoreType) Descriptor() protoreflect.EnumDescriptor {
	return file_rpc_proto_enumTypes[0].Descriptor()
}

func (StoreType) Type() protoreflect.EnumType {
	return &file_rpc_proto_enumTypes[0]
}

func (x StoreType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use StoreType.Descriptor instead.
func (StoreType) EnumDescriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{0}
}

type Aggr int32

const (
	Aggr_RAW     Aggr = 0
	Aggr_COUNT   Aggr = 1
	Aggr_SUM     Aggr = 2
	Aggr_MIN     Aggr = 3
	Aggr_MAX     Aggr = 4
	Aggr_COUNTER Aggr = 5
)

// Enum value maps for Aggr.
var (
	Aggr_name = map[int32]string{
		0: "RAW",
		1: "COUNT",
		2: "SUM",
		3: "MIN",
		4: "MAX",
		5: "COUNTER",
	}
	Aggr_value = map[string]int32{
		"RAW":     0,
		"COUNT":   1,
		"SUM":     2,
		"MIN":     3,
		"MAX":     4,
		"COUNTER": 5,
	}
)

func (x Aggr) Enum() *Aggr {
	p := new(Aggr)
	*p = x
	return p
}

func (x Aggr) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Aggr) Descriptor() protoreflect.EnumDescriptor {
	return file_rpc_proto_enumTypes[1].Descriptor()
}

func (Aggr) Type() protoreflect.EnumType {
	return &file_rpc_proto_enumTypes[1]
}

func (x Aggr) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Aggr.Descriptor instead.
func (Aggr) EnumDescriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{1}
}

/// PartialResponseStrategy controls partial response handling.
type PartialResponseStrategy int32

const (
	/// WARN strategy tells server to treat any error that will related to single StoreAPI (e.g missing chunk series because of underlying
	/// storeAPI is temporarily not available) as warning which will not fail the whole query (still OK response).
	/// Server should produce those as a warnings field in response.
	PartialResponseStrategy_WARN PartialResponseStrategy = 0
	/// ABORT strategy tells server to treat any error that will related to single StoreAPI (e.g missing chunk series because of underlying
	/// storeAPI is temporarily not available) as the gRPC error that aborts the query.
	///
	/// This is especially useful for any rule/alert evaluations on top of StoreAPI which usually does not tolerate partial
	/// errors.
	PartialResponseStrategy_ABORT PartialResponseStrategy = 1
)

// Enum value maps for PartialResponseStrategy.
var (
	PartialResponseStrategy_name = map[int32]string{
		0: "WARN",
		1: "ABORT",
	}
	PartialResponseStrategy_value = map[string]int32{
		"WARN":  0,
		"ABORT": 1,
	}
)

func (x PartialResponseStrategy) Enum() *PartialResponseStrategy {
	p := new(PartialResponseStrategy)
	*p = x
	return p
}

func (x PartialResponseStrategy) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PartialResponseStrategy) Descriptor() protoreflect.EnumDescriptor {
	return file_rpc_proto_enumTypes[2].Descriptor()
}

func (PartialResponseStrategy) Type() protoreflect.EnumType {
	return &file_rpc_proto_enumTypes[2]
}

func (x PartialResponseStrategy) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PartialResponseStrategy.Descriptor instead.
func (PartialResponseStrategy) EnumDescriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{2}
}

type Chunk_Encoding int32

const (
	Chunk_XOR Chunk_Encoding = 0
)

// Enum value maps for Chunk_Encoding.
var (
	Chunk_Encoding_name = map[int32]string{
		0: "XOR",
	}
	Chunk_Encoding_value = map[string]int32{
		"XOR": 0,
	}
)

func (x Chunk_Encoding) Enum() *Chunk_Encoding {
	p := new(Chunk_Encoding)
	*p = x
	return p
}

func (x Chunk_Encoding) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Chunk_Encoding) Descriptor() protoreflect.EnumDescriptor {
	return file_rpc_proto_enumTypes[3].Descriptor()
}

func (Chunk_Encoding) Type() protoreflect.EnumType {
	return &file_rpc_proto_enumTypes[3]
}

func (x Chunk_Encoding) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Chunk_Encoding.Descriptor instead.
func (Chunk_Encoding) EnumDescriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{2, 0}
}

type LabelMatcher_Type int32

const (
	LabelMatcher_EQ  LabelMatcher_Type = 0 // =
	LabelMatcher_NEQ LabelMatcher_Type = 1 // !=
	LabelMatcher_RE  LabelMatcher_Type = 2 // =~
	LabelMatcher_NRE LabelMatcher_Type = 3 // !~
)

// Enum value maps for LabelMatcher_Type.
var (
	LabelMatcher_Type_name = map[int32]string{
		0: "EQ",
		1: "NEQ",
		2: "RE",
		3: "NRE",
	}
	LabelMatcher_Type_value = map[string]int32{
		"EQ":  0,
		"NEQ": 1,
		"RE":  2,
		"NRE": 3,
	}
)

func (x LabelMatcher_Type) Enum() *LabelMatcher_Type {
	p := new(LabelMatcher_Type)
	*p = x
	return p
}

func (x LabelMatcher_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (LabelMatcher_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_rpc_proto_enumTypes[4].Descriptor()
}

func (LabelMatcher_Type) Type() protoreflect.EnumType {
	return &file_rpc_proto_enumTypes[4]
}

func (x LabelMatcher_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use LabelMatcher_Type.Descriptor instead.
func (LabelMatcher_Type) EnumDescriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{5, 0}
}

type SeriesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MinTime             int64           `protobuf:"varint,1,opt,name=min_time,json=minTime,proto3" json:"min_time,omitempty"`
	MaxTime             int64           `protobuf:"varint,2,opt,name=max_time,json=maxTime,proto3" json:"max_time,omitempty"`
	Matchers            []*LabelMatcher `protobuf:"bytes,3,rep,name=matchers,proto3" json:"matchers,omitempty"`
	MaxResolutionWindow int64           `protobuf:"varint,4,opt,name=max_resolution_window,json=maxResolutionWindow,proto3" json:"max_resolution_window,omitempty"`
	Aggregates          []Aggr          `protobuf:"varint,5,rep,packed,name=aggregates,proto3,enum=thanos.Aggr" json:"aggregates,omitempty"`
	// TODO(bwplotka): Move Thanos components to use strategy instead. Including QueryAPI.
	PartialResponseStrategy PartialResponseStrategy `protobuf:"varint,7,opt,name=partial_response_strategy,json=partialResponseStrategy,proto3,enum=thanos.PartialResponseStrategy" json:"partial_response_strategy,omitempty"`
	// skip_chunks controls whether sending chunks or not in series responses.
	SkipChunks bool `protobuf:"varint,8,opt,name=skip_chunks,json=skipChunks,proto3" json:"skip_chunks,omitempty"`
	// hints is an opaque data structure that can be used to carry additional information.
	// The content of this field and whether it's supported depends on the
	// implementation of a specific store.
	Hints *anypb.Any `protobuf:"bytes,9,opt,name=hints,proto3" json:"hints,omitempty"`
}

func (x *SeriesRequest) Reset() {
	*x = SeriesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SeriesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SeriesRequest) ProtoMessage() {}

func (x *SeriesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SeriesRequest.ProtoReflect.Descriptor instead.
func (*SeriesRequest) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{0}
}

func (x *SeriesRequest) GetMinTime() int64 {
	if x != nil {
		return x.MinTime
	}
	return 0
}

func (x *SeriesRequest) GetMaxTime() int64 {
	if x != nil {
		return x.MaxTime
	}
	return 0
}

func (x *SeriesRequest) GetMatchers() []*LabelMatcher {
	if x != nil {
		return x.Matchers
	}
	return nil
}

func (x *SeriesRequest) GetMaxResolutionWindow() int64 {
	if x != nil {
		return x.MaxResolutionWindow
	}
	return 0
}

func (x *SeriesRequest) GetAggregates() []Aggr {
	if x != nil {
		return x.Aggregates
	}
	return nil
}

func (x *SeriesRequest) GetPartialResponseStrategy() PartialResponseStrategy {
	if x != nil {
		return x.PartialResponseStrategy
	}
	return PartialResponseStrategy_WARN
}

func (x *SeriesRequest) GetSkipChunks() bool {
	if x != nil {
		return x.SkipChunks
	}
	return false
}

func (x *SeriesRequest) GetHints() *anypb.Any {
	if x != nil {
		return x.Hints
	}
	return nil
}

type SeriesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Result:
	//	*SeriesResponse_Series
	//	*SeriesResponse_Warning
	//	*SeriesResponse_Hints
	Result isSeriesResponse_Result `protobuf_oneof:"result"`
}

func (x *SeriesResponse) Reset() {
	*x = SeriesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SeriesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SeriesResponse) ProtoMessage() {}

func (x *SeriesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SeriesResponse.ProtoReflect.Descriptor instead.
func (*SeriesResponse) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{1}
}

func (m *SeriesResponse) GetResult() isSeriesResponse_Result {
	if m != nil {
		return m.Result
	}
	return nil
}

func (x *SeriesResponse) GetSeries() *Series {
	if x, ok := x.GetResult().(*SeriesResponse_Series); ok {
		return x.Series
	}
	return nil
}

func (x *SeriesResponse) GetWarning() string {
	if x, ok := x.GetResult().(*SeriesResponse_Warning); ok {
		return x.Warning
	}
	return ""
}

func (x *SeriesResponse) GetHints() *anypb.Any {
	if x, ok := x.GetResult().(*SeriesResponse_Hints); ok {
		return x.Hints
	}
	return nil
}

type isSeriesResponse_Result interface {
	isSeriesResponse_Result()
}

type SeriesResponse_Series struct {
	/// series contains 1 response series. The series labels are sorted by name.
	Series *Series `protobuf:"bytes,1,opt,name=series,proto3,oneof"`
}

type SeriesResponse_Warning struct {
	/// warning is considered an information piece in place of series for warning purposes.
	/// It is used to warn store API user about suspicious cases or partial response (if enabled).
	Warning string `protobuf:"bytes,2,opt,name=warning,proto3,oneof"`
}

type SeriesResponse_Hints struct {
	/// hints is an opaque data structure that can be used to carry additional information from
	/// the store. The content of this field and whether it's supported depends on the
	/// implementation of a specific store. It's also implementation specific if it's allowed that
	/// multiple SeriesResponse frames contain hints for a single Series() request and how should they
	/// be handled in such case (ie. merged vs keep the first/last one).
	Hints *anypb.Any `protobuf:"bytes,3,opt,name=hints,proto3,oneof"`
}

func (*SeriesResponse_Series) isSeriesResponse_Result() {}

func (*SeriesResponse_Warning) isSeriesResponse_Result() {}

func (*SeriesResponse_Hints) isSeriesResponse_Result() {}

type Chunk struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type Chunk_Encoding `protobuf:"varint,1,opt,name=type,proto3,enum=thanos.Chunk_Encoding" json:"type,omitempty"`
	Data []byte         `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Chunk) Reset() {
	*x = Chunk{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Chunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Chunk) ProtoMessage() {}

func (x *Chunk) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Chunk.ProtoReflect.Descriptor instead.
func (*Chunk) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{2}
}

func (x *Chunk) GetType() Chunk_Encoding {
	if x != nil {
		return x.Type
	}
	return Chunk_XOR
}

func (x *Chunk) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type Series struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Labels []*Label     `protobuf:"bytes,1,rep,name=labels,proto3" json:"labels,omitempty"`
	Chunks []*AggrChunk `protobuf:"bytes,2,rep,name=chunks,proto3" json:"chunks,omitempty"`
}

func (x *Series) Reset() {
	*x = Series{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Series) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Series) ProtoMessage() {}

func (x *Series) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Series.ProtoReflect.Descriptor instead.
func (*Series) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{3}
}

func (x *Series) GetLabels() []*Label {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *Series) GetChunks() []*AggrChunk {
	if x != nil {
		return x.Chunks
	}
	return nil
}

type AggrChunk struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MinTime int64  `protobuf:"varint,1,opt,name=min_time,json=minTime,proto3" json:"min_time,omitempty"`
	MaxTime int64  `protobuf:"varint,2,opt,name=max_time,json=maxTime,proto3" json:"max_time,omitempty"`
	Raw     *Chunk `protobuf:"bytes,3,opt,name=raw,proto3" json:"raw,omitempty"`
	Count   *Chunk `protobuf:"bytes,4,opt,name=count,proto3" json:"count,omitempty"`
	Sum     *Chunk `protobuf:"bytes,5,opt,name=sum,proto3" json:"sum,omitempty"`
	Min     *Chunk `protobuf:"bytes,6,opt,name=min,proto3" json:"min,omitempty"`
	Max     *Chunk `protobuf:"bytes,7,opt,name=max,proto3" json:"max,omitempty"`
	Counter *Chunk `protobuf:"bytes,8,opt,name=counter,proto3" json:"counter,omitempty"`
}

func (x *AggrChunk) Reset() {
	*x = AggrChunk{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AggrChunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AggrChunk) ProtoMessage() {}

func (x *AggrChunk) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AggrChunk.ProtoReflect.Descriptor instead.
func (*AggrChunk) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{4}
}

func (x *AggrChunk) GetMinTime() int64 {
	if x != nil {
		return x.MinTime
	}
	return 0
}

func (x *AggrChunk) GetMaxTime() int64 {
	if x != nil {
		return x.MaxTime
	}
	return 0
}

func (x *AggrChunk) GetRaw() *Chunk {
	if x != nil {
		return x.Raw
	}
	return nil
}

func (x *AggrChunk) GetCount() *Chunk {
	if x != nil {
		return x.Count
	}
	return nil
}

func (x *AggrChunk) GetSum() *Chunk {
	if x != nil {
		return x.Sum
	}
	return nil
}

func (x *AggrChunk) GetMin() *Chunk {
	if x != nil {
		return x.Min
	}
	return nil
}

func (x *AggrChunk) GetMax() *Chunk {
	if x != nil {
		return x.Max
	}
	return nil
}

func (x *AggrChunk) GetCounter() *Chunk {
	if x != nil {
		return x.Counter
	}
	return nil
}

// Matcher specifies a rule, which can match or set of labels or not.
type LabelMatcher struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type  LabelMatcher_Type `protobuf:"varint,1,opt,name=type,proto3,enum=thanos.LabelMatcher_Type" json:"type,omitempty"`
	Name  string            `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Value string            `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *LabelMatcher) Reset() {
	*x = LabelMatcher{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LabelMatcher) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LabelMatcher) ProtoMessage() {}

func (x *LabelMatcher) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LabelMatcher.ProtoReflect.Descriptor instead.
func (*LabelMatcher) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{5}
}

func (x *LabelMatcher) GetType() LabelMatcher_Type {
	if x != nil {
		return x.Type
	}
	return LabelMatcher_EQ
}

func (x *LabelMatcher) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *LabelMatcher) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type Label struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name  string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Label) Reset() {
	*x = Label{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Label) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Label) ProtoMessage() {}

func (x *Label) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Label.ProtoReflect.Descriptor instead.
func (*Label) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{6}
}

func (x *Label) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Label) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type LabelSet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Labels []*Label `protobuf:"bytes,1,rep,name=labels,proto3" json:"labels,omitempty"`
}

func (x *LabelSet) Reset() {
	*x = LabelSet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LabelSet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LabelSet) ProtoMessage() {}

func (x *LabelSet) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LabelSet.ProtoReflect.Descriptor instead.
func (*LabelSet) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{7}
}

func (x *LabelSet) GetLabels() []*Label {
	if x != nil {
		return x.Labels
	}
	return nil
}

type ZLabelSet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Labels []*Label `protobuf:"bytes,1,rep,name=labels,proto3" json:"labels,omitempty"`
}

func (x *ZLabelSet) Reset() {
	*x = ZLabelSet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ZLabelSet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ZLabelSet) ProtoMessage() {}

func (x *ZLabelSet) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ZLabelSet.ProtoReflect.Descriptor instead.
func (*ZLabelSet) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{8}
}

func (x *ZLabelSet) GetLabels() []*Label {
	if x != nil {
		return x.Labels
	}
	return nil
}

var File_rpc_proto protoreflect.FileDescriptor

var file_rpc_proto_rawDesc = []byte{
	0x0a, 0x09, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x74, 0x68, 0x61,
	0x6e, 0x6f, 0x73, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x83,
	0x03, 0x0a, 0x0d, 0x53, 0x65, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x19, 0x0a, 0x08, 0x6d, 0x69, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x07, 0x6d, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x6d,
	0x61, 0x78, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x6d,
	0x61, 0x78, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x30, 0x0a, 0x08, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65,
	0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x74, 0x68, 0x61, 0x6e, 0x6f,
	0x73, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x52, 0x08,
	0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x73, 0x12, 0x32, 0x0a, 0x15, 0x6d, 0x61, 0x78, 0x5f,
	0x72, 0x65, 0x73, 0x6f, 0x6c, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x77, 0x69, 0x6e, 0x64, 0x6f,
	0x77, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x13, 0x6d, 0x61, 0x78, 0x52, 0x65, 0x73, 0x6f,
	0x6c, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x57, 0x69, 0x6e, 0x64, 0x6f, 0x77, 0x12, 0x2c, 0x0a, 0x0a,
	0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0e,
	0x32, 0x0c, 0x2e, 0x74, 0x68, 0x61, 0x6e, 0x6f, 0x73, 0x2e, 0x41, 0x67, 0x67, 0x72, 0x52, 0x0a,
	0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x73, 0x12, 0x5b, 0x0a, 0x19, 0x70, 0x61,
	0x72, 0x74, 0x69, 0x61, 0x6c, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x73,
	0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e,
	0x74, 0x68, 0x61, 0x6e, 0x6f, 0x73, 0x2e, 0x50, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x53, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x52, 0x17,
	0x70, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x53,
	0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x6b, 0x69, 0x70, 0x5f,
	0x63, 0x68, 0x75, 0x6e, 0x6b, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x73, 0x6b,
	0x69, 0x70, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x73, 0x12, 0x2a, 0x0a, 0x05, 0x68, 0x69, 0x6e, 0x74,
	0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x05, 0x68,
	0x69, 0x6e, 0x74, 0x73, 0x22, 0x8e, 0x01, 0x0a, 0x0e, 0x53, 0x65, 0x72, 0x69, 0x65, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a, 0x06, 0x73, 0x65, 0x72, 0x69, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x74, 0x68, 0x61, 0x6e, 0x6f, 0x73,
	0x2e, 0x53, 0x65, 0x72, 0x69, 0x65, 0x73, 0x48, 0x00, 0x52, 0x06, 0x73, 0x65, 0x72, 0x69, 0x65,
	0x73, 0x12, 0x1a, 0x0a, 0x07, 0x77, 0x61, 0x72, 0x6e, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x00, 0x52, 0x07, 0x77, 0x61, 0x72, 0x6e, 0x69, 0x6e, 0x67, 0x12, 0x2c, 0x0a,
	0x05, 0x68, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41,
	0x6e, 0x79, 0x48, 0x00, 0x52, 0x05, 0x68, 0x69, 0x6e, 0x74, 0x73, 0x42, 0x08, 0x0a, 0x06, 0x72,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x5c, 0x0a, 0x05, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x2a,
	0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x74,
	0x68, 0x61, 0x6e, 0x6f, 0x73, 0x2e, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x2e, 0x45, 0x6e, 0x63, 0x6f,
	0x64, 0x69, 0x6e, 0x67, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x13,
	0x0a, 0x08, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x12, 0x07, 0x0a, 0x03, 0x58, 0x4f,
	0x52, 0x10, 0x00, 0x22, 0x5a, 0x0a, 0x06, 0x53, 0x65, 0x72, 0x69, 0x65, 0x73, 0x12, 0x25, 0x0a,
	0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e,
	0x74, 0x68, 0x61, 0x6e, 0x6f, 0x73, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x52, 0x06, 0x6c, 0x61,
	0x62, 0x65, 0x6c, 0x73, 0x12, 0x29, 0x0a, 0x06, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x74, 0x68, 0x61, 0x6e, 0x6f, 0x73, 0x2e, 0x41, 0x67,
	0x67, 0x72, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x52, 0x06, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x73, 0x22,
	0x93, 0x02, 0x0a, 0x09, 0x41, 0x67, 0x67, 0x72, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x19, 0x0a,
	0x08, 0x6d, 0x69, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x07, 0x6d, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x6d, 0x61, 0x78, 0x5f,
	0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x6d, 0x61, 0x78, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x03, 0x72, 0x61, 0x77, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0d, 0x2e, 0x74, 0x68, 0x61, 0x6e, 0x6f, 0x73, 0x2e, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x52,
	0x03, 0x72, 0x61, 0x77, 0x12, 0x23, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x74, 0x68, 0x61, 0x6e, 0x6f, 0x73, 0x2e, 0x43, 0x68, 0x75,
	0x6e, 0x6b, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1f, 0x0a, 0x03, 0x73, 0x75, 0x6d,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x74, 0x68, 0x61, 0x6e, 0x6f, 0x73, 0x2e,
	0x43, 0x68, 0x75, 0x6e, 0x6b, 0x52, 0x03, 0x73, 0x75, 0x6d, 0x12, 0x1f, 0x0a, 0x03, 0x6d, 0x69,
	0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x74, 0x68, 0x61, 0x6e, 0x6f, 0x73,
	0x2e, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x52, 0x03, 0x6d, 0x69, 0x6e, 0x12, 0x1f, 0x0a, 0x03, 0x6d,
	0x61, 0x78, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x74, 0x68, 0x61, 0x6e, 0x6f,
	0x73, 0x2e, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x52, 0x03, 0x6d, 0x61, 0x78, 0x12, 0x27, 0x0a, 0x07,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e,
	0x74, 0x68, 0x61, 0x6e, 0x6f, 0x73, 0x2e, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x52, 0x07, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x65, 0x72, 0x22, 0x91, 0x01, 0x0a, 0x0c, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x4d,
	0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x12, 0x2d, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x74, 0x68, 0x61, 0x6e, 0x6f, 0x73, 0x2e, 0x4c, 0x61,
	0x62, 0x65, 0x6c, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22,
	0x28, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x06, 0x0a, 0x02, 0x45, 0x51, 0x10, 0x00, 0x12,
	0x07, 0x0a, 0x03, 0x4e, 0x45, 0x51, 0x10, 0x01, 0x12, 0x06, 0x0a, 0x02, 0x52, 0x45, 0x10, 0x02,
	0x12, 0x07, 0x0a, 0x03, 0x4e, 0x52, 0x45, 0x10, 0x03, 0x22, 0x31, 0x0a, 0x05, 0x4c, 0x61, 0x62,
	0x65, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x31, 0x0a, 0x08,
	0x4c, 0x61, 0x62, 0x65, 0x6c, 0x53, 0x65, 0x74, 0x12, 0x25, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65,
	0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x74, 0x68, 0x61, 0x6e, 0x6f,
	0x73, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x22,
	0x32, 0x0a, 0x09, 0x5a, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x53, 0x65, 0x74, 0x12, 0x25, 0x0a, 0x06,
	0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x74,
	0x68, 0x61, 0x6e, 0x6f, 0x73, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x52, 0x06, 0x6c, 0x61, 0x62,
	0x65, 0x6c, 0x73, 0x2a, 0x5d, 0x0a, 0x09, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x09, 0x0a,
	0x05, 0x51, 0x55, 0x45, 0x52, 0x59, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x52, 0x55, 0x4c, 0x45,
	0x10, 0x02, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x49, 0x44, 0x45, 0x43, 0x41, 0x52, 0x10, 0x03, 0x12,
	0x09, 0x0a, 0x05, 0x53, 0x54, 0x4f, 0x52, 0x45, 0x10, 0x04, 0x12, 0x0b, 0x0a, 0x07, 0x52, 0x45,
	0x43, 0x45, 0x49, 0x56, 0x45, 0x10, 0x05, 0x12, 0x09, 0x0a, 0x05, 0x44, 0x45, 0x42, 0x55, 0x47,
	0x10, 0x06, 0x2a, 0x42, 0x0a, 0x04, 0x41, 0x67, 0x67, 0x72, 0x12, 0x07, 0x0a, 0x03, 0x52, 0x41,
	0x57, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x10, 0x01, 0x12, 0x07,
	0x0a, 0x03, 0x53, 0x55, 0x4d, 0x10, 0x02, 0x12, 0x07, 0x0a, 0x03, 0x4d, 0x49, 0x4e, 0x10, 0x03,
	0x12, 0x07, 0x0a, 0x03, 0x4d, 0x41, 0x58, 0x10, 0x04, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x4f, 0x55,
	0x4e, 0x54, 0x45, 0x52, 0x10, 0x05, 0x2a, 0x2e, 0x0a, 0x17, 0x50, 0x61, 0x72, 0x74, 0x69, 0x61,
	0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x53, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67,
	0x79, 0x12, 0x08, 0x0a, 0x04, 0x57, 0x41, 0x52, 0x4e, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x41,
	0x42, 0x4f, 0x52, 0x54, 0x10, 0x01, 0x32, 0x42, 0x0a, 0x05, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x12,
	0x39, 0x0a, 0x06, 0x53, 0x65, 0x72, 0x69, 0x65, 0x73, 0x12, 0x15, 0x2e, 0x74, 0x68, 0x61, 0x6e,
	0x6f, 0x73, 0x2e, 0x53, 0x65, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x16, 0x2e, 0x74, 0x68, 0x61, 0x6e, 0x6f, 0x73, 0x2e, 0x53, 0x65, 0x72, 0x69, 0x65, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x42, 0x3c, 0x5a, 0x3a, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x66, 0x66, 0x69, 0x63, 0x69, 0x65,
	0x6e, 0x74, 0x67, 0x6f, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x70, 0x61, 0x72, 0x71, 0x75, 0x65, 0x74, 0x2d, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74,
	0x2f, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x32, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_proto_rawDescOnce sync.Once
	file_rpc_proto_rawDescData = file_rpc_proto_rawDesc
)

func file_rpc_proto_rawDescGZIP() []byte {
	file_rpc_proto_rawDescOnce.Do(func() {
		file_rpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_proto_rawDescData)
	})
	return file_rpc_proto_rawDescData
}

var file_rpc_proto_enumTypes = make([]protoimpl.EnumInfo, 5)
var file_rpc_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_rpc_proto_goTypes = []interface{}{
	(StoreType)(0),               // 0: thanos.StoreType
	(Aggr)(0),                    // 1: thanos.Aggr
	(PartialResponseStrategy)(0), // 2: thanos.PartialResponseStrategy
	(Chunk_Encoding)(0),          // 3: thanos.Chunk.Encoding
	(LabelMatcher_Type)(0),       // 4: thanos.LabelMatcher.Type
	(*SeriesRequest)(nil),        // 5: thanos.SeriesRequest
	(*SeriesResponse)(nil),       // 6: thanos.SeriesResponse
	(*Chunk)(nil),                // 7: thanos.Chunk
	(*Series)(nil),               // 8: thanos.Series
	(*AggrChunk)(nil),            // 9: thanos.AggrChunk
	(*LabelMatcher)(nil),         // 10: thanos.LabelMatcher
	(*Label)(nil),                // 11: thanos.Label
	(*LabelSet)(nil),             // 12: thanos.LabelSet
	(*ZLabelSet)(nil),            // 13: thanos.ZLabelSet
	(*anypb.Any)(nil),            // 14: google.protobuf.Any
}
var file_rpc_proto_depIdxs = []int32{
	10, // 0: thanos.SeriesRequest.matchers:type_name -> thanos.LabelMatcher
	1,  // 1: thanos.SeriesRequest.aggregates:type_name -> thanos.Aggr
	2,  // 2: thanos.SeriesRequest.partial_response_strategy:type_name -> thanos.PartialResponseStrategy
	14, // 3: thanos.SeriesRequest.hints:type_name -> google.protobuf.Any
	8,  // 4: thanos.SeriesResponse.series:type_name -> thanos.Series
	14, // 5: thanos.SeriesResponse.hints:type_name -> google.protobuf.Any
	3,  // 6: thanos.Chunk.type:type_name -> thanos.Chunk.Encoding
	11, // 7: thanos.Series.labels:type_name -> thanos.Label
	9,  // 8: thanos.Series.chunks:type_name -> thanos.AggrChunk
	7,  // 9: thanos.AggrChunk.raw:type_name -> thanos.Chunk
	7,  // 10: thanos.AggrChunk.count:type_name -> thanos.Chunk
	7,  // 11: thanos.AggrChunk.sum:type_name -> thanos.Chunk
	7,  // 12: thanos.AggrChunk.min:type_name -> thanos.Chunk
	7,  // 13: thanos.AggrChunk.max:type_name -> thanos.Chunk
	7,  // 14: thanos.AggrChunk.counter:type_name -> thanos.Chunk
	4,  // 15: thanos.LabelMatcher.type:type_name -> thanos.LabelMatcher.Type
	11, // 16: thanos.LabelSet.labels:type_name -> thanos.Label
	11, // 17: thanos.ZLabelSet.labels:type_name -> thanos.Label
	5,  // 18: thanos.Store.Series:input_type -> thanos.SeriesRequest
	6,  // 19: thanos.Store.Series:output_type -> thanos.SeriesResponse
	19, // [19:20] is the sub-list for method output_type
	18, // [18:19] is the sub-list for method input_type
	18, // [18:18] is the sub-list for extension type_name
	18, // [18:18] is the sub-list for extension extendee
	0,  // [0:18] is the sub-list for field type_name
}

func init() { file_rpc_proto_init() }
func file_rpc_proto_init() {
	if File_rpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SeriesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SeriesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Chunk); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Series); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AggrChunk); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LabelMatcher); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Label); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LabelSet); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ZLabelSet); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_rpc_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*SeriesResponse_Series)(nil),
		(*SeriesResponse_Warning)(nil),
		(*SeriesResponse_Hints)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rpc_proto_rawDesc,
			NumEnums:      5,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rpc_proto_goTypes,
		DependencyIndexes: file_rpc_proto_depIdxs,
		EnumInfos:         file_rpc_proto_enumTypes,
		MessageInfos:      file_rpc_proto_msgTypes,
	}.Build()
	File_rpc_proto = out.File
	file_rpc_proto_rawDesc = nil
	file_rpc_proto_goTypes = nil
	file_rpc_proto_depIdxs = nil
}