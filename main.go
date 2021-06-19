package main

import (
	"io/ioutil"
	"log"
	"os"

	ipa "github.com/ubccr/goipa"
	"gopkg.in/yaml.v2"
)

func main() {
	client := getIpaClient()

	groups := readGroupFromYaml()
	for _, group := range groups.Groups {
		err := ensureGroupExist(group.Name, client)
		if err != nil {
			log.Fatal(err)
		}
		groupRecord, err := client.GroupShow(group.Name)
		if err != nil {
			log.Fatal(err)
		}

		usersRemote, err := groupRecord.GetUsers()
		if err != nil {
			log.Fatal(err)
		}

		for _, userRemote := range usersRemote {
			// if userRemote is not in group.Users - remove
			userShouldBeRemoved := true
			for _, userYaml := range group.Users {
				if userYaml == userRemote {
					userShouldBeRemoved = false
				}
			}

			if userShouldBeRemoved {
				client.RemoveUserFromGroup(group.Name, userRemote)
			}
		}

		for _, userYaml := range group.Users {
			// if user is not in the group - add user to the group
			userExist, err := client.CheckUserExist(userYaml)
			if err != nil {
				log.Fatal(err)
			}
			if !userExist {
				log.Printf("user does not exist: %s, this user will be ignored", userYaml)
			}
			userIsMember, err := client.CheckUserMemberOfGroup(userYaml, group.Name)
			if err != nil {
				log.Fatal(err)
			}
			if !userIsMember {
				client.AddUserToGroup(group.Name, userYaml)
			}
		}
	}
}

func getIpaClient() *ipa.Client {
	host := os.Getenv("IPA_HOST")
	realm := os.Getenv("IPA_REALM")
	username := os.Getenv("IPA_USERNAME")
	password := os.Getenv("IPA_PASSWORD")

	client := ipa.NewClient(host, realm)

	// To get a keytab file first run kinit, then klist
	err := client.RemoteLogin(username, password)
	if err != nil {
		panic(err)
	}

	return client
}

type GroupsYaml struct {
	Groups []Group
}

type Group struct {
	Name  string
	Users []string
}

func readGroupFromYaml() GroupsYaml {
	yamlFilePath := os.Getenv("IPA_GROUPS_YAML_PATH")
	data, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	groupsYaml := GroupsYaml{}
	err = yaml.Unmarshal([]byte(data), &groupsYaml)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return groupsYaml
}

func ensureGroupExist(cn string, c *ipa.Client) error {
	groupExists, err := c.CheckGroupExist(cn)
	if err != nil {
		return err
	}

	if !groupExists {
		_, err := c.GroupAdd(cn)
		if err != nil {
			return err
		}
	}
	return nil
}
