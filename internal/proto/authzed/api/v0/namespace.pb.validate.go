// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: authzed/api/v0/namespace.proto

package v0

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
)

// Validate checks the field values on Metadata with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *Metadata) Validate() error {
	if m == nil {
		return nil
	}

	if len(m.GetMetadataMessage()) < 1 {
		return MetadataValidationError{
			field:  "MetadataMessage",
			reason: "value must contain at least 1 item(s)",
		}
	}

	for idx, item := range m.GetMetadataMessage() {
		_, _ = idx, item

		if item == nil {
			return MetadataValidationError{
				field:  fmt.Sprintf("MetadataMessage[%v]", idx),
				reason: "value is required",
			}
		}

		if a := item; a != nil {

			if _, ok := _Metadata_MetadataMessage_InLookup[a.GetTypeUrl()]; !ok {
				return MetadataValidationError{
					field:  fmt.Sprintf("MetadataMessage[%v]", idx),
					reason: "type URL must be in list [type.googleapis.com/impl.v1.DocComment type.googleapis.com/impl.v1.RelationMetadata]",
				}
			}

		}

	}

	return nil
}

// MetadataValidationError is the validation error returned by
// Metadata.Validate if the designated constraints aren't met.
type MetadataValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MetadataValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MetadataValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MetadataValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MetadataValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MetadataValidationError) ErrorName() string { return "MetadataValidationError" }

// Error satisfies the builtin error interface
func (e MetadataValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMetadata.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MetadataValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MetadataValidationError{}

// Validate checks the field values on NamespaceDefinition with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *NamespaceDefinition) Validate() error {
	if m == nil {
		return nil
	}

	if len(m.GetName()) > 128 {
		return NamespaceDefinitionValidationError{
			field:  "Name",
			reason: "value length must be at most 128 bytes",
		}
	}

	if !_NamespaceDefinition_Name_Pattern.MatchString(m.GetName()) {
		return NamespaceDefinitionValidationError{
			field:  "Name",
			reason: "value does not match regex pattern \"^([a-z][a-z0-9_]{2,62}[a-z0-9]/)?[a-z][a-z0-9_]{2,62}[a-z0-9]$\"",
		}
	}

	for idx, item := range m.GetRelation() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return NamespaceDefinitionValidationError{
					field:  fmt.Sprintf("Relation[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if v, ok := interface{}(m.GetMetadata()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return NamespaceDefinitionValidationError{
				field:  "Metadata",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// NamespaceDefinitionValidationError is the validation error returned by
// NamespaceDefinition.Validate if the designated constraints aren't met.
type NamespaceDefinitionValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e NamespaceDefinitionValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e NamespaceDefinitionValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e NamespaceDefinitionValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e NamespaceDefinitionValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e NamespaceDefinitionValidationError) ErrorName() string {
	return "NamespaceDefinitionValidationError"
}

// Error satisfies the builtin error interface
func (e NamespaceDefinitionValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sNamespaceDefinition.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = NamespaceDefinitionValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = NamespaceDefinitionValidationError{}

var _NamespaceDefinition_Name_Pattern = regexp.MustCompile("^([a-z][a-z0-9_]{2,62}[a-z0-9]/)?[a-z][a-z0-9_]{2,62}[a-z0-9]$")

// Validate checks the field values on Relation with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *Relation) Validate() error {
	if m == nil {
		return nil
	}

	if len(m.GetName()) > 64 {
		return RelationValidationError{
			field:  "Name",
			reason: "value length must be at most 64 bytes",
		}
	}

	if !_Relation_Name_Pattern.MatchString(m.GetName()) {
		return RelationValidationError{
			field:  "Name",
			reason: "value does not match regex pattern \"^[a-z][a-z0-9_]{2,62}[a-z0-9]$\"",
		}
	}

	if v, ok := interface{}(m.GetUsersetRewrite()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RelationValidationError{
				field:  "UsersetRewrite",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetTypeInformation()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RelationValidationError{
				field:  "TypeInformation",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetMetadata()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RelationValidationError{
				field:  "Metadata",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// RelationValidationError is the validation error returned by
// Relation.Validate if the designated constraints aren't met.
type RelationValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RelationValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RelationValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RelationValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RelationValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RelationValidationError) ErrorName() string { return "RelationValidationError" }

// Error satisfies the builtin error interface
func (e RelationValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRelation.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RelationValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RelationValidationError{}

var _Relation_Name_Pattern = regexp.MustCompile("^[a-z][a-z0-9_]{2,62}[a-z0-9]$")

// Validate checks the field values on TypeInformation with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *TypeInformation) Validate() error {
	if m == nil {
		return nil
	}

	for idx, item := range m.GetAllowedDirectRelations() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return TypeInformationValidationError{
					field:  fmt.Sprintf("AllowedDirectRelations[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// TypeInformationValidationError is the validation error returned by
// TypeInformation.Validate if the designated constraints aren't met.
type TypeInformationValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TypeInformationValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TypeInformationValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TypeInformationValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TypeInformationValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TypeInformationValidationError) ErrorName() string { return "TypeInformationValidationError" }

// Error satisfies the builtin error interface
func (e TypeInformationValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTypeInformation.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TypeInformationValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TypeInformationValidationError{}

// Validate checks the field values on UsersetRewrite with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *UsersetRewrite) Validate() error {
	if m == nil {
		return nil
	}

	switch m.RewriteOperation.(type) {

	case *UsersetRewrite_Union:

		if m.GetUnion() == nil {
			return UsersetRewriteValidationError{
				field:  "Union",
				reason: "value is required",
			}
		}

		if v, ok := interface{}(m.GetUnion()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UsersetRewriteValidationError{
					field:  "Union",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *UsersetRewrite_Intersection:

		if m.GetIntersection() == nil {
			return UsersetRewriteValidationError{
				field:  "Intersection",
				reason: "value is required",
			}
		}

		if v, ok := interface{}(m.GetIntersection()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UsersetRewriteValidationError{
					field:  "Intersection",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *UsersetRewrite_Exclusion:

		if m.GetExclusion() == nil {
			return UsersetRewriteValidationError{
				field:  "Exclusion",
				reason: "value is required",
			}
		}

		if v, ok := interface{}(m.GetExclusion()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UsersetRewriteValidationError{
					field:  "Exclusion",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	default:
		return UsersetRewriteValidationError{
			field:  "RewriteOperation",
			reason: "value is required",
		}

	}

	return nil
}

// UsersetRewriteValidationError is the validation error returned by
// UsersetRewrite.Validate if the designated constraints aren't met.
type UsersetRewriteValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UsersetRewriteValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UsersetRewriteValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UsersetRewriteValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UsersetRewriteValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UsersetRewriteValidationError) ErrorName() string { return "UsersetRewriteValidationError" }

// Error satisfies the builtin error interface
func (e UsersetRewriteValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUsersetRewrite.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UsersetRewriteValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UsersetRewriteValidationError{}

// Validate checks the field values on SetOperation with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *SetOperation) Validate() error {
	if m == nil {
		return nil
	}

	if len(m.GetChild()) < 1 {
		return SetOperationValidationError{
			field:  "Child",
			reason: "value must contain at least 1 item(s)",
		}
	}

	for idx, item := range m.GetChild() {
		_, _ = idx, item

		if item == nil {
			return SetOperationValidationError{
				field:  fmt.Sprintf("Child[%v]", idx),
				reason: "value is required",
			}
		}

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return SetOperationValidationError{
					field:  fmt.Sprintf("Child[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// SetOperationValidationError is the validation error returned by
// SetOperation.Validate if the designated constraints aren't met.
type SetOperationValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SetOperationValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SetOperationValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SetOperationValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SetOperationValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SetOperationValidationError) ErrorName() string { return "SetOperationValidationError" }

// Error satisfies the builtin error interface
func (e SetOperationValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSetOperation.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SetOperationValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SetOperationValidationError{}

// Validate checks the field values on TupleToUserset with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *TupleToUserset) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetTupleset() == nil {
		return TupleToUsersetValidationError{
			field:  "Tupleset",
			reason: "value is required",
		}
	}

	if v, ok := interface{}(m.GetTupleset()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TupleToUsersetValidationError{
				field:  "Tupleset",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if m.GetComputedUserset() == nil {
		return TupleToUsersetValidationError{
			field:  "ComputedUserset",
			reason: "value is required",
		}
	}

	if v, ok := interface{}(m.GetComputedUserset()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TupleToUsersetValidationError{
				field:  "ComputedUserset",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// TupleToUsersetValidationError is the validation error returned by
// TupleToUserset.Validate if the designated constraints aren't met.
type TupleToUsersetValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TupleToUsersetValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TupleToUsersetValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TupleToUsersetValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TupleToUsersetValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TupleToUsersetValidationError) ErrorName() string { return "TupleToUsersetValidationError" }

// Error satisfies the builtin error interface
func (e TupleToUsersetValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTupleToUserset.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TupleToUsersetValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TupleToUsersetValidationError{}

// Validate checks the field values on ComputedUserset with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *ComputedUserset) Validate() error {
	if m == nil {
		return nil
	}

	if _, ok := ComputedUserset_Object_name[int32(m.GetObject())]; !ok {
		return ComputedUsersetValidationError{
			field:  "Object",
			reason: "value must be one of the defined enum values",
		}
	}

	if len(m.GetRelation()) > 64 {
		return ComputedUsersetValidationError{
			field:  "Relation",
			reason: "value length must be at most 64 bytes",
		}
	}

	if !_ComputedUserset_Relation_Pattern.MatchString(m.GetRelation()) {
		return ComputedUsersetValidationError{
			field:  "Relation",
			reason: "value does not match regex pattern \"^[a-z][a-z0-9_]{2,62}[a-z0-9]$\"",
		}
	}

	return nil
}

// ComputedUsersetValidationError is the validation error returned by
// ComputedUserset.Validate if the designated constraints aren't met.
type ComputedUsersetValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ComputedUsersetValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ComputedUsersetValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ComputedUsersetValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ComputedUsersetValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ComputedUsersetValidationError) ErrorName() string { return "ComputedUsersetValidationError" }

// Error satisfies the builtin error interface
func (e ComputedUsersetValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sComputedUserset.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ComputedUsersetValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ComputedUsersetValidationError{}

var _ComputedUserset_Relation_Pattern = regexp.MustCompile("^[a-z][a-z0-9_]{2,62}[a-z0-9]$")

// Validate checks the field values on SetOperation_Child with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *SetOperation_Child) Validate() error {
	if m == nil {
		return nil
	}

	switch m.ChildType.(type) {

	case *SetOperation_Child_XThis:

		if v, ok := interface{}(m.GetXThis()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return SetOperation_ChildValidationError{
					field:  "XThis",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *SetOperation_Child_ComputedUserset:

		if m.GetComputedUserset() == nil {
			return SetOperation_ChildValidationError{
				field:  "ComputedUserset",
				reason: "value is required",
			}
		}

		if v, ok := interface{}(m.GetComputedUserset()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return SetOperation_ChildValidationError{
					field:  "ComputedUserset",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *SetOperation_Child_TupleToUserset:

		if m.GetTupleToUserset() == nil {
			return SetOperation_ChildValidationError{
				field:  "TupleToUserset",
				reason: "value is required",
			}
		}

		if v, ok := interface{}(m.GetTupleToUserset()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return SetOperation_ChildValidationError{
					field:  "TupleToUserset",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *SetOperation_Child_UsersetRewrite:

		if m.GetUsersetRewrite() == nil {
			return SetOperation_ChildValidationError{
				field:  "UsersetRewrite",
				reason: "value is required",
			}
		}

		if v, ok := interface{}(m.GetUsersetRewrite()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return SetOperation_ChildValidationError{
					field:  "UsersetRewrite",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	default:
		return SetOperation_ChildValidationError{
			field:  "ChildType",
			reason: "value is required",
		}

	}

	return nil
}

// SetOperation_ChildValidationError is the validation error returned by
// SetOperation_Child.Validate if the designated constraints aren't met.
type SetOperation_ChildValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SetOperation_ChildValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SetOperation_ChildValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SetOperation_ChildValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SetOperation_ChildValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SetOperation_ChildValidationError) ErrorName() string {
	return "SetOperation_ChildValidationError"
}

// Error satisfies the builtin error interface
func (e SetOperation_ChildValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSetOperation_Child.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SetOperation_ChildValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SetOperation_ChildValidationError{}

// Validate checks the field values on SetOperation_Child_This with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *SetOperation_Child_This) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// SetOperation_Child_ThisValidationError is the validation error returned by
// SetOperation_Child_This.Validate if the designated constraints aren't met.
type SetOperation_Child_ThisValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SetOperation_Child_ThisValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SetOperation_Child_ThisValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SetOperation_Child_ThisValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SetOperation_Child_ThisValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SetOperation_Child_ThisValidationError) ErrorName() string {
	return "SetOperation_Child_ThisValidationError"
}

// Error satisfies the builtin error interface
func (e SetOperation_Child_ThisValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSetOperation_Child_This.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SetOperation_Child_ThisValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SetOperation_Child_ThisValidationError{}

// Validate checks the field values on TupleToUserset_Tupleset with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *TupleToUserset_Tupleset) Validate() error {
	if m == nil {
		return nil
	}

	if len(m.GetRelation()) > 64 {
		return TupleToUserset_TuplesetValidationError{
			field:  "Relation",
			reason: "value length must be at most 64 bytes",
		}
	}

	if !_TupleToUserset_Tupleset_Relation_Pattern.MatchString(m.GetRelation()) {
		return TupleToUserset_TuplesetValidationError{
			field:  "Relation",
			reason: "value does not match regex pattern \"^[a-z][a-z0-9_]{2,62}[a-z0-9]$\"",
		}
	}

	return nil
}

// TupleToUserset_TuplesetValidationError is the validation error returned by
// TupleToUserset_Tupleset.Validate if the designated constraints aren't met.
type TupleToUserset_TuplesetValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TupleToUserset_TuplesetValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TupleToUserset_TuplesetValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TupleToUserset_TuplesetValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TupleToUserset_TuplesetValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TupleToUserset_TuplesetValidationError) ErrorName() string {
	return "TupleToUserset_TuplesetValidationError"
}

// Error satisfies the builtin error interface
func (e TupleToUserset_TuplesetValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTupleToUserset_Tupleset.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TupleToUserset_TuplesetValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TupleToUserset_TuplesetValidationError{}

var _TupleToUserset_Tupleset_Relation_Pattern = regexp.MustCompile("^[a-z][a-z0-9_]{2,62}[a-z0-9]$")
