// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.7
// source: api/EventService.proto

package event

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type EventIdV1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *EventIdV1) Reset() {
	*x = EventIdV1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_EventService_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventIdV1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventIdV1) ProtoMessage() {}

func (x *EventIdV1) ProtoReflect() protoreflect.Message {
	mi := &file_api_EventService_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventIdV1.ProtoReflect.Descriptor instead.
func (*EventIdV1) Descriptor() ([]byte, []int) {
	return file_api_EventService_proto_rawDescGZIP(), []int{0}
}

func (x *EventIdV1) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type EventV1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Title       string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Description *string                `protobuf:"bytes,3,opt,name=description,proto3,oneof" json:"description,omitempty"`
	StartAt     *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=start_at,json=startAt,proto3" json:"start_at,omitempty"`
	EndAt       *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=end_at,json=endAt,proto3" json:"end_at,omitempty"`
	UserId      string                 `protobuf:"bytes,6,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	RemindFor   *uint32                `protobuf:"varint,7,opt,name=remind_for,json=remindFor,proto3,oneof" json:"remind_for,omitempty"`
}

func (x *EventV1) Reset() {
	*x = EventV1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_EventService_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventV1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventV1) ProtoMessage() {}

func (x *EventV1) ProtoReflect() protoreflect.Message {
	mi := &file_api_EventService_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventV1.ProtoReflect.Descriptor instead.
func (*EventV1) Descriptor() ([]byte, []int) {
	return file_api_EventService_proto_rawDescGZIP(), []int{1}
}

func (x *EventV1) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *EventV1) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *EventV1) GetDescription() string {
	if x != nil && x.Description != nil {
		return *x.Description
	}
	return ""
}

func (x *EventV1) GetStartAt() *timestamppb.Timestamp {
	if x != nil {
		return x.StartAt
	}
	return nil
}

func (x *EventV1) GetEndAt() *timestamppb.Timestamp {
	if x != nil {
		return x.EndAt
	}
	return nil
}

func (x *EventV1) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *EventV1) GetRemindFor() uint32 {
	if x != nil && x.RemindFor != nil {
		return *x.RemindFor
	}
	return 0
}

type EventResponseV1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*EventV1 `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *EventResponseV1) Reset() {
	*x = EventResponseV1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_EventService_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventResponseV1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventResponseV1) ProtoMessage() {}

func (x *EventResponseV1) ProtoReflect() protoreflect.Message {
	mi := &file_api_EventService_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventResponseV1.ProtoReflect.Descriptor instead.
func (*EventResponseV1) Descriptor() ([]byte, []int) {
	return file_api_EventService_proto_rawDescGZIP(), []int{2}
}

func (x *EventResponseV1) GetItems() []*EventV1 {
	if x != nil {
		return x.Items
	}
	return nil
}

var File_api_EventService_proto protoreflect.FileDescriptor

var file_api_EventService_proto_rawDesc = []byte{
	0x0a, 0x16, 0x61, 0x70, 0x69, 0x2f, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x1b, 0x0a, 0x09, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x49,
	0x64, 0x56, 0x31, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x22, 0x9c, 0x02, 0x0a, 0x07, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x56, 0x31, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x25, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0b, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x35, 0x0a, 0x08,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x41, 0x74, 0x12, 0x31, 0x0a, 0x06, 0x65, 0x6e, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x05, 0x65, 0x6e, 0x64, 0x41, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x22, 0x0a, 0x0a, 0x72, 0x65, 0x6d, 0x69, 0x6e, 0x64, 0x5f, 0x66, 0x6f, 0x72, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x0d, 0x48, 0x01, 0x52, 0x09, 0x72, 0x65, 0x6d, 0x69, 0x6e, 0x64, 0x46, 0x6f, 0x72,
	0x88, 0x01, 0x01, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x72, 0x65, 0x6d, 0x69, 0x6e, 0x64, 0x5f, 0x66,
	0x6f, 0x72, 0x22, 0x31, 0x0a, 0x0f, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x56, 0x31, 0x12, 0x1e, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x56, 0x31, 0x52, 0x05,
	0x69, 0x74, 0x65, 0x6d, 0x73, 0x32, 0xe8, 0x02, 0x0a, 0x08, 0x43, 0x61, 0x6c, 0x65, 0x6e, 0x64,
	0x61, 0x72, 0x12, 0x31, 0x0a, 0x0d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x56, 0x31, 0x12, 0x08, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x56, 0x31, 0x1a, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x31, 0x0a, 0x0d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x56, 0x31, 0x12, 0x08, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x56, 0x31,
	0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x33, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x56, 0x31, 0x12, 0x0a, 0x2e, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x49, 0x64, 0x56, 0x31, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x3e, 0x0a,
	0x0e, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x44, 0x61, 0x79, 0x56, 0x31, 0x12,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x1a, 0x10, 0x2e, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x56, 0x31, 0x12, 0x3f, 0x0a,
	0x0f, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x57, 0x65, 0x65, 0x6b, 0x56, 0x31,
	0x12, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x1a, 0x10, 0x2e, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x56, 0x31, 0x12, 0x40,
	0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4d, 0x6f, 0x6e, 0x74, 0x68,
	0x56, 0x31, 0x12, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x1a, 0x10,
	0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x56, 0x31,
	0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x3b, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_EventService_proto_rawDescOnce sync.Once
	file_api_EventService_proto_rawDescData = file_api_EventService_proto_rawDesc
)

func file_api_EventService_proto_rawDescGZIP() []byte {
	file_api_EventService_proto_rawDescOnce.Do(func() {
		file_api_EventService_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_EventService_proto_rawDescData)
	})
	return file_api_EventService_proto_rawDescData
}

var file_api_EventService_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_api_EventService_proto_goTypes = []interface{}{
	(*EventIdV1)(nil),             // 0: EventIdV1
	(*EventV1)(nil),               // 1: EventV1
	(*EventResponseV1)(nil),       // 2: EventResponseV1
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),         // 4: google.protobuf.Empty
}
var file_api_EventService_proto_depIdxs = []int32{
	3, // 0: EventV1.start_at:type_name -> google.protobuf.Timestamp
	3, // 1: EventV1.end_at:type_name -> google.protobuf.Timestamp
	1, // 2: EventResponseV1.items:type_name -> EventV1
	1, // 3: Calendar.CreateEventV1:input_type -> EventV1
	1, // 4: Calendar.UpdateEventV1:input_type -> EventV1
	0, // 5: Calendar.DeleteEventV1:input_type -> EventIdV1
	3, // 6: Calendar.ListEventDayV1:input_type -> google.protobuf.Timestamp
	3, // 7: Calendar.ListEventWeekV1:input_type -> google.protobuf.Timestamp
	3, // 8: Calendar.ListEventMonthV1:input_type -> google.protobuf.Timestamp
	4, // 9: Calendar.CreateEventV1:output_type -> google.protobuf.Empty
	4, // 10: Calendar.UpdateEventV1:output_type -> google.protobuf.Empty
	4, // 11: Calendar.DeleteEventV1:output_type -> google.protobuf.Empty
	2, // 12: Calendar.ListEventDayV1:output_type -> EventResponseV1
	2, // 13: Calendar.ListEventWeekV1:output_type -> EventResponseV1
	2, // 14: Calendar.ListEventMonthV1:output_type -> EventResponseV1
	9, // [9:15] is the sub-list for method output_type
	3, // [3:9] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_api_EventService_proto_init() }
func file_api_EventService_proto_init() {
	if File_api_EventService_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_EventService_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventIdV1); i {
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
		file_api_EventService_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventV1); i {
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
		file_api_EventService_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventResponseV1); i {
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
	file_api_EventService_proto_msgTypes[1].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_EventService_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_EventService_proto_goTypes,
		DependencyIndexes: file_api_EventService_proto_depIdxs,
		MessageInfos:      file_api_EventService_proto_msgTypes,
	}.Build()
	File_api_EventService_proto = out.File
	file_api_EventService_proto_rawDesc = nil
	file_api_EventService_proto_goTypes = nil
	file_api_EventService_proto_depIdxs = nil
}
