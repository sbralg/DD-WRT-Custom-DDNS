package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lightsail"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("usage: lightsail-update domainname hostname ip")
		os.Exit(64)
	}

	svc := lightsail.New(session.New())

	inputGetDomain := &lightsail.GetDomainInput{
		DomainName: aws.String(os.Args[1]),
	}

	fmt.Printf("Collecting list of DNS entries from domain \"%s\" ...", os.Args[1])
	resultGetDomain, err := svc.GetDomain(inputGetDomain)
	if err != nil {
		fmt.Printf(" FAILED\n")
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Printf(" Done!\n")
	fmt.Printf("Looking for entry with name \"%s\" and type \"A\" ...", os.Args[2])

	var entryId string
	for _, entry := range resultGetDomain.Domain.DomainEntries {
		if *entry.Name == os.Args[2] && *entry.Type == "A" {
			fmt.Printf(" Done!\n")
			entryId = *entry.Id
			break
		}
	}

	if entryId != "" {
		inputUpdateDomain := &lightsail.UpdateDomainEntryInput{
			DomainName: aws.String(os.Args[1]),
			DomainEntry: &lightsail.DomainEntry{
				Id:      aws.String(entryId),
				IsAlias: aws.Bool(false),
				Name:    aws.String(os.Args[2]),
				Target:  aws.String(os.Args[3]),
				Type:    aws.String("A"),
			},
		}

		fmt.Printf("Updating entry ...")
		_, err := svc.UpdateDomainEntry(inputUpdateDomain)
		if err != nil {
			fmt.Printf(" FAILED\n")
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		} else {
			fmt.Printf(" Done!\n\n")
			fmt.Printf("Entry \"%s\" of type \"A\" updated with IP address \"%s\"\n\nUpdate Successful!\n", os.Args[2], os.Args[3])
			os.Exit(0)
		}
	} else {
		fmt.Printf(" FAILED\n\n")
		fmt.Println(resultGetDomain)
		fmt.Printf("\nEntry \"%s\" of type \"A\" not found!\n\nUPDATE FAILED!\n", os.Args[2])
		os.Exit(1)
	}
}