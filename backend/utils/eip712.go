package utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
)

func init() {
	// Types those are trivially valid
	for _, t := range []string{
		"address", "address[]", "bool", "bool[]", "string", "string[]",
		"bytes", "bytes[]", "int", "int[]", "uint", "uint[]",
	} {
		validPrimitiveTypes[t] = struct{}{}
	}
	// For 'bytesN', 'bytesN[]', we allow N from 1 to 32
	for n := 1; n <= 32; n++ {
		validPrimitiveTypes[fmt.Sprintf("bytes%d", n)] = struct{}{}
		validPrimitiveTypes[fmt.Sprintf("bytes%d[]", n)] = struct{}{}
	}
	// For 'intN','intN[]' and 'uintN','uintN[]' we allow N in increments of 8, from 8 up to 256
	for n := 8; n <= 256; n += 8 {
		validPrimitiveTypes[fmt.Sprintf("int%d", n)] = struct{}{}
		validPrimitiveTypes[fmt.Sprintf("int%d[]", n)] = struct{}{}
		validPrimitiveTypes[fmt.Sprintf("uint%d", n)] = struct{}{}
		validPrimitiveTypes[fmt.Sprintf("uint%d[]", n)] = struct{}{}
	}
}
func TypedDataAndHash(typedData TypedData) ([]byte, string, error) {
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, "", err
	}
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, "", err
	}
	rawData := fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash))
	return crypto.Keccak256([]byte(rawData)), rawData, nil
}

// TypedData is a type to encapsulate EIP-712 typed messages
type TypedData struct {
	Types       Types            `json:"types"`
	PrimaryType string           `json:"primaryType"`
	Domain      TypedDataDomain  `json:"domain"`
	Message     TypedDataMessage `json:"message"`
}

type Types map[string][]Type
type TypedDataDomain struct {
	Name              string                `json:"name"`
	Version           string                `json:"version"`
	ChainId           *math.HexOrDecimal256 `json:"chainId"`
	VerifyingContract string                `json:"verifyingContract"`
	Salt              string                `json:"salt"`
}

func (domain *TypedDataDomain) validate() error {
	if domain.ChainId == nil && len(domain.Name) == 0 && len(domain.Version) == 0 && len(domain.VerifyingContract) == 0 && len(domain.Salt) == 0 {
		return errors.New("domain is undefined")
	}

	return nil
}
func (domain *TypedDataDomain) Map() map[string]interface{} {
	dataMap := map[string]interface{}{}

	if domain.ChainId != nil {
		dataMap["chainId"] = domain.ChainId
	}

	if len(domain.Name) > 0 {
		dataMap["name"] = domain.Name
	}

	if len(domain.Version) > 0 {
		dataMap["version"] = domain.Version
	}

	if len(domain.VerifyingContract) > 0 {
		dataMap["verifyingContract"] = domain.VerifyingContract
	}

	if len(domain.Salt) > 0 {
		dataMap["salt"] = domain.Salt
	}
	return dataMap
}

type TypedDataMessage = map[string]interface{}

func (typedData *TypedData) HashStruct(primaryType string, data TypedDataMessage) (hexutil.Bytes, error) {
	encodedData, err := typedData.EncodeData(primaryType, data, 1)
	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(encodedData), nil
}
func (typedData *TypedData) validate() error {
	if err := typedData.Types.validate(); err != nil {
		return err
	}
	if err := typedData.Domain.validate(); err != nil {
		return err
	}
	return nil
}
func (typedData *TypedData) TypeHash(primaryType string) hexutil.Bytes {
	return crypto.Keccak256(typedData.EncodeType(primaryType))
}
func (typedData *TypedData) EncodeData(primaryType string, data map[string]interface{}, depth int) (hexutil.Bytes, error) {
	if err := typedData.validate(); err != nil {
		return nil, err
	}

	buffer := bytes.Buffer{}

	// Verify extra data
	if exp, got := len(typedData.Types[primaryType]), len(data); exp < got {
		return nil, fmt.Errorf("there is extra data provided in the message (%d < %d)", exp, got)
	}

	// Add typehash
	buffer.Write(typedData.TypeHash(primaryType))

	// Add field contents. Structs and arrays have special handlers.
	for _, field := range typedData.Types[primaryType] {
		encType := field.Type
		encValue := data[field.Name]
		if encType[len(encType)-1:] == "]" {
			encodedData, err := typedData.encodeArrayValue(encValue, encType, depth)
			if err != nil {
				return nil, err
			}
			buffer.Write(encodedData)
		} else if typedData.Types[field.Type] != nil {
			mapValue, ok := encValue.(map[string]interface{})
			if !ok {
				return nil, dataMismatchError(encType, encValue)
			}
			encodedData, err := typedData.EncodeData(field.Type, mapValue, depth+1)
			if err != nil {
				return nil, err
			}
			buffer.Write(crypto.Keccak256(encodedData))
		} else {
			byteValue, err := typedData.EncodePrimitiveValue(encType, encValue, depth)
			if err != nil {
				return nil, err
			}
			buffer.Write(byteValue)
		}
	}
	return buffer.Bytes(), nil
}
func (typedData *TypedData) encodeArrayValue(encValue interface{}, encType string, depth int) (hexutil.Bytes, error) {
	arrayValue, err := convertDataToSlice(encValue)
	if err != nil {
		return nil, dataMismatchError(encType, encValue)
	}

	arrayBuffer := new(bytes.Buffer)
	parsedType := strings.Split(encType, "[")[0]
	for _, item := range arrayValue {
		if reflect.TypeOf(item).Kind() == reflect.Slice ||
			reflect.TypeOf(item).Kind() == reflect.Array {
			encodedData, err := typedData.encodeArrayValue(item, parsedType, depth+1)
			if err != nil {
				return nil, err
			}
			arrayBuffer.Write(encodedData)
		} else {
			if typedData.Types[parsedType] != nil {
				mapValue, ok := item.(map[string]interface{})
				if !ok {
					return nil, dataMismatchError(parsedType, item)
				}
				encodedData, err := typedData.EncodeData(parsedType, mapValue, depth+1)
				if err != nil {
					return nil, err
				}
				digest := crypto.Keccak256(encodedData)
				arrayBuffer.Write(digest)
			} else {
				bytesValue, err := typedData.EncodePrimitiveValue(parsedType, item, depth)
				if err != nil {
					return nil, err
				}
				arrayBuffer.Write(bytesValue)
			}
		}
	}
	return crypto.Keccak256(arrayBuffer.Bytes()), nil
}
func convertDataToSlice(encValue interface{}) ([]interface{}, error) {
	var outEncValue []interface{}
	rv := reflect.ValueOf(encValue)
	if rv.Kind() == reflect.Slice {
		for i := 0; i < rv.Len(); i++ {
			outEncValue = append(outEncValue, rv.Index(i).Interface())
		}
	} else {
		return outEncValue, fmt.Errorf("provided data '%v' is not slice", encValue)
	}
	return outEncValue, nil
}

func (typedData *TypedData) EncodePrimitiveValue(encType string, encValue interface{}, depth int) ([]byte, error) {
	switch encType {
	case "address":
		retval := make([]byte, 32)
		switch val := encValue.(type) {
		case string:
			if common.IsHexAddress(val) {
				copy(retval[12:], common.HexToAddress(val).Bytes())
				return retval, nil
			}
		case []byte:
			if len(val) == 20 {
				copy(retval[12:], val)
				return retval, nil
			}
		case [20]byte:
			copy(retval[12:], val[:])
			return retval, nil
		}
		return nil, dataMismatchError(encType, encValue)
	case "bool":
		boolValue, ok := encValue.(bool)
		if !ok {
			return nil, dataMismatchError(encType, encValue)
		}
		if boolValue {
			return math.PaddedBigBytes(common.Big1, 32), nil
		}
		return math.PaddedBigBytes(common.Big0, 32), nil
	case "string":
		strVal, ok := encValue.(string)
		if !ok {
			return nil, dataMismatchError(encType, encValue)
		}
		return crypto.Keccak256([]byte(strVal)), nil
	case "bytes":
		bytesValue, ok := parseBytes(encValue)
		if !ok {
			return nil, dataMismatchError(encType, encValue)
		}
		return crypto.Keccak256(bytesValue), nil
	}
	if strings.HasPrefix(encType, "bytes") {
		lengthStr := strings.TrimPrefix(encType, "bytes")
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return nil, fmt.Errorf("invalid size on bytes: %v", lengthStr)
		}
		if length < 0 || length > 32 {
			return nil, fmt.Errorf("invalid size on bytes: %d", length)
		}
		if byteValue, ok := parseBytes(encValue); !ok || len(byteValue) != length {
			return nil, dataMismatchError(encType, encValue)
		} else {
			// Right-pad the bits
			dst := make([]byte, 32)
			copy(dst, byteValue)
			return dst, nil
		}
	}
	if strings.HasPrefix(encType, "int") || strings.HasPrefix(encType, "uint") {
		b, err := parseInteger(encType, encValue)
		if err != nil {
			return nil, err
		}
		return math.U256Bytes(new(big.Int).Set(b)), nil
	}
	return nil, fmt.Errorf("unrecognized type '%s'", encType)
}
func parseBytes(encType interface{}) ([]byte, bool) {
	// Handle array types.
	val := reflect.ValueOf(encType)
	if val.Kind() == reflect.Array && val.Type().Elem().Kind() == reflect.Uint8 {
		v := reflect.MakeSlice(reflect.TypeOf([]byte{}), val.Len(), val.Len())
		reflect.Copy(v, val)
		return v.Bytes(), true
	}

	switch v := encType.(type) {
	case []byte:
		return v, true
	case hexutil.Bytes:
		return v, true
	case string:
		bytes, err := hexutil.Decode(v)
		if err != nil {
			return nil, false
		}
		return bytes, true
	default:
		return nil, false
	}
}
func parseInteger(encType string, encValue interface{}) (*big.Int, error) {
	var (
		length int
		signed = strings.HasPrefix(encType, "int")
		b      *big.Int
	)
	if encType == "int" || encType == "uint" {
		length = 256
	} else {
		lengthStr := ""
		if strings.HasPrefix(encType, "uint") {
			lengthStr = strings.TrimPrefix(encType, "uint")
		} else {
			lengthStr = strings.TrimPrefix(encType, "int")
		}
		atoiSize, err := strconv.Atoi(lengthStr)
		if err != nil {
			return nil, fmt.Errorf("invalid size on integer: %v", lengthStr)
		}
		length = atoiSize
	}
	switch v := encValue.(type) {
	case *math.HexOrDecimal256:
		b = (*big.Int)(v)
	case *big.Int:
		b = v
	case string:
		var hexIntValue math.HexOrDecimal256
		if err := hexIntValue.UnmarshalText([]byte(v)); err != nil {
			return nil, err
		}
		b = (*big.Int)(&hexIntValue)
	case float64:
		// JSON parses non-strings as float64. Fail if we cannot
		// convert it losslessly
		if float64(int64(v)) == v {
			b = big.NewInt(int64(v))
		} else {
			return nil, fmt.Errorf("invalid float value %v for type %v", v, encType)
		}
	}
	if b == nil {
		return nil, fmt.Errorf("invalid integer value %v/%v for type %v", encValue, reflect.TypeOf(encValue), encType)
	}
	if b.BitLen() > length {
		return nil, fmt.Errorf("integer larger than '%v'", encType)
	}
	if !signed && b.Sign() == -1 {
		return nil, fmt.Errorf("invalid negative value for unsigned type %v", encType)
	}
	return b, nil
}
func dataMismatchError(encType string, encValue interface{}) error {
	return fmt.Errorf("provided data '%v' doesn't match type '%s'", encValue, encType)
}

func (typedData *TypedData) EncodeType(primaryType string) hexutil.Bytes {
	// Get dependencies primary first, then alphabetical
	deps := typedData.Dependencies(primaryType, []string{})
	if len(deps) > 0 {
		slicedDeps := deps[1:]
		sort.Strings(slicedDeps)
		deps = append([]string{primaryType}, slicedDeps...)
	}

	// Format as a string with fields
	var buffer bytes.Buffer
	for _, dep := range deps {
		buffer.WriteString(dep)
		buffer.WriteString("(")
		for _, obj := range typedData.Types[dep] {
			buffer.WriteString(obj.Type)
			buffer.WriteString(" ")
			buffer.WriteString(obj.Name)
			buffer.WriteString(",")
		}
		buffer.Truncate(buffer.Len() - 1)
		buffer.WriteString(")")
	}
	return buffer.Bytes()
}
func (typedData *TypedData) Dependencies(primaryType string, found []string) []string {
	primaryType = strings.Split(primaryType, "[")[0]

	if slices.Contains(found, primaryType) {
		return found
	}
	if typedData.Types[primaryType] == nil {
		return found
	}
	found = append(found, primaryType)
	for _, field := range typedData.Types[primaryType] {
		for _, dep := range typedData.Dependencies(field.Type, found) {
			if !slices.Contains(found, dep) {
				found = append(found, dep)
			}
		}
	}
	return found
}

// Type is the inner type of an EIP-712 message
type Type struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// validate checks if the types object is conformant to the specs
func (t Types) validate() error {
	for typeKey, typeArr := range t {
		if len(typeKey) == 0 {
			return errors.New("empty type key")
		}
		for i, typeObj := range typeArr {
			if len(typeObj.Type) == 0 {
				return fmt.Errorf("type %q:%d: empty Type", typeKey, i)
			}
			if len(typeObj.Name) == 0 {
				return fmt.Errorf("type %q:%d: empty Name", typeKey, i)
			}
			if typeKey == typeObj.Type {
				return fmt.Errorf("type %q cannot reference itself", typeObj.Type)
			}
			if isPrimitiveTypeValid(typeObj.Type) {
				continue
			}
			// Must be reference type
			if _, exist := t[typeObj.typeName()]; !exist {
				return fmt.Errorf("reference type %q is undefined", typeObj.Type)
			}
			if !typedDataReferenceTypeRegexp.MatchString(typeObj.Type) {
				return fmt.Errorf("unknown reference type %q", typeObj.Type)
			}
		}
	}
	return nil
}
func isPrimitiveTypeValid(primitiveType string) bool {
	input := strings.Split(primitiveType, "[")[0]
	_, ok := validPrimitiveTypes[input]
	return ok
}

var validPrimitiveTypes = map[string]struct{}{}
var typedDataReferenceTypeRegexp = regexp.MustCompile(`^[A-Za-z](\w*)(\[\d*\])*$`)

// isArray returns true if the type is a fixed or variable sized array.
// This method may return false positives, in case the Type is not a valid
// expression, e.g. "fooo[[[[".
func (t *Type) isArray() bool {
	return strings.IndexByte(t.Type, '[') > 0
}

// typeName returns the canonical name of the type. If the type is 'Person[]' or 'Person[2]', then
// this method returns 'Person'
func (t *Type) typeName() string {
	return strings.Split(t.Type, "[")[0]
}
