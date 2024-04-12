package types

import (
	"fmt"
	"strconv"
)

type ObjectType struct {
	name      string
	anonymous bool
	fields    map[string]RuntimeType
	defaults  map[string]interface{}
}

var objId int = 0

func (o *ObjectType) GetName() string {
	return o.name
}

func (o *ObjectType) SetName(name string) {
	o.anonymous = false
	o.name = name
}

func (o *ObjectType) IsAnonymous() bool {
	return o.anonymous
}

func (o *ObjectType) GetFieldType(name string) (RuntimeType, error) {
	v, ok := o.fields[name]
	if !ok {
		return v, fmt.Errorf("trying to access undefined field [%s]", name)
	}
	return v, nil
}

func (o *ObjectType) GetFieldDefault(name string) (interface{}, error) {
	v, ok := o.defaults[name]
	if !ok {
		return v, fmt.Errorf("trying to access undefined field [%s]", name)
	}
	return v, nil
}

func (o *ObjectType) GetDefaults() map[string]interface{} {
	return o.defaults
}

func (t *ObjectType) Match(t2 RuntimeType) error {
	t2Obj, err := ToObjectType(t2)
	if err != nil {
		return err
	}
	if t2Obj.name != t.GetName() {
		return fmt.Errorf("expected type %s, got %s", t.GetName(), t2.GetName())
	}
	return nil
}

func (t *ObjectType) AddField(name string, fieldType RuntimeType, defaultValue interface{}) error {
	_, ok := t.fields[name]
	if ok {
		return fmt.Errorf("cannot redeclare field %s", name)
	}
	t.fields[name] = fieldType
	t.defaults[name] = defaultValue
	return nil
}

func NewObjectType() *ObjectType {
	id := objId
	objId++
	return &ObjectType{
		name:      "AnonymousObject" + strconv.Itoa(id),
		anonymous: true,
		fields:    make(map[string]RuntimeType),
		defaults:  make(map[string]interface{}),
	}
}

func IsObjectType(t RuntimeType) bool {
	_, err := ToObjectType(t)
	return err == nil
}

func ToObjectType(t RuntimeType) (*ObjectType, error) {
	tObj, ok := t.(*ObjectType)
	if !ok {
		return nil, fmt.Errorf("object type expected")
	}
	return tObj, nil
}
