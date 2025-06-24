package main

import "github.com/stxkxs/ok-cli/cmd"
import _ "github.com/stxkxs/ok-cli/cmd/tidy"
import _ "github.com/stxkxs/ok-cli/cmd/prep"

func main() {
	cmd.Execute()
}
