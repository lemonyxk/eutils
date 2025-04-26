/**
* @program: eutils
*
* @create: 2025-04-25 23:11
**/

package types

type Config struct {
	Addresses []string `json:"addresses"`
	Username  string   `json:"user_name"`
	Password  string   `json:"password"`
}
