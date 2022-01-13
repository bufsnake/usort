package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
)

// 排序URL
// 获取全部路径，同时获取路径协议-http/https
// 获取全部路径数组，排序数组，输出
func main() {
	target := flag.String("target", "", "")
	targets := flag.String("target-list", "", "")
	flag.Parse()
	probes := make(map[string]bool)
	file := make([]byte, 0)
	var err error
	if *target != "" {
		probes[*target] = true
	} else if *targets != "" {
		file, err = os.ReadFile(*targets)
		if err != nil {
			fmt.Println(err)
			return
		}
		urls := strings.Split(string(file), "\n")
		for i := 0; i < len(urls); i++ {
			urls[i] = strings.Trim(urls[i], " \t\r")
			if urls[i] == "" {
				continue
			}
			probes[urls[i]] = true
		}
	} else {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			probes[sc.Text()] = true
		}
		if err = sc.Err(); err != nil {
			fmt.Println("failed to read input:", err)
			return
		}
	}
	if len(probes) == 0 {
		flag.Usage()
		return
	}
	hostpath_schemas := make(map[string][]string)
	hostpaths := make([]string, 0)
	for urlstr, _ := range probes {
		var schema = ""
		schema, err = urlparse(urlstr)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s %s", urlstr, err)
			continue
		}
		if schema == "" {
			continue
		}
		if _, ok := hostpath_schemas[strings.TrimLeft(urlstr, schema)]; !ok {
			hostpath_schemas[strings.TrimLeft(urlstr, schema)] = make([]string, 0)
			hostpaths = append(hostpaths, strings.TrimLeft(urlstr, schema))
		}
		hostpath_schemas[strings.TrimLeft(urlstr, schema)] = append(hostpath_schemas[strings.TrimLeft(urlstr, schema)], schema)
	}
	sort.Strings(hostpaths)
	for i := 0; i < len(hostpaths); i++ {
		sort.Strings(hostpath_schemas[hostpaths[i]])
		for j := 0; j < len(hostpath_schemas[hostpaths[i]]); j++ {
			fmt.Println(hostpath_schemas[hostpaths[i]][j] + hostpaths[i])
		}
	}
}

func urlparse(urlstr string) (schema string, err error) {
	var parse *url.URL
	parse, err = url.Parse(urlstr)
	if err != nil {
		return
	}
	schema = parse.Scheme
	return
}
