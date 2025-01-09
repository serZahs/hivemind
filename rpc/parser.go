package main

import (
	"bufio"
	"os"
	"io"
	"strings"
)

func parseProcedureFile(filename string) ([]rpc, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	var procedures []rpc
	for {
		library, err := reader.ReadString(' ')
		if err != nil {
			if err == io.EOF {
				return procedures, nil
			}
			return nil, err
		}
		library = strings.TrimSpace(library)
		return_type, err := reader.ReadString(' ')
		if err != nil {
			return nil, err
		}
		return_type = strings.TrimSpace(return_type)
		name, err := reader.ReadString('(')
		if err != nil {
			return nil, err
		}
		name = strings.TrimSpace(strings.TrimSuffix(name, "("))

		param_types, err := reader.ReadString(')')
		if err != nil {
			return nil, err
		}
		param_types = strings.TrimSpace(strings.TrimSuffix(param_types, ")"))
		params := strings.Split(param_types, ",")

		final_params := make([]rpcParam, len(params))
		for i, v := range params {
			var the_type rpcParamType
			v = strings.TrimSpace(v)
			if v == "i32" {
				the_type = INTEGER
			} else if v == "cstr" {
				the_type = C_STRING
			}
			final_params[i] = rpcParam{
				the_type,
				nil,
			}
		}
		procedures = append(procedures, rpc{
			name: name,
			library: library,
			params: final_params,
		})
	}
}
