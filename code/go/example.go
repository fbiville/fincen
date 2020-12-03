// go mod init main
// go run example.go
package main
import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j" //Go 1.8
)
func main() {
	s, err := runQuery("bolt://<HOST>:<BOLTPORT>", "<USERNAME>", "<PASSWORD>")
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
}
func runQuery(uri, username, password string) ([]string, error) {
	configForNeo4j4 := func(conf *neo4j.Config) { conf.Encrypted = false }
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), configForNeo4j4)
	if err != nil {
		return nil, err
	}
	defer driver.Close()
	sessionConfig := neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead, DatabaseName: "neo4j"}
	session, err := driver.NewSession(sessionConfig)
	if err != nil {
		return nil, err
	}
	defer session.Close()
	results, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			`
			MATCH (from:Entity)<-[:ORIGINATOR]-(f:Filing)-[:BENEFITS]->(to:Entity)-[:COUNTRY]->(c:Country {name:$country})
			WITH from, to, round(sum(f.amount)) as sum
			ORDER BY sum DESC LIMIT 10
			RETURN from.name as originator
			`, map[string]interface{}{
				"country": "Russia",
			})
		if err != nil {
			return nil, err
		}
		var arr []string
		for result.Next() {
			value, found := result.Record().Get("originator")
			if found {
			  arr = append(arr, value.(string))
			}
		}
		if err = result.Err(); err != nil {
			return nil, err
		}
		return arr, nil
	})
	if err != nil {
		return nil, err
	}
	return results.([]string), err
}