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
	testFilePath := filepath.Join(path.Dir(currentFile), dataDirRelPath, "storage/3f2888f8-075d-49fe-9d09-ea7e951000dc/848a1dea-5ddd-42d6-a00d-030547bffcfe/8fe61d65-6a2e-4f11-87b3-d8a3170bfd4f.txt")
	testImgPath := filepath.Join(path.Dir(currentFile), dataDirRelPath, "storage/3f2888f8-075d-49fe-9d09-ea7e951000dc/577bd676-aacb-4072-b7da-99d00ee210a4/4881bdef-06b4-4dea-8d97-6125ad242677.png")
	testThumbCropCenterPath := filepath.Join(path.Dir(currentFile), dataDirRelPath, "storage/3f2888f8-075d-49fe-9d09-ea7e951000dc/577bd676-aacb-4072-b7da-99d00ee210a4/thumbs_4881bdef-06b4-4dea-8d97-6125ad242677.png/70x50_4881bdef-06b4-4dea-8d97-6125ad242677.png")
	testThumbCropTopPath := filepath.Join(path.Dir(currentFile), dataDirRelPath, "storage/3f2888f8-075d-49fe-9d09-ea7e951000dc/577bd676-aacb-4072-b7da-99d00ee210a4/thumbs_4881bdef-06b4-4dea-8d97-6125ad242677.png/70x50t_4881bdef-06b4-4dea-8d97-6125ad242677.png")
	testThumbCropBottomPath := filepath.Join(path.Dir(currentFile), dataDirRelPath, "storage/3f2888f8-075d-49fe-9d09-ea7e951000dc/577bd676-aacb-4072-b7da-99d00ee210a4/thumbs_4881bdef-06b4-4dea-8d97-6125ad242677.png/70x50b_4881bdef-06b4-4dea-8d97-6125ad242677.png")
	testThumbFitPath := filepath.Join(path.Dir(currentFile), dataDirRelPath, "storage/3f2888f8-075d-49fe-9d09-ea7e951000dc/577bd676-aacb-4072-b7da-99d00ee210a4/thumbs_4881bdef-06b4-4dea-8d97-6125ad242677.png/70x50f_4881bdef-06b4-4dea-8d97-6125ad242677.png")
	testThumbZeroWidthPath := filepath.Join(path.Dir(currentFile), dataDirRelPath, "storage/3f2888f8-075d-49fe-9d09-ea7e951000dc/577bd676-aacb-4072-b7da-99d00ee210a4/thumbs_4881bdef-06b4-4dea-8d97-6125ad242677.png/0x50_4881bdef-06b4-4dea-8d97-6125ad242677.png")
	testThumbZeroHeightPath := filepath.Join(path.Dir(currentFile), dataDirRelPath, "storage/3f2888f8-075d-49fe-9d09-ea7e951000dc/577bd676-aacb-4072-b7da-99d00ee210a4/thumbs_4881bdef-06b4-4dea-8d97-6125ad242677.png/70x0_4881bdef-06b4-4dea-8d97-6125ad242677.png")

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
			Url:             "/api/files/missing/577bd676-aacb-4072-b7da-99d00ee210a4/4881bdef-06b4-4dea-8d97-6125ad242677.png",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "missing record",
			Method:          http.MethodGet,
			Url:             "/api/files/demo/00000000-aacb-4072-b7da-99d00ee210a4/4881bdef-06b4-4dea-8d97-6125ad242677.png",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "missing file",
			Method:          http.MethodGet,
			Url:             "/api/files/demo/577bd676-aacb-4072-b7da-99d00ee210a4/00000000-06b4-4dea-8d97-6125ad242677.png",
			ExpectedStatus:  404,
			ExpectedContent: []string{`"data":{}`},
		},
		{
			Name:            "existing image",
			Method:          http.MethodGet,
			Url:             "/api/files/demo/577bd676-aacb-4072-b7da-99d00ee210a4/4881bdef-06b4-4dea-8d97-6125ad242677.png",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testImg)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "existing image - missing thumb (should fallback to the original)",
			Method:          http.MethodGet,
			Url:             "/api/files/demo/577bd676-aacb-4072-b7da-99d00ee210a4/4881bdef-06b4-4dea-8d97-6125ad242677.png?thumb=999x999",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testImg)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "existing image - existing thumb (crop center)",
			Method:          http.MethodGet,
			Url:             "/api/files/demo/577bd676-aacb-4072-b7da-99d00ee210a4/4881bdef-06b4-4dea-8d97-6125ad242677.png?thumb=70x50",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testThumbCropCenter)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "existing image - existing thumb (crop top)",
			Method:          http.MethodGet,
			Url:             "/api/files/demo/577bd676-aacb-4072-b7da-99d00ee210a4/4881bdef-06b4-4dea-8d97-6125ad242677.png?thumb=70x50t",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testThumbCropTop)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "existing image - existing thumb (crop bottom)",
			Method:          http.MethodGet,
			Url:             "/api/files/demo/577bd676-aacb-4072-b7da-99d00ee210a4/4881bdef-06b4-4dea-8d97-6125ad242677.png?thumb=70x50b",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testThumbCropBottom)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "existing image - existing thumb (fit)",
			Method:          http.MethodGet,
			Url:             "/api/files/demo/577bd676-aacb-4072-b7da-99d00ee210a4/4881bdef-06b4-4dea-8d97-6125ad242677.png?thumb=70x50f",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testThumbFit)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "existing image - existing thumb (zero width)",
			Method:          http.MethodGet,
			Url:             "/api/files/demo/577bd676-aacb-4072-b7da-99d00ee210a4/4881bdef-06b4-4dea-8d97-6125ad242677.png?thumb=0x50",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testThumbZeroWidth)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "existing image - existing thumb (zero height)",
			Method:          http.MethodGet,
			Url:             "/api/files/demo/577bd676-aacb-4072-b7da-99d00ee210a4/4881bdef-06b4-4dea-8d97-6125ad242677.png?thumb=70x0",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testThumbZeroHeight)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
		{
			Name:            "existing non image file - thumb parameter should be ignored",
			Method:          http.MethodGet,
			Url:             "/api/files/demo/848a1dea-5ddd-42d6-a00d-030547bffcfe/8fe61d65-6a2e-4f11-87b3-d8a3170bfd4f.txt?thumb=100x100",
			ExpectedStatus:  200,
			ExpectedContent: []string{string(testFile)},
			ExpectedEvents: map[string]int{
				"OnFileDownloadRequest": 1,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
