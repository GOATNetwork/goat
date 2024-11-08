// Code generated by protoc-gen-go-pulsar. DO NOT EDIT.
package bitcoinv1

import (
	_ "cosmossdk.io/api/amino"
	fmt "fmt"
	runtime "github.com/cosmos/cosmos-proto/runtime"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoiface "google.golang.org/protobuf/runtime/protoiface"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	io "io"
	reflect "reflect"
	sync "sync"
)

var (
	md_Params                      protoreflect.MessageDescriptor
	fd_Params_network_name         protoreflect.FieldDescriptor
	fd_Params_confirmation_number  protoreflect.FieldDescriptor
	fd_Params_min_deposit_amount   protoreflect.FieldDescriptor
	fd_Params_deposit_magic_prefix protoreflect.FieldDescriptor
	fd_Params_deposit_tax_rate     protoreflect.FieldDescriptor
	fd_Params_max_deposit_tax      protoreflect.FieldDescriptor
)

func init() {
	file_goat_bitcoin_v1_params_proto_init()
	md_Params = File_goat_bitcoin_v1_params_proto.Messages().ByName("Params")
	fd_Params_network_name = md_Params.Fields().ByName("network_name")
	fd_Params_confirmation_number = md_Params.Fields().ByName("confirmation_number")
	fd_Params_min_deposit_amount = md_Params.Fields().ByName("min_deposit_amount")
	fd_Params_deposit_magic_prefix = md_Params.Fields().ByName("deposit_magic_prefix")
	fd_Params_deposit_tax_rate = md_Params.Fields().ByName("deposit_tax_rate")
	fd_Params_max_deposit_tax = md_Params.Fields().ByName("max_deposit_tax")
}

var _ protoreflect.Message = (*fastReflection_Params)(nil)

type fastReflection_Params Params

func (x *Params) ProtoReflect() protoreflect.Message {
	return (*fastReflection_Params)(x)
}

func (x *Params) slowProtoReflect() protoreflect.Message {
	mi := &file_goat_bitcoin_v1_params_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

var _fastReflection_Params_messageType fastReflection_Params_messageType
var _ protoreflect.MessageType = fastReflection_Params_messageType{}

type fastReflection_Params_messageType struct{}

func (x fastReflection_Params_messageType) Zero() protoreflect.Message {
	return (*fastReflection_Params)(nil)
}
func (x fastReflection_Params_messageType) New() protoreflect.Message {
	return new(fastReflection_Params)
}
func (x fastReflection_Params_messageType) Descriptor() protoreflect.MessageDescriptor {
	return md_Params
}

// Descriptor returns message descriptor, which contains only the protobuf
// type information for the message.
func (x *fastReflection_Params) Descriptor() protoreflect.MessageDescriptor {
	return md_Params
}

// Type returns the message type, which encapsulates both Go and protobuf
// type information. If the Go type information is not needed,
// it is recommended that the message descriptor be used instead.
func (x *fastReflection_Params) Type() protoreflect.MessageType {
	return _fastReflection_Params_messageType
}

// New returns a newly allocated and mutable empty message.
func (x *fastReflection_Params) New() protoreflect.Message {
	return new(fastReflection_Params)
}

// Interface unwraps the message reflection interface and
// returns the underlying ProtoMessage interface.
func (x *fastReflection_Params) Interface() protoreflect.ProtoMessage {
	return (*Params)(x)
}

// Range iterates over every populated field in an undefined order,
// calling f for each field descriptor and value encountered.
// Range returns immediately if f returns false.
// While iterating, mutating operations may only be performed
// on the current field descriptor.
func (x *fastReflection_Params) Range(f func(protoreflect.FieldDescriptor, protoreflect.Value) bool) {
	if x.NetworkName != "" {
		value := protoreflect.ValueOfString(x.NetworkName)
		if !f(fd_Params_network_name, value) {
			return
		}
	}
	if x.ConfirmationNumber != uint64(0) {
		value := protoreflect.ValueOfUint64(x.ConfirmationNumber)
		if !f(fd_Params_confirmation_number, value) {
			return
		}
	}
	if x.MinDepositAmount != uint64(0) {
		value := protoreflect.ValueOfUint64(x.MinDepositAmount)
		if !f(fd_Params_min_deposit_amount, value) {
			return
		}
	}
	if len(x.DepositMagicPrefix) != 0 {
		value := protoreflect.ValueOfBytes(x.DepositMagicPrefix)
		if !f(fd_Params_deposit_magic_prefix, value) {
			return
		}
	}
	if x.DepositTaxRate != uint64(0) {
		value := protoreflect.ValueOfUint64(x.DepositTaxRate)
		if !f(fd_Params_deposit_tax_rate, value) {
			return
		}
	}
	if x.MaxDepositTax != uint64(0) {
		value := protoreflect.ValueOfUint64(x.MaxDepositTax)
		if !f(fd_Params_max_deposit_tax, value) {
			return
		}
	}
}

// Has reports whether a field is populated.
//
// Some fields have the property of nullability where it is possible to
// distinguish between the default value of a field and whether the field
// was explicitly populated with the default value. Singular message fields,
// member fields of a oneof, and proto2 scalar fields are nullable. Such
// fields are populated only if explicitly set.
//
// In other cases (aside from the nullable cases above),
// a proto3 scalar field is populated if it contains a non-zero value, and
// a repeated field is populated if it is non-empty.
func (x *fastReflection_Params) Has(fd protoreflect.FieldDescriptor) bool {
	switch fd.FullName() {
	case "goat.bitcoin.v1.Params.network_name":
		return x.NetworkName != ""
	case "goat.bitcoin.v1.Params.confirmation_number":
		return x.ConfirmationNumber != uint64(0)
	case "goat.bitcoin.v1.Params.min_deposit_amount":
		return x.MinDepositAmount != uint64(0)
	case "goat.bitcoin.v1.Params.deposit_magic_prefix":
		return len(x.DepositMagicPrefix) != 0
	case "goat.bitcoin.v1.Params.deposit_tax_rate":
		return x.DepositTaxRate != uint64(0)
	case "goat.bitcoin.v1.Params.max_deposit_tax":
		return x.MaxDepositTax != uint64(0)
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: goat.bitcoin.v1.Params"))
		}
		panic(fmt.Errorf("message goat.bitcoin.v1.Params does not contain field %s", fd.FullName()))
	}
}

// Clear clears the field such that a subsequent Has call reports false.
//
// Clearing an extension field clears both the extension type and value
// associated with the given field number.
//
// Clear is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_Params) Clear(fd protoreflect.FieldDescriptor) {
	switch fd.FullName() {
	case "goat.bitcoin.v1.Params.network_name":
		x.NetworkName = ""
	case "goat.bitcoin.v1.Params.confirmation_number":
		x.ConfirmationNumber = uint64(0)
	case "goat.bitcoin.v1.Params.min_deposit_amount":
		x.MinDepositAmount = uint64(0)
	case "goat.bitcoin.v1.Params.deposit_magic_prefix":
		x.DepositMagicPrefix = nil
	case "goat.bitcoin.v1.Params.deposit_tax_rate":
		x.DepositTaxRate = uint64(0)
	case "goat.bitcoin.v1.Params.max_deposit_tax":
		x.MaxDepositTax = uint64(0)
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: goat.bitcoin.v1.Params"))
		}
		panic(fmt.Errorf("message goat.bitcoin.v1.Params does not contain field %s", fd.FullName()))
	}
}

// Get retrieves the value for a field.
//
// For unpopulated scalars, it returns the default value, where
// the default value of a bytes scalar is guaranteed to be a copy.
// For unpopulated composite types, it returns an empty, read-only view
// of the value; to obtain a mutable reference, use Mutable.
func (x *fastReflection_Params) Get(descriptor protoreflect.FieldDescriptor) protoreflect.Value {
	switch descriptor.FullName() {
	case "goat.bitcoin.v1.Params.network_name":
		value := x.NetworkName
		return protoreflect.ValueOfString(value)
	case "goat.bitcoin.v1.Params.confirmation_number":
		value := x.ConfirmationNumber
		return protoreflect.ValueOfUint64(value)
	case "goat.bitcoin.v1.Params.min_deposit_amount":
		value := x.MinDepositAmount
		return protoreflect.ValueOfUint64(value)
	case "goat.bitcoin.v1.Params.deposit_magic_prefix":
		value := x.DepositMagicPrefix
		return protoreflect.ValueOfBytes(value)
	case "goat.bitcoin.v1.Params.deposit_tax_rate":
		value := x.DepositTaxRate
		return protoreflect.ValueOfUint64(value)
	case "goat.bitcoin.v1.Params.max_deposit_tax":
		value := x.MaxDepositTax
		return protoreflect.ValueOfUint64(value)
	default:
		if descriptor.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: goat.bitcoin.v1.Params"))
		}
		panic(fmt.Errorf("message goat.bitcoin.v1.Params does not contain field %s", descriptor.FullName()))
	}
}

// Set stores the value for a field.
//
// For a field belonging to a oneof, it implicitly clears any other field
// that may be currently set within the same oneof.
// For extension fields, it implicitly stores the provided ExtensionType.
// When setting a composite type, it is unspecified whether the stored value
// aliases the source's memory in any way. If the composite value is an
// empty, read-only value, then it panics.
//
// Set is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_Params) Set(fd protoreflect.FieldDescriptor, value protoreflect.Value) {
	switch fd.FullName() {
	case "goat.bitcoin.v1.Params.network_name":
		x.NetworkName = value.Interface().(string)
	case "goat.bitcoin.v1.Params.confirmation_number":
		x.ConfirmationNumber = value.Uint()
	case "goat.bitcoin.v1.Params.min_deposit_amount":
		x.MinDepositAmount = value.Uint()
	case "goat.bitcoin.v1.Params.deposit_magic_prefix":
		x.DepositMagicPrefix = value.Bytes()
	case "goat.bitcoin.v1.Params.deposit_tax_rate":
		x.DepositTaxRate = value.Uint()
	case "goat.bitcoin.v1.Params.max_deposit_tax":
		x.MaxDepositTax = value.Uint()
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: goat.bitcoin.v1.Params"))
		}
		panic(fmt.Errorf("message goat.bitcoin.v1.Params does not contain field %s", fd.FullName()))
	}
}

// Mutable returns a mutable reference to a composite type.
//
// If the field is unpopulated, it may allocate a composite value.
// For a field belonging to a oneof, it implicitly clears any other field
// that may be currently set within the same oneof.
// For extension fields, it implicitly stores the provided ExtensionType
// if not already stored.
// It panics if the field does not contain a composite type.
//
// Mutable is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_Params) Mutable(fd protoreflect.FieldDescriptor) protoreflect.Value {
	switch fd.FullName() {
	case "goat.bitcoin.v1.Params.network_name":
		panic(fmt.Errorf("field network_name of message goat.bitcoin.v1.Params is not mutable"))
	case "goat.bitcoin.v1.Params.confirmation_number":
		panic(fmt.Errorf("field confirmation_number of message goat.bitcoin.v1.Params is not mutable"))
	case "goat.bitcoin.v1.Params.min_deposit_amount":
		panic(fmt.Errorf("field min_deposit_amount of message goat.bitcoin.v1.Params is not mutable"))
	case "goat.bitcoin.v1.Params.deposit_magic_prefix":
		panic(fmt.Errorf("field deposit_magic_prefix of message goat.bitcoin.v1.Params is not mutable"))
	case "goat.bitcoin.v1.Params.deposit_tax_rate":
		panic(fmt.Errorf("field deposit_tax_rate of message goat.bitcoin.v1.Params is not mutable"))
	case "goat.bitcoin.v1.Params.max_deposit_tax":
		panic(fmt.Errorf("field max_deposit_tax of message goat.bitcoin.v1.Params is not mutable"))
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: goat.bitcoin.v1.Params"))
		}
		panic(fmt.Errorf("message goat.bitcoin.v1.Params does not contain field %s", fd.FullName()))
	}
}

// NewField returns a new value that is assignable to the field
// for the given descriptor. For scalars, this returns the default value.
// For lists, maps, and messages, this returns a new, empty, mutable value.
func (x *fastReflection_Params) NewField(fd protoreflect.FieldDescriptor) protoreflect.Value {
	switch fd.FullName() {
	case "goat.bitcoin.v1.Params.network_name":
		return protoreflect.ValueOfString("")
	case "goat.bitcoin.v1.Params.confirmation_number":
		return protoreflect.ValueOfUint64(uint64(0))
	case "goat.bitcoin.v1.Params.min_deposit_amount":
		return protoreflect.ValueOfUint64(uint64(0))
	case "goat.bitcoin.v1.Params.deposit_magic_prefix":
		return protoreflect.ValueOfBytes(nil)
	case "goat.bitcoin.v1.Params.deposit_tax_rate":
		return protoreflect.ValueOfUint64(uint64(0))
	case "goat.bitcoin.v1.Params.max_deposit_tax":
		return protoreflect.ValueOfUint64(uint64(0))
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: goat.bitcoin.v1.Params"))
		}
		panic(fmt.Errorf("message goat.bitcoin.v1.Params does not contain field %s", fd.FullName()))
	}
}

// WhichOneof reports which field within the oneof is populated,
// returning nil if none are populated.
// It panics if the oneof descriptor does not belong to this message.
func (x *fastReflection_Params) WhichOneof(d protoreflect.OneofDescriptor) protoreflect.FieldDescriptor {
	switch d.FullName() {
	default:
		panic(fmt.Errorf("%s is not a oneof field in goat.bitcoin.v1.Params", d.FullName()))
	}
	panic("unreachable")
}

// GetUnknown retrieves the entire list of unknown fields.
// The caller may only mutate the contents of the RawFields
// if the mutated bytes are stored back into the message with SetUnknown.
func (x *fastReflection_Params) GetUnknown() protoreflect.RawFields {
	return x.unknownFields
}

// SetUnknown stores an entire list of unknown fields.
// The raw fields must be syntactically valid according to the wire format.
// An implementation may panic if this is not the case.
// Once stored, the caller must not mutate the content of the RawFields.
// An empty RawFields may be passed to clear the fields.
//
// SetUnknown is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_Params) SetUnknown(fields protoreflect.RawFields) {
	x.unknownFields = fields
}

// IsValid reports whether the message is valid.
//
// An invalid message is an empty, read-only value.
//
// An invalid message often corresponds to a nil pointer of the concrete
// message type, but the details are implementation dependent.
// Validity is not part of the protobuf data model, and may not
// be preserved in marshaling or other operations.
func (x *fastReflection_Params) IsValid() bool {
	return x != nil
}

// ProtoMethods returns optional fastReflectionFeature-path implementations of various operations.
// This method may return nil.
//
// The returned methods type is identical to
// "google.golang.org/protobuf/runtime/protoiface".Methods.
// Consult the protoiface package documentation for details.
func (x *fastReflection_Params) ProtoMethods() *protoiface.Methods {
	size := func(input protoiface.SizeInput) protoiface.SizeOutput {
		x := input.Message.Interface().(*Params)
		if x == nil {
			return protoiface.SizeOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Size:              0,
			}
		}
		options := runtime.SizeInputToOptions(input)
		_ = options
		var n int
		var l int
		_ = l
		l = len(x.NetworkName)
		if l > 0 {
			n += 1 + l + runtime.Sov(uint64(l))
		}
		if x.ConfirmationNumber != 0 {
			n += 1 + runtime.Sov(uint64(x.ConfirmationNumber))
		}
		if x.MinDepositAmount != 0 {
			n += 1 + runtime.Sov(uint64(x.MinDepositAmount))
		}
		l = len(x.DepositMagicPrefix)
		if l > 0 {
			n += 1 + l + runtime.Sov(uint64(l))
		}
		if x.DepositTaxRate != 0 {
			n += 1 + runtime.Sov(uint64(x.DepositTaxRate))
		}
		if x.MaxDepositTax != 0 {
			n += 1 + runtime.Sov(uint64(x.MaxDepositTax))
		}
		if x.unknownFields != nil {
			n += len(x.unknownFields)
		}
		return protoiface.SizeOutput{
			NoUnkeyedLiterals: input.NoUnkeyedLiterals,
			Size:              n,
		}
	}

	marshal := func(input protoiface.MarshalInput) (protoiface.MarshalOutput, error) {
		x := input.Message.Interface().(*Params)
		if x == nil {
			return protoiface.MarshalOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Buf:               input.Buf,
			}, nil
		}
		options := runtime.MarshalInputToOptions(input)
		_ = options
		size := options.Size(x)
		dAtA := make([]byte, size)
		i := len(dAtA)
		_ = i
		var l int
		_ = l
		if x.unknownFields != nil {
			i -= len(x.unknownFields)
			copy(dAtA[i:], x.unknownFields)
		}
		if x.MaxDepositTax != 0 {
			i = runtime.EncodeVarint(dAtA, i, uint64(x.MaxDepositTax))
			i--
			dAtA[i] = 0x30
		}
		if x.DepositTaxRate != 0 {
			i = runtime.EncodeVarint(dAtA, i, uint64(x.DepositTaxRate))
			i--
			dAtA[i] = 0x28
		}
		if len(x.DepositMagicPrefix) > 0 {
			i -= len(x.DepositMagicPrefix)
			copy(dAtA[i:], x.DepositMagicPrefix)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.DepositMagicPrefix)))
			i--
			dAtA[i] = 0x22
		}
		if x.MinDepositAmount != 0 {
			i = runtime.EncodeVarint(dAtA, i, uint64(x.MinDepositAmount))
			i--
			dAtA[i] = 0x18
		}
		if x.ConfirmationNumber != 0 {
			i = runtime.EncodeVarint(dAtA, i, uint64(x.ConfirmationNumber))
			i--
			dAtA[i] = 0x10
		}
		if len(x.NetworkName) > 0 {
			i -= len(x.NetworkName)
			copy(dAtA[i:], x.NetworkName)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.NetworkName)))
			i--
			dAtA[i] = 0xa
		}
		if input.Buf != nil {
			input.Buf = append(input.Buf, dAtA...)
		} else {
			input.Buf = dAtA
		}
		return protoiface.MarshalOutput{
			NoUnkeyedLiterals: input.NoUnkeyedLiterals,
			Buf:               input.Buf,
		}, nil
	}
	unmarshal := func(input protoiface.UnmarshalInput) (protoiface.UnmarshalOutput, error) {
		x := input.Message.Interface().(*Params)
		if x == nil {
			return protoiface.UnmarshalOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Flags:             input.Flags,
			}, nil
		}
		options := runtime.UnmarshalInputToOptions(input)
		_ = options
		dAtA := input.Buf
		l := len(dAtA)
		iNdEx := 0
		for iNdEx < l {
			preIndex := iNdEx
			var wire uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
				}
				if iNdEx >= l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				wire |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			fieldNum := int32(wire >> 3)
			wireType := int(wire & 0x7)
			if wireType == 4 {
				return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: Params: wiretype end group for non-group")
			}
			if fieldNum <= 0 {
				return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
			}
			switch fieldNum {
			case 1:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field NetworkName", wireType)
				}
				var stringLen uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					stringLen |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				intStringLen := int(stringLen)
				if intStringLen < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + intStringLen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				x.NetworkName = string(dAtA[iNdEx:postIndex])
				iNdEx = postIndex
			case 2:
				if wireType != 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field ConfirmationNumber", wireType)
				}
				x.ConfirmationNumber = 0
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					x.ConfirmationNumber |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
			case 3:
				if wireType != 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field MinDepositAmount", wireType)
				}
				x.MinDepositAmount = 0
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					x.MinDepositAmount |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
			case 4:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field DepositMagicPrefix", wireType)
				}
				var byteLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					byteLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if byteLen < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + byteLen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				x.DepositMagicPrefix = append(x.DepositMagicPrefix[:0], dAtA[iNdEx:postIndex]...)
				if x.DepositMagicPrefix == nil {
					x.DepositMagicPrefix = []byte{}
				}
				iNdEx = postIndex
			case 5:
				if wireType != 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field DepositTaxRate", wireType)
				}
				x.DepositTaxRate = 0
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					x.DepositTaxRate |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
			case 6:
				if wireType != 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field MaxDepositTax", wireType)
				}
				x.MaxDepositTax = 0
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					x.MaxDepositTax |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
			default:
				iNdEx = preIndex
				skippy, err := runtime.Skip(dAtA[iNdEx:])
				if err != nil {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, err
				}
				if (skippy < 0) || (iNdEx+skippy) < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if (iNdEx + skippy) > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				if !options.DiscardUnknown {
					x.unknownFields = append(x.unknownFields, dAtA[iNdEx:iNdEx+skippy]...)
				}
				iNdEx += skippy
			}
		}

		if iNdEx > l {
			return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
		}
		return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, nil
	}
	return &protoiface.Methods{
		NoUnkeyedLiterals: struct{}{},
		Flags:             protoiface.SupportMarshalDeterministic | protoiface.SupportUnmarshalDiscardUnknown,
		Size:              size,
		Marshal:           marshal,
		Unmarshal:         unmarshal,
		Merge:             nil,
		CheckInitialized:  nil,
	}
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.0
// 	protoc        (unknown)
// source: goat/bitcoin/v1/params.proto

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Params defines the parameters for the module.
type Params struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NetworkName string `protobuf:"bytes,1,opt,name=network_name,json=networkName,proto3" json:"network_name,omitempty"`
	// A block has the number should be considered as a finalized block
	ConfirmationNumber uint64 `protobuf:"varint,2,opt,name=confirmation_number,json=confirmationNumber,proto3" json:"confirmation_number,omitempty"`
	// min deposit amount in satoshi
	MinDepositAmount   uint64 `protobuf:"varint,3,opt,name=min_deposit_amount,json=minDepositAmount,proto3" json:"min_deposit_amount,omitempty"`
	DepositMagicPrefix []byte `protobuf:"bytes,4,opt,name=deposit_magic_prefix,json=depositMagicPrefix,proto3" json:"deposit_magic_prefix,omitempty"`
	DepositTaxRate     uint64 `protobuf:"varint,5,opt,name=deposit_tax_rate,json=depositTaxRate,proto3" json:"deposit_tax_rate,omitempty"`
	// max deposit tax in satoshi
	MaxDepositTax uint64 `protobuf:"varint,6,opt,name=max_deposit_tax,json=maxDepositTax,proto3" json:"max_deposit_tax,omitempty"`
}

func (x *Params) Reset() {
	*x = Params{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goat_bitcoin_v1_params_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Params) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Params) ProtoMessage() {}

// Deprecated: Use Params.ProtoReflect.Descriptor instead.
func (*Params) Descriptor() ([]byte, []int) {
	return file_goat_bitcoin_v1_params_proto_rawDescGZIP(), []int{0}
}

func (x *Params) GetNetworkName() string {
	if x != nil {
		return x.NetworkName
	}
	return ""
}

func (x *Params) GetConfirmationNumber() uint64 {
	if x != nil {
		return x.ConfirmationNumber
	}
	return 0
}

func (x *Params) GetMinDepositAmount() uint64 {
	if x != nil {
		return x.MinDepositAmount
	}
	return 0
}

func (x *Params) GetDepositMagicPrefix() []byte {
	if x != nil {
		return x.DepositMagicPrefix
	}
	return nil
}

func (x *Params) GetDepositTaxRate() uint64 {
	if x != nil {
		return x.DepositTaxRate
	}
	return 0
}

func (x *Params) GetMaxDepositTax() uint64 {
	if x != nil {
		return x.MaxDepositTax
	}
	return 0
}

var File_goat_bitcoin_v1_params_proto protoreflect.FileDescriptor

var file_goat_bitcoin_v1_params_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x67, 0x6f, 0x61, 0x74, 0x2f, 0x62, 0x69, 0x74, 0x63, 0x6f, 0x69, 0x6e, 0x2f, 0x76,
	0x31, 0x2f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f,
	0x67, 0x6f, 0x61, 0x74, 0x2e, 0x62, 0x69, 0x74, 0x63, 0x6f, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x1a,
	0x11, 0x61, 0x6d, 0x69, 0x6e, 0x6f, 0x2f, 0x61, 0x6d, 0x69, 0x6e, 0x6f, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xaa, 0x02, 0x0a, 0x06, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x21, 0x0a,
	0x0c, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x2f, 0x0a, 0x13, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x12, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x12, 0x2c, 0x0a, 0x12, 0x6d, 0x69, 0x6e, 0x5f, 0x64, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74,
	0x5f, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x10, 0x6d,
	0x69, 0x6e, 0x44, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12,
	0x30, 0x0a, 0x14, 0x64, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x5f, 0x6d, 0x61, 0x67, 0x69, 0x63,
	0x5f, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x12, 0x64,
	0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x4d, 0x61, 0x67, 0x69, 0x63, 0x50, 0x72, 0x65, 0x66, 0x69,
	0x78, 0x12, 0x28, 0x0a, 0x10, 0x64, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x5f, 0x74, 0x61, 0x78,
	0x5f, 0x72, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0e, 0x64, 0x65, 0x70,
	0x6f, 0x73, 0x69, 0x74, 0x54, 0x61, 0x78, 0x52, 0x61, 0x74, 0x65, 0x12, 0x26, 0x0a, 0x0f, 0x6d,
	0x61, 0x78, 0x5f, 0x64, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x5f, 0x74, 0x61, 0x78, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x0d, 0x6d, 0x61, 0x78, 0x44, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74,
	0x54, 0x61, 0x78, 0x3a, 0x1a, 0x8a, 0xe7, 0xb0, 0x2a, 0x15, 0x67, 0x6f, 0x61, 0x74, 0x2f, 0x78,
	0x2f, 0x62, 0x69, 0x74, 0x63, 0x6f, 0x69, 0x6e, 0x2f, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x42,
	0xbb, 0x01, 0x0a, 0x13, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x6f, 0x61, 0x74, 0x2e, 0x62, 0x69, 0x74,
	0x63, 0x6f, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x42, 0x0b, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x39, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x61, 0x74, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x67,
	0x6f, 0x61, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x6f, 0x61, 0x74, 0x2f, 0x62, 0x69, 0x74,
	0x63, 0x6f, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x3b, 0x62, 0x69, 0x74, 0x63, 0x6f, 0x69, 0x6e, 0x76,
	0x31, 0xa2, 0x02, 0x03, 0x47, 0x42, 0x58, 0xaa, 0x02, 0x0f, 0x47, 0x6f, 0x61, 0x74, 0x2e, 0x42,
	0x69, 0x74, 0x63, 0x6f, 0x69, 0x6e, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0f, 0x47, 0x6f, 0x61, 0x74,
	0x5c, 0x42, 0x69, 0x74, 0x63, 0x6f, 0x69, 0x6e, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1b, 0x47, 0x6f,
	0x61, 0x74, 0x5c, 0x42, 0x69, 0x74, 0x63, 0x6f, 0x69, 0x6e, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50,
	0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x11, 0x47, 0x6f, 0x61, 0x74,
	0x3a, 0x3a, 0x42, 0x69, 0x74, 0x63, 0x6f, 0x69, 0x6e, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_goat_bitcoin_v1_params_proto_rawDescOnce sync.Once
	file_goat_bitcoin_v1_params_proto_rawDescData = file_goat_bitcoin_v1_params_proto_rawDesc
)

func file_goat_bitcoin_v1_params_proto_rawDescGZIP() []byte {
	file_goat_bitcoin_v1_params_proto_rawDescOnce.Do(func() {
		file_goat_bitcoin_v1_params_proto_rawDescData = protoimpl.X.CompressGZIP(file_goat_bitcoin_v1_params_proto_rawDescData)
	})
	return file_goat_bitcoin_v1_params_proto_rawDescData
}

var file_goat_bitcoin_v1_params_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_goat_bitcoin_v1_params_proto_goTypes = []interface{}{
	(*Params)(nil), // 0: goat.bitcoin.v1.Params
}
var file_goat_bitcoin_v1_params_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_goat_bitcoin_v1_params_proto_init() }
func file_goat_bitcoin_v1_params_proto_init() {
	if File_goat_bitcoin_v1_params_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_goat_bitcoin_v1_params_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Params); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_goat_bitcoin_v1_params_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_goat_bitcoin_v1_params_proto_goTypes,
		DependencyIndexes: file_goat_bitcoin_v1_params_proto_depIdxs,
		MessageInfos:      file_goat_bitcoin_v1_params_proto_msgTypes,
	}.Build()
	File_goat_bitcoin_v1_params_proto = out.File
	file_goat_bitcoin_v1_params_proto_rawDesc = nil
	file_goat_bitcoin_v1_params_proto_goTypes = nil
	file_goat_bitcoin_v1_params_proto_depIdxs = nil
}
