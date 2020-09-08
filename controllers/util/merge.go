//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package util

import (
	"encoding/json"
	"reflect"

	"k8s.io/klog"
)

// MergeCR deep merge two custom resource spec
func MergeCR(defaultCR, changedCR []byte) map[string]interface{} {
	var defaultCRDecoded map[string]interface{}
	defaultCRUnmarshalErr := json.Unmarshal(defaultCR, &defaultCRDecoded)
	if defaultCRUnmarshalErr != nil {
		klog.Error("error unmarshalling CR Template: ", defaultCRUnmarshalErr)
	}
	var changedCRDecoded map[string]interface{}
	changedCRUnmarshalErr := json.Unmarshal(changedCR, &changedCRDecoded)
	if changedCRUnmarshalErr != nil {
		klog.Error("error unmarshalling OperandConfig service spec: ", changedCRUnmarshalErr)
	}
	for key := range defaultCRDecoded {
		checkKeyBeforeMerging(key, defaultCRDecoded[key], changedCRDecoded[key], changedCRDecoded)
	}
	return changedCRDecoded
}

func checkKeyBeforeMerging(key string, defaultMap interface{}, changedMap interface{}, finalMap map[string]interface{}) {
	if !reflect.DeepEqual(defaultMap, changedMap) {
		switch defaultMap := defaultMap.(type) {
		case map[string]interface{}:
			//Check that the changed map value doesn't contain this map at all and is nil
			if changedMap == nil {
				finalMap[key] = defaultMap
			} else if _, ok := changedMap.(map[string]interface{}); ok { //Check that the changed map value is also a map[string]interface
				defaultMapRef := defaultMap
				changedMapRef := changedMap.(map[string]interface{})
				for newKey := range defaultMapRef {
					checkKeyBeforeMerging(newKey, defaultMapRef[newKey], changedMapRef[newKey], finalMap[key].(map[string]interface{}))
				}
			}
		default:
			//Check if the value was set, otherwise set it
			if changedMap == nil {
				finalMap[key] = defaultMap
			}
		}
	}
}