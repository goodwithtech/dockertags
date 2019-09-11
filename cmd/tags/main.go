package main

import (
	"fmt"
	"time"

	"github.com/goodwithtech/image-tag-sorter/pkg/types"
	"github.com/k0kubun/pp"

	"github.com/goodwithtech/image-tag-sorter/pkg/provider"
)

func main() {

	image := "goodwithtech/dockle"
	// image := "gcr.io/google_containers/hyperkube"
	// image := "335149406041.dkr.ecr.ap-northeast-1.amazonaws.com/api-streaming"
	opt := types.AuthOption{
		Timeout: time.Second * 10,
	}
	tags, err := provider.Exec(image, opt)
	if err != nil {
		fmt.Println("err", err)
	}
	pp.Print(tags)
}
