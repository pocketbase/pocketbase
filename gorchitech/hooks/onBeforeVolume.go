package hooks

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	Gtype "github.com/pocketbase/pocketbase/gorchitech/types"
	"github.com/pocketbase/pocketbase/models"
	"log"
	"strings"
)

func SetupOnBeforeVolumeHooks(app *pocketbase.PocketBase) {
	onBeforeCreateVolume(app)
	onBeforeDeleteVolume(app)
	onBeforeUpdateVolume(app)
}

func onBeforeCreateVolume(app *pocketbase.PocketBase) {
	app.OnModelBeforeCreate("volumen").Add(func(e *core.ModelEvent) error {
		log.Println("onBeforeCreateVolume", e.Dao, e.Model, e.Tags())

		gVolumen := Gtype.GVolumenRecord{}
		log.Println("on  FromRecordToVolumen")
		err := gVolumen.FromRecordToVolumen(e.Model.(*models.Record))

		log.Println("on  Create CLI docker")
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}
		defer cli.Close()

		log.Println("on  Create network Labels", gVolumen.Label)
		LabelsString := strings.Split(gVolumen.Label, ",")
		log.Println("on  Create network Labels", LabelsString)
		Labels := make(map[string]string)

		for i, label := range LabelsString {
			log.Println("on  label", label, i)
			Labels["key"+fmt.Sprintf("%d", i)] = label
		}

		log.Println("on Labels", Labels)

		log.Println("on VolumeCreate")

		options := volume.CreateOptions{
			Name:   gVolumen.Name,
			Driver: gVolumen.Driver,
			Labels: Labels,
		}

		// Create the volume
		volume, err := cli.VolumeCreate(context.Background(), options)
		if err != nil {
			panic(err)
		}
		log.Printf("Volume %s has been created with ID %s.\n", volume)
		gVolumen.VolumenId = volume.Name
		Gtype.FromGVolumenToRecord(e.Model.(*models.Record), &gVolumen)
		log.Println("on  FromGVolumenToRecord")
		if err != nil {
			panic(err)
		}
		return nil
	})
}

func onBeforeDeleteVolume(app *pocketbase.PocketBase) {
	app.OnModelBeforeDelete("volumen").Add(func(e *core.ModelEvent) error {
		log.Println("onBeforeCreateVolume", e.Dao, e.Model, e.Tags())

		gVolumen := Gtype.GVolumenRecord{}
		log.Println("on  FromRecordToVolumen")
		err := gVolumen.FromRecordToVolumen(e.Model.(*models.Record))

		log.Println("on  Create CLI volumen")
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}
		defer cli.Close()

		// Delete the volume
		err = cli.VolumeRemove(context.Background(), gVolumen.VolumenId, false)
		if err != nil {
			panic(err)
		}
		return nil
	})
}

func onBeforeUpdateVolume(app *pocketbase.PocketBase) {
	app.OnModelBeforeUpdate("volumen").Add(func(e *core.ModelEvent) error {
		log.Println("onBeforeCreateVolume", e.Dao, e.Model, e.Tags())

		gVolumen := Gtype.GVolumenRecord{}
		err := gVolumen.FromRecordToVolumen(e.Model.(*models.Record))

		log.Println("on  Create CLI volumen")
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}
		defer cli.Close()

		if gVolumen.Deleted != "" {
			// Delete the volume
			err = cli.VolumeRemove(context.Background(), gVolumen.VolumenId, false)
			if err != nil {
				panic(err)
			}
			return nil
		}

		err = cli.VolumeRemove(context.Background(), gVolumen.VolumenId, false)
		if err != nil {
			panic(err)
		}

		log.Println("on  Create network Labels", gVolumen.Label)
		LabelsString := strings.Split(gVolumen.Label, ",")
		log.Println("on  Create network Labels", LabelsString)
		Labels := make(map[string]string)

		for i, label := range LabelsString {
			log.Println("on  label", label, i)
			Labels["key"+fmt.Sprintf("%d", i)] = label
		}

		log.Println("on Labels", Labels)

		log.Println("on VolumeCreate")

		options := volume.CreateOptions{
			Name:   gVolumen.Name,
			Driver: gVolumen.Driver,
			Labels: Labels,
		}

		// Create the volume
		volume, err := cli.VolumeCreate(context.Background(), options)
		if err != nil {
			panic(err)
		}
		log.Printf("Volume %s has been created with ID %s.\n", volume)
		gVolumen.VolumenId = volume.Name
		Gtype.FromGVolumenToRecord(e.Model.(*models.Record), &gVolumen)
		log.Println("on  FromGVolumenToRecord")
		if err != nil {
			panic(err)
		}
		return nil
	})
}
