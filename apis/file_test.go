package apis_test

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/pocketbase/pocketbase/tests"
)

func TestFileDownload(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	dataDirRelPath := "../tests/data/"

	testFilePath := filepath.Join(path.Dir(currentFile), dataDirRelPath, "storage/_pb_users_auth_/oap640cot4yru2s/test_kfd2wYLxkz.txt")
	testImgPath := filepath.Join(path.Dir(currentFile), dataDirRelPath, "storage/_pb_users_auth_/4q1xlclmfloku33/300_1SEi6Q6U72.png")
	testThumbCropCenterPath := filepath.Join(path.Dir(currentFile), dataDirRelPath, "storage/_pb_users_auth_/4q1xlclmfloku33/thumbs_300_1SEi6Q6U72.png/70x50_300_1SEi6Q6U72.png")
	testThumbCropTopPath := filepath.Join(path.Dir(currentFile), dataDirRelPath, "storage/_pb_users_auth_/4q1xlclmfloku33/thumbs_300_1SEi6Q6U72.png/70x50t_300_1SEi6Q6U72.png")
	testThumbCropBottomPath := filepath.Join(path.Dir(currentFile), dataDirRelPath, "storage/_pb_users_auth_/4q1xlclmfloku33/thumbs_300_1SEi6Q6U72.png/70x50b_300_1SEi6Q6U72.png")
	testThumbFitPath := filepath.Join(path.Dir(currentFile), dataDirRelPath, "storage/_pb_users_auth_/4q1xlclmfloku33/thumbs_300_1SEi6Q6U72.png/70x50f_300_1SEi6Q6U72.png")
	testThumbZeroWidthPath := filepath.Join(path.Dir(currentFile), dataDirRelPath, "storage/_pb_users_auth_/4q1xlclmfloku33/thumbs_300_1SEi6Q6U72.png/0x50_300_1SEi6Q6U72.png")
	testThumbZeroHeightPath := filepath.Join(path.Dir(currentFile), dataDirRelPath, "storage/_pb_users_auth_/4q1xlclmfloku33/thumbs_300_1SEi6Q6U72.png/70x0_300_1SEi6Q6U72.png")

	testFile, fileErr := os.ReadFile(testFilePath)
	if fileErr != nil {
		t.Fatal(fileErr)
	}

	testImg, imgErr := os.ReadFile(testImgPath)
	if imgErr != nil {
		t.Fatal(imgErr)
	}

	testThumbCropCenter, thumbErr := os.ReadFile(testThumbCropCenterPath)
	if thumbErr != nil {
		t.Fatal(thumbErr)
	}

	testThumbCropTop, thumbErr := os.ReadFile(testThumbCropTopPath)
	if thumbErr != nil {
		t.Fatal(thumbErr)
	}

	testThumbCropBottom, thumbErr := os.ReadFile(testThumbCropBottomPath)
	if thumbErr != nil {
		t.Fatal(thumbErr)
	}

	testThumbFit, thumbErr := os.ReadFile(testThumbFitPath)
	if thumbErr != nil {
		t.Fatal(thumbErr)
	}

	testThumbZeroWidth, thumbErr := os.ReadFile(testThumbZeroWidthPath)
	if thumbErr != nil {
		t.Fatal(thumbErr)
	}

	testThumbZeroHeight, thumbErr := os.ReadFile(testThumbZeroHeightPath)
	if thumbErr != nil {
		t.Fatal(thumbErr)
	}

	scenarios := []tests.ApiScenario{
		{
			Name:            "missing collection",
			Method:          http.MethodGet,
			Url:             "/api/files/missing/4q1xlclmfloku33/300_1SEi6Q6U72.png",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "missing record",
			Method:          http.MethodGet,
			Url:             "/api/files/_pb_users_auth_/missing/300_1SEi6Q6U72.png",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "missing file",
			Method:          http.MethodGet,
			Url:             "/api/files/_pb_users_auth_/4q1xlclmfloku33/missing.png",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "existing image",
			Method:          http.MethodGet,
			Url:             "/api/files/_pb_users_auth_/4q1xlclmfloku33/300_1SEi6Q6U72.png",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testImg)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "existing image - missing thumb (should fallback to the original)",
			Method:          http.MethodGet,
			Url:             "/api/files/_pb_users_auth_/4q1xlclmfloku33/300_1SEi6Q6U72.png?thumb=999x999",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testImg)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "existing image - existing thumb (crop center)",
			Method:          http.MethodGet,
			Url:             "/api/files/_pb_users_auth_/4q1xlclmfloku33/300_1SEi6Q6U72.png?thumb=70x50",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testThumbCropCenter)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "existing image - existing thumb (crop top)",
			Method:          http.MethodGet,
			Url:             "/api/files/_pb_users_auth_/4q1xlclmfloku33/300_1SEi6Q6U72.png?thumb=70x50t",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testThumbCropTop)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "existing image - existing thumb (crop bottom)",
			Method:          http.MethodGet,
			Url:             "/api/files/_pb_users_auth_/4q1xlclmfloku33/300_1SEi6Q6U72.png?thumb=70x50b",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testThumbCropBottom)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "existing image - existing thumb (fit)",
			Method:          http.MethodGet,
			Url:             "/api/files/_pb_users_auth_/4q1xlclmfloku33/300_1SEi6Q6U72.png?thumb=70x50f",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testThumbFit)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "existing image - existing thumb (zero width)",
			Method:          http.MethodGet,
			Url:             "/api/files/_pb_users_auth_/4q1xlclmfloku33/300_1SEi6Q6U72.png?thumb=0x50",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testThumbZeroWidth)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "existing image - existing thumb (zero height)",
			Method:          http.MethodGet,
			Url:             "/api/files/_pb_users_auth_/4q1xlclmfloku33/300_1SEi6Q6U72.png?thumb=70x0",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testThumbZeroHeight)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "existing non image file - thumb parameter should be ignored",
			Method:          http.MethodGet,
			Url:             "/api/files/_pb_users_auth_/oap640cot4yru2s/test_kfd2wYLxkz.txt?thumb=100x100",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testFile)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
	}

	for _, scenario := range scenarios {
		// clone for the HEAD test (the same as the original scenario but without body)
		head := scenario
		head.Method = http.MethodHead
		head.Name = ("(HEAD) " + scenario.Name)
		head.ExpectedContent = nil
		head.Test(t)

		// regular request test
		scenario.Test(t)
	}
}
