//   Copyright 2009-2012 Joubin Houshyar
// 
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//    
//   http://www.apache.org/licenses/LICENSE-2.0
//    
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//
package main

import (
	"os"
	"redis"
	"log"
	"fmt"
	"reflect"
    "io/ioutil"
    "strings"
)

type error struct {
    msg     string
    cause   os.Error
}

// Checks the redis.Client interface methods against a
// a list of canonical Redis methods obtained from a
// spec file and reports the missing methods
func main() {
    specfname, e := getSpecFileName("checksync.txt");
    if e != nil {
        log.Println("error -", e)
    }
    specmap, e1 := readCommandsFromSpecFile(specfname)
    if e1 != nil {
        log.Println("error -", e1)
        return
    }

    defined, e2 := getDefinedMethods()
    if e2 != nil {
        log.Println ("error -", e2)
    }

    reportCompliance(defined, specmap)
}

// check method name map against spec'd method names
// and report undefined methods
func reportCompliance(mmap map[string]string, specms []string){
    fmt.Println("=== compliance report =========================")
    var nccnt = 0
    for _, method := range specms {
        if mmap[method] == "" {
            nccnt++
            fmt.Printf("not defined - [%d]: %s\n", nccnt, method)
        }
    }
    if nccnt == 0 {
        fmt.Printf ("client is compliant")
    }
}

// Fully reads the file named (specfile) and converts content
// to a []string.  It is expected that the file is a simple list
// of redis commands, each on a single line.
func readCommandsFromSpecFile (specfile string) ([]string, os.Error) {

    spec, e := ioutil.ReadFile(specfile)
    if e != nil {
        return nil, e
    }
    slen := len(spec)
    if spec[slen-1] == 0x0a {
        spec = spec[:slen-1]
    }
    commands := strings.Split(string(spec), "\n")

    return commands, nil
}

// Reflect over the methods defined in redis.Client
// and send back as []string (tolowercase)
func getDefinedMethods () (map[string]string, *error) {

    var mmap = map[string]string {}

	spec := redis.DefaultSpec().Db(13).Password("go-redis")

	client, e := redis.NewSynchClientWithSpec(spec)
	if e != nil {
        log.Println("ignoring - ", e)
	}
	if client == nil {
	    return mmap, &error{"client is nil", nil}
	}
	defer client.Quit()

    tc := reflect.TypeOf(client)
    nm := tc.NumMethod()

    for i := 0; i < nm; i++ {
        m := tc.Method(i)
        mname := strings.ToLower(m.Name)
        mmap[mname] = mname
    }
	return mmap, nil
}

// Reads the spec file to use for the check
// csmetafile is the name of the expected prop file and
// should contain just the name of another file.
func getSpecFileName(csmetafile string) (string, os.Error) {
    var fname string
    buff, e := ioutil.ReadFile(csmetafile)
    if e != nil {
        return fname, e
    }
    fname = string(buff)
    return fname[:len(fname)-1], nil
}