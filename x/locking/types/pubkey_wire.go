package types

import "fmt"

// RequestKeyType maps the stored name onto the wire tag, the inverse of
// KeyTypeFromRequest. Anything building a RotateRequest needs this direction.
func RequestKeyType(keyType string) (uint8, error) {
	switch KeyTypeName(keyType) {
	case KeyTypeSecp256k1:
		return WireKeyTypeSecp256k1, nil
	case KeyTypeMlDsa65:
		return WireKeyTypeMlDsa65, nil
	default:
		return 0, fmt.Errorf("unsupported consensus key type %q", keyType)
	}
}
