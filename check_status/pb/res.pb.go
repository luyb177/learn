// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v4.25.6
// source: res.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Res struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Ret           int32                  `protobuf:"varint,1,opt,name=ret,proto3" json:"ret,omitempty"`
	Act           string                 `protobuf:"bytes,2,opt,name=act,proto3" json:"act,omitempty"`
	Msg           string                 `protobuf:"bytes,3,opt,name=msg,proto3" json:"msg,omitempty"`
	Data          []*Info                `protobuf:"bytes,4,rep,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Res) Reset() {
	*x = Res{}
	mi := &file_res_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Res) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Res) ProtoMessage() {}

func (x *Res) ProtoReflect() protoreflect.Message {
	mi := &file_res_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Res.ProtoReflect.Descriptor instead.
func (*Res) Descriptor() ([]byte, []int) {
	return file_res_proto_rawDescGZIP(), []int{0}
}

func (x *Res) GetRet() int32 {
	if x != nil {
		return x.Ret
	}
	return 0
}

func (x *Res) GetAct() string {
	if x != nil {
		return x.Act
	}
	return ""
}

func (x *Res) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *Res) GetData() []*Info {
	if x != nil {
		return x.Data
	}
	return nil
}

type Info struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RoomName      string                 `protobuf:"bytes,1,opt,name=room_name,json=roomName,proto3" json:"room_name,omitempty"`
	Title         string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Ts            []*T                   `protobuf:"bytes,3,rep,name=ts,proto3" json:"ts,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Info) Reset() {
	*x = Info{}
	mi := &file_res_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Info) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Info) ProtoMessage() {}

func (x *Info) ProtoReflect() protoreflect.Message {
	mi := &file_res_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Info.ProtoReflect.Descriptor instead.
func (*Info) Descriptor() ([]byte, []int) {
	return file_res_proto_rawDescGZIP(), []int{1}
}

func (x *Info) GetRoomName() string {
	if x != nil {
		return x.RoomName
	}
	return ""
}

func (x *Info) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Info) GetTs() []*T {
	if x != nil {
		return x.Ts
	}
	return nil
}

type T struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Start         string                 `protobuf:"bytes,1,opt,name=start,proto3" json:"start,omitempty"`
	End           string                 `protobuf:"bytes,2,opt,name=end,proto3" json:"end,omitempty"`
	State         string                 `protobuf:"bytes,3,opt,name=state,proto3" json:"state,omitempty"`
	Title         string                 `protobuf:"bytes,4,opt,name=title,proto3" json:"title,omitempty"`
	Owner         string                 `protobuf:"bytes,5,opt,name=owner,proto3" json:"owner,omitempty"`
	Occupy        bool                   `protobuf:"varint,6,opt,name=occupy,proto3" json:"occupy,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *T) Reset() {
	*x = T{}
	mi := &file_res_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *T) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*T) ProtoMessage() {}

func (x *T) ProtoReflect() protoreflect.Message {
	mi := &file_res_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use T.ProtoReflect.Descriptor instead.
func (*T) Descriptor() ([]byte, []int) {
	return file_res_proto_rawDescGZIP(), []int{2}
}

func (x *T) GetStart() string {
	if x != nil {
		return x.Start
	}
	return ""
}

func (x *T) GetEnd() string {
	if x != nil {
		return x.End
	}
	return ""
}

func (x *T) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *T) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *T) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

func (x *T) GetOccupy() bool {
	if x != nil {
		return x.Occupy
	}
	return false
}

var File_res_proto protoreflect.FileDescriptor

const file_res_proto_rawDesc = "" +
	"\n" +
	"\tres.proto\x12\x02pb\"Y\n" +
	"\x03Res\x12\x10\n" +
	"\x03ret\x18\x01 \x01(\x05R\x03ret\x12\x10\n" +
	"\x03act\x18\x02 \x01(\tR\x03act\x12\x10\n" +
	"\x03msg\x18\x03 \x01(\tR\x03msg\x12\x1c\n" +
	"\x04data\x18\x04 \x03(\v2\b.pb.InfoR\x04data\"P\n" +
	"\x04Info\x12\x1b\n" +
	"\troom_name\x18\x01 \x01(\tR\broomName\x12\x14\n" +
	"\x05title\x18\x02 \x01(\tR\x05title\x12\x15\n" +
	"\x02ts\x18\x03 \x03(\v2\x05.pb.TR\x02ts\"\x85\x01\n" +
	"\x01T\x12\x14\n" +
	"\x05start\x18\x01 \x01(\tR\x05start\x12\x10\n" +
	"\x03end\x18\x02 \x01(\tR\x03end\x12\x14\n" +
	"\x05state\x18\x03 \x01(\tR\x05state\x12\x14\n" +
	"\x05title\x18\x04 \x01(\tR\x05title\x12\x14\n" +
	"\x05owner\x18\x05 \x01(\tR\x05owner\x12\x16\n" +
	"\x06occupy\x18\x06 \x01(\bR\x06occupyB\x06Z\x04.;pbb\x06proto3"

var (
	file_res_proto_rawDescOnce sync.Once
	file_res_proto_rawDescData []byte
)

func file_res_proto_rawDescGZIP() []byte {
	file_res_proto_rawDescOnce.Do(func() {
		file_res_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_res_proto_rawDesc), len(file_res_proto_rawDesc)))
	})
	return file_res_proto_rawDescData
}

var file_res_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_res_proto_goTypes = []any{
	(*Res)(nil),  // 0: pb.Res
	(*Info)(nil), // 1: pb.Info
	(*T)(nil),    // 2: pb.T
}
var file_res_proto_depIdxs = []int32{
	1, // 0: pb.Res.data:type_name -> pb.Info
	2, // 1: pb.Info.ts:type_name -> pb.T
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_res_proto_init() }
func file_res_proto_init() {
	if File_res_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_res_proto_rawDesc), len(file_res_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_res_proto_goTypes,
		DependencyIndexes: file_res_proto_depIdxs,
		MessageInfos:      file_res_proto_msgTypes,
	}.Build()
	File_res_proto = out.File
	file_res_proto_goTypes = nil
	file_res_proto_depIdxs = nil
}
