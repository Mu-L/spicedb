// Code generated by "stringer -type=NodeType -output zz_generated.nodetype_string.go"; DO NOT EDIT.

package dslshape

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NodeTypeError-0]
	_ = x[NodeTypeFile-1]
	_ = x[NodeTypeComment-2]
	_ = x[NodeTypeDefinition-3]
	_ = x[NodeTypeCaveatDefinition-4]
	_ = x[NodeTypeCaveatParameter-5]
	_ = x[NodeTypeCaveatExpression-6]
	_ = x[NodeTypeRelation-7]
	_ = x[NodeTypePermission-8]
	_ = x[NodeTypeTypeReference-9]
	_ = x[NodeTypeSpecificTypeReference-10]
	_ = x[NodeTypeCaveatReference-11]
	_ = x[NodeTypeUnionExpression-12]
	_ = x[NodeTypeIntersectExpression-13]
	_ = x[NodeTypeExclusionExpression-14]
	_ = x[NodeTypeArrowExpression-15]
	_ = x[NodeTypeIdentifier-16]
	_ = x[NodeTypeNilExpression-17]
	_ = x[NodeTypeCaveatTypeReference-18]
	_ = x[NodeTypeImport-19]
}

const _NodeType_name = "NodeTypeErrorNodeTypeFileNodeTypeCommentNodeTypeDefinitionNodeTypeCaveatDefinitionNodeTypeCaveatParameterNodeTypeCaveatExpressionNodeTypeRelationNodeTypePermissionNodeTypeTypeReferenceNodeTypeSpecificTypeReferenceNodeTypeCaveatReferenceNodeTypeUnionExpressionNodeTypeIntersectExpressionNodeTypeExclusionExpressionNodeTypeArrowExpressionNodeTypeIdentifierNodeTypeNilExpressionNodeTypeCaveatTypeReferenceNodeTypeImport"

var _NodeType_index = [...]uint16{0, 13, 25, 40, 58, 82, 105, 129, 145, 163, 184, 213, 236, 259, 286, 313, 336, 354, 375, 402, 416}

func (i NodeType) String() string {
	if i < 0 || i >= NodeType(len(_NodeType_index)-1) {
		return "NodeType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _NodeType_name[_NodeType_index[i]:_NodeType_index[i+1]]
}
