/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package model

import (
	log "github.com/cihub/seelog"
	"github.com/nixgnehc/infini-framework/core/errors"
	"github.com/nixgnehc/infini-framework/core/orm"
	"github.com/nixgnehc/infini-framework/core/util"
	"time"
)

type Request struct {
	ID                 string    `gorm:"not null;unique;primary_key" json:"id" elastic_meta:"_id"`
	Url                string    `json:"url"`
	Method             string    `json:"method"`
	Body               string    `json:"body"`
	Upstream           string    `json:"upstream"`
	Response           string    `json:"response"`
	ResponseSize       int64     `json:"response_size"`
	ResponseStatusCode int       `json:"response_code"`
	Created            time.Time `json:"created"`
	Updated            time.Time `json:"updated"`
	Status             int       `json:"status"`
	Message            string    `json:"message"`
}

const Created = 1
const Ignored = 2
const ReplayedSuccess = 3
const ReplayedFailure = 4

func CreateRequest(request *Request) error {
	time := time.Now().UTC()
	request.ID = util.GetUUID()
	request.Status = Created
	request.Created = time
	request.Updated = time

	err := orm.Save(request)
	if err != nil {
		log.Error(request, ", ", err)
	}
	return err
}

func UpdateRequest(request *Request) error {
	time := time.Now().UTC()
	request.Updated = time
	if request.Url == "" {
		return errors.New("url can't be nil")
	}
	return orm.Update(request)
}

func DeleteRequest(id string) error {
	request := Request{ID: id}
	err := orm.Delete(&request)
	if err != nil {
		log.Error(id, ", ", err)
	}
	return err
}

func GetRequest(id string) (Request, error) {
	request := Request{}
	request.ID = id
	err := orm.Get(&request)
	if err != nil {
		log.Error(id, ", ", err)
	}

	if len(request.ID) == 0 || request.Updated.IsZero() {
		err = errors.New("not found," + id)
	}

	return request, err
}

func GetRequestList(from, size int, upstream string, status int) (int, []Request, error) {

	var tasks []Request
	sort := []orm.Sort{}
	sort = append(sort, orm.Sort{Field: "created", SortType: orm.DESC})
	queryO := orm.Query{Sort: &sort, From: from, Size: size}
	if upstream != "" {
		queryO.Conds = orm.And(orm.Eq("upstream", upstream))
	}

	if status >= 0 {
		queryO.Conds = orm.Combine(queryO.Conds, orm.And(orm.Eq("status", status)))
	}

	err, result := orm.Search(Request{}, &tasks, &queryO)
	if err != nil {
		log.Error(err)
		return 0, tasks, err
	}
	if result.Result != nil && tasks == nil || len(tasks) == 0 {
		convertRequest(result, &tasks)
	}
	return result.Total, tasks, err
}

func GetRequestByField(k, v string) ([]Request, error) {
	request := Request{}
	requests := []Request{}
	err, result := orm.GetBy(k, v, request, &requests)

	if err != nil {
		log.Error(k, ", ", err)
		return requests, err
	}
	if result.Result != nil && requests == nil || len(requests) == 0 {
		convertRequest(result, &requests)
	}

	return requests, err
}

func convertRequest(result orm.Result, requests *[]Request) {
	if result.Result == nil {
		return
	}

	t, ok := result.Result.([]interface{})
	if ok {
		for _, i := range t {
			js := util.ToJson(i, false)
			t := Request{}
			util.FromJson(js, &t)
			*requests = append(*requests, t)
		}
	}
}
