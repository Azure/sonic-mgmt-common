////////////////////////////////////////////////////////////////////////////////
//                                                                            //
//  Copyright 2019 Dell, Inc.                                                 //
//                                                                            //
//  Licensed under the Apache License, Version 2.0 (the "License");           //
//  you may not use this file except in compliance with the License.          //
//  You may obtain a copy of the License at                                   //
//                                                                            //
//  http://www.apache.org/licenses/LICENSE-2.0                                //
//                                                                            //
//  Unless required by applicable law or agreed to in writing, software       //
//  distributed under the License is distributed on an "AS IS" BASIS,         //
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  //
//  See the License for the specific language governing permissions and       //
//  limitations under the License.                                            //
//                                                                            //
////////////////////////////////////////////////////////////////////////////////

package transformer

import (
	"github.com/openconfig/ygot/ygot"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"sync"
)

type RedisDbMap = map[db.DBNum]map[string]map[string]db.Value

// XfmrParams represents input parameters for table-transformer, key-transformer, field-transformer & subtree-transformer
type XfmrParams struct {
	d *db.DB
	dbs [db.MaxDB]*db.DB
	curDb db.DBNum
	ygRoot *ygot.GoStruct
	uri string
	requestUri string //original uri using which a curl/NBI request is made
	oper int
	table string
	key string
	dbDataMap *map[db.DBNum]map[string]map[string]db.Value
	subOpDataMap map[int]*RedisDbMap // used to add an in-flight data with a sub-op
	param interface{}
	txCache *sync.Map
	skipOrdTblChk *bool
	isVirtualTbl *bool
    pCascadeDelTbl *[] string //used to populate list of tables needed cascade delete by subtree overloaded methods
    yangDefValMap map[string]map[string]db.Value
}

// SubscProcType represents subcription process type identifying the type of subscription request made from translib.
type SubscProcType int
const (
    TRANSLATE_SUBSCRIBE SubscProcType = iota
    PROCESS_SUBSCRIBE
)

/*susbcription sampling interval and subscription preference type*/
type notificationOpts struct {
	mInterval int
	pType     NotificationType
}

// XfmrSubscInParams represents input to subscribe subtree callbacks - request uri, DBs info access-pointers, DB info for request uri and subscription process type from translib. 
type XfmrSubscInParams struct {
    uri string
    dbs [db.MaxDB]*db.DB
    dbDataMap RedisDbMap
    subscProc SubscProcType
}

// XfmrSubscOutParams represents output from subscribe subtree callback - DB data for request uri, Need cache, OnChange, subscription preference and interval.
type XfmrSubscOutParams struct {
    dbDataMap RedisDbMap
    needCache bool
    onChange bool
    nOpts *notificationOpts  //these can be set regardless of error 
    isVirtualTbl bool //used for RFC parent table check, set to true when no Redis Mapping
}

// XfmrDbParams represents input paraDeters for value-transformer
type XfmrDbParams struct {
	oper           int
	dbNum          db.DBNum
	tableName      string
	key            string
	fieldName      string
	value          string
}


// KeyXfmrYangToDb type is defined to use for conversion of Yang key to DB Key,
// Transformer function definition.
// Param: XfmrParams structure having Database info, YgotRoot, operation, Xpath
// Return: Database keys to access db entry, error
type KeyXfmrYangToDb func (inParams XfmrParams) (string, error)

// KeyXfmrDbToYang type is defined to use for conversion of DB key to Yang key,
// Transformer function definition.
// Param: XfmrParams structure having Database info, operation, Database keys to access db entry
// Return: multi dimensional map to hold the yang key attributes of complete xpath, error */
type KeyXfmrDbToYang func (inParams XfmrParams) (map[string]interface{}, error)

// FieldXfmrYangToDb type is defined to use for conversion of yang Field to DB field
// Transformer function definition.
// Param: Database info, YgotRoot, operation, Xpath
// Return: multi dimensional map to hold the DB data, error
type FieldXfmrYangToDb func (inParams XfmrParams) (map[string]string, error)

// FieldXfmrDbtoYang type is defined to use for conversion of DB field to Yang field
// Transformer function definition.
// Param: XfmrParams structure having Database info, operation, DB data in multidimensional map, output param YgotRoot
// Return: error
type FieldXfmrDbtoYang func (inParams XfmrParams)  (map[string]interface{}, error)

// SubTreeXfmrYangToDb type is defined to use for handling the yang subtree to DB
// Transformer function definition.
// Param: XfmrParams structure having Database info, YgotRoot, operation, Xpath
// Return: multi dimensional map to hold the DB data, error
type SubTreeXfmrYangToDb func (inParams XfmrParams) (map[string]map[string]db.Value, error)

// SubTreeXfmrDbToYang type is defined to use for handling the DB to Yang subtree
// Transformer function definition.
// Param : XfmrParams structure having Database pointers, current db, operation, DB data in multidimensional map, output param YgotRoot, uri
// Return :  error
type SubTreeXfmrDbToYang func (inParams XfmrParams) (error)

// SubTreeXfmrSubscribe type is defined to use for handling subscribe(translateSubscribe & processSubscribe) subtree
// Transformer function definition.
// Param : XfmrSubscInParams structure having uri, database pointers,  subcribe process(translate/processSusbscribe), DB data in multidimensional map 
// Return :  XfmrSubscOutParams structure (db data in multiD map, needCache, pType, onChange, minInterval), error
type SubTreeXfmrSubscribe func (inParams XfmrSubscInParams) (XfmrSubscOutParams, error)

// ValidateCallpoint is used to validate a YANG node during data translation back to YANG as a response to GET
// Param : XfmrParams structure having Database pointers, current db, operation, DB data in multidimensional map, output param YgotRoot, uri
// Return :  bool
type ValidateCallpoint func (inParams XfmrParams) (bool)

// RpcCallpoint is used to invoke a callback for action
// Param : []byte input payload, dbi indices
// Return :  []byte output payload, error
type RpcCallpoint func (body []byte, dbs [db.MaxDB]*db.DB) ([]byte, error)

// PostXfmrFunc type is defined to use for handling any default handling operations required as part of the CREATE
// Transformer function definition.
// Param: XfmrParams structure having database pointers, current db, operation, DB data in multidimensional map, YgotRoot, uri
// Return: Multi dimensional map to hold the DB data Map (tblName, key and Fields), error
type PostXfmrFunc func (inParams XfmrParams) (map[string]map[string]db.Value, error)

// TableXfmrFunc type is defined to use for table transformer function for dynamic derviation of redis table.
// Param: XfmrParams structure having database pointers, current db, operation, DB data in multidimensional map, YgotRoot, uri
// Return: List of table names, error
type TableXfmrFunc func (inParams XfmrParams) ([]string, error)

// ValueXfmrFunc type is defined to use for conversion of DB field value from one forma to another
// Transformer function definition.
// Param: XfmrDbParams structure having Database info, operation, db-number, table, key, field, value
// Return: value string, error
type ValueXfmrFunc func (inParams XfmrDbParams)  (string, error)

 // PreXfmrFunc type is defined to use for handling any default handling operations required as part of the CREATE, UPDATE, REPLACE, DELETE & GET
 // Transformer function definition.
 // Param: XfmrParams structure having database pointers, current db, operation, DB data in multidimensional map, YgotRoot, uri
 // Return: error
type PreXfmrFunc func (inParams XfmrParams) (error)

// XfmrInterface is a validation interface for validating the callback registration of app modules 
// transformer methods.
type XfmrInterface interface {
    xfmrInterfaceValiidate()
}

func (KeyXfmrYangToDb) xfmrInterfaceValiidate () {
    xfmrLogInfo("xfmrInterfaceValiidate for KeyXfmrYangToDb")
}
func (KeyXfmrDbToYang) xfmrInterfaceValiidate () {
    xfmrLogInfo("xfmrInterfaceValiidate for KeyXfmrDbToYang")
}
func (FieldXfmrYangToDb) xfmrInterfaceValiidate () {
    xfmrLogInfo("xfmrInterfaceValiidate for FieldXfmrYangToDb")
}
func (FieldXfmrDbtoYang) xfmrInterfaceValiidate () {
    xfmrLogInfo("xfmrInterfaceValiidate for FieldXfmrDbtoYang")
}
func (SubTreeXfmrYangToDb) xfmrInterfaceValiidate () {
    xfmrLogInfo("xfmrInterfaceValiidate for SubTreeXfmrYangToDb")
}
func (SubTreeXfmrDbToYang) xfmrInterfaceValiidate () {
    xfmrLogInfo("xfmrInterfaceValiidate for SubTreeXfmrDbToYang")
}
func (SubTreeXfmrSubscribe) xfmrInterfaceValiidate () {
    xfmrLogInfo("xfmrInterfaceValiidate for SubTreeXfmrSubscribe")
}
func (TableXfmrFunc) xfmrInterfaceValiidate () {
    xfmrLogInfo("xfmrInterfaceValiidate for TableXfmrFunc")
}
