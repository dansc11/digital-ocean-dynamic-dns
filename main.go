package main

import (
	"context"
	"fmt"
	"github.com/digitalocean/godo"
	"io/ioutil"
	"net/http"
	"strings"
)

const CheckIpUrl = "http://checkip.amazonaws.com/"

const Token = ""
const Domain = ""

const RecordType = "A"
const RecordName = "@"

func main() {
	client := godo.NewFromToken(Token)
	ctx := context.Background()

	ip := externalIp()

	records, _, err := client.Domains.Records(ctx, Domain, nil)

	if err != nil {
		panic(err)
	}

	for _, value := range records {
		if value.Type == RecordType && value.Name == RecordName && value.Data != ip {
			edit := godo.DomainRecordEditRequest{
				Data: ip,
			}

			record, _, err := client.Domains.EditRecord(ctx, Domain, value.ID, &edit)

			if err != nil {
				panic(err)
			}

			fmt.Println(record)
		}
	}
}

func externalIp() string {
	resp, err := http.Get(CheckIpUrl)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	return strings.TrimSuffix(string(body), "\n")
}
