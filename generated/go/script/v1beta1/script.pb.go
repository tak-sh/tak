// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        (unknown)
// source: script/v1beta1/script.proto

package v1beta1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Component_Input_Type int32

const (
	Component_Input_TEXT     Component_Input_Type = 0
	Component_Input_PASSWORD Component_Input_Type = 1
)

// Enum value maps for Component_Input_Type.
var (
	Component_Input_Type_name = map[int32]string{
		0: "TEXT",
		1: "PASSWORD",
	}
	Component_Input_Type_value = map[string]int32{
		"TEXT":     0,
		"PASSWORD": 1,
	}
)

func (x Component_Input_Type) Enum() *Component_Input_Type {
	p := new(Component_Input_Type)
	*p = x
	return p
}

func (x Component_Input_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Component_Input_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_script_v1beta1_script_proto_enumTypes[0].Descriptor()
}

func (Component_Input_Type) Type() protoreflect.EnumType {
	return &file_script_v1beta1_script_proto_enumTypes[0]
}

func (x Component_Input_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Component_Input_Type.Descriptor instead.
func (Component_Input_Type) EnumDescriptor() ([]byte, []int) {
	return file_script_v1beta1_script_proto_rawDescGZIP(), []int{4, 1, 0}
}

// A way to programmatically control what the headless browser should do.
type Script struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Steps []*Step `protobuf:"bytes,1,rep,name=steps,proto3" json:"steps,omitempty"`
}

func (x *Script) Reset() {
	*x = Script{}
	if protoimpl.UnsafeEnabled {
		mi := &file_script_v1beta1_script_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Script) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Script) ProtoMessage() {}

func (x *Script) ProtoReflect() protoreflect.Message {
	mi := &file_script_v1beta1_script_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Script.ProtoReflect.Descriptor instead.
func (*Script) Descriptor() ([]byte, []int) {
	return file_script_v1beta1_script_proto_rawDescGZIP(), []int{0}
}

func (x *Script) GetSteps() []*Step {
	if x != nil {
		return x.Steps
	}
	return nil
}

// A single line within a Script.
type Step struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A top-level referential key for the Step. If not provided, the result is not
	// stored nor is it accessible by following Steps. If it is provided, any following
	// Step can refer to the value returned by the Action.
	Id *string `protobuf:"bytes,1,opt,name=id,proto3,oneof" json:"id,omitempty"`
	// Provide the action that should be taken for this Step.
	Action *Action `protobuf:"bytes,2,opt,name=action,proto3" json:"action,omitempty"`
	// Execute the Step if a truthy value is returned.
	If *string `protobuf:"bytes,3,opt,name=if,proto3,oneof" json:"if,omitempty"`
}

func (x *Step) Reset() {
	*x = Step{}
	if protoimpl.UnsafeEnabled {
		mi := &file_script_v1beta1_script_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Step) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Step) ProtoMessage() {}

func (x *Step) ProtoReflect() protoreflect.Message {
	mi := &file_script_v1beta1_script_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Step.ProtoReflect.Descriptor instead.
func (*Step) Descriptor() ([]byte, []int) {
	return file_script_v1beta1_script_proto_rawDescGZIP(), []int{1}
}

func (x *Step) GetId() string {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return ""
}

func (x *Step) GetAction() *Action {
	if x != nil {
		return x.Action
	}
	return nil
}

func (x *Step) GetIf() string {
	if x != nil && x.If != nil {
		return *x.If
	}
	return ""
}

type Action struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Click on an element.
	MouseClick *Action_MouseClick `protobuf:"bytes,1,opt,name=mouse_click,proto3,oneof" json:"mouse_click,omitempty"`
	// Input some text into an element.
	Input *Action_Input `protobuf:"bytes,2,opt,name=input,proto3,oneof" json:"input,omitempty"`
	// Ask the user for some information.
	Ask *Action_PromptUser `protobuf:"bytes,4,opt,name=ask,proto3,oneof" json:"ask,omitempty"`
}

func (x *Action) Reset() {
	*x = Action{}
	if protoimpl.UnsafeEnabled {
		mi := &file_script_v1beta1_script_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Action) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Action) ProtoMessage() {}

func (x *Action) ProtoReflect() protoreflect.Message {
	mi := &file_script_v1beta1_script_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Action.ProtoReflect.Descriptor instead.
func (*Action) Descriptor() ([]byte, []int) {
	return file_script_v1beta1_script_proto_rawDescGZIP(), []int{2}
}

func (x *Action) GetMouseClick() *Action_MouseClick {
	if x != nil {
		return x.MouseClick
	}
	return nil
}

func (x *Action) GetInput() *Action_Input {
	if x != nil {
		return x.Input
	}
	return nil
}

func (x *Action) GetAsk() *Action_PromptUser {
	if x != nil {
		return x.Ask
	}
	return nil
}

type Prompt struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A user-friendly description for what you're asking. This is displayed above the component.
	Description *string `protobuf:"bytes,1,opt,name=description,proto3,oneof" json:"description,omitempty"`
	// A user-friendly and short title for the prompt.
	Title string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	// The component to display to the user.
	Component *Component `protobuf:"bytes,3,opt,name=component,proto3" json:"component,omitempty"`
}

func (x *Prompt) Reset() {
	*x = Prompt{}
	if protoimpl.UnsafeEnabled {
		mi := &file_script_v1beta1_script_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Prompt) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Prompt) ProtoMessage() {}

func (x *Prompt) ProtoReflect() protoreflect.Message {
	mi := &file_script_v1beta1_script_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Prompt.ProtoReflect.Descriptor instead.
func (*Prompt) Descriptor() ([]byte, []int) {
	return file_script_v1beta1_script_proto_rawDescGZIP(), []int{3}
}

func (x *Prompt) GetDescription() string {
	if x != nil && x.Description != nil {
		return *x.Description
	}
	return ""
}

func (x *Prompt) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Prompt) GetComponent() *Component {
	if x != nil {
		return x.Component
	}
	return nil
}

// Mutually exclusive set of components that can be displayed to the user. Only a single field
// can be set.
type Component struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Dropdown *Component_Dropdown `protobuf:"bytes,1,opt,name=dropdown,proto3,oneof" json:"dropdown,omitempty"`
	Input    *Component_Input    `protobuf:"bytes,2,opt,name=input,proto3,oneof" json:"input,omitempty"`
}

func (x *Component) Reset() {
	*x = Component{}
	if protoimpl.UnsafeEnabled {
		mi := &file_script_v1beta1_script_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Component) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Component) ProtoMessage() {}

func (x *Component) ProtoReflect() protoreflect.Message {
	mi := &file_script_v1beta1_script_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Component.ProtoReflect.Descriptor instead.
func (*Component) Descriptor() ([]byte, []int) {
	return file_script_v1beta1_script_proto_rawDescGZIP(), []int{4}
}

func (x *Component) GetDropdown() *Component_Dropdown {
	if x != nil {
		return x.Dropdown
	}
	return nil
}

func (x *Component) GetInput() *Component_Input {
	if x != nil {
		return x.Input
	}
	return nil
}

// A field representing the value of an Action.
type Action_Value struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Hardcode the value for the action.
	Raw string `protobuf:"bytes,1,opt,name=raw,proto3" json:"raw,omitempty"`
}

func (x *Action_Value) Reset() {
	*x = Action_Value{}
	if protoimpl.UnsafeEnabled {
		mi := &file_script_v1beta1_script_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Action_Value) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Action_Value) ProtoMessage() {}

func (x *Action_Value) ProtoReflect() protoreflect.Message {
	mi := &file_script_v1beta1_script_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Action_Value.ProtoReflect.Descriptor instead.
func (*Action_Value) Descriptor() ([]byte, []int) {
	return file_script_v1beta1_script_proto_rawDescGZIP(), []int{2, 0}
}

func (x *Action_Value) GetRaw() string {
	if x != nil {
		return x.Raw
	}
	return ""
}

// Click on an element.
type Action_MouseClick struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Selector string `protobuf:"bytes,1,opt,name=selector,proto3" json:"selector,omitempty"`
}

func (x *Action_MouseClick) Reset() {
	*x = Action_MouseClick{}
	if protoimpl.UnsafeEnabled {
		mi := &file_script_v1beta1_script_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Action_MouseClick) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Action_MouseClick) ProtoMessage() {}

func (x *Action_MouseClick) ProtoReflect() protoreflect.Message {
	mi := &file_script_v1beta1_script_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Action_MouseClick.ProtoReflect.Descriptor instead.
func (*Action_MouseClick) Descriptor() ([]byte, []int) {
	return file_script_v1beta1_script_proto_rawDescGZIP(), []int{2, 1}
}

func (x *Action_MouseClick) GetSelector() string {
	if x != nil {
		return x.Selector
	}
	return ""
}

// Input text into a field.
type Action_Input struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Selector string        `protobuf:"bytes,1,opt,name=selector,proto3" json:"selector,omitempty"`
	Value    *Action_Value `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Action_Input) Reset() {
	*x = Action_Input{}
	if protoimpl.UnsafeEnabled {
		mi := &file_script_v1beta1_script_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Action_Input) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Action_Input) ProtoMessage() {}

func (x *Action_Input) ProtoReflect() protoreflect.Message {
	mi := &file_script_v1beta1_script_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Action_Input.ProtoReflect.Descriptor instead.
func (*Action_Input) Descriptor() ([]byte, []int) {
	return file_script_v1beta1_script_proto_rawDescGZIP(), []int{2, 2}
}

func (x *Action_Input) GetSelector() string {
	if x != nil {
		return x.Selector
	}
	return ""
}

func (x *Action_Input) GetValue() *Action_Value {
	if x != nil {
		return x.Value
	}
	return nil
}

// Prompt the user for some data.
type Action_PromptUser struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Prompt *Prompt `protobuf:"bytes,1,opt,name=prompt,proto3" json:"prompt,omitempty"`
}

func (x *Action_PromptUser) Reset() {
	*x = Action_PromptUser{}
	if protoimpl.UnsafeEnabled {
		mi := &file_script_v1beta1_script_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Action_PromptUser) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Action_PromptUser) ProtoMessage() {}

func (x *Action_PromptUser) ProtoReflect() protoreflect.Message {
	mi := &file_script_v1beta1_script_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Action_PromptUser.ProtoReflect.Descriptor instead.
func (*Action_PromptUser) Descriptor() ([]byte, []int) {
	return file_script_v1beta1_script_proto_rawDescGZIP(), []int{2, 3}
}

func (x *Action_PromptUser) GetPrompt() *Prompt {
	if x != nil {
		return x.Prompt
	}
	return nil
}

// Choose from a set of predefined options.
type Component_Dropdown struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Options []*Component_Dropdown_Option `protobuf:"bytes,1,rep,name=options,proto3" json:"options,omitempty"`
}

func (x *Component_Dropdown) Reset() {
	*x = Component_Dropdown{}
	if protoimpl.UnsafeEnabled {
		mi := &file_script_v1beta1_script_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Component_Dropdown) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Component_Dropdown) ProtoMessage() {}

func (x *Component_Dropdown) ProtoReflect() protoreflect.Message {
	mi := &file_script_v1beta1_script_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Component_Dropdown.ProtoReflect.Descriptor instead.
func (*Component_Dropdown) Descriptor() ([]byte, []int) {
	return file_script_v1beta1_script_proto_rawDescGZIP(), []int{4, 0}
}

func (x *Component_Dropdown) GetOptions() []*Component_Dropdown_Option {
	if x != nil {
		return x.Options
	}
	return nil
}

type Component_Input struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Defaults to TEXT.
	Type Component_Input_Type `protobuf:"varint,1,opt,name=type,proto3,enum=tak.sh.script.v1beta1.Component_Input_Type" json:"type,omitempty"`
}

func (x *Component_Input) Reset() {
	*x = Component_Input{}
	if protoimpl.UnsafeEnabled {
		mi := &file_script_v1beta1_script_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Component_Input) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Component_Input) ProtoMessage() {}

func (x *Component_Input) ProtoReflect() protoreflect.Message {
	mi := &file_script_v1beta1_script_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Component_Input.ProtoReflect.Descriptor instead.
func (*Component_Input) Descriptor() ([]byte, []int) {
	return file_script_v1beta1_script_proto_rawDescGZIP(), []int{4, 1}
}

func (x *Component_Input) GetType() Component_Input_Type {
	if x != nil {
		return x.Type
	}
	return Component_Input_TEXT
}

type Component_Dropdown_Option struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Component_Dropdown_Option) Reset() {
	*x = Component_Dropdown_Option{}
	if protoimpl.UnsafeEnabled {
		mi := &file_script_v1beta1_script_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Component_Dropdown_Option) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Component_Dropdown_Option) ProtoMessage() {}

func (x *Component_Dropdown_Option) ProtoReflect() protoreflect.Message {
	mi := &file_script_v1beta1_script_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Component_Dropdown_Option.ProtoReflect.Descriptor instead.
func (*Component_Dropdown_Option) Descriptor() ([]byte, []int) {
	return file_script_v1beta1_script_proto_rawDescGZIP(), []int{4, 0, 0}
}

func (x *Component_Dropdown_Option) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

var File_script_v1beta1_script_proto protoreflect.FileDescriptor

var file_script_v1beta1_script_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31,
	0x2f, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x74,
	0x61, 0x6b, 0x2e, 0x73, 0x68, 0x2e, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x2e, 0x76, 0x31, 0x62,
	0x65, 0x74, 0x61, 0x31, 0x22, 0x3b, 0x0a, 0x06, 0x53, 0x63, 0x72, 0x69, 0x70, 0x74, 0x12, 0x31,
	0x0a, 0x05, 0x73, 0x74, 0x65, 0x70, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e,
	0x74, 0x61, 0x6b, 0x2e, 0x73, 0x68, 0x2e, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x2e, 0x76, 0x31,
	0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x53, 0x74, 0x65, 0x70, 0x52, 0x05, 0x73, 0x74, 0x65, 0x70,
	0x73, 0x22, 0x75, 0x0a, 0x04, 0x53, 0x74, 0x65, 0x70, 0x12, 0x13, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x02, 0x69, 0x64, 0x88, 0x01, 0x01, 0x12, 0x35,
	0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d,
	0x2e, 0x74, 0x61, 0x6b, 0x2e, 0x73, 0x68, 0x2e, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x2e, 0x76,
	0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x13, 0x0a, 0x02, 0x69, 0x66, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x01, 0x52, 0x02, 0x69, 0x66, 0x88, 0x01, 0x01, 0x42, 0x05, 0x0a, 0x03, 0x5f, 0x69,
	0x64, 0x42, 0x05, 0x0a, 0x03, 0x5f, 0x69, 0x66, 0x22, 0xe6, 0x03, 0x0a, 0x06, 0x41, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x4f, 0x0a, 0x0b, 0x6d, 0x6f, 0x75, 0x73, 0x65, 0x5f, 0x63, 0x6c, 0x69,
	0x63, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x74, 0x61, 0x6b, 0x2e, 0x73,
	0x68, 0x2e, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31,
	0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4d, 0x6f, 0x75, 0x73, 0x65, 0x43, 0x6c, 0x69,
	0x63, 0x6b, 0x48, 0x00, 0x52, 0x0b, 0x6d, 0x6f, 0x75, 0x73, 0x65, 0x5f, 0x63, 0x6c, 0x69, 0x63,
	0x6b, 0x88, 0x01, 0x01, 0x12, 0x3e, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x74, 0x61, 0x6b, 0x2e, 0x73, 0x68, 0x2e, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x41, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x48, 0x01, 0x52, 0x05, 0x69, 0x6e, 0x70, 0x75,
	0x74, 0x88, 0x01, 0x01, 0x12, 0x3f, 0x0a, 0x03, 0x61, 0x73, 0x6b, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x28, 0x2e, 0x74, 0x61, 0x6b, 0x2e, 0x73, 0x68, 0x2e, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x50, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x55, 0x73, 0x65, 0x72, 0x48, 0x02, 0x52, 0x03, 0x61,
	0x73, 0x6b, 0x88, 0x01, 0x01, 0x1a, 0x19, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x10,
	0x0a, 0x03, 0x72, 0x61, 0x77, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x72, 0x61, 0x77,
	0x1a, 0x28, 0x0a, 0x0a, 0x4d, 0x6f, 0x75, 0x73, 0x65, 0x43, 0x6c, 0x69, 0x63, 0x6b, 0x12, 0x1a,
	0x0a, 0x08, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x1a, 0x5e, 0x0a, 0x05, 0x49, 0x6e,
	0x70, 0x75, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x12,
	0x39, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23,
	0x2e, 0x74, 0x61, 0x6b, 0x2e, 0x73, 0x68, 0x2e, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x2e, 0x76,
	0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x43, 0x0a, 0x0a, 0x50, 0x72,
	0x6f, 0x6d, 0x70, 0x74, 0x55, 0x73, 0x65, 0x72, 0x12, 0x35, 0x0a, 0x06, 0x70, 0x72, 0x6f, 0x6d,
	0x70, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x74, 0x61, 0x6b, 0x2e, 0x73,
	0x68, 0x2e, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31,
	0x2e, 0x50, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x52, 0x06, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x42,
	0x0e, 0x0a, 0x0c, 0x5f, 0x6d, 0x6f, 0x75, 0x73, 0x65, 0x5f, 0x63, 0x6c, 0x69, 0x63, 0x6b, 0x42,
	0x08, 0x0a, 0x06, 0x5f, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x61, 0x73,
	0x6b, 0x22, 0x95, 0x01, 0x0a, 0x06, 0x50, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x12, 0x25, 0x0a, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x00, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x88, 0x01, 0x01, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x3e, 0x0a, 0x09, 0x63, 0x6f, 0x6d,
	0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x74,
	0x61, 0x6b, 0x2e, 0x73, 0x68, 0x2e, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x2e, 0x76, 0x31, 0x62,
	0x65, 0x74, 0x61, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x52, 0x09,
	0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x93, 0x03, 0x0a, 0x09, 0x43, 0x6f,
	0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x12, 0x4a, 0x0a, 0x08, 0x64, 0x72, 0x6f, 0x70, 0x64,
	0x6f, 0x77, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x74, 0x61, 0x6b, 0x2e,
	0x73, 0x68, 0x2e, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61,
	0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x2e, 0x44, 0x72, 0x6f, 0x70,
	0x64, 0x6f, 0x77, 0x6e, 0x48, 0x00, 0x52, 0x08, 0x64, 0x72, 0x6f, 0x70, 0x64, 0x6f, 0x77, 0x6e,
	0x88, 0x01, 0x01, 0x12, 0x41, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x26, 0x2e, 0x74, 0x61, 0x6b, 0x2e, 0x73, 0x68, 0x2e, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x6f,
	0x6e, 0x65, 0x6e, 0x74, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x48, 0x01, 0x52, 0x05, 0x69, 0x6e,
	0x70, 0x75, 0x74, 0x88, 0x01, 0x01, 0x1a, 0x76, 0x0a, 0x08, 0x44, 0x72, 0x6f, 0x70, 0x64, 0x6f,
	0x77, 0x6e, 0x12, 0x4a, 0x0a, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x30, 0x2e, 0x74, 0x61, 0x6b, 0x2e, 0x73, 0x68, 0x2e, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x70,
	0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x2e, 0x44, 0x72, 0x6f, 0x70, 0x64, 0x6f, 0x77, 0x6e, 0x2e, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1a, 0x1e,
	0x0a, 0x06, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x68,
	0x0a, 0x05, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x3f, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2b, 0x2e, 0x74, 0x61, 0x6b, 0x2e, 0x73, 0x68, 0x2e, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x43, 0x6f,
	0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x2e, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x1e, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x08, 0x0a, 0x04, 0x54, 0x45, 0x58, 0x54, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x50, 0x41,
	0x53, 0x53, 0x57, 0x4f, 0x52, 0x44, 0x10, 0x01, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x64, 0x72, 0x6f,
	0x70, 0x64, 0x6f, 0x77, 0x6e, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x42,
	0xd2, 0x01, 0x0a, 0x19, 0x63, 0x6f, 0x6d, 0x2e, 0x74, 0x61, 0x6b, 0x2e, 0x73, 0x68, 0x2e, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x42, 0x0b, 0x53,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x31, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x61, 0x6b, 0x2d, 0x73, 0x68, 0x2f,
	0x74, 0x61, 0x6b, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2f, 0x67, 0x6f,
	0x2f, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0xa2,
	0x02, 0x03, 0x54, 0x53, 0x53, 0xaa, 0x02, 0x15, 0x54, 0x61, 0x6b, 0x2e, 0x53, 0x68, 0x2e, 0x53,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x2e, 0x56, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0xca, 0x02, 0x15,
	0x54, 0x61, 0x6b, 0x5c, 0x53, 0x68, 0x5c, 0x53, 0x63, 0x72, 0x69, 0x70, 0x74, 0x5c, 0x56, 0x31,
	0x62, 0x65, 0x74, 0x61, 0x31, 0xe2, 0x02, 0x21, 0x54, 0x61, 0x6b, 0x5c, 0x53, 0x68, 0x5c, 0x53,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x5c, 0x56, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x5c, 0x47, 0x50,
	0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x18, 0x54, 0x61, 0x6b, 0x3a,
	0x3a, 0x53, 0x68, 0x3a, 0x3a, 0x53, 0x63, 0x72, 0x69, 0x70, 0x74, 0x3a, 0x3a, 0x56, 0x31, 0x62,
	0x65, 0x74, 0x61, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_script_v1beta1_script_proto_rawDescOnce sync.Once
	file_script_v1beta1_script_proto_rawDescData = file_script_v1beta1_script_proto_rawDesc
)

func file_script_v1beta1_script_proto_rawDescGZIP() []byte {
	file_script_v1beta1_script_proto_rawDescOnce.Do(func() {
		file_script_v1beta1_script_proto_rawDescData = protoimpl.X.CompressGZIP(file_script_v1beta1_script_proto_rawDescData)
	})
	return file_script_v1beta1_script_proto_rawDescData
}

var file_script_v1beta1_script_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_script_v1beta1_script_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_script_v1beta1_script_proto_goTypes = []interface{}{
	(Component_Input_Type)(0),         // 0: tak.sh.script.v1beta1.Component.Input.Type
	(*Script)(nil),                    // 1: tak.sh.script.v1beta1.Script
	(*Step)(nil),                      // 2: tak.sh.script.v1beta1.Step
	(*Action)(nil),                    // 3: tak.sh.script.v1beta1.Action
	(*Prompt)(nil),                    // 4: tak.sh.script.v1beta1.Prompt
	(*Component)(nil),                 // 5: tak.sh.script.v1beta1.Component
	(*Action_Value)(nil),              // 6: tak.sh.script.v1beta1.Action.Value
	(*Action_MouseClick)(nil),         // 7: tak.sh.script.v1beta1.Action.MouseClick
	(*Action_Input)(nil),              // 8: tak.sh.script.v1beta1.Action.Input
	(*Action_PromptUser)(nil),         // 9: tak.sh.script.v1beta1.Action.PromptUser
	(*Component_Dropdown)(nil),        // 10: tak.sh.script.v1beta1.Component.Dropdown
	(*Component_Input)(nil),           // 11: tak.sh.script.v1beta1.Component.Input
	(*Component_Dropdown_Option)(nil), // 12: tak.sh.script.v1beta1.Component.Dropdown.Option
}
var file_script_v1beta1_script_proto_depIdxs = []int32{
	2,  // 0: tak.sh.script.v1beta1.Script.steps:type_name -> tak.sh.script.v1beta1.Step
	3,  // 1: tak.sh.script.v1beta1.Step.action:type_name -> tak.sh.script.v1beta1.Action
	7,  // 2: tak.sh.script.v1beta1.Action.mouse_click:type_name -> tak.sh.script.v1beta1.Action.MouseClick
	8,  // 3: tak.sh.script.v1beta1.Action.input:type_name -> tak.sh.script.v1beta1.Action.Input
	9,  // 4: tak.sh.script.v1beta1.Action.ask:type_name -> tak.sh.script.v1beta1.Action.PromptUser
	5,  // 5: tak.sh.script.v1beta1.Prompt.component:type_name -> tak.sh.script.v1beta1.Component
	10, // 6: tak.sh.script.v1beta1.Component.dropdown:type_name -> tak.sh.script.v1beta1.Component.Dropdown
	11, // 7: tak.sh.script.v1beta1.Component.input:type_name -> tak.sh.script.v1beta1.Component.Input
	6,  // 8: tak.sh.script.v1beta1.Action.Input.value:type_name -> tak.sh.script.v1beta1.Action.Value
	4,  // 9: tak.sh.script.v1beta1.Action.PromptUser.prompt:type_name -> tak.sh.script.v1beta1.Prompt
	12, // 10: tak.sh.script.v1beta1.Component.Dropdown.options:type_name -> tak.sh.script.v1beta1.Component.Dropdown.Option
	0,  // 11: tak.sh.script.v1beta1.Component.Input.type:type_name -> tak.sh.script.v1beta1.Component.Input.Type
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_script_v1beta1_script_proto_init() }
func file_script_v1beta1_script_proto_init() {
	if File_script_v1beta1_script_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_script_v1beta1_script_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Script); i {
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
		file_script_v1beta1_script_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Step); i {
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
		file_script_v1beta1_script_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Action); i {
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
		file_script_v1beta1_script_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Prompt); i {
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
		file_script_v1beta1_script_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Component); i {
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
		file_script_v1beta1_script_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Action_Value); i {
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
		file_script_v1beta1_script_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Action_MouseClick); i {
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
		file_script_v1beta1_script_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Action_Input); i {
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
		file_script_v1beta1_script_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Action_PromptUser); i {
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
		file_script_v1beta1_script_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Component_Dropdown); i {
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
		file_script_v1beta1_script_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Component_Input); i {
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
		file_script_v1beta1_script_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Component_Dropdown_Option); i {
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
	file_script_v1beta1_script_proto_msgTypes[1].OneofWrappers = []interface{}{}
	file_script_v1beta1_script_proto_msgTypes[2].OneofWrappers = []interface{}{}
	file_script_v1beta1_script_proto_msgTypes[3].OneofWrappers = []interface{}{}
	file_script_v1beta1_script_proto_msgTypes[4].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_script_v1beta1_script_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_script_v1beta1_script_proto_goTypes,
		DependencyIndexes: file_script_v1beta1_script_proto_depIdxs,
		EnumInfos:         file_script_v1beta1_script_proto_enumTypes,
		MessageInfos:      file_script_v1beta1_script_proto_msgTypes,
	}.Build()
	File_script_v1beta1_script_proto = out.File
	file_script_v1beta1_script_proto_rawDesc = nil
	file_script_v1beta1_script_proto_goTypes = nil
	file_script_v1beta1_script_proto_depIdxs = nil
}