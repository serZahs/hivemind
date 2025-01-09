package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"syscall"
	"unsafe"
)

type rpc struct {
	name    string
	params  []rpcParam // For now, these are not pointers for serialization.
	library string
}

type rpcParamType int

const (
	C_STRING rpcParamType = iota
	INTEGER
	POINTER
	COMPOSITE // value could be []rpcParam
)

type rpcParam struct {
	param_type rpcParamType
	value      interface{}
}

type rpcHost struct {
	address string
}

func (host rpcHost) send(proc rpc) error {
	conn, err := net.Dial("tcp", host.address)
	if err != nil {
		return err
	}
	defer conn.Close()

	serialized, err := serialize(proc)
	if err != nil {
		return err
	}
	size := uint32(len(serialized))
	err = binary.Write(conn, binary.BigEndian, size)
	if err != nil {
		return err
	}
	_, err = conn.Write(serialized)
	if err != nil {
		return err
	}
	return nil
}

func callRPC(proc rpc) error {
	lib := syscall.NewLazyDLL(proc.library)
	callable := lib.NewProc(proc.name)
	actual_params := make([]uintptr, len(proc.params))
	for i, v := range proc.params {
		if v.param_type == INTEGER {
			actual_params[i] = uintptr(v.value.(int))
		} else if v.param_type == C_STRING {
			byte_value, _ := syscall.BytePtrFromString(v.value.(string))
			actual_params[i] = uintptr(unsafe.Pointer(byte_value))
		} else {
			return errors.New("found unknown RPC parameter type")
		}
	}
	callable.Call(actual_params...)
	return nil
}

func (proc *rpc) fill(args ...interface{}) error {
	if len(proc.params) != len(args) {
		return errors.New("incorrect number of arguments passed")
	}
	for i, v := range args {
		proc.params[i].value = v
	}
	return nil
}

func append_string(list []interface{}, str string) []interface{} {
	list = append(list, uint64(len(str)))
	list = append(list, []byte(str))
	return list
}

func serialize(proc rpc) ([]byte, error) {
	buffer := new(bytes.Buffer)

	var fields []interface{}
	fields = append_string(fields, proc.name)
	fields = append(fields, uint64(len(proc.params)))
	for _, v := range proc.params {
		fields = append(fields, uint64(v.param_type));
		if v.param_type == INTEGER {
			fields = append(fields, int64(v.value.(int)))
		} else if v.param_type == C_STRING {
			fields = append_string(fields, v.value.(string))
		} else {
			return nil, errors.New("encountered unknown parameter type")
		}
	}
	fields = append_string(fields, proc.library)

	for _, v := range fields {
		if err := binary.Write(buffer, binary.LittleEndian, v); err != nil {
			return nil, err
		}
	}
	return buffer.Bytes(), nil
}

func deserialize(data []byte) (rpc, error) {
	reader := bytes.NewReader(data)
	var name_size uint64
	if err := binary.Read(reader, binary.LittleEndian, &name_size); err != nil {
		return rpc{}, err
	}
	name_bytes := make([]byte, name_size)
	if err := binary.Read(reader, binary.LittleEndian, &name_bytes); err != nil {
		return rpc{}, err
	}

	var params_size uint64
	if err := binary.Read(reader, binary.LittleEndian, &params_size); err != nil {
		return rpc{}, err
	}
	params := make([]rpcParam, params_size)
	for i := 0; i < int(params_size); i++ {
		var param rpcParam
		var param_type uint64
		if err := binary.Read(reader, binary.LittleEndian, &param_type); err != nil {
			return rpc{}, err
		}
		param.param_type = rpcParamType(param_type)
		if param.param_type == INTEGER {
			var value int64
			if err := binary.Read(reader, binary.LittleEndian, &value); err != nil {
				return rpc{}, err
			}
			param.value = int(value)
		} else if param.param_type == C_STRING {
			var value_size uint64
			if err := binary.Read(reader, binary.LittleEndian, &value_size); err != nil {
				return rpc{}, err
			}
			value_bytes := make([]byte, value_size)
			if err := binary.Read(reader, binary.LittleEndian, &value_bytes); err != nil {
				return rpc{}, err
			}
			param.value = string(value_bytes)
		}
		params[i] = param
	}

	var library_size uint64
	if err := binary.Read(reader, binary.LittleEndian, &library_size); err != nil {
		return rpc{}, err
	}
	library_bytes := make([]byte, library_size)
	if err := binary.Read(reader, binary.LittleEndian, &library_bytes); err != nil {
		return rpc{}, err
	}

	return rpc{string(name_bytes), params, string(library_bytes)}, nil
}
