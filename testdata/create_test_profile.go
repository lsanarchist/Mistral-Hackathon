package main

import (
	"log"
	"os"

	"github.com/google/pprof/profile"
)

func main() {
	// Create a simple test profile
	prof := &profile.Profile{
		SampleType: []*profile.ValueType{
			{Type: "cpu", Unit: "nanoseconds"},
		},
		PeriodType: &profile.ValueType{Type: "cpu", Unit: "nanoseconds"},
		Period:     1000000, // 1ms
		Sample: []*profile.Sample{
			{
				Value: []int64{1000000000}, // 1 second
				Location: []*profile.Location{
					{
						ID:      1,
						Mapping: &profile.Mapping{ID: 1},
						Line: []profile.Line{
							{
								Function: &profile.Function{ID: 1, Name: "main.main"},
								Line:     10,
							},
						},
					},
					{
						ID:      2,
						Mapping: &profile.Mapping{ID: 2},
						Line: []profile.Line{
							{
								Function: &profile.Function{ID: 2, Name: "runtime.allocm"},
								Line:     2276,
							},
						},
					},
				},
			},
			{
				Value: []int64{500000000}, // 0.5 seconds
				Location: []*profile.Location{
					{
						ID:      3,
						Mapping: &profile.Mapping{ID: 3},
						Line: []profile.Line{
							{
								Function: &profile.Function{ID: 3, Name: "runtime.mallocgc"},
								Line:     1045,
							},
						},
					},
				},
			},
		},
		Function: []*profile.Function{
			{ID: 1, Name: "main.main"},
			{ID: 2, Name: "runtime.allocm"},
			{ID: 3, Name: "runtime.mallocgc"},
		},
		Mapping: []*profile.Mapping{
			{ID: 1},
			{ID: 2},
			{ID: 3},
		},
	}

	// Write the profile to a file
	data, err := profile.Encode(prof)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("test_cpu.pb.gz", data, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Test profile created: test_cpu.pb.gz")
}