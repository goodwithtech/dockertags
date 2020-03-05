package dockerhub

import (
	"encoding/json"
	"io/ioutil"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/goodwithtech/dockertags/internal/types"
)

func TestScanImage(t *testing.T) {
	testcases := map[string]struct {
		filePath  string
		filterOpt types.FilterOption
		expected  types.ImageTags
	}{
		"debian page1": {
			filePath:  "./testdata/page1.json",
			filterOpt: types.FilterOption{},
			expected: types.ImageTags{
				types.ImageTag{
					Tags: []string{"unstable-slim", "unstable-20191118-slim"},
					Data: []types.TagAttr{
						{
							Os:     "linux",
							Arch:   "ppc64le",
							Digest: "sha256:af6a41a4147be448b99befff0216877a5f982caa7b94b17689ddf783dfdb13fb",
							Byte:   31299332,
						},
						{
							Os:     "linux",
							Arch:   "arm",
							Digest: "sha256:8c279d26352cbb75f7afdfdbbd8fa1503d70c21e29ed697a315f28024cea60e7",
							Byte:   25710653,
						},
						{
							Os:     "linux",
							Arch:   "386",
							Digest: "sha256:0a69ee0bb8d4bb1fab07d985813fce2f55c9c3911d30b85b1da8c556851ba98f",
							Byte:   28737447,
						},
						{
							Os:     "linux",
							Arch:   "amd64",
							Digest: "sha256:06dcd5bd1693f192a0a0b5d8451c43e8b294931fbdd4077f258ff900dd788a93",
							Byte:   27770486,
						},
						{
							Os:     "linux",
							Arch:   "arm64",
							Digest: "sha256:7c1890a916f5a4572c5941eed2952fe518f0f962ee8b51883a80fa99c3521074",
							Byte:   26670752,
						},
						{
							Os:     "linux",
							Arch:   "arm",
							Digest: "sha256:b5c3b48ea6ae08ffbabac80e1ee017d551368e9ae028493257d37984b9131f4b",
							Byte:   23567986,
						},
						{
							Os:     "linux",
							Arch:   "s390x",
							Digest: "sha256:4398583092b4bff413cb67471d2703eaead5f9e2405dbf80c4f515d7a1eb24fb",
							Byte:   26457563,
						},
					},
				},
				types.ImageTag{
					Tags: []string{"unstable", "unstable-20191118"},
					Data: []types.TagAttr{
						{
							Os:     "linux",
							Arch:   "386",
							Digest: "sha256:11e3b8f71aecd9267c1e844278afa01c2fd6aad24d7a09327ad7ff08c8451e4a",
							Byte:   52411371,
						},
						{
							Os:     "linux",
							Arch:   "ppc64le",
							Digest: "sha256:8347e968bbd6d84238eeb3badd44c565009294556fdf615b9e37c0193e7bf7b7",
							Byte:   55128391,
						},
						{
							Os:     "linux",
							Arch:   "s390x",
							Digest: "sha256:931fbc8e7e61eab27b6907d28872fb837c4cc2d105ea69d8e40790b6c04557a1",
							Byte:   49968440,
						},
						{
							Os:     "linux",
							Arch:   "arm",
							Digest: "sha256:74606051b9104909274b3aa01714484c8b149d8e77a4abde33e6db34d19761ac",
							Byte:   49263130,
						},
						{
							Os:     "linux",
							Arch:   "arm",
							Digest: "sha256:6d73fad4b05882d055376f98e5453b599ec6a20de1d1ac938420802ada9d8130",
							Byte:   47015922,
						},
						{
							Os:     "linux",
							Arch:   "arm64",
							Digest: "sha256:cf0a802cc88f3596bc16ee702d9a790683e50d71e36c7391e35af637660a28e5",
							Byte:   50259191,
						},
						{
							Os:     "linux",
							Arch:   "amd64",
							Digest: "sha256:73df2931068bb33252c3c86594da818e23ca298cc580b8b378022f9fa74f74b8",
							Byte:   51302998,
						},
					},
				},
				types.ImageTag{
					Tags: []string{"testing-slim", "testing-20191118-slim"},
					Data: []types.TagAttr{
						{
							Os:     "linux",
							Arch:   "arm64",
							Digest: "sha256:50fdf0820936f2da4f3810565dc4c202e39190ea4aed93c015e9f5515d3bcf18",
							Byte:   26662818,
						},
						{
							Os:     "linux",
							Arch:   "ppc64le",
							Digest: "sha256:dfebd9f2b9962f1daf3be5c81a2f941a203c30121ab9c4dd10b925a456024396",
							Byte:   31287771,
						},
						{
							Os:     "linux",
							Arch:   "s390x",
							Digest: "sha256:6d10a72d2f03d7e487e18ea2e50da0a3be26674d863731a4a8a49e1ed0e476bd",
							Byte:   26444802,
						},
						{
							Os:     "linux",
							Arch:   "arm",
							Digest: "sha256:e8f91639723e09cd06b67d7350941d9259bba797c8e6369391baa72140626b01",
							Byte:   23562693,
						},
						{
							Os:     "linux",
							Arch:   "amd64",
							Digest: "sha256:9493179cd3d54f466c8819a2a1bfec74f02571d4466c53b0f3a7d0d8f93c8193",
							Byte:   27764368,
						},
						{
							Os:     "linux",
							Arch:   "386",
							Digest: "sha256:a3ed7808f955565ce8306d12457b55142d00082322f200929106df0a47be6a24",
							Byte:   28735414,
						},
						{
							Os:     "linux",
							Arch:   "arm",
							Digest: "sha256:883f0ddfeb966293073cc11214ceb3e50d6468f4c4fcaae1bdd00c1fca1907b5",
							Byte:   25708037,
						},
					},
				},
				types.ImageTag{
					Tags: []string{"testing-backports"},
					Data: []types.TagAttr{
						{
							Os:     "linux",
							Arch:   "arm64",
							Digest: "sha256:2402d0314bf3027a1d1e6a8c94778224c463f5b5f6f24ba2ceb271ae6bee87a1",
							Byte:   50254334,
						},
						{
							Os:     "linux",
							Arch:   "amd64",
							Digest: "sha256:408a340acf4b20ecd24bea9e0347068f3808fef0c2ad64aebd5918b5ca29a49f",
							Byte:   51291130,
						},
						{
							Os:     "linux",
							Arch:   "386",
							Digest: "sha256:41230cb5ffe4898b071871373d109ceee3983506a4d67ce76c8f539b09728ac4",
							Byte:   52411733,
						},
						{
							Os:     "linux",
							Arch:   "arm",
							Digest: "sha256:e7808905c7c28f8855f71ecb8851e08d74755df26d62454642e2a404a3d4b9e5",
							Byte:   47005040,
						},
						{
							Os:     "linux",
							Arch:   "ppc64le",
							Digest: "sha256:a33cb9703eb32ae90f6229d7102033789a2cc971601a03a5b8661e695b1088fd",
							Byte:   55119411,
						},
						{
							Os:     "linux",
							Arch:   "s390x",
							Digest: "sha256:a7e7321022ac83f1bbebf598cb3c67cdf293b7c4c374231046e9157986efa120",
							Byte:   49940365,
						},
						{
							Os:     "linux",
							Arch:   "arm",
							Digest: "sha256:7c44df1bb8db120b456e8057ba0041dbdf0b78f8c2e66fb486721ddee789561c",
							Byte:   49268237,
						},
					},
				},
				types.ImageTag{
					Tags: []string{"testing", "testing-20191118"},
					Data: []types.TagAttr{
						{
							Os:     "linux",
							Arch:   "s390x",
							Digest: "sha256:868be1721e5c655b557ab4c391b3758158ee0a6794f8eec7f4304cdec2142733",
							Byte:   49940141,
						},
						{
							Os:     "linux",
							Arch:   "386",
							Digest: "sha256:e9137b989e1d00a08cd3a30f1668cf1cb01a4d99f47cce3b02443d8dc3e71e16",
							Byte:   52411510,
						},
						{
							Os:     "linux",
							Arch:   "arm",
							Digest: "sha256:49621f8580714781d422f6d3dbb1bcbdf1e8be44ad7cfef76ad5f8327ed8e332",
							Byte:   47004816,
						},
						{
							Os:     "linux",
							Arch:   "arm64",
							Digest: "sha256:dc428d7fb5c8a96b65ab272b0a7e5f1c47aa89cc3f1ba55a7513be8b62498259",
							Byte:   50254108,
						},
						{
							Os:     "linux",
							Arch:   "arm",
							Digest: "sha256:5e994cd9d64ce17b03708525141f7911dc2b05646247009c652b6800a71da7ad",
							Byte:   49268013,
						},
						{
							Os:     "linux",
							Arch:   "ppc64le",
							Digest: "sha256:3d07e60ad8916daca6016eb44e693b3c65b018ec2ba2ed4088892c549a561243",
							Byte:   55119187,
						},
						{
							Os:     "linux",
							Arch:   "amd64",
							Digest: "sha256:1e9ae682e239cb6b539c4d2e3945c8f6d62c87e55df24be2d6b049e0cd6efdfd",
							Byte:   51290905,
						},
					},
				},
				types.ImageTag{
					Tags: []string{"stretch-slim"},
					Data: []types.TagAttr{
						{
							Os:     "linux",
							Arch:   "arm64",
							Digest: "sha256:2e24d88457527d82b1cab0a9327eab35de3e1902691a346bc474af36462df043",
							Byte:   20385759,
						},
						{
							Os:     "linux",
							Arch:   "ppc64le",
							Digest: "sha256:91ba1e73b204bf281e2cd13a4ff673ce763303b5b1106824b3fd4becc357667f",
							Byte:   22800737,
						},
						{
							Os:     "linux",
							Arch:   "arm",
							Digest: "sha256:14bfaabb2f662fd51c251b6d27671ee78f022b9f5c0cf9bb63a0b849cd5faea8",
							Byte:   21202864,
						},
						{
							Os:     "linux",
							Arch:   "386",
							Digest: "sha256:e5763d6be644027e414d06907e214ed5aa3899ab7ac1a9d10427347eab587353",
							Byte:   23152070,
						},
						{
							Os:     "linux",
							Arch:   "s390x",
							Digest: "sha256:bef911ad7b423dcaa98615040e63ee750f8b78d97587e60b48d4a22b12e06f08",
							Byte:   22380089,
						},
						{
							Os:     "linux",
							Arch:   "arm",
							Digest: "sha256:37a60ef0aad40394224a94aaec8ea65edceb52cb8bc44c069cf931b656d94a66",
							Byte:   19311578,
						},
						{
							Os:     "linux",
							Arch:   "amd64",
							Digest: "sha256:df522a8ee081e81521d499440e2620dd39ee75434940e4890b0424d83f6332c2",
							Byte:   22524572,
						},
					},
				},
			},
		},
		"debian filter slim": {
			filePath:  "./testdata/page1.json",
			filterOpt: types.FilterOption{Contain: []string{"slim"}},
			expected: types.ImageTags{
				{
					Tags: []string{"unstable-slim", "unstable-20191118-slim"},
					Data: []types.TagAttr{
						{
							Os:     "linux",
							Arch:   "amd64",
							Digest: "sha256:06dcd5bd1693f192a0a0b5d8451c43e8b294931fbdd4077f258ff900dd788a93",
							Byte:   27770486,
						},
						{
							Os:     "linux",
							Arch:   "ppc64le",
							Digest: "sha256:af6a41a4147be448b99befff0216877a5f982caa7b94b17689ddf783dfdb13fb",
							Byte:   31299332,
						},
						{
							Os:     "linux",
							Arch:   "386",
							Digest: "sha256:0a69ee0bb8d4bb1fab07d985813fce2f55c9c3911d30b85b1da8c556851ba98f",
							Byte:   28737447,
						},
						{
							Os:     "linux",
							Arch:   "arm",
							Digest: "sha256:8c279d26352cbb75f7afdfdbbd8fa1503d70c21e29ed697a315f28024cea60e7",
							Byte:   25710653,
						},
						{
							Os:     "linux",
							Arch:   "s390x",
							Digest: "sha256:4398583092b4bff413cb67471d2703eaead5f9e2405dbf80c4f515d7a1eb24fb",
							Byte:   26457563,
						},
						{
							Os:     "linux",
							Arch:   "arm64",
							Digest: "sha256:7c1890a916f5a4572c5941eed2952fe518f0f962ee8b51883a80fa99c3521074",
							Byte:   26670752,
						},
						{
							Os:     "linux",
							Arch:   "arm",
							Digest: "sha256:b5c3b48ea6ae08ffbabac80e1ee017d551368e9ae028493257d37984b9131f4b",
							Byte:   23567986,
						},
					},
				},
				{
					Tags: []string{"testing-slim", "testing-20191118-slim"},
					Data: []types.TagAttr{
						{
							Os:     "linux",
							Arch:   "arm",
							Digest: "sha256:e8f91639723e09cd06b67d7350941d9259bba797c8e6369391baa72140626b01",
							Byte:   23562693,
						},
						{
							Os:     "linux",
							Arch:   "ppc64le",
							Digest: "sha256:dfebd9f2b9962f1daf3be5c81a2f941a203c30121ab9c4dd10b925a456024396",
							Byte:   31287771,
						},
						{
							Os:     "linux",
							Arch:   "s390x",
							Digest: "sha256:6d10a72d2f03d7e487e18ea2e50da0a3be26674d863731a4a8a49e1ed0e476bd",
							Byte:   26444802,
						},
						{
							Os:     "linux",
							Arch:   "386",
							Digest: "sha256:a3ed7808f955565ce8306d12457b55142d00082322f200929106df0a47be6a24",
							Byte:   28735414,
						},
						{
							Os:     "linux",
							Arch:   "arm",
							Digest: "sha256:883f0ddfeb966293073cc11214ceb3e50d6468f4c4fcaae1bdd00c1fca1907b5",
							Byte:   25708037,
						},
						{
							Os:     "linux",
							Arch:   "arm64",
							Digest: "sha256:50fdf0820936f2da4f3810565dc4c202e39190ea4aed93c015e9f5515d3bcf18",
							Byte:   26662818,
						},
						{
							Os:     "linux",
							Arch:   "amd64",
							Digest: "sha256:9493179cd3d54f466c8819a2a1bfec74f02571d4466c53b0f3a7d0d8f93c8193",
							Byte:   27764368,
						},
					},
				},
				{
					Tags: []string{"stretch-slim"},
					Data: []types.TagAttr{
						{
							Os:     "linux",
							Arch:   "amd64",
							Digest: "sha256:df522a8ee081e81521d499440e2620dd39ee75434940e4890b0424d83f6332c2",
							Byte:   22524572,
						},
						{
							Os:     "linux",
							Arch:   "arm64",
							Digest: "sha256:2e24d88457527d82b1cab0a9327eab35de3e1902691a346bc474af36462df043",
							Byte:   20385759,
						},
						{
							Os:     "linux",
							Arch:   "ppc64le",
							Digest: "sha256:91ba1e73b204bf281e2cd13a4ff673ce763303b5b1106824b3fd4becc357667f",
							Byte:   22800737,
						},
						{
							Os:     "linux",
							Arch:   "s390x",
							Digest: "sha256:bef911ad7b423dcaa98615040e63ee750f8b78d97587e60b48d4a22b12e06f08",
							Byte:   22380089,
						},
						{
							Os:     "linux",
							Arch:   "arm",
							Digest: "sha256:37a60ef0aad40394224a94aaec8ea65edceb52cb8bc44c069cf931b656d94a66",
							Byte:   19311578,
						},
						{
							Os:     "linux",
							Arch:   "arm",
							Digest: "sha256:14bfaabb2f662fd51c251b6d27671ee78f022b9f5c0cf9bb63a0b849cd5faea8",
							Byte:   21202864,
						},
						{
							Os:     "linux",
							Arch:   "386",
							Digest: "sha256:e5763d6be644027e414d06907e214ed5aa3899ab7ac1a9d10427347eab587353",
							Byte:   23152070,
						},
					},
				},
			},
		},
	}

	for tc, v := range testcases {
		dockerHub := DockerHub{filterOpt: &v.filterOpt}
		var data tagsResponse
		file, err := ioutil.ReadFile(v.filePath)
		if err != nil {
			t.Errorf("readfile error: %w", err)
			continue
		}
		json.Unmarshal(file, &data)
		actual := dockerHub.convertResultToTag(data.Results)
		opts := []cmp.Option{
			cmp.Transformer("Sort", func(in []string) []string {
				out := append([]string{}, in...) // Copy input to avoid mutating it
				sort.Strings(out)
				return out
			}),
			cmp.Transformer("Sort", func(in []types.TagAttr) []types.TagAttr {
				out := append([]types.TagAttr{}, in...) // Copy input to avoid mutating it
				sort.Sort(types.TagAttrs(out))
				return out
			}),

			cmpopts.IgnoreFields(types.ImageTag{}, "UploadedAt"),
		}
		sort.Sort(actual)
		if diff := cmp.Diff(v.expected, actual, opts...); diff != "" {
			t.Errorf("%s: diff %v", tc, diff)
		}
	}
}

func TestSummarizeByHash(t *testing.T) {
	testcases := map[string]struct {
		tags     []tagSummary
		expected map[string]types.ImageTag
	}{
		"OK": {
			tags: []tagSummary{
				{
					Name:        "a",
					LastUpdated: "2019-12-02T00:00:00.00000Z",
					Images: []image{
						{Digest: "001", Architecture: "999test"},
						{Digest: "100", Architecture: "998test"},
						{Digest: "200", Architecture: "997test"},
						{Digest: "300", Architecture: "996test"},
						{Digest: "400", Architecture: "995test"},
					},
				},
				{
					Name:        "b",
					LastUpdated: "2019-12-01T00:00:00.00000Z",
					Images: []image{
						{Digest: "400b", Architecture: "995test"},
						{Digest: "001b", Architecture: "999test"},
						{Digest: "100b", Architecture: "998test"},
						{Digest: "200b", Architecture: "997test"},
						{Digest: "300b", Architecture: "996test"},
					},
				},
				{
					Name:        "c",
					LastUpdated: "2019-12-03T00:00:00.00000Z",
					Images: []image{
						{Digest: "400", Architecture: "995test"},
						{Digest: "300", Architecture: "996test"},
						{Digest: "001", Architecture: "999test"},
						{Digest: "100", Architecture: "998test"},
						{Digest: "200", Architecture: "997test"},
					},
				},
			},
			expected: map[string]types.ImageTag{
				"001": {
					Tags:       []string{"a", "c"},
					Data:       []types.TagAttr{{Arch: "999test", Digest: "001"}},
					UploadedAt: time.Date(2019, time.December, 3, 0, 0, 0, 0, time.UTC),
				},
				"001b": {
					Tags:       []string{"b"},
					Data:       []types.TagAttr{{Arch: "999test", Digest: "001b"}},
					UploadedAt: time.Date(2019, time.December, 1, 0, 0, 0, 0, time.UTC),
				},
				"100": {
					Tags:       []string{"a", "c"},
					Data:       []types.TagAttr{{Arch: "998test", Digest: "100"}},
					UploadedAt: time.Date(2019, time.December, 3, 0, 0, 0, 0, time.UTC),
				},
				"100b": {
					Tags:       []string{"b"},
					Data:       []types.TagAttr{{Arch: "998test", Digest: "100b"}},
					UploadedAt: time.Date(2019, time.December, 1, 0, 0, 0, 0, time.UTC),
				},
				"200": {
					Tags:       []string{"a", "c"},
					Data:       []types.TagAttr{{Arch: "997test", Digest: "200"}},
					UploadedAt: time.Date(2019, time.December, 3, 0, 0, 0, 0, time.UTC),
				},
				"200b": {
					Tags:       []string{"b"},
					Data:       []types.TagAttr{{Arch: "997test", Digest: "200b"}},
					UploadedAt: time.Date(2019, time.December, 1, 0, 0, 0, 0, time.UTC),
				},
				"300": {
					Tags:       []string{"a", "c"},
					Data:       []types.TagAttr{{Arch: "996test", Digest: "300"}},
					UploadedAt: time.Date(2019, time.December, 3, 0, 0, 0, 0, time.UTC),
				},
				"300b": {
					Tags:       []string{"b"},
					Data:       []types.TagAttr{{Arch: "996test", Digest: "300b"}},
					UploadedAt: time.Date(2019, time.December, 1, 0, 0, 0, 0, time.UTC),
				},
				"400": {
					Tags:       []string{"a", "c"},
					Data:       []types.TagAttr{{Arch: "995test", Digest: "400"}},
					UploadedAt: time.Date(2019, time.December, 3, 0, 0, 0, 0, time.UTC),
				},
				"400b": {
					Tags:       []string{"b"},
					Data:       []types.TagAttr{{Arch: "995test", Digest: "400b"}},
					UploadedAt: time.Date(2019, time.December, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		"LoadUpdatedAt": {
			tags: []tagSummary{
				{
					Name:        "a",
					LastUpdated: "2019-12-02T00:00:00.00000Z",
					Images: []image{
						{Digest: "001", Architecture: "999test"},
						{Digest: "100", Architecture: "998test"},
						{Digest: "200", Architecture: "997test"},
						{Digest: "300", Architecture: "996test"},
						{Digest: "400", Architecture: "995test"},
					},
				},
				{
					Name:        "b",
					LastUpdated: "2019-12-01T00:00:00.00000Z",
					Images: []image{
						{Digest: "400b", Architecture: "995test"},
						{Digest: "001b", Architecture: "999test"},
						{Digest: "100b", Architecture: "998test"},
						{Digest: "200b", Architecture: "997test"},
						{Digest: "300b", Architecture: "996test"},
					},
				},
				{
					Name:        "c",
					LastUpdated: "2019-12-01T00:00:00.00000Z",
					Images: []image{
						{Digest: "400", Architecture: "995test"},
						{Digest: "300", Architecture: "996test"},
						{Digest: "001", Architecture: "999test"},
						{Digest: "100", Architecture: "998test"},
						{Digest: "200", Architecture: "997test"},
					},
				},
			},
			expected: map[string]types.ImageTag{
				"001": {
					Tags:       []string{"a", "c"},
					Data:       []types.TagAttr{{Arch: "999test", Digest: "001"}},
					UploadedAt: time.Date(2019, time.December, 2, 0, 0, 0, 0, time.UTC),
				},
				"001b": {
					Tags:       []string{"b"},
					Data:       []types.TagAttr{{Arch: "999test", Digest: "001b"}},
					UploadedAt: time.Date(2019, time.December, 1, 0, 0, 0, 0, time.UTC),
				},
				"100": {
					Tags:       []string{"a", "c"},
					Data:       []types.TagAttr{{Arch: "998test", Digest: "100"}},
					UploadedAt: time.Date(2019, time.December, 2, 0, 0, 0, 0, time.UTC),
				},
				"100b": {
					Tags:       []string{"b"},
					Data:       []types.TagAttr{{Arch: "998test", Digest: "100b"}},
					UploadedAt: time.Date(2019, time.December, 1, 0, 0, 0, 0, time.UTC),
				},
				"200": {
					Tags:       []string{"a", "c"},
					Data:       []types.TagAttr{{Arch: "997test", Digest: "200"}},
					UploadedAt: time.Date(2019, time.December, 2, 0, 0, 0, 0, time.UTC),
				},
				"200b": {
					Tags:       []string{"b"},
					Data:       []types.TagAttr{{Arch: "997test", Digest: "200b"}},
					UploadedAt: time.Date(2019, time.December, 1, 0, 0, 0, 0, time.UTC),
				},
				"300": {
					Tags:       []string{"a", "c"},
					Data:       []types.TagAttr{{Arch: "996test", Digest: "300"}},
					UploadedAt: time.Date(2019, time.December, 2, 0, 0, 0, 0, time.UTC),
				},
				"300b": {
					Tags:       []string{"b"},
					Data:       []types.TagAttr{{Arch: "996test", Digest: "300b"}},
					UploadedAt: time.Date(2019, time.December, 1, 0, 0, 0, 0, time.UTC),
				},
				"400": {
					Tags:       []string{"a", "c"},
					Data:       []types.TagAttr{{Arch: "995test", Digest: "400"}},
					UploadedAt: time.Date(2019, time.December, 2, 0, 0, 0, 0, time.UTC),
				},
				"400b": {
					Tags:       []string{"b"},
					Data:       []types.TagAttr{{Arch: "995test", Digest: "400b"}},
					UploadedAt: time.Date(2019, time.December, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}
	for tc, v := range testcases {
		actual := summarizeByHash(v.tags)
		if diff := cmp.Diff(v.expected, actual); diff != "" {
			t.Errorf("%s: diff %v", tc, diff)
		}
	}
}
